package repository

import (
	"backend/internal/domain"
	"backend/internal/utils"
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
	GetTopicsByNames(ctx context.Context, names []string) ([]*domain.Topic, error)
	BulkWrite(ctx context.Context, topics []*domain.Topic) error
}

type ICourseRepository interface {
	GetAllCourses(ctx context.Context, pagination utils.Pagination) ([]*domain.Course, *utils.Pagination, error)
	UpsertCourse(ctx context.Context, course *domain.Course) error
	GetCourseByID(ctx context.Context, courseID string) (*domain.Course, error)
	GetBulkByUrl(ctx context.Context, urls []string) ([]*domain.Course, error)
	BulkInsert(ctx context.Context, courses []*domain.Course) error
}

type IRoadmapRepository interface {
	GetAllRoadmaps(ctx context.Context) ([]*domain.Roadmap, error)
	UpsertRoadmap(ctx context.Context, roadmap *domain.Roadmap) error
	GetRoadmap(ctx context.Context, roadmapID string) (*domain.Roadmap, error)
	GetRoadmapsByUser(ctx context.Context, userID string, userRepo IUserRepository) ([]*domain.Roadmap, error)
	GetByTopic(ctx context.Context, topic string) ([]domain.Roadmap, error)
}
