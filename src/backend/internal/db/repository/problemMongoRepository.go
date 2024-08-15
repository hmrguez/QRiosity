package repository

import (
	"backend/internal/graphql/models"
	"backend/internal/services"
)

type MongoDBProblemRepository struct {
	services.LLMService
}

func NewMongoDBProblemRepository(llmService services.LLMService) *MongoDBProblemRepository {

	return &MongoDBProblemRepository{
		LLMService: llmService,
	}

}

func (r *MongoDBProblemRepository) UpsertProblem(problem models.Problem) models.Problem {
	panic("Not implemented")
}

func (r *MongoDBProblemRepository) GetProblems() []*models.Problem {
	panic("Not implemented")
}

func (r *MongoDBProblemRepository) GetDailyChallenge(category string) (*models.Problem, error) {
	problem, err := r.GetQuestion(category)
	return &problem, err
}

func (r *MongoDBProblemRepository) SubmitChallengeResponse(userId, question, answer string) (models.ChallengeResponse, error) {
	return r.LLMService.RateQuestion(question, answer)
}
