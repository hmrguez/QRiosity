package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	s3Client   *s3.S3
	bucketName = os.Getenv("BUCKET_NAME")
)

func init() {
	sess := session.Must(session.NewSession())
	s3Client = s3.New(sess)
	rand.Seed(time.Now().UnixNano())
}

func generateRandomFilename() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileData, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Failed to decode file data: %s", err),
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	fileName := generateRandomFilename()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   strings.NewReader(string(fileData)),
	})

	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to upload file to S3: %s", err),
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("{filename: \"%s\"}", fileName),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
