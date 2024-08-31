package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	maxFileSize    = 5 * 1024 * 1024 // 5MB
	fileNameLength = 15
	bucketName     = "roadmap-images"
)

var (
	s3Client *s3.S3
)

func init() {
	sess := session.Must(session.NewSession())
	s3Client = s3.New(sess)
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Decode base64 encoded file
	fileContent, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		return createResponse(http.StatusBadRequest, "Invalid file content"), nil
	}

	// Check file size
	if len(fileContent) > maxFileSize {
		return createResponse(http.StatusBadRequest, "File size exceeds 5MB limit"), nil
	}

	// Detect file type
	fileType := http.DetectContentType(fileContent)
	if !isAllowedFileType(fileType) {
		return createResponse(http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, and WebP are allowed"), nil
	}

	// Generate random file name
	fileName := generateRandomFileName(fileType)

	// Upload file to S3
	fileURL, err := uploadToS3(fileContent, fileName, fileType)
	if err != nil {
		return createResponse(http.StatusInternalServerError, "Failed to upload file"), nil
	}

	return createResponse(http.StatusOK, fileURL), nil
}

func isAllowedFileType(fileType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	return allowedTypes[fileType]
}

func generateRandomFileName(fileType string) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomBytes := make([]byte, fileNameLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	for i, b := range randomBytes {
		randomBytes[i] = letters[b%byte(len(letters))]
	}

	extension := strings.TrimPrefix(fileType, "image/")
	if extension == "jpeg" {
		extension = "jpg"
	}

	return string(randomBytes) + "." + extension
}

func uploadToS3(fileContent []byte, fileName, fileType string) (string, error) {
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(fileContent),
		ContentType:   aws.String(fileType),
		ContentLength: aws.Int64(int64(len(fileContent))),
	})

	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName)
	return fileURL, nil
}

func createResponse(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                 "text/plain",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}
}

func main() {
	lambda.Start(handleRequest)
}
