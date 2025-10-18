package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type DeepSeekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

// Основная функция для запроса к DeepSeek
func AskDeepSeek(prompt string, personality string) (string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY not set")
	}

	// Формируем полный промпт с личностью судьи
	fullPrompt := fmt.Sprintf("Ты - %s. Ответь на вопрос: %s", personality, prompt)

	requestBody := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "user",
				Content: fullPrompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %v", err)
	}

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request creation error: %v", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("response read error: %v", err)
	}

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	// Парсим JSON ответ
	var deepseekResp DeepSeekResponse
	err = json.Unmarshal(body, &deepseekResp)
	if err != nil {
		return "", fmt.Errorf("JSON unmarshal error: %v", err)
	}

	// Извлекаем ответ
	if len(deepseekResp.Choices) > 0 {
		return deepseekResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}
