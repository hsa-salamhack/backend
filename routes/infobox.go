package routes

import (
	"context"
	"encoding/json"
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
		Name:   "/infobox",
		Method: "POST",
		Run:    InfoHandler,
	})
}

type ArticleO struct {
	Lang string `json:"lang" example:"en"`
	Wiki string `json:"wiki" example:"French_Revolution"`
}

type InfoboxRes struct {
	Infobox string `json:"infobox" example:"{name: 'French', leader: 'Bonaparte'}"`
}

// @Summary Generate article infobox
// @Description Generates a concise infobox using AI
// @Tags Wiki
// @Accept json
// @Produce json
// @Param request body ArticleO true "Article Request Body"
// @Success 200 {object} InfoboxRes "Successful response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /infobox [post]
func InfoHandler(c *fiber.Ctx) error {
	godotenv.Load()
	apiKey := os.Getenv("GEMINI_API_KEY")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gemini initialization failed")
	}
	defer client.Close()

	var wiki Article
	if err := c.BodyParser(&wiki); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	url := fmt.Sprintf("http://localhost:5050/wiki/%s/%s", wiki.Lang, wiki.Wiki)
	respo, err := http.Get(url)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch article")
	}
	defer respo.Body.Close()

	body, err := io.ReadAll(respo.Body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read article response")
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Invalid JSON response from wiki service",
		)
	}

	articleContent, ok := data["full_body"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid article content format")
	}

	sysint := fmt.Sprintf(
		`Create a comprehensive Wikipedia-style infobox in JSON format based on the following article content. Extract only the most significant and interesting information that would typically appear in a Wikipedia infobox.

Format the JSON with:
- Precise, concise labels
- Nested objects for related data
- Arrays for multiple related items
- UNIX timestamps for dates
- Language: %s

If a Wikimedia image is available, include the URL. Example JSON structure:

{
  "countryInfo": {
    "anthem": "Bilady, Bilady, Bilady",
    "capital": "Cairo",
    "officialLanguage": "Arabic",
    "currency": "Egyptian pound (EGP)"
  }
}

I'll now provide the article overview.`,
		wiki.Lang,
	)

	modelName := "gemini-1.5-flash"
	model := client.GenerativeModel(modelName)
	model.SystemInstruction = genai.NewUserContent(genai.Text(sysint))

	resp, err := model.GenerateContent(ctx, genai.Text(articleContent))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate infobox")
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return fiber.NewError(fiber.StatusInternalServerError, "No response from AI model")
	}

	rawResponse := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	cleanedJSON := strings.TrimSpace(rawResponse)
	if strings.HasPrefix(cleanedJSON, "```json") {
		cleanedJSON = strings.TrimPrefix(cleanedJSON, "```json")
	}
	if strings.HasSuffix(cleanedJSON, "```") {
		cleanedJSON = strings.TrimSuffix(cleanedJSON, "```")
	}
	cleanedJSON = strings.TrimSpace(cleanedJSON)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanedJSON), &result); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error parsing AI response JSON")
	}

	return c.JSON(
		fiber.Map{"infobox": result, "images": data["images"], "thumbnail": data["thumbnail"]},
	)
}
