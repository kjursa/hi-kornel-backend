package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-go-backend/models"
	"net/http"
)

type AiClient struct {
	config AiClientConfig
}

func NewAiClient(config AiClientConfig) *AiClient {
	return &AiClient{
		config: config,
	}
}

type AiBody struct {
	Model    string      `json:"model"`
	Messages []AiMessage `json:"messages"`
	Stream   bool        `json:"stream"`
}

type AiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var RoleSystem = "system"
var RoleUser = "user"
var RoleAssistant = "assistant"

func (a *AiClient) Send(system string, context string, history []models.Chat, question string) (string, error) {
	messages := buildAiMessages(system, context, history, question)
	for _, res := range messages {
		fmt.Println(res.Role + "->" + res.Content)
	}
	fmt.Println("---")
	answer, err := completions(a.config, messages)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func completions(config AiClientConfig, messages []AiMessage) (string, error) {
	body := AiBody{
		Model:    config.Model,
		Messages: messages,
		Stream:   false,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var aiResp AiResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return aiResp.Choices[0].Message.Content, nil
}

func buildAiMessages(system string, context string, history []models.Chat, question string) []AiMessage {
	var response []AiMessage
	response = append(response, AiMessage{
		Role:    RoleSystem,
		Content: system,
	})
	response = append(response, AiMessage{
		Role:    RoleUser,
		Content: context,
	})

	for _, chat := range reverse(history) {
		response = append(response, AiMessage{
			Role:    RoleUser,
			Content: chat.Question,
		})
		response = append(response, AiMessage{
			Role:    RoleAssistant,
			Content: chat.Answer,
		})
	}

	response = append(response, AiMessage{
		Role:    RoleUser,
		Content: question,
	})

	return response
}

func reverse(input []models.Chat) []models.Chat {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
	return input
}
