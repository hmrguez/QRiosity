package main

import (
	"backend/internal/repository"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

func TestGetBulkByUrl() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})

	if err != nil {
		log.Fatal(err)
	}

	courseRepository := repository.NewDynamoDBCourseRepository(sess, "Qriosity-Courses")

	urls := []string{"https://www.coursera.org/learn/learning-how-to-learn", "https://www.coursera.org/learn/machine-learning"}

	courses, err := courseRepository.GetBulkByUrl(context.Background(), urls)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(courses)

}

func main() {
	TestGetBulkByUrl()
}
