package repository

import (
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBRoadmapRepository struct {
	db        *dynamodb.DynamoDB
	topicRepo ITopicRepository
	tableName string
}

func NewDynamoDBRoadmapRepository(sess *session.Session, tableName string, topicRepo ITopicRepository) *DynamoDBRoadmapRepository {
	return &DynamoDBRoadmapRepository{
		db:        dynamodb.New(sess),
		topicRepo: topicRepo,
		tableName: tableName,
	}
}

func (r *DynamoDBRoadmapRepository) GetAllRoadmaps(ctx context.Context) ([]*domain.Roadmap, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var roadmaps []*domain.Roadmap
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &roadmaps)
	if err != nil {
		return nil, err
	}

	return roadmaps, nil
}

func (r *DynamoDBRoadmapRepository) UpsertRoadmap(ctx context.Context, roadmap *domain.Roadmap) error {
	// Fetch all topics
	topics, err := r.topicRepo.GetTopicsByNames(ctx, roadmap.Topics)
	if err != nil {
		return err
	}

	// Add this roadmap to the roadmapIds of each topic
	for _, topic := range topics {
		topic.RoadmapIds = append(topic.RoadmapIds, roadmap.ID)
	}

	// Update the topics in a separate goroutine
	errChan := make(chan error)
	go func() {
		errChan <- r.topicRepo.BulkWrite(ctx, topics)
	}()

	item, err := dynamodbattribute.MarshalMap(roadmap)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.db.PutItemWithContext(ctx, input)

	// Wait for the goroutine to finish and check for errors
	if err := <-errChan; err != nil {
		return err
	}

	return err
}

func (r *DynamoDBRoadmapRepository) GetRoadmap(ctx context.Context, roadmapID string) (*domain.Roadmap, error) {
	// Fetch the roadmap by ID
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(roadmapID)},
		},
	}

	result, err := r.db.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("roadmap not found")
	}

	var roadmap domain.Roadmap
	err = dynamodbattribute.UnmarshalMap(result.Item, &roadmap)
	if err != nil {
		return nil, err
	}

	// Collect all course IDs
	courseIDs := make(map[string]struct{})
	for _, courseID := range roadmap.CourseIDs {
		courseIDs[courseID] = struct{}{}
	}

	// Prepare keys for BatchGetItem
	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(courseIDs))
	for courseID := range courseIDs {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(courseID)},
		})
	}

	// Perform BatchGetItem
	batchGetInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			"Qriosity-Courses": {
				Keys: keys,
			},
		},
	}

	batchGetResult, err := r.db.BatchGetItemWithContext(ctx, batchGetInput)
	if err != nil {
		return nil, err
	}

	// Unmarshal courses
	var courses []domain.Course
	err = dynamodbattribute.UnmarshalListOfMaps(batchGetResult.Responses["Qriosity-Courses"], &courses)
	if err != nil {
		return nil, err
	}

	// Populate Courses field in roadmap
	roadmap.Courses = courses

	return &roadmap, nil
}

func (r *DynamoDBRoadmapRepository) GetRoadmapsByUser(ctx context.Context, userID string, userRepo IUserRepository) ([]*domain.Roadmap, error) {
	// Fetch the user
	user, err := userRepo.GetUserByName(userID)
	if err != nil {
		return nil, err
	}

	// Now fetch all roadmaps using a batch operation from the user.RoadmapIds
	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(user.Roadmaps))
	for _, roadmapID := range user.Roadmaps {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(roadmapID)},
		})
	}

	batchGetInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			r.tableName: {
				Keys: keys,
			},
		},
	}

	batchGetResult, err := r.db.BatchGetItemWithContext(ctx, batchGetInput)
	if err != nil {
		return nil, nil
	}

	// Unmarshal roadmaps
	var roadmaps []*domain.Roadmap
	err = dynamodbattribute.UnmarshalListOfMaps(batchGetResult.Responses[r.tableName], &roadmaps)
	if err != nil {
		return nil, err
	}

	return roadmaps, nil
}

func (r *DynamoDBRoadmapRepository) GetByTopic(ctx context.Context, topic string) ([]domain.Roadmap, error) {
	// Use topic repo to fetch by name
	topics, err := r.topicRepo.GetTopicsByNames(ctx, []string{topic})
	if err != nil {
		return nil, err
	}

	if len(topics) == 0 {
		return nil, fmt.Errorf("topic %s not found", topic)
	}

	// Fetch all roadmaps by topic
	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(topics[0].RoadmapIds))
	for _, roadmapID := range topics[0].RoadmapIds {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(roadmapID)},
		})
	}

	batchGetInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			r.tableName: {
				Keys: keys,
			},
		},
	}

	batchGetResult, err := r.db.BatchGetItemWithContext(ctx, batchGetInput)
	if err != nil {
		return nil, err
	}

	// Unmarshal roadmaps
	var roadmaps []domain.Roadmap
	err = dynamodbattribute.UnmarshalListOfMaps(batchGetResult.Responses[r.tableName], &roadmaps)
	if err != nil {
		return nil, err
	}

	return roadmaps, nil
}
