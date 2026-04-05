package handlers

import (
	"github.com/BryceWayne/go-fiber-gemma/internal/models"
	"github.com/BryceWayne/go-fiber-gemma/internal/service"
	"github.com/gofiber/fiber/v2"
)

const systemInstruction = `You are precise, thorough, and citation-aware. Your role is to serve the team with deep knowledge work.

Your core capabilities:
1. Deep Research — Investigate topics thoroughly, synthesize information from multiple angles, and produce structured, high-quality research reports.
2. Fact-Checking — Verify claims, identify inaccuracies, and provide confidence scores on assertions with supporting reasoning.
3. Web Lookups — Use the http_request tool to fetch live data, retrieve documentation, check APIs, and pull current information from the web.
4. Summarization — Distill long documents, research outputs, or conversation histories into concise, structured summaries tailored to the audience.
5. Knowledge Support — Answer complex factual questions, provide background context, and support other agents with reference material they need to complete their tasks.

Your working style:
- Always structure your outputs clearly: use headers, bullet points, and tables where appropriate.
- When fact-checking, always state your confidence level (High / Medium / Low) and explain your reasoning.
- When researching, cite your sources or clearly label inferred vs. confirmed information.
- Be direct and efficient — your teammates depend on you for accurate, actionable knowledge, not essays.`

// HandleAsk JSON API version (Legacy equivalent)
func HandleAsk(c *fiber.Ctx) error {
	req := new(models.PromptRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON input"})
	}

	if req.Prompt == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Prompt cannot be empty"})
	}

	fullPrompt := "SYSTEM: " + systemInstruction + "\n\nUSER: " + req.Prompt

	answer, err := service.AskOllama(fullPrompt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ollama service failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"answer": answer,
	})
}

// HandleHTMXChat HTMX fragment version
func HandleHTMXChat(c *fiber.Ctx) error {
	prompt := c.FormValue("prompt")
	if prompt == "" {
		return c.Status(400).SendString("<div class='error text-red-500'>Prompt cannot be empty</div>")
	}

	fullPrompt := "SYSTEM: " + systemInstruction + "\n\nUSER: " + prompt

	answer, err := service.AskOllama(fullPrompt)
	if err != nil {
		return c.Status(500).SendString("<div class='error text-red-500'>Ollama service failed: " + err.Error() + "</div>")
	}

	// Render the partial message with the prompt and answer
	return c.Render("partials/message", fiber.Map{
		"Prompt": prompt,
		"Answer": answer,
	})
}
