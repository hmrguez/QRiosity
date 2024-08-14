package resolvers

import (
	"backend/internal/db/repository"
	"backend/internal/services"
)

type Resolver struct {
	AuthService services.CognitoAuthService
	UserRepo    repository.UserRepository
	ProblemRepo repository.ProblemRepository
	TopicRepo   repository.TopicRepository
}
