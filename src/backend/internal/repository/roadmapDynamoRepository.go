package repository

import (
	"backend/internal/domain"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBRoadmapRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBRoadmapRepository(sess *session.Session, tableName string) *DynamoDBRoadmapRepository {
	return &DynamoDBRoadmapRepository{
		db:        dynamodb.New(sess),
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
	item, err := dynamodbattribute.MarshalMap(roadmap)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.db.PutItemWithContext(ctx, input)
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
	courses := make(map[string]*domain.Course)
	err = dynamodbattribute.UnmarshalListOfMaps(batchGetResult.Responses["Qriosity-Courses"], &courses)
	if err != nil {
		return nil, err
	}

	// Populate Courses field in roadmap
	for _, courseID := range roadmap.CourseIDs {
		if course, ok := courses[courseID]; ok {
			roadmap.Courses = append(roadmap.Courses, *course)
		}
	}

	return &roadmap, nil
}
