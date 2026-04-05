package models

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatRequest struct {
	Model    string                 `json:"model"`
	Messages []ChatMessage          `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options"`
}

type OllamaChatResponse struct {
	Message ChatMessage `json:"message"`
}
