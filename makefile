graph:
	go get github.com/99designs/gqlgen@v0.17.49 && go run github.com/99designs/gqlgen generate

deploy-auth:
	cd src/backend/cmd/auth && \
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go && \
	zip -r auth.zip bootstrap && \
	aws lambda update-function-code --function-name auth-appsync --zip-file fileb://auth.zip

deploy-daily:
	cd src/backend/cmd/daily && \
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go && \
	zip -r daily.zip bootstrap && \
	aws lambda update-function-code --function-name daily-appsync --zip-file fileb://daily.zip

deploy-learning:
	cd src/backend/cmd/learning && \
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go && \
	zip -r learning.zip bootstrap && \
	aws lambda update-function-code --function-name learning-appsync --zip-file fileb://learning.zip

deploy-s3-lambda:
	cd src/backend/cmd/image-upload && \
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go && \
	zip -r image-upload.zip bootstrap && \
	aws lambda update-function-code --function-name image-upload --zip-file fileb://image-upload.zip

deploy-python-lambda:
	@echo "Usage: make deploy-python-lambda NAME=<name>"
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required"; \
		exit 1; \
	fi
	cd src/ai/$(NAME) && \
	zip -r $(NAME).zip lambda_function.py && \
	aws lambda update-function-code --function-name $(NAME) --zip-file fileb://$(NAME).zip