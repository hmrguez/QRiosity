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

type DynamoDBCourseRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBCourseRepository(sess *session.Session, tableName string) *DynamoDBCourseRepository {
	return &DynamoDBCourseRepository{
		db:        dynamodb.New(sess),
		tableName: tableName,
	}
}

func (r *DynamoDBCourseRepository) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var courses []*domain.Course
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &courses)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *DynamoDBCourseRepository) UpsertCourse(ctx context.Context, course *domain.Course) error {
	item, err := dynamodbattribute.MarshalMap(course)
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

func (r *DynamoDBCourseRepository) GetCourseByID(ctx context.Context, courseID string) (*domain.Course, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(courseID)},
		},
	}

	result, err := r.db.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("course not found")
	}

	var course domain.Course
	err = dynamodbattribute.UnmarshalMap(result.Item, &course)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *DynamoDBCourseRepository) GetBulkByUrl(ctx context.Context, urls []string) ([]*domain.Course, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(urls))
	for i, url := range urls {
		keys[i] = map[string]*dynamodb.AttributeValue{
			"url": {S: aws.String(url)},
		}
	}

	batchGetInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			r.tableName: {
				Keys: keys,
			},
		},
	}

	result, err := r.db.BatchGetItemWithContext(ctx, batchGetInput)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get items: %w", err)
	}

	courses := make([]*domain.Course, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Responses[r.tableName], &courses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal courses: %w", err)
	}

	return courses, nil
}

func (r *DynamoDBCourseRepository) BulkInsert(ctx context.Context, courses []*domain.Course) error {
	writeRequests := make([]*dynamodb.WriteRequest, len(courses))
	for i, course := range courses {
		item, err := dynamodbattribute.MarshalMap(course)
		if err != nil {
			return fmt.Errorf("failed to marshal course: %w", err)
		}
		writeRequests[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		}
	}

	batchWriteInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			r.tableName: writeRequests,
		},
	}

	_, err := r.db.BatchWriteItemWithContext(ctx, batchWriteInput)
	if err != nil {
		return fmt.Errorf("failed to batch write items: %w", err)
	}

	return nil
}
