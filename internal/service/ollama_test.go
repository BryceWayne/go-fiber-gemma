package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BryceWayne/go-fiber-gemma/internal/models"
)

func TestAskOllama_Success(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		var req models.OllamaChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Ensure we received the expected prompt
		if len(req.Messages) == 0 || req.Messages[0].Content != "SYSTEM: context\n\nUSER: hello" {
			t.Errorf("Unexpected prompt sent to Ollama: %v", req.Messages)
		}
		
		// Return a mock success response
		resp := models.OllamaChatResponse{
			Message: models.ChatMessage{
				Role:    "assistant",
				Content: "Hello from mock Gemma!",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	// Override global endpoint
	originalEndpoint := OllamaEndpoint
	OllamaEndpoint = mockServer.URL
	defer func() { OllamaEndpoint = originalEndpoint }()

	// Call the function
	answer, err := AskOllama("SYSTEM: context\n\nUSER: hello")
	if err != nil {
		t.Fatalf("AskOllama returned unexpected error: %v", err)
	}

	if answer != "Hello from mock Gemma!" {
		t.Errorf("Expected 'Hello from mock Gemma!', got '%s'", answer)
	}
}

func TestAskOllama_HTTPError(t *testing.T) {
	// Create a mock server that returns 500
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	originalEndpoint := OllamaEndpoint
	OllamaEndpoint = mockServer.URL
	defer func() { OllamaEndpoint = originalEndpoint }()

	_, err := AskOllama("test prompt")
	// Since json decode will fail on "internal server error" body
	if err == nil {
		t.Fatal("Expected error when server returns 500, got nil")
	}
}
