package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var (
	oncce     sync.Once
	aiClient2 *genai.Client
)

func init() {
	oncce.Do(func() {
		godotenv.Load()
		var err error
		aiClient2, err = genai.NewClient(
			context.Background(),
			option.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize AI client: %v", err))
		}
	})
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
	var wiki Article
	if err := c.BodyParser(&wiki); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Fetch Wikipedia content
	respo, err := http.Get("http://localhost:5050/wiki/" + wiki.Lang + "/" + wiki.Wiki)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch Wikipedia data")
	}
	defer respo.Body.Close()

	body, err := io.ReadAll(respo.Body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read Wikipedia response")
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse Wikipedia data")
	}

	// Construct AI prompt
	var sysint string
	if wiki.Section == nil {
		sysint = "Summarize this article in under 3 paragraphs in " + wiki.Lang + ":\n\n" + data["full_body"].(string)
	} else {
		sections, ok := data["sections"].(map[string]interface{})
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "Invalid sections format")
		}

		sectionKey := strings.TrimSpace(*wiki.Section)
		sectionText, exists := sections[sectionKey].(map[string]interface{})
		if !exists {
			return fiber.NewError(fiber.StatusBadRequest, "Section is invalid")
		}

		sysint = "You're a Wikipedia summarizer. Summarize the given article section while keeping all info intact. Keep it under 3 paragraphs. Use " +
			wiki.Lang + ":\n\n" + sectionText["body"].(string)
	}

	// Generate AI response
	model := aiClient2.GenerativeModel("gemini-2.0-flash")
	resp, err := model.GenerateContent(context.Background(), genai.Text(sysint))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return c.JSON(fiber.Map{"message": content})
}
