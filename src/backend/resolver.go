package main

import (
	"context"
	_ "github.com/99designs/gqlgen/graphql"
)

type Resolver struct {
	userRepo UserRepository
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) UpsertUser(ctx context.Context, input UserInput) (*User, error) {
	user := User{
		Name:  input.Name,
		Email: input.Email,
	}
	upsertedUser := r.userRepo.UpsertUser(user)
	return &upsertedUser, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetUsers(ctx context.Context) ([]*User, error) {
	users := r.userRepo.GetUsers()
	return users, nil
}
