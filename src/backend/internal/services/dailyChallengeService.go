package services

import (
	"backend/internal/graphql/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DailyChallengeService struct {
	baseURL string
}

func NewDailyChallengeService() *DailyChallengeService {
	return &DailyChallengeService{
		baseURL: os.Getenv("AI_API_URL"),
	}
}

func (s *DailyChallengeService) GetQuestion(category string) (models.Problem, error) {
	var problem models.Problem
	url := fmt.Sprintf("%s/ask_question?topic=%s", s.baseURL, category)
	resp, err := http.Get(url)
	if err != nil {
		return problem, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return problem, fmt.Errorf("failed to get question: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&problem)
	if err != nil {
		return problem, err
	}

	return problem, nil
}

func (s *DailyChallengeService) RateQuestion(question, answer string) (models.ChallengeResponse, error) {
	var challengeResponse models.ChallengeResponse
	url := fmt.Sprintf("%s/rate_question", s.baseURL)

	payload := map[string]string{
		"question": question,
		"answer":   answer,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return challengeResponse, err
	}

	fmt.Println("Url:", url)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return challengeResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return challengeResponse, fmt.Errorf("failed to rate question: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&challengeResponse)
	if err != nil {
		return challengeResponse, err
	}

	return challengeResponse, nil
}
