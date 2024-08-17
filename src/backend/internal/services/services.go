package services

import "backend/internal/domain"

type IDailyChallengeService interface {
	GetQuestion(username string) (domain.Problem, error)
	RateQuestion(question, answer string) (domain.ChallengeResponse, error)
}
