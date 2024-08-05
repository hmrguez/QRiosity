package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type LLMService interface {
	RateQuestion(question, answer string) (ChallengeResponse, error)
}

// LMStudioService implementation
type LMStudioService struct {
	endpoint string
}

func NewLDefaultMStudioService() *LMStudioService {
	return &LMStudioService{endpoint: "http://localhost:1234/v1/chat/completions"}
}

func (s *LMStudioService) RateQuestion(question, answer string) (ChallengeResponse, error) {
	// Define the request payload
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "system", "content": "Rate the answer to the following question in a scale from 1 to 10."},
			{"role": "user", "content": fmt.Sprintf("Question: %s, Answer: %s", question, answer)},
		},
		"temperature": 0.7,
		"max_tokens":  -1,
		"stream":      true,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return ChallengeResponse{}, err
	}

	// Make the HTTP request
	resp, err := http.Post(s.endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ChallengeResponse{}, err
	}
	defer resp.Body.Close()

	// Read and print the response body for debugging
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return ChallengeResponse{}, err
	//}
	//fmt.Println("Response Body:", string(body))

	// Read the streaming response
	var insightBuilder strings.Builder
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF || err.Error() == "unexpected end of JSON input" {
				break
			}

			//return ChallengeResponse{}, err
		}

		line = strings.TrimPrefix(line, "data: ")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse the JSON line
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			return ChallengeResponse{}, err
		}

		// Extract the content
		if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
			if delta, ok := choices[0].(map[string]interface{})["delta"].(map[string]interface{}); ok {
				if content, ok := delta["content"].(string); ok {
					insightBuilder.WriteString(content)
				}
			}
		}
	}

	insight := insightBuilder.String()

	// Extract the first number found in the insight as the rating
	re := regexp.MustCompile(`\d+`)
	rating := "0"
	if match := re.FindString(insight); match != "" {
		rating = match
	}

	return ChallengeResponse{
		Question: question,
		Answer:   answer,
		Rating:   rating,
		Insight:  insight,
	}, nil
}
