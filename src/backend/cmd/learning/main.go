package main

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

	topicRepository = repository.NewDynamoDBTopicRepository(sess, "Qriosity-Topics")

	lambda.Start(Handler)
}

func Handler(ctx context.Context, event AppSyncEvent) (json.RawMessage, error) {
	switch event.TypeName {
	case "Query":
		switch event.FieldName {
		case "getAllTopics":
			return handleGetAllTopics(ctx)
		case "getQuizes":
			return handleGetQuizes(ctx)
		case "getLessons":
			return handleGetLessons(ctx)
		case "getCourses":
			return handleGetCourses(ctx)
		case "getRoadmaps":
			return handleGetRoadmaps(ctx)
		}
	case "Mutation":
		switch event.FieldName {
		case "addTopics":
			return handleAddTopics(ctx, event.Arguments)
		case "upsertCourse":
			return handleUpsertCourse(ctx, event.Arguments)
		case "upsertRoadmap":
			return handleUpsertRoadmap(ctx, event.Arguments)
		case "courseAddedToRoadmap":
			return handleCourseAddedToRoadmap(ctx, event.Arguments)
		case "userLikedRoadmap":
			return handleUserLikedRoadmap(ctx, event.Arguments)
		}
	}

	return nil, errors.New("unhandled operation")
}

var (
	topicRepository repository.ITopicRepository
)

type AppSyncEvent struct {
	TypeName  string          `json:"parentTypeName"`
	FieldName string          `json:"fieldName"`
	Arguments json.RawMessage `json:"arguments"`
}

type AddTopicsArguments struct {
	Names []string `json:"names"`
}

func handleGetAllTopics(ctx context.Context) (json.RawMessage, error) {
	topics, err := topicRepository.GetAllTopics(ctx)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(topics)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleGetQuizes(ctx context.Context) (json.RawMessage, error) {
	panic("Not Implemented")
}

func handleGetLessons(ctx context.Context) (json.RawMessage, error) {
	panic("Not Implemented")
}

func handleGetCourses(ctx context.Context) (json.RawMessage, error) {
	panic("Not Implemented")
}

func handleGetRoadmaps(ctx context.Context) (json.RawMessage, error) {
	panic("Not Implemented")
}

func handleCourseAddedToRoadmap(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	panic("Not Implemented")
}

func handleUserLikedRoadmap(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	panic("Not Implemented")
}

func handleAddTopics(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var addTopicsArgs AddTopicsArguments
	if err := json.Unmarshal(args, &addTopicsArgs); err != nil {
		return nil, err
	}

	var topics []*domain.Topic
	for _, name := range addTopicsArgs.Names {
		topics = append(topics, &domain.Topic{Name: name})
	}

	if err := topicRepository.Insert(ctx, topics); err != nil {
		return nil, err
	}

	response, err := json.Marshal(topics)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleUpsertCourse(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	panic("Not Implemented")

}

func handleUpsertRoadmap(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	panic("Not Implemented")

}
