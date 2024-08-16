package server

//
//import (
//	"backend/internal/config"
//	"backend/internal/graphql/generated"
//	"backend/internal/graphql/resolvers"
//	"context"
//	h "github.com/99designs/gqlgen/graphql/handler"
//	"github.com/aws/aws-lambda-go/events"
//	"github.com/aws/aws-lambda-go/lambda"
//	"github.com/gorilla/handlers"
//	"net/http"
//	"net/http/httptest"
//)
//
//func graphqlHandler(config config.Config) http.Handler {
//	srv := h.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{
//		UserRepo:    config.UserRepo,
//		ProblemRepo: config.ProblemRepo,
//		TopicRepo:   config.TopicRepo,
//		AuthService: config.AuthService,
//	}}))
//
//	return handlers.CORS(
//		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // Adjust as needed for production
//		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
//		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
//	)(srv)
//}
//
//func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	// Create a new HTTP request from the API Gateway event
//	request, err := http.NewRequest(req.HTTPMethod, req.Path, nil)
//	if err != nil {
//		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
//	}
//
//	// Create a ResponseRecorder to capture the response
//	recorder := httptest.NewRecorder()
//
//	// Serve the request
//	graphqlHandler(*config.NewProdConfig()).ServeHTTP(recorder, request)
//
//	// Return the response back to API Gateway
//	return events.APIGatewayProxyResponse{
//		StatusCode: recorder.Code,
//		Body:       recorder.Body.String(),
//		Headers:    map[string]string{"Content-Type": "application/json"},
//	}, nil
//}
//
//func LambdaServer() {
//	lambda.Start(handler)
//}
