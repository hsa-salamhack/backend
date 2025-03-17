package routes

import (
	"context"
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
		Run: func(c *fiber.Ctx) error {
			godotenv.Load()
			database.Connect()

			ctx := context.Background()

			client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
			if err != nil {
				return err
			}

			defer client.Close()

			type Chat struct {
				ID      string `json:"uuid"`
				Message string `json:"message"`
				Model   string `json:"model"`
				Sum     string `json:"summary"`
			}

			var chat Chat
			if err := c.BodyParser(&chat); err != nil {
				return err
			}

			user_message := database.Message{ChatID: chat.ID, Content: chat.Message, Role: "User"}
			database.DB.Create(&user_message)

			modelName := "gemini-1.5-flash"
			model := client.GenerativeModel(modelName)
			model.SystemInstruction = genai.NewUserContent(
				genai.Text(
					"You are an expert at this topic:\n\n" + chat.Sum + "\n, talk from the prespective of " + chat.Model,
				),
			)

			resp, err := model.GenerateContent(ctx, genai.Text(chat.Message))
			if err != nil {
				return err
			}

			airesp := resp.Candidates[0].Content.Parts[0]

			aimessage := database.Message{
				ChatID:  chat.ID,
				Content: airesp,
				Role:    "AI",
			}
			database.DB.Create(&aimessage)

			return c.JSON(fiber.Map{"message": airesp})
		},
	})
}
