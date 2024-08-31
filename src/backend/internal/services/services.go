package services

import "backend/internal/domain"

type IDailyChallengeService interface {
	GetQuestion(username string) (domain.Problem, error)
	RateQuestion(question, answer string) (domain.ChallengeResponse, error)
}

type IRoadmapService interface {
	GetCustomRoadmap(topic string) (*domain.Roadmap, error)
}
