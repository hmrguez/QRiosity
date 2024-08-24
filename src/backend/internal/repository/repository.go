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

type ICourseRepository interface {
	GetAllCourses(ctx context.Context) ([]*domain.Course, error)
	UpsertCourse(ctx context.Context, course *domain.Course) error
}

type IRoadmapRepository interface {
	GetAllRoadmaps(ctx context.Context) ([]*domain.Roadmap, error)
	UpsertRoadmap(ctx context.Context, roadmap *domain.Roadmap) error
}
