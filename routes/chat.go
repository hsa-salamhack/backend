package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/hsa-salamhack/backend/database"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func init() {
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
	godotenv.Load()
	database.Connect()

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer client.Close()

	var chat Chat
	if err := c.BodyParser(&chat); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userMessage := database.Message{
		ChatID:  chat.ID,
		Content: chat.Message,
		Role:    "User",
		Topic:   chat.Wiki,
	}
	database.DB.Create(&userMessage)

	respo, _ := http.Get("http://localhost:5050/wiki/" + chat.Lang + "/" + chat.Wiki)
	defer respo.Body.Close()
	body, _ := io.ReadAll(respo.Body)

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	sysint := fmt.Sprintf(
		"You are discussing the topic of %s with the user.\n"+
			"Please adopt the perspective of %s throughout our conversation...\n\n"+
			"The article is:\n\n %s",
		chat.Wiki,
		chat.Model,
		data["full_body"].(string),
	)

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SystemInstruction = genai.NewUserContent(genai.Text(sysint))

	resp, err := model.GenerateContent(ctx, genai.Text(chat.Message))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	aiMessage := database.Message{
		ChatID:  chat.ID,
		Content: content,
		Role:    "AI",
	}
	database.DB.Create(&aiMessage)

	return c.JSON(fiber.Map{"message": content})
}
