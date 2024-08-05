package main

import (
	"context"
	_ "github.com/99designs/gqlgen/graphql"
)

type Resolver struct {
	userRepo    UserRepository
	problemRepo ProblemRepository
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) UpsertProblem(ctx context.Context, input ProblemInput) (*Problem, error) {
	problem := Problem{
		Question:   input.Question,
		Categories: input.Categories,
		Type:       input.Type,
	}
	upsertedProblem := r.problemRepo.UpsertProblem(problem)
	return &upsertedProblem, nil
}

func (r *mutationResolver) UpsertUser(ctx context.Context, input UserInput) (*User, error) {
	user := User{
		Name:  input.Name,
		Email: input.Email,
	}
	upsertedUser := r.userRepo.UpsertUser(user)
	return &upsertedUser, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetProblems(ctx context.Context) ([]*Problem, error) {
	problems := r.problemRepo.GetProblems()
	return problems, nil
}

func (r *queryResolver) GetUsers(ctx context.Context) ([]*User, error) {
	users := r.userRepo.GetUsers()
	return users, nil
}
