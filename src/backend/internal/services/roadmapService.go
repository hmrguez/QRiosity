package services

import (
	"backend/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type RoadmapService struct {
	baseURL string
}

func NewRoadmapService() *RoadmapService {
	return &RoadmapService{
		baseURL: os.Getenv("API_ENDPOINT"),
	}
}

func (s *RoadmapService) GetCustomRoadmap(topic string) (*domain.Roadmap, error) {
	var roadmap domain.Roadmap
	encodedTopic := url.QueryEscape(topic)
	apiRequest := fmt.Sprintf("%s/get-roadmap?topic=%s", s.baseURL, encodedTopic)

	fmt.Println("URL: ", apiRequest)

	resp, err := http.Get(apiRequest)
	if err != nil {
		return &roadmap, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &roadmap, fmt.Errorf("failed to get roadmap: %s, %v", resp.Status, resp.Body)
	}

	err = json.NewDecoder(resp.Body).Decode(&roadmap)
	if err != nil {
		return &roadmap, err
	}

	return &roadmap, nil
}
