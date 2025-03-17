package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func init() {
	Register(Route{
		Name:   "/wiki/:lang/:wiki",
		Method: "GET",
		Run:    wikiHandler,
	})
}

// WikiResponse represents the structure of the Wikipedia article response
type WikiResponse struct {
	Title    string `json:"title"     example:"French Revolution"`
	FullBody string `json:"full_body" example:"The French Revolution was a period of radical political and societal change in France..."`
}

// @Summary Get Wikipedia article
// @Description Retrieves a Wikipedia article by language and article name
// @Tags Wiki
// @Accept json
// @Produce json
// @Param lang path string true "Language code" example:"en"
// @Param wiki path string true "Wikipedia article name" example:"French_Revolution"
// @Success 200 {object} WikiResponse "Wikipedia article content"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Article not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /wiki/{lang}/{wiki} [get]
func wikiHandler(c *fiber.Ctx) error {
	lang := c.Params("lang")
	wiki := c.Params("wiki")

	respo, err := http.Get("http://localhost:5050/wiki/" + lang + "/" + wiki)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch article",
		})
	}
	defer respo.Body.Close()

	if respo.StatusCode == 404 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Article not found",
		})
	}

	body, err := io.ReadAll(respo.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse response",
		})
	}

	return c.JSON(data)
}
