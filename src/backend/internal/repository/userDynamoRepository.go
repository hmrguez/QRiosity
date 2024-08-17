package repository

import (
	"backend/internal/domain"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBUserRepository struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBUserRepository(session *session.Session, tableName string) (*DynamoDBUserRepository, error) {
	client := dynamodb.New(session)
	return &DynamoDBUserRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (r *DynamoDBUserRepository) UpsertUser(user domain.User) (domain.User, error) {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return domain.User{}, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.tableName),
	}

	_, err = r.client.PutItem(input)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *DynamoDBUserRepository) GetUsers() []*domain.User {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.client.Scan(input)
	if err != nil {
		fmt.Println("Error getting users:", err)
		return nil
	}

	var users []*domain.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		fmt.Println("Error unmarshalling users:", err)
		return nil
	}

	return users
}

func (r *DynamoDBUserRepository) GetUserByName(name string) (*domain.User, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"name": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(name),
					},
				},
			},
		},
	}

	result, err := r.client.Query(input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user domain.User
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
