package repository

import (
	"backend/internal/graphql/models"
	"context"
)

type IUserRepository interface {
	UpsertUser(user models.User) (models.User, error)
	GetUsers() []*models.User
	GetUserByName(name string) (*models.User, error)
}

type ITopicRepository interface {
	Insert(ctx context.Context, names []*models.Topic) error
	GetAllTopics(ctx context.Context) ([]*models.Topic, error)
}
