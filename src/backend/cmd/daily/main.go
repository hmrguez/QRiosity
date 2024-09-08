package main

import (
	"backend/internal/repository"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"math/rand"
	"os"
	"time"
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

	dailyChallengeService = services.NewDailyChallengeService()

	lambda.Start(Handler)
}

func Handler(ctx context.Context, event utils.AppSyncEvent) (json.RawMessage, error) {

	if err := utils.CheckAuthorization(ctx, event); err != nil {
		return nil, err
	}

	switch event.TypeName {
	case "Query":
		if event.FieldName == "dailyChallenge" {
			return handleDailyChallengeQuery(ctx, event.Arguments)
		}
	case "Mutation":
		if event.FieldName == "dailyChallenge" {
			return handleDailyChallengeMutation(ctx, event.Arguments)
		}
	}

	return nil, errors.New("unhandled operation")
}

var (
	userRepository        repository.IUserRepository
	dailyChallengeService services.IDailyChallengeService
)

type QueryArguments struct {
	UserID string `json:"userId"`
}

type MutationArguments struct {
	Username string `json:"username"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func handleDailyChallengeQuery(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {

	var queryArgs QueryArguments
	if err := json.Unmarshal(args, &queryArgs); err != nil {
		return nil, err
	}

	user, err := userRepository.GetUserByName(queryArgs.UserID)
	if err != nil {
		return nil, err
	}

	i := rand.Int() % len(user.Topics)
	problem, err := dailyChallengeService.GetQuestion(user.Topics[i])
	if err != nil {
		return nil, err
	}

	if problem.Categories == nil {
		problem.Categories = make([]string, 0)
	}

	response, err := json.Marshal(problem)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleDailyChallengeMutation(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var mutationArgs MutationArguments
	if err := json.Unmarshal(args, &mutationArgs); err != nil {
		return nil, err
	}

	response, err := dailyChallengeService.RateQuestion(mutationArgs.Question, mutationArgs.Answer)
	if err != nil {
		return nil, err
	}

	user, err := userRepository.GetUserByName(mutationArgs.Username)
	if err != nil {
		return nil, err
	}

	user.DailyChallengesRemaining -= 1

	if user.DailyChallengesRemaining == 0 {
		user.DailyChallengeAvailable = false
	}

	if response.Rating > 5 && user.LastDailyChallenge.Day() == time.Now().Day()-1 {
		user.DailyChallengeStreak += 1
		user.LastDailyChallenge = time.Now()
	}

	_, err = userRepository.UpsertUser(*user)
	if err != nil {
		return nil, err
	}

	response.Left = user.DailyChallengesRemaining

	resp, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
