package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func init() {
	Register(Route{
		Name:   "/summary",
		Method: "POST",
		Run:    summaryHandler,
	})
}

type Article struct {
	Lang    string  `json:"lang"    example:"en"`
	Wiki    string  `json:"wiki"    example:"French_Revolution"`
	Section *string `json:"section" example:"Causes"`
}

type SummaryResponse struct {
	Message string `json:"message" example:"The French Revolution was a period of radical political and societal change in France..."`
}

// @Summary Generate article summary
// @Description Generates a concise summary of a Wikipedia article
// @Tags Summary
// @Accept json
// @Produce json
// @Param request body Article true "Article Request Body"
// @Success 200 {object} SummaryResponse "Successful response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /summary [post]
func summaryHandler(c *fiber.Ctx) error {
	godotenv.Load()
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gemini cant get up")
	}

	defer client.Close()

	var wiki Article
	if err := c.BodyParser(&wiki); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	respo, _ := http.Get("http://localhost:5050/wiki/" + wiki.Lang + "/" + wiki.Wiki)
	defer respo.Body.Close()
	body, _ := io.ReadAll(respo.Body)

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	var sysint string

	if wiki.Section == nil {
		sysint = "Summarize this article in under 3 paragraphs in " + wiki.Lang + ":\n\n" + data["full_body"].(string)
	} else {
		sections, ok := data["sections"].(map[string]interface{})
		if !ok {
			return errors.New("invalid sections format")
		}

		sectionKey := strings.TrimSpace(*wiki.Section)
		sectionText, exists := sections[sectionKey].(map[string]interface{})
		if !exists {
			return fiber.NewError(fiber.StatusBadRequest, "Section is invalid")
		}

		sysint = "You're a Wikipedia summarizer. Summarize the given article section while keeping all info intact. Keep it under 3 paragraphs. Use " +
			wiki.Lang + ":\n\n" + sectionText["body"].(string)
	}
	modelName := "gemini-2.0-flash"
	model := client.GenerativeModel(modelName)

	resp, err := model.GenerateContent(ctx, genai.Text(sysint))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return c.JSON(fiber.Map{"message": content})
}
