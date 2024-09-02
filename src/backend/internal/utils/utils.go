package utils

import "github.com/aws/aws-sdk-go/service/dynamodb"

type Pagination struct {
	Page             int                                 `json:"page"`
	PerPage          int                                 `json:"perPage"`
	LastEvaluatedKey map[string]*dynamodb.AttributeValue `json:"lastEvaluatedKey,omitempty"`
}
