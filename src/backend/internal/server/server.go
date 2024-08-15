package server

import (
	"backend/internal/config"
	"backend/internal/graphql/generated"
	"backend/internal/graphql/resolvers"
	h "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func StartServer(config config.Config) {
	srv := h.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{
		UserRepo:    config.UserRepo,
		ProblemRepo: config.ProblemRepo,
		TopicRepo:   config.TopicRepo,
		AuthService: config.AuthService,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(srv))

	log.Printf("connect to http://localhost%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
