package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type AIService struct {
	apiKey  string
	baseURL string
	model   string
}

func NewAIService(apiKey, baseURL, model string) *AIService {
	return &AIService{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: strings.TrimSpace(baseURL),
		model:   strings.TrimSpace(model),
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *AIService) GenerateReport(prompt string) (string, error) {
	if s.apiKey == "" || s.baseURL == "" {
		return "", errors.New("AI service is not configured")
	}

	reqBody := chatRequest{
		Model: s.model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Printf("AI request failed with status=%d body=%s", resp.StatusCode, string(bodyText))
		return "", fmt.Errorf("AI request failed with status %d", resp.StatusCode)
	}

	var aiResp chatResponse
	if err := json.Unmarshal(bodyText, &aiResp); err != nil {
		return "", err
	}
	if aiResp.Error != nil && aiResp.Error.Message != "" {
		return "", errors.New(aiResp.Error.Message)
	}
	if len(aiResp.Choices) == 0 {
		return "", errors.New("no response from AI")
	}
	return aiResp.Choices[0].Message.Content, nil
}
