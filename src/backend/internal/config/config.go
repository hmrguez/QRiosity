package config

import "backend/internal/db/repository"

type Config struct {
	Port        string
	UserRepo    repository.UserRepository
	ProblemRepo repository.ProblemRepository
	TopicRepo   repository.TopicRepository
}

func NewLocalConfig() *Config {

	userRepo, err := repository.NewMongoDBUserRepository("mongodb://localhost:27017", "saas", "users")
	if err != nil {
		panic("Couldn't connnect to mongodb: " + err.Error())
	}

	topicRepo, err := repository.NewMongoDBTopicRepository("mongodb://localhost:27017", "saas", "topics")
	if err != nil {
		panic("Couldn't connnect to mongodb: " + err.Error())
	}

	return &Config{
		Port:        ":9000",
		UserRepo:    userRepo,
		ProblemRepo: repository.NewMongoDBProblemRepository(),
		TopicRepo:   topicRepo,
	}
}
