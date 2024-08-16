package services

import "backend/internal/graphql/models"

type IDailyChallengeService interface {
	GetQuestion(username string) (models.Problem, error)
	RateQuestion(question, answer string) (models.ChallengeResponse, error)
}
