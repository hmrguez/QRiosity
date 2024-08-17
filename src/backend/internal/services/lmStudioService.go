package services

import (
	"backend/internal/domain"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// LMStudioService implementation
type LMStudioService struct {
	endpoint string
}

func NewLDefaultLMStudioService() *LMStudioService {
	return &LMStudioService{endpoint: "http://localhost:1234/v1/chat/completions"}
}

func (s *LMStudioService) RateQuestion(question, answer string) (domain.ChallengeResponse, error) {
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

	insight, err := AskLMStudio(payload, s.endpoint)
	if err != nil {
		return domain.ChallengeResponse{}, err
	}

	// Extract the first number found in the insight as the rating
	re := regexp.MustCompile(`\d+`)
	rating := 0
	if match := re.FindString(insight); match != "" {
		rating, _ = strconv.Atoi(match)
	}

	return domain.ChallengeResponse{
		Question: question,
		Answer:   answer,
		Rating:   rating,
		Insight:  insight,
	}, nil
}

func (s *LMStudioService) GetQuestion(category string) (domain.Problem, error) {
	// Define the request payload
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": fmt.Sprintf("Give me a question for the following category: %s", category)},
		},
		"temperature": 0.7,
		"max_tokens":  -1,
		"stream":      true,
	}

	question, err := AskLMStudio(payload, s.endpoint)
	if err != nil {
		return domain.Problem{}, err
	}

	return domain.Problem{
		Question:   question,
		Categories: append(make([]string, 0), category),
		Type:       "LLM Asking",
	}, nil
}

func AskLMStudio(payload map[string]interface{}, endpoint string) (string, error) {
	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Make the HTTP request
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the streaming response
	var builder strings.Builder
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF || err.Error() == "unexpected end of JSON input" {
				break
			}
		}

		line = strings.TrimPrefix(line, "data: ")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse the JSON line
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			return "", err
		}

		// Extract the content
		if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
			if delta, ok := choices[0].(map[string]interface{})["delta"].(map[string]interface{}); ok {
				if content, ok := delta["content"].(string); ok {
					builder.WriteString(content)
				}
			}
		}
	}

	return builder.String(), err
}
