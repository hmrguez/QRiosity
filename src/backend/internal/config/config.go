package config

import (
	"backend/internal/db/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

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

func NewProdConfig() *Config {
	// Create a new AWS session using credentials from .env
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	userRepo, err := repository.NewMongoDBUserRepository("mongodb://localhost:27017", "saas", "users")
	if err != nil {
		panic("Couldn't connnect to mongodb: " + err.Error())
	}

	return &Config{
		Port:        ":9000",
		UserRepo:    userRepo,
		ProblemRepo: repository.NewMongoDBProblemRepository(),
		TopicRepo:   repository.NewDynamoDBTopicRepository(sess, "Qriosity-Topics"),
	}
}
