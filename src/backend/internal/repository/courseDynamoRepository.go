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
	courses := make([]*domain.Course, 0)

	for _, url := range urls {
		queryInput := &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			IndexName:              aws.String("YourSecondaryIndexName"), // Replace with your secondary index name
			KeyConditionExpression: aws.String("#u = :url"),
			ExpressionAttributeNames: map[string]*string{
				"#u": aws.String("url"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":url": {S: aws.String(url)},
			},
		}

		result, err := r.db.QueryWithContext(ctx, queryInput)
		if err != nil {
			return nil, fmt.Errorf("failed to query items: %w", err)
		}

		var batchCourses []*domain.Course
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &batchCourses)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal courses: %w", err)
		}

		courses = append(courses, batchCourses...)
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
