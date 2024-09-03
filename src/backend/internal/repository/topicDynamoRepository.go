package repository

import (
	"backend/internal/domain"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBTopicRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBTopicRepository(sess *session.Session, tableName string) *DynamoDBTopicRepository {
	return &DynamoDBTopicRepository{
		db:        dynamodb.New(sess),
		tableName: tableName,
	}
}

func (r *DynamoDBTopicRepository) Insert(ctx context.Context, topics []*domain.Topic) error {
	for _, topic := range topics {
		av, err := dynamodbattribute.MarshalMap(topic)
		if err != nil {
			return err
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(r.tableName),
		}

		_, err = r.db.PutItemWithContext(ctx, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DynamoDBTopicRepository) GetAllTopics(ctx context.Context) ([]*domain.Topic, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var topics []*domain.Topic
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &topics)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *DynamoDBTopicRepository) GetTopicsByNames(ctx context.Context, names []string) ([]*domain.Topic, error) {
	// Check if names slice is empty
	if len(names) == 0 {
		return nil, fmt.Errorf("names slice is empty")
	}

	// Do a batch operation and fetch all roadmap.Topics
	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(names))
	for _, topic := range names {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"name": {S: aws.String(topic)},
		})
	}

	batchGetInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			"Qriosity-Topics": {
				Keys: keys,
			},
		},
	}

	batchGetResult, err := r.db.BatchGetItemWithContext(ctx, batchGetInput)
	if err != nil {
		return nil, err
	}

	// Unmarshal topics
	var topics []*domain.Topic
	err = dynamodbattribute.UnmarshalListOfMaps(batchGetResult.Responses["Qriosity-Topics"], &topics)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *DynamoDBTopicRepository) BulkWrite(ctx context.Context, topics []*domain.Topic) error {
	// Prepare the write requests
	writeRequests := make([]*dynamodb.WriteRequest, 0, len(topics))
	for _, topic := range topics {
		av, err := dynamodbattribute.MarshalMap(topic)
		if err != nil {
			return err
		}
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		})
	}

	// Create the batch write input
	batchWriteInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			r.tableName: writeRequests,
		},
	}

	// Execute the batch write operation
	_, err := r.db.BatchWriteItemWithContext(ctx, batchWriteInput)
	return err
}
