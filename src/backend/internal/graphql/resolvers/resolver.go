package resolvers

import (
	"backend/internal/db/repository"
	"backend/internal/services"
)

type Resolver struct {
	AuthService services.CognitoAuthService
	UserRepo    repository.IUserRepository
	ProblemRepo repository.IProblemRepository
	TopicRepo   repository.ITopicRepository
}
