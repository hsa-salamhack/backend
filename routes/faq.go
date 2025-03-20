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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var (
	gClient    *genai.Client
	onceee     sync.Once
	wikiCachez = make(map[string]string)
)

func init() {
	onceee.Do(func() {
		godotenv.Load()
		var err error
		gClient, err = genai.NewClient(
			context.Background(),
			option.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
		)
		if err != nil {
			fmt.Sprintf("Failed to initialize AI client: %v", err)
		}
		go resetCachee()
	})
	Register(Route{
		Name:   "/prompts",
		Method: "POST",
		Run:    promptHandler,
	})
}

// WikiArticle represents the request payload for the wiki endpoint
type WikiArticle struct {
	Lang string `json:"lang" example:"en"`
	Wiki string `json:"wiki" example:"French_Revolution"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// WikiRespons represents the response from the wiki endpoint
type WikiRespons struct {
	Questions []string `json:"questions" example:"["What were the main causes of the French Revolution?", "How did the Revolution impact European politics?", "Was the Reign of Terror justified?"]"`
}

type cachezStruct struct {
	FullBody string `json:"full_body"`
}

// resetCachee clears the cache every 5 minutes
func resetCachee() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		wikiCachez = make(map[string]string)
	}
}

// fetchWikii retrieves wiki content either from cache or from the wiki service
func fetchWikii(lang, topic string) (string, error) {
	cacheKey := lang + ":" + topic
	if data, exists := wikiCachez[cacheKey]; exists {
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

	var data cachezStruct
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	wikiCachez[cacheKey] = data.FullBody
	return data.FullBody, nil
}

// @Summary Generate discussion questions
// @Description Generates 3 thought-provoking questions based on a Wikipedia article
// @Tags Wiki
// @Accept json
// @Produce json
// @Param request body WikiArticle true "Wiki Request Body"
// @Success 200 {object} WikiRespons "Successful response with generated questions"
// @Failure 400 {object} ErrorResponse "Bad request when payload is invalid"
// @Failure 500 {object} ErrorResponse "Internal server error when AI processing fails"
// @Router /prompts [post]
func promptHandler(c *fiber.Ctx) error {
	var wiki WikiArticle
	if err := c.BodyParser(&wiki); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	wikiContent, err := fetchWikii(wiki.Lang, wiki.Wiki)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Failed to fetch Wikipedia data " + err.Error()})
	}

	sysint := fmt.Sprintf(
		`
Create a focused "You may ask" JSON array with exactly 3 thought-provoking questions based on the following article content. Extract only the most significant questions that would spark meaningful discussion about the core themes or controversies in the text.

Requirements:
- Include exactly 3 questions (no more, no less)
- Focus on questions that directly relate to the article's central arguments or claims
- Include at least one question that invites debate or critical analysis
- Formulate questions that would naturally arise for readers seeking deeper understanding
- Please KEEP the questions SHORT, 5-8 words maximum.
- Avoid basic factual questions in favor of those requiring analysis or opinion
- Use %s langauge
Format the JSON array as follows:
{ questions: {
["Question that addresses a key argument in the article?", "Question that examines a controversial claim or position?", "Question that explores broader implications of the article's content?"] }}

Article content will be provided
`, wiki.Lang)

	model := gClient.GenerativeModel("gemini-2.0-flash")
	model.SystemInstruction = genai.NewUserContent(genai.Text(sysint))

	resp, err := model.GenerateContent(context.Background(), genai.Text(wikiContent))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	return c.JSON(result)
}
