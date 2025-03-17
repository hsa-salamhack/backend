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
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func init() {
	Register(Route{
		Name:   "/summary",
		Method: "POST",
		Run: func(c *fiber.Ctx) error {
			godotenv.Load()
			ctx := context.Background()

			client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
			if err != nil {
				return err
			}

			defer client.Close()

			type Article struct {
				Lang string `json:"lang"`
				Wiki string `json:"wiki"`
			}

			var wiki Article
			if err := c.BodyParser(&wiki); err != nil {
				return err
			}

			respo, _ := http.Get("http://localhost:5050/wiki/" + wiki.Lang + "/" + wiki.Wiki)
			defer respo.Body.Close()
			body, _ := io.ReadAll(respo.Body)

			var data map[string]interface{}
			json.Unmarshal(body, &data)

			sysint := "You're a wikipedia summerizer, summerize the info inside the article text i give you while keeping all info intact, Keep it short under 3 paragraphs don't do more. avoid removing any or editing any:\n\n" + data["full_body"].(string)
			modelName := "gemini-1.5-flash"
			model := client.GenerativeModel(modelName)

			resp, err := model.GenerateContent(ctx, genai.Text(sysint))
			if err != nil {
				return err
			}

			content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
			return c.JSON(fiber.Map{"message": content})
		},
	})
}
