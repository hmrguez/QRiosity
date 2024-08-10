package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

const defaultPort = "8080"

func main() {

	userRepo, err := NewMongoDBUserRepository("mongodb://localhost:27017", "saas", "users")
	if err != nil {
		panic("Couldn't connnect to mongodb: " + err.Error())
	}

	topicRepo, err := NewMongoDBTopicRepository("mongodb://localhost:27017", "saas", "topics")
	if err != nil {
		panic("Couldn't connnect to mongodb: " + err.Error())
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{
		userRepo:    userRepo,
		problemRepo: NewInMemoryProblemRepository(),
		topicRepo:   topicRepo,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", "9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
