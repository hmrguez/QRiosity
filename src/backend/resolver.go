package main

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/99designs/gqlgen/graphql"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Resolver struct {
	userRepo    UserRepository
	problemRepo ProblemRepository
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) DailyChallenge(ctx context.Context, userID string, question string, answer string) (*ChallengeResponse, error) {
	response, err := r.problemRepo.SubmitChallengeResponse(userID, question, answer)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *mutationResolver) UpsertProblem(ctx context.Context, input ProblemInput) (*Problem, error) {
	problem := Problem{
		Question:   input.Question,
		Categories: input.Categories,
		Type:       input.Type,
	}
	upsertedProblem := r.problemRepo.UpsertProblem(problem)
	return &upsertedProblem, nil
}

func (r *mutationResolver) UpsertUser(ctx context.Context, input UserInput) (*User, error) {
	user := User{
		ID:       input.ID,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
	upsertedUser := r.userRepo.UpsertUser(user)
	return &upsertedUser, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Login(ctx context.Context, username string, password string) (*AuthPayload, error) {
	user, err := r.userRepo.GetUserByName(username)
	if err != nil || user.Password != password {
		fmt.Println(username)
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("User: %v\n", user)
		fmt.Printf("Password: %v\n", password)
		return nil, errors.New("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Token: tokenString,
		User:  user,
	}, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*User, error) {
	return r.userRepo.GetUserByID(id)
}

func (r *queryResolver) DailyChallenge(ctx context.Context, category string) (*Problem, error) {
	problem, err := r.problemRepo.GetDailyChallenge(category)
	if err != nil {
		return nil, err
	}
	return problem, nil

}

func (r *queryResolver) GetProblems(ctx context.Context) ([]*Problem, error) {
	problems := r.problemRepo.GetProblems()
	return problems, nil
}

func (r *queryResolver) GetUsers(ctx context.Context) ([]*User, error) {
	users := r.userRepo.GetUsers()
	return users, nil
}

var jwtSecret = []byte("your_secret_key")
