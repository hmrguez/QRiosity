package main

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"backend/internal/services"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/gomail.v2"
	"log"
	"os"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REPO_AWS_REGION")),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	userRepository, err = repository.NewDynamoDBUserRepository(sess, "Qriosity-Users")
	if err != nil {
		panic("Couldn't connect to dynamo: " + err.Error())
	}

	authService = *services.NewCognitoAuthService(
		os.Getenv("COGNITO_APP_CLIENT_ID"),
		os.Getenv("COGNITO_USER_POOL_ID"),
		os.Getenv("COGNITO_USER_POOL_REGION"),
	)

	lambda.Start(Handler)
}

func Handler(ctx context.Context, event AppSyncEvent) (json.RawMessage, error) {
	switch event.TypeName {
	case "Query":
		switch event.FieldName {
		case "login":
			return handleLogin(ctx, event.Arguments)
		case "getUsers":
			return handleGetUsers(ctx)
		case "resendConfirmationEmail":
			return handleResendConfirmationEmail(ctx, event.Arguments)
		case "getUserByName":
			return handleGetUserByName(ctx, event.Arguments)
		}
	case "Mutation":
		switch event.FieldName {
		case "register":
			return handleRegister(ctx, event.Arguments)
		case "confirmEmail":
			return handleConfirmEmail(ctx, event.Arguments)
		case "sendFeedback":
			return handleSendFeedback(ctx, event.Arguments)
		}
	}

	return nil, errors.New("unhandled operation")
}

var (
	userRepository repository.IUserRepository
	authService    services.CognitoAuthService
)

type AppSyncEvent struct {
	TypeName  string          `json:"parentTypeName"`
	FieldName string          `json:"fieldName"`
	Arguments json.RawMessage `json:"arguments"`
}

type LoginArguments struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterArguments struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Topics   []string `json:"topics"`
}

type ConfirmEmailArguments struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResendConfirmationEmailArguments struct {
	Email string `json:"email"`
}

func handleLogin(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var loginArgs LoginArguments
	if err := json.Unmarshal(args, &loginArgs); err != nil {
		return nil, err
	}

	cognitoResponse, err := authService.Login(loginArgs.Username, loginArgs.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	tokenString := cognitoResponse.AuthenticationResult.IdToken

	response, err := json.Marshal(map[string]string{"token": *tokenString})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleRegister(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var registerArgs RegisterArguments
	if err := json.Unmarshal(args, &registerArgs); err != nil {
		return nil, err
	}

	cognitoResponse, err := authService.SignUp(registerArgs.Email, registerArgs.Username, registerArgs.Password)
	if err != nil {
		return nil, err
	}

	registeredUsername := cognitoResponse.UserSub

	user := domain.User{
		Name:                    *registeredUsername,
		Email:                   registerArgs.Email,
		Topics:                  registerArgs.Topics,
		DailyChallengeAvailable: true,
	}

	user, err = userRepository.UpsertUser(user)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(map[string]string{"username": *registeredUsername})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleGetUsers(ctx context.Context) (json.RawMessage, error) {
	users := userRepository.GetUsers()

	response, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleConfirmEmail(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var confirmArgs ConfirmEmailArguments
	if err := json.Unmarshal(args, &confirmArgs); err != nil {
		return nil, err
	}

	_, err := authService.ConfirmSignUp(confirmArgs.Email, confirmArgs.Token)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(map[string]bool{"success": true})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleResendConfirmationEmail(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var resendArgs ResendConfirmationEmailArguments
	if err := json.Unmarshal(args, &resendArgs); err != nil {
		return nil, err
	}

	_, err := authService.ResendConfirmationCode(resendArgs.Email)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(map[string]bool{"success": true})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleGetUserByName(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var getUserArgs struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(args, &getUserArgs); err != nil {
		return nil, err
	}

	user, err := userRepository.GetUserByName(getUserArgs.Name)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleSendFeedback(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var feedbackArgs struct {
		Feedback string `json:"feedback"`
		From     string `json:"from"`
	}
	if err := json.Unmarshal(args, &feedbackArgs); err != nil {
		return nil, err
	}

	email := os.Getenv("EMAIL")

	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "New Feedback Received")
	m.SetBody("text/plain", "Feedback from: "+feedbackArgs.From+"\n\n"+feedbackArgs.Feedback)

	d := gomail.NewDialer("smtp.gmail.com", 587, email, os.Getenv("EMAIL_PASS"))

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email: %v", err)
		return nil, err
	}

	response, err := json.Marshal(map[string]bool{"success": true})
	if err != nil {
		return nil, err
	}

	return response, nil
}
