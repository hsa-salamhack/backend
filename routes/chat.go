package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/hsa-salamhack/backend/database"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var (
	aiClient  *genai.Client
	once      sync.Once
	wikiCache = make(map[string]string) // Simple in-memory cache for Wikipedia articles
)

func init() {
	once.Do(func() {
		godotenv.Load()
		database.Connect()
		var err error
		aiClient, err = genai.NewClient(
			context.Background(),
			option.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize AI client: %v", err))
		}
	})
	Register(Route{
		Name:   "/chat",
		Method: "POST",
		Run:    chatHandler,
	})
}

// Chat represents the request payload for the chat endpoint
type Chat struct {
	ID      string `json:"uuid"    example:"123e4567-e89b-12d3-a456-426614174000"`
	Message string `json:"message" example:"What is the capital of France?"`
	Model   string `json:"model"   example:"conservative"`
	Lang    string `json:"lang"    example:"en"`
	Wiki    string `json:"wiki"    example:"French_Revolution"`
}

// ChatResponse represents the response from the chat endpoint
type ChatResponse struct {
	Message string `json:"message" example:"Paris is the capital of France."`
}

type CacheStruct struct {
	FullBody string `json:"full_body"`
}

// @Summary Chat with AI
// @Description Sends a message to the AI model with a Wikipedia context.
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body Chat true "Chat Request Body"
// @Success 200 {object} ChatResponse "Successful response"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /chat [post]
func chatHandler(c *fiber.Ctx) error {
	var chat Chat
	if err := c.BodyParser(&chat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	database.DB.Create(&database.Message{
		ChatID:  chat.ID,
		Content: chat.Message,
		Role:    "User",
		Topic:   chat.Wiki,
	})

	wikiContent, err := fetchWiki(chat.Lang, chat.Wiki)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Failed to fetch Wikipedia data"})
	}

	sysint := fmt.Sprintf(
		"You are discussing the topic of %s with the user.\n"+
			"Please adopt the perspective of %s throughout our conversation...\n\n"+
			"The article is:\n\n %s",
		chat.Wiki,
		chat.Model,
		wikiContent,
	)

	model := aiClient.GenerativeModel("gemini-2.0-flash")
	model.SystemInstruction = genai.NewUserContent(genai.Text(sysint))

	resp, err := model.GenerateContent(context.Background(), genai.Text(chat.Message))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	database.DB.Create(&database.Message{
		ChatID:  chat.ID,
		Content: content,
		Role:    "AI",
	})

	return c.JSON(fiber.Map{"message": content})
}

func fetchWiki(lang, topic string) (string, error) {
	cacheKey := lang + ":" + topic
	if data, exists := wikiCache[cacheKey]; exists {
		return data, nil
	}

	respo, err := http.Get("http://localhost:5050/wiki/" + lang + "/" + topic)
	if err != nil {
		return "", err
	}
	defer respo.Body.Close()

	body, err := io.ReadAll(respo.Body)
	if err != nil {
		return "", err
	}

	var data CacheStruct
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	wikiCache[cacheKey] = data.FullBody
	return data.FullBody, nil
}
