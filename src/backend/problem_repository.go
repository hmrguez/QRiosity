package main

import "fmt"

// ProblemRepository interface
type ProblemRepository interface {
	UpsertProblem(problem Problem) Problem
	GetProblems() []*Problem
	GetDailyChallenge(category string) (*Problem, error)
	SubmitChallengeResponse(userId, question, answer string) (ChallengeResponse, error)
}

// InMemoryProblemRepository implementation
type InMemoryProblemRepository struct {
	problems []Problem
}

func NewInMemoryProblemRepository() *InMemoryProblemRepository {
	//return &InMemoryProblemRepository{
	//	problems: make([]Problem, 0),
	//}

	// Prepopulate the problems with 1 random value
	problems := make([]Problem, 1)
	problems[0] = Problem{
		ID:         "1",
		Question:   "What AWS Kinesis used for?",
		Categories: []string{"Solution Architecture"},
	}

	return &InMemoryProblemRepository{
		problems: problems,
	}

}

func (r *InMemoryProblemRepository) UpsertProblem(problem Problem) Problem {
	for i, p := range r.problems {
		if p.ID == problem.ID {
			r.problems[i] = problem
			return problem
		}
	}
	problem.ID = generatePID(len(r.problems))
	r.problems = append(r.problems, problem)
	return problem
}

func (r *InMemoryProblemRepository) GetProblems() []*Problem {
	answer := make([]*Problem, len(r.problems))
	for i, p := range r.problems {
		answer[i] = &p
	}

	return answer
}

func (r *InMemoryProblemRepository) GetDailyChallenge(category string) (*Problem, error) {
	return &r.problems[0], nil
}

func (r *InMemoryProblemRepository) SubmitChallengeResponse(userId, question, answer string) (ChallengeResponse, error) {
	return ChallengeResponse{
		Question: question,
		Answer:   answer,
		Rating:   "Good answer",
		Insight:  "Nothing to improve",
	}, nil
}

func generatePID(count int) string {
	return fmt.Sprintf("%d", count+1)
}
