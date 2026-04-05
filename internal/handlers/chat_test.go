package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/BryceWayne/go-fiber-gemma/internal/models"
	"github.com/BryceWayne/go-fiber-gemma/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func mockOllamaServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.OllamaChatResponse{
			Message: models.ChatMessage{
				Role:    "assistant",
				Content: "Mocked response",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestHandleAsk_Success(t *testing.T) {
	mockServer := mockOllamaServer()
	defer mockServer.Close()
	service.OllamaEndpoint = mockServer.URL

	app := fiber.New()
	app.Post("/api/ask", HandleAsk)

	body := `{"prompt": "hello world"}`
	req := httptest.NewRequest("POST", "/api/ask", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, int(2*time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(respBody), "Mocked response") {
		t.Errorf("Expected 'Mocked response' in JSON, got %s", string(respBody))
	}
}

func TestHandleHTMXChat_Success(t *testing.T) {
	mockServer := mockOllamaServer()
	defer mockServer.Close()
	service.OllamaEndpoint = mockServer.URL

	// NOTE: We need templates mapped correctly because HandleHTMXChat calls c.Render
	// Wait, the test runs inside internal/handlers, so the path is ../../web/templates.
	engine := html.New("../../web/templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Post("/chat", HandleHTMXChat)

	body := "prompt=tell%20me%20a%20joke"
	req := httptest.NewRequest("POST", "/chat", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, int(2*time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(respBody), "Mocked response") {
		t.Errorf("Expected 'Mocked response' in HTML, got %s", string(respBody))
	}
}

func TestHandleHTMXChat_EmptyPrompt(t *testing.T) {
	app := fiber.New()
	app.Post("/chat", HandleHTMXChat)

	req := httptest.NewRequest("POST", "/chat", bytes.NewBufferString("prompt="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, int(2*time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", resp.StatusCode)
	}
}
