package main

import "fmt"

// ProblemRepository interface
type ProblemRepository interface {
	UpsertProblem(problem Problem) Problem
	GetProblems() []*Problem
}

// InMemoryProblemRepository implementation
type InMemoryProblemRepository struct {
	problems []Problem
}

func NewInMemoryProblemRepository() *InMemoryProblemRepository {
	return &InMemoryProblemRepository{
		problems: make([]Problem, 0),
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

func generatePID(count int) string {
	return fmt.Sprintf("%d", count+1)
}
