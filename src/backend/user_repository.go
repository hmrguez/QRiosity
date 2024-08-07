package main

import "fmt"

type UserRepository interface {
	UpsertUser(user User) (User, error)
	GetUsers() []*User
	GetUserByID(id string) (*User, error)
	GetUserByName(name string) (*User, error)
}

type InMemoryUserRepository struct {
	users []User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make([]User, 0),
	}
}

func (r *InMemoryUserRepository) UpsertUser(user User) User {
	for i, u := range r.users {
		if u.ID == user.ID {
			r.users[i] = user
			return user
		}
	}
	user.ID = generateID(r) // function to generate a new ID
	r.users = append(r.users, user)
	return user
}

func (r *InMemoryUserRepository) GetUsers() []*User {
	// Return a new array with the pointer to each user
	users := make([]*User, len(r.users))
	for i, user := range r.users {
		users[i] = &user
	}

	return users
}

func generateID(r *InMemoryUserRepository) string {
	return fmt.Sprintf("%d", len(r.users)+1) // simple ID generation
}
