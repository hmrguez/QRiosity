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

var (
	topicRepository   repository.ITopicRepository
	courseRepository  repository.ICourseRepository
	roadmapRepository repository.IRoadmapRepository
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REPO_AWS_REGION")),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	topicRepository = repository.NewDynamoDBTopicRepository(sess, "Qriosity-Topics")
	courseRepository = repository.NewDynamoDBCourseRepository(sess, "Qriosity-Courses")
	roadmapRepository = repository.NewDynamoDBRoadmapRepository(sess, "Qriosity-Roadmaps")

	lambda.Start(Handler)
}

func Handler(ctx context.Context, event AppSyncEvent) (json.RawMessage, error) {
	switch event.TypeName {
	case "Query":
		switch event.FieldName {
		case "getRoadmapById":
			return handleGetRoadmapById(ctx, event.Arguments)
		case "getCourseById":
			return handleGetCourseById(ctx, event.Arguments)
		case "getAllTopics":
			return handleGetAllTopics(ctx)
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

type AppSyncEvent struct {
	TypeName  string          `json:"parentTypeName"`
	FieldName string          `json:"fieldName"`
	Arguments json.RawMessage `json:"arguments"`
}

type AddTopicsArguments struct {
	Names []string `json:"names"`
}

type CourseAddedToRoadmapArguments struct {
	CourseID  string `json:"courseId"`
	RoadmapID string `json:"roadmapId"`
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

func handleGetCourses(ctx context.Context) (json.RawMessage, error) {
	courses, err := courseRepository.GetAllCourses(ctx)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(courses)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleGetRoadmaps(ctx context.Context) (json.RawMessage, error) {
	roadmaps, err := roadmapRepository.GetAllRoadmaps(ctx)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(roadmaps)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleCourseAddedToRoadmap(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var courseAddedToRoadmapArgs CourseAddedToRoadmapArguments
	if err := json.Unmarshal(args, &courseAddedToRoadmapArgs); err != nil {
		return nil, err
	}

	roadmap, err := roadmapRepository.GetRoadmap(ctx, courseAddedToRoadmapArgs.RoadmapID)
	if err != nil {
		return nil, err
	}

	roadmap.CourseIDs = append(roadmap.CourseIDs, courseAddedToRoadmapArgs.CourseID)

	if err := roadmapRepository.UpsertRoadmap(ctx, roadmap); err != nil {
		return nil, err
	}

	return json.RawMessage(`{"success": true}`), nil
}

func handleUserLikedRoadmap(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	panic("not implemented")
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

	// Unmarshal into map[string]interface{}
	var tempMap map[string]interface{}
	if err := json.Unmarshal(args, &tempMap); err != nil {
		return nil, err
	}

	// Extract the input field
	inputData, ok := tempMap["input"]
	if !ok {
		return nil, errors.New("input field is missing")
	}

	// Marshal the input field back to JSON
	inputJSON, err := json.Marshal(inputData)
	if err != nil {
		return nil, err
	}

	// Unmarshal the input JSON into the domain.Course struct
	var course domain.Course
	if err := json.Unmarshal(inputJSON, &course); err != nil {
		return nil, err
	}

	if err := courseRepository.UpsertCourse(ctx, &course); err != nil {
		return nil, err
	}

	response, err := json.Marshal(course)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleUpsertRoadmap(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	// Unmarshal into map[string]interface{}
	var tempMap map[string]interface{}
	if err := json.Unmarshal(args, &tempMap); err != nil {
		return nil, err
	}

	// Extract the input field
	inputData, ok := tempMap["input"]
	if !ok {
		return nil, errors.New("input field is missing")
	}

	// Marshal the input field back to JSON
	inputJSON, err := json.Marshal(inputData)
	if err != nil {
		return nil, err
	}

	// Unmarshal the input JSON into the domain.Roadmap struct
	var roadmap domain.Roadmap
	if err := json.Unmarshal(inputJSON, &roadmap); err != nil {
		return nil, err
	}

	log.Println("Upserting unmarshalled roadmap ", roadmap)

	if err := roadmapRepository.UpsertRoadmap(ctx, &roadmap); err != nil {
		return nil, err
	}

	response, err := json.Marshal(roadmap)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleGetRoadmapById(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var input struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	roadmap, err := roadmapRepository.GetRoadmap(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(roadmap)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleGetCourseById(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var input struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	course, err := courseRepository.GetCourseByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(course)
	if err != nil {
		return nil, err
	}

	return response, nil
}
