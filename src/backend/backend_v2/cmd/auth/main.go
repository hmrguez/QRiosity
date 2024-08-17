package main

import (
	"backend/backend_v2/internal/repository"
	"backend/internal/services"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"math/rand"
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

	dailyChallengeService = services.NewDailyChallengeService()

	lambda.Start(Handler)
}

func Handler(ctx context.Context, event AppSyncEvent) (json.RawMessage, error) {

	//log.Printf("Event %+v\n", event)
	//log.Printf("Variables %+v\n", event["variables"])
	//
	//log.Printf("Raw args: %+v\n", event["variables"])

	//var queryArgs QueryArguments
	//if err := json.Unmarshal(event["variables"].(byte), &queryArgs); err != nil {
	//	return nil, err
	//}
	//
	//log.Printf("Parsed args: %v", queryArgs)

	//return nil, nil

	//log.Printf("Event %v", event)
	//
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

type AppSyncEvent struct {
	TypeName  string          `json:"parentTypeName"`
	FieldName string          `json:"fieldName"`
	Arguments json.RawMessage `json:"arguments"`
}

type QueryArguments struct {
	UserID string `json:"userId"`
}

type MutationArguments struct {
	Username string `json:"username"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func handleDailyChallengeQuery(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {

	log.Printf("Raw args: %s", args)

	var queryArgs QueryArguments
	if err := json.Unmarshal(args, &queryArgs); err != nil {
		return nil, err
	}

	log.Printf("Parsed args: %v", queryArgs)

	user, err := userRepository.GetUserByName(queryArgs.UserID)
	if err != nil {
		return nil, err
	}

	i := rand.Int() % len(user.Topics)
	problem, err := dailyChallengeService.GetQuestion(user.Topics[i])
	if err != nil {
		return nil, err
	}
	log.Printf("Response %v", problem)

	if problem.Categories == nil {
		problem.Categories = make([]string, 0)
	}

	response, err := json.Marshal(problem)
	if err != nil {
		return nil, err
	}

	log.Printf("Response json %s", response)

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

	user.DailyChallengeAvailable = false
	_, err = userRepository.UpsertUser(*user)
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
