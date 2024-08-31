package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(ctx context.Context) (interface{}, error) {
	// Create an S3 client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	s3Client := s3.New(sess)

	// Get the request body
	requestBody, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	// Parse the multipart form data
	var formData struct {
		File []struct {
			Content string `form:"file"`
		} `form:"file"`
	}
	if err := json.Unmarshal(requestBody, &formData); err != nil {
		return nil, fmt.Errorf("error parsing form data: %v", err)
	}

	// Check if there's at least one file
	if len(formData.File) == 0 {
		return nil, fmt.Errorf("no file provided")
	}

	// Process each file
	for _, file := range formData.File {
		// Extract the content from the base64 encoded string
		imgData, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return nil, fmt.Errorf("error decoding base64: %v", err)
		}

		// Validate the file type
		contentType := http.DetectContentType(imgData)
		if !isValidImageType(contentType) {
			return nil, fmt.Errorf("invalid image type: %s", contentType)
		}

		// Get extension from content type
		var ext = strings.Split(contentType, "/")[1]

		// Validate the file size
		if float64(len(imgData)) > 5*1024*1024 { // 5 MB
			return nil, fmt.Errorf("file exceeds 5MB limit")
		}

		// Generate a random filename
		rand.Seed(time.Now().UnixNano())
		fileName := generateRandomName() + "." + ext

		// Upload the file to S3
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String("roadmap-images"),
			Key:    aws.String(fileName),
			Body:   bytes.NewReader(imgData),
		})
		if err != nil {
			return nil, fmt.Errorf("error uploading file to S3: %v", err)
		}

		// Return the URL of the uploaded file
		url := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", "your-bucket-name", fileName)
		return map[string]string{"url": url}, nil
	}

	return nil, nil
}

// Helper functions

func isValidImageType(contentType string) bool {
	validTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}
	for _, t := range validTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

func generateRandomName() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 15)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	lambda.Start(handler)
}
