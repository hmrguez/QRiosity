package repository

import (
	"backend/internal/graphql/models"
	"context"
)

type UserRepository interface {
	UpsertUser(user models.User) (models.User, error)
	GetUsers() []*models.User
	GetUserByID(id string) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
}

type ProblemRepository interface {
	UpsertProblem(problem models.Problem) models.Problem
	GetProblems() []*models.Problem
	GetDailyChallenge(category string) (*models.Problem, error)
	SubmitChallengeResponse(userId, question, answer string) (models.ChallengeResponse, error)
}

type TopicRepository interface {
	Insert(ctx context.Context, names []*models.Topic) error
	GetAllTopics(ctx context.Context) ([]*models.Topic, error)
}
