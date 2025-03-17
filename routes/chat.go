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
				Lang    string `json:"lang"`
				Wiki    string `json:"wiki"`
			}

			var chat Chat
			if err := c.BodyParser(&chat); err != nil {
				return err
			}

			user_message := database.Message{ChatID: chat.ID, Content: chat.Message, Role: "User"}
			database.DB.Create(&user_message)

			respo, _ := http.Get("http://localhost:5050/wiki/" + chat.Lang + "/" + chat.Wiki)
			defer respo.Body.Close()
			body, _ := io.ReadAll(respo.Body)

			var data map[string]interface{}
			json.Unmarshal(body, &data)

			sysint := fmt.Sprintf(
				"You are discussing the topic of %s with the user.\n"+
					"Please adopt the perspective of %s throughout our conversation, including this perspective's political views and ideological framework.\n"+
					"Present arguments consistent with this viewpoint, and defend the positions naturally while remaining open to changing your stance if I present compelling counterarguments.\n\n"+
					"This conversation should feel like I'm speaking with someone who genuinely holds these views, not just someone roleplaying.\n"+
					"Respond to my points directly, ask thoughtful follow-up questions, and engage with the nuances of this perspective.\n\n"+
					"Stay in character even if the discussion becomes challenging, but maintain respectful discourse throughout.\n\n"+
					"Note that this information was retrieved from Wikipedia, which may contain certain biases or limitations in perspective.\n"+
					"Feel free to acknowledge these potential biases when relevant to our discussion.\n\n"+
					"The article is:\n\n %s",
				chat.Wiki,
				chat.Model,
				data["full_body"].(string))
			modelName := "gemini-1.5-flash"
			model := client.GenerativeModel(modelName)
			model.SystemInstruction = genai.NewUserContent(
				genai.Text(sysint),
			)

			resp, err := model.GenerateContent(ctx, genai.Text(chat.Message))
			if err != nil {
				return err
			}

			content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

			aimessage := database.Message{
				ChatID:  chat.ID,
				Content: content,
				Role:    "AI",
			}
			database.DB.Create(&aimessage)

			return c.JSON(fiber.Map{"message": content})
		},
	})
}
