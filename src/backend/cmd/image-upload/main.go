package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strings"
)

var (
	bucketName = os.Getenv("BUCKET_NAME")
	s3Client   *s3.S3
)

func init() {
	sess := session.Must(session.NewSession())
	s3Client = s3.New(sess)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileData, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("Failed to decode file data: %s", err)}, nil
	}

	fileName := request.Headers["filename"]

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   strings.NewReader(string(fileData)),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Failed to upload file to S3: %s", err)}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "File uploaded successfully!"}, nil
}

func main() {
	lambda.Start(handler)
}
