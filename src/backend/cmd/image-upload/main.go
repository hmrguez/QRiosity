package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	maxFileSize = 5 * 1024 * 1024 // 5MB
	bucketName  = "roadmap-images"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the incoming request to extract the image file
	contentType := request.Headers["Content-Type"]
	if contentType == "" {
		contentType = request.Headers["content-type"]
	}

	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid content type",
		}, nil
	}

	reader := multipart.NewReader(bytes.NewReader([]byte(request.Body)), contentType)
	part, err := reader.NextPart()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Failed to read the image file",
		}, nil
	}

	defer part.Close()

	// Validate file type
	fileType := strings.ToLower(part.Header.Get("Content-Type"))
	if fileType != "image/jpeg" && fileType != "image/jpg" && fileType != "image/png" && fileType != "image/webp" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid file type",
		}, nil
	}

	// Check file size
	buf := new(bytes.Buffer)
	fileSize, err := io.Copy(buf, part)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to read the file",
		}, nil
	}

	if fileSize > maxFileSize {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "File size exceeds the 5MB limit",
		}, nil
	}

	// Generate a random file name
	fileName, err := generateRandomFileName(part.FileName())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to generate file name",
		}, nil
	}

	// Upload the file to S3
	sess := session.Must(session.NewSession())
	uploader := s3.New(sess)
	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to upload file to S3",
		}, nil
	}

	// Return success response with the S3 URL as plain text
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fileURL,
	}, nil
}

func generateRandomFileName(fileName string) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 15)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}

	ext := strings.ToLower(strings.TrimPrefix(fileName, filepath.Ext(fileName)))

	return fmt.Sprintf("%s.%s", string(b), ext), nil
}

func main() {
	lambda.Start(handler)
}
