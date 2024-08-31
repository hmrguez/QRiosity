package services

import (
	"backend/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
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

func (s *RoadmapService) GetCustomRoadmap(topic string) (domain.Roadmap, error) {
	var roadmap domain.Roadmap
	url := fmt.Sprintf("%s/get-roadmap?topic=%s", s.baseURL, topic)
	resp, err := http.Get(url)
	if err != nil {
		return roadmap, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return roadmap, fmt.Errorf("failed to get roadmap: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&roadmap)
	if err != nil {
		return roadmap, err
	}

	return roadmap, nil
}
