package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BryceWayne/go-fiber-gemma/internal/models"
)

// OllamaEndpoint allows unit tests to override the target URL
var OllamaEndpoint = "http://localhost:11434/api/chat"

// AskOllama sends the given full prompt to the local Ollama instance (gemma4:e2b)
func AskOllama(fullPrompt string) (string, error) {
	chatPayload := models.OllamaChatRequest{
		Model: "gemma4:e2b",
		Messages: []models.ChatMessage{
			{Role: "user", Content: fullPrompt},
		},
		Stream: false,
		Options: map[string]interface{}{
			"num_gpu":     35,   // Pin all layers to RTX 4070
			"num_ctx":     8192, // High enough for code, small enough for speed
			"temperature": 1.0,  // Gemma 4 recommended temp
		},
	}

	jsonData, err := json.Marshal(chatPayload)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", OllamaEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result models.OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}
