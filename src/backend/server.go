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

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{
		userRepo:    userRepo,
		problemRepo: NewInMemoryProblemRepository(),
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", "9000")
	log.Fatal(http.ListenAndServe(":9000", nil))

	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = defaultPort
	//}
	//
	//var userRepo UserRepository = NewInMemoryUserRepository()
	//var problemRepo ProblemRepository = NewInMemoryProblemRepository()
	//
	//srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{
	//	userRepo:    userRepo,
	//	problemRepo: problemRepo,
	//}}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}
