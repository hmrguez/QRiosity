package services

import "backend/internal/graphql/models"

type LLMService interface {
	GetQuestion(category string) (models.Problem, error)
	RateQuestion(question, answer string) (models.ChallengeResponse, error)
}
