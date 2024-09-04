package main

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"os"
	"sync"
)

var (
	userRepository    repository.IUserRepository
	topicRepository   repository.ITopicRepository
	courseRepository  repository.ICourseRepository
	roadmapRepository repository.IRoadmapRepository
	roadmapService    services.IRoadmapService
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
	roadmapRepository = repository.NewDynamoDBRoadmapRepository(sess, "Qriosity-Roadmaps", topicRepository)
	userRepository, _ = repository.NewDynamoDBUserRepository(sess, "Qriosity-Users")
	roadmapService = services.NewRoadmapService()

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
			return handleGetCourses(ctx, event.Arguments)
		case "getRoadmaps":
			return handleGetRoadmaps(ctx)
		case "getRoadmapsByUser":
			return handleGetRoadmapsByUser(ctx, event.Arguments)
		case "getRoadmapFeed":
			return handleGetRoadmapFeed(ctx, event.Arguments)
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
		case "customRoadmapRequested":
			return handleCustomRoadmapRequested(ctx, event.Arguments)
		}
	}

	return nil, errors.New("unhandled operation")
}

func handleCustomRoadmapRequested(ctx context.Context, arguments json.RawMessage) (json.RawMessage, error) {
	var input struct {
		UserID string `json:"userId"`
		Prompt string `json:"prompt"`
	}

	if err := json.Unmarshal(arguments, &input); err != nil {
		return nil, err
	}

	log.Printf("User %s requested a custom roadmap for prompt %s", input.UserID, input.Prompt)

	// Fetch roadmap from the roadmap service
	roadmap, err := roadmapService.GetCustomRoadmap(input.Prompt)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	// In a go routine fetch the user, decrease its gen uses by 1 and insert it back
	go func() {
		defer wg.Done()
		user, err := userRepository.GetUserByName(input.UserID)
		if err != nil {
			log.Printf("Failed to fetch user %s: %v", input.UserID, err)
		}

		user.GenUsagesRemaining = user.GenUsagesRemaining - 1
		if _, err := userRepository.UpsertUser(*user); err != nil {
			log.Printf("Failed to update user %s: %v", input.UserID, err)
		}
	}()

	// Extract URLs from the courses in the roadmap
	urls := make([]string, len(roadmap.Courses))
	for i, course := range roadmap.Courses {
		urls[i] = course.URL
	}

	// Check which courses already exist in DynamoDB
	existingCourses, err := courseRepository.GetBulkByUrl(ctx, urls)
	if err != nil {
		return nil, err
	}

	// Create a map of existing course URLs for quick lookup
	existingCourseMap := make(map[string]bool)
	for _, course := range existingCourses {
		existingCourseMap[course.URL] = true
	}

	// Log map
	for key, value := range existingCourseMap {
		log.Printf("Key: %s, Value: %t", key, value)
	}

	var allCourses []domain.Course
	var newCourses []*domain.Course
	for _, course := range roadmap.Courses {
		if !existingCourseMap[course.URL] {
			randomGuid, err := uuid.NewRandom()
			if err != nil {
				continue
			}

			newCourse := course // Create a new instance of course
			newCourse.ID = randomGuid.String()
			newCourse.Author = "Qriosity-AI"
			newCourses = append(newCourses, &newCourse)
			allCourses = append(allCourses, newCourse)
		} else {
			// Attach ID
			for _, existingCourse := range existingCourses {
				if existingCourse.URL == course.URL {
					oldCourse := course
					oldCourse.ID = existingCourse.ID
					oldCourse.Author = existingCourse.Author
					allCourses = append(allCourses, oldCourse)
					break
				}
			}
		}
	}

	// Print newCourses
	for _, course := range newCourses {
		log.Printf("New course with id %s: %v ", course.ID, course)
	}

	// Add non-existing courses to DynamoDB
	if len(newCourses) > 0 {
		if err := courseRepository.BulkInsert(ctx, newCourses); err != nil {
			return nil, err
		}
	}

	roadmap.Courses = allCourses

	// Return the updated roadmap
	response, err := json.Marshal(roadmap)
	if err != nil {
		return nil, err
	}

	wg.Wait()

	return response, nil
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

func handleGetCourses(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var input struct {
		UserID     string `json:"userId"`
		Pagination struct {
			Page             int    `json:"page"`
			PerPage          int    `json:"perPage"`
			LastEvaluatedKey string `json:"lastEvaluatedKey,omitempty"`
		} `json:"pagination"`
	}

	if err := json.Unmarshal(args, &input); err != nil {
		log.Printf("Error unmarshalling input: %v", err)
		return nil, err
	}

	log.Printf("Input: %+v", input)

	// Assuming LastEvaluatedKey is a base64 encoded string or a simple key
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue
	if input.Pagination.LastEvaluatedKey != "" {
		// Decode the LastEvaluatedKey string back into the appropriate map structure
		decodedKey, err := base64.StdEncoding.DecodeString(input.Pagination.LastEvaluatedKey)
		if err != nil {
			log.Printf("Error decoding LastEvaluatedKey: %v", err)
			return nil, err
		}
		if err := json.Unmarshal(decodedKey, &lastEvaluatedKey); err != nil {
			log.Printf("Error unmarshalling LastEvaluatedKey: %v", err)
			return nil, err
		}
	} else {
		lastEvaluatedKey = nil
	}

	// Create a new pagination object
	parsedPagination := utils.Pagination{
		Page:             input.Pagination.Page,
		PerPage:          input.Pagination.PerPage,
		LastEvaluatedKey: lastEvaluatedKey,
	}

	// Pass the decoded LastEvaluatedKey to your repository function
	courses, pagination, err := courseRepository.GetAllCourses(ctx, parsedPagination)
	if err != nil {
		log.Printf("Error fetching courses: %v", err)
		return nil, err
	}
	log.Printf("Fetched courses: %+v", courses)
	log.Printf("Pagination: %+v", pagination)

	// If there's a LastEvaluatedKey in the pagination, re-encode it as a string
	var encodedLastEvaluatedKey string
	if pagination.LastEvaluatedKey != nil {
		marshaledKey, err := json.Marshal(pagination.LastEvaluatedKey)
		if err != nil {
			log.Printf("Error marshalling LastEvaluatedKey: %v", err)
			return nil, err
		}
		encodedLastEvaluatedKey = base64.StdEncoding.EncodeToString(marshaledKey)
	}

	output := struct {
		Courses    []*domain.Course `json:"courses"`
		Pagination struct {
			Page             int    `json:"page"`
			PerPage          int    `json:"perPage"`
			LastEvaluatedKey string `json:"lastEvaluatedKey,omitempty"`
		} `json:"pagination"`
	}{
		Courses: courses,
		Pagination: struct {
			Page             int    `json:"page"`
			PerPage          int    `json:"perPage"`
			LastEvaluatedKey string `json:"lastEvaluatedKey,omitempty"`
		}{
			Page:             pagination.Page,
			PerPage:          pagination.PerPage,
			LastEvaluatedKey: encodedLastEvaluatedKey,
		},
	}

	response, err := json.Marshal(output)
	if err != nil {
		log.Printf("Error marshalling output: %v", err)
		return nil, err
	}
	log.Printf("Response: %s", response)

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
	// Fetch the user ID and roadmap ID from the arguments
	var input struct {
		UserID    string `json:"userId"`
		RoadmapID string `json:"roadmapId"`
	}

	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	var user *domain.User
	var roadmap *domain.Roadmap

	wg.Add(2)

	// Handle user operations in a goroutine
	go func() {
		defer wg.Done()
		var err error
		user, err = userRepository.GetUserByName(input.UserID)
		if err != nil {
			errChan <- err
			return
		}
		user.Roadmaps = append(user.Roadmaps, input.RoadmapID)
		if _, err := userRepository.UpsertUser(*user); err != nil {
			errChan <- err
		}
	}()

	// Handle roadmap operations in a goroutine
	go func() {
		defer wg.Done()
		var err error
		roadmap, err = roadmapRepository.GetRoadmap(ctx, input.RoadmapID)
		if err != nil {
			errChan <- err
			return
		}
		roadmap.Likes++
		if err := roadmapRepository.UpsertRoadmap(ctx, roadmap); err != nil {
			errChan <- err
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return json.RawMessage(`{"success": true}`), nil
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

func handleGetRoadmapsByUser(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var input struct {
		UserID string `json:"userId"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	roadmaps, err := roadmapRepository.GetRoadmapsByUser(ctx, input.UserID, userRepository)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(roadmaps)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleGetRoadmapFeed(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var input struct {
		UserID string `json:"userId"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	// Fetch the user by ID
	user, err := userRepository.GetUserByName(input.UserID)
	if err != nil {
		return nil, err
	}

	// Create a set of the user's topics
	userTopics := make(map[string]struct{})
	for _, topic := range user.Topics {
		userTopics[topic] = struct{}{}
	}

	// Fetch roadmaps by each topic and ensure no duplicates
	roadmapMap := make(map[string]domain.Roadmap)
	for topic := range userTopics {
		roadmaps, err := roadmapRepository.GetByTopic(ctx, topic)
		if err != nil {
			return nil, err
		}
		for _, roadmap := range roadmaps {
			roadmapMap[roadmap.ID] = *roadmap
		}
	}

	// Convert map to slice
	var filteredRoadmaps []domain.Roadmap
	for _, roadmap := range roadmapMap {
		filteredRoadmaps = append(filteredRoadmaps, roadmap)
	}

	// Mark roadmaps liked by the user
	for i, roadmap := range filteredRoadmaps {
		for _, roadmapID := range user.Roadmaps {
			if roadmap.ID == roadmapID {
				filteredRoadmaps[i].Liked = true
				break
			}
		}
	}

	response, err := json.Marshal(filteredRoadmaps)
	if err != nil {
		return nil, err
	}

	return response, nil
}
