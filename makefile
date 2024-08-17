graph:
	go get github.com/99designs/gqlgen@v0.17.49 && go run github.com/99designs/gqlgen generate

deploy-auth:
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go && \
	zip -r auth.zip bootstrap && \
	aws lambda update-function-code --function-name auth-appsync --zip-file fileb://auth.zip