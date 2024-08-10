package resolvers

import "backend/internal/db/repository"

type Resolver struct {
	UserRepo    repository.UserRepository
	ProblemRepo repository.ProblemRepository
	TopicRepo   repository.TopicRepository
}
