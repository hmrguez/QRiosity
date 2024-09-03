package repository

import (
	"backend/internal/domain"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"sync"
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
	if len(names) == 0 {
		return nil, fmt.Errorf("names slice is empty")
	}

	var wg sync.WaitGroup
	topicChan := make(chan *domain.Topic, len(names))
	errChan := make(chan error, len(names))

	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			input := &dynamodb.GetItemInput{
				TableName: aws.String("Qriosity-Topics"),
				Key: map[string]*dynamodb.AttributeValue{
					"name": {S: aws.String(name)},
				},
			}

			result, err := r.db.GetItemWithContext(ctx, input)
			if err != nil {
				errChan <- err
				return
			}

			if result.Item == nil {
				// Create a new topic if not found
				newTopic := &domain.Topic{
					Name:       name,
					RoadmapIds: []string{},
				}

				topicChan <- newTopic
				return
			}

			var topic domain.Topic
			err = dynamodbattribute.UnmarshalMap(result.Item, &topic)
			if err != nil {
				errChan <- err
				return
			}

			topicChan <- &topic
		}(name)
	}

	wg.Wait()
	close(topicChan)
	close(errChan)

	var topics []*domain.Topic
	for topic := range topicChan {
		topics = append(topics, topic)
	}

	if len(errChan) > 0 {
		return nil, <-errChan
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

	fmt.Printf("writeRequests: %v\n", writeRequests)

	// Create the batch write input
	batchWriteInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			r.tableName: writeRequests,
		},
	}

	fmt.Printf("batchWriteInput: %v\n", batchWriteInput)

	// Execute the batch write operation
	_, err := r.db.BatchWriteItemWithContext(ctx, batchWriteInput)
	return err
}
