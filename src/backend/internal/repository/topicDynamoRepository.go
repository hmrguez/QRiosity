package repository

import (
	"backend/internal/domain"
	"context"
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
