package repository

import (
	"backend/internal/domain"
	"context"
)

type IUserRepository interface {
	UpsertUser(user domain.User) (domain.User, error)
	GetUsers() []*domain.User
	GetUserByName(name string) (*domain.User, error)
}

type ITopicRepository interface {
	Insert(ctx context.Context, names []*domain.Topic) error
	GetAllTopics(ctx context.Context) ([]*domain.Topic, error)
}
