package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoAuthService struct {
	ClientId   string
	UserPoolId string
	AwsRegion  string
}

func NewCognitoAuthService(clientId, userPoolId, region string) *CognitoAuthService {
	return &CognitoAuthService{
		ClientId:   clientId,
		UserPoolId: userPoolId,
		AwsRegion:  region,
	}
}

func (s CognitoAuthService) SignUp(email, username, password string) (*cognitoidentityprovider.SignUpOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.AwsRegion), // Replace with your region
	}))

	svc := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.ClientId),
		Username: aws.String(email), // Email is used as the username
		Password: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("name"), // Provide the required name.formatted attribute
				Value: aws.String(username),
			},
		},
	}

	result, err := svc.SignUp(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s CognitoAuthService) Login(email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.AwsRegion), // Replace with your region
	}))

	svc := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(email), // Email is used as the username
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String(s.ClientId),
	}

	result, err := svc.InitiateAuth(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ConfirmSignUp confirms the sign-up using the verification code sent to the user's email.
func (s CognitoAuthService) ConfirmSignUp(email, confirmationCode string) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.AwsRegion), // Replace with your region
	}))

	svc := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.ClientId),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
	}

	result, err := svc.ConfirmSignUp(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s CognitoAuthService) ResendConfirmationCode(email string) (*cognitoidentityprovider.ResendConfirmationCodeOutput, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.AwsRegion), // Replace with your region
	}))

	svc := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(s.ClientId),
		Username: aws.String(email),
	}

	result, err := svc.ResendConfirmationCode(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
