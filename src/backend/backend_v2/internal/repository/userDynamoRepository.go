package repository

import (
	"backend/internal/graphql/models"
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

func (r *DynamoDBUserRepository) UpsertUser(user models.User) (models.User, error) {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return models.User{}, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.tableName),
	}

	_, err = r.client.PutItem(input)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *DynamoDBUserRepository) GetUsers() []*models.User {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := r.client.Scan(input)
	if err != nil {
		fmt.Println("Error getting users:", err)
		return nil
	}

	var users []*models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		fmt.Println("Error unmarshalling users:", err)
		return nil
	}

	return users
}

func (r *DynamoDBUserRepository) GetUserByName(name string) (*models.User, error) {
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

	var user models.User
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
