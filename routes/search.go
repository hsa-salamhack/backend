package routes

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	gowiki "github.com/trietmn/go-wiki"
)

func init() {
	Register(Route{
		Name:   "/search",
		Method: "GET",
		Run:    searchHandler,
	})
}

type SearchResult struct {
	Title string `json:"title"   example:"French Revolution"`
	Sum   string `json:"summary" example:"The French Revolution was a period of radical political and societal change in France..."`
	URL   string `json:"url"     example:"/wiki/en/French_Revolution"`
}

// @Summary Search Wikipedia
// @Description Searches Wikipedia for articles matching the query
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param lang query string false "Language code (default: en)" default:"en"
// @Success 200 {array} SearchResult "Search results"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /search [get]
func searchHandler(c *fiber.Ctx) error {
	query := c.Query("q")
	lang := c.Query("lang")

	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	gowiki.SetLanguage(lang)
	search_result, _, err := gowiki.Search(query, 3, false)
	if err != nil {
		fmt.Println(err)
	}

	var results []SearchResult
	for _, result := range search_result {
		if strings.Contains(result, "(disambiguation)") {
			continue
		}
		sum, _ := gowiki.Summary(result, 2, -1, false, true)
		cutSum := strings.Split(sum, "==")[0]

		results = append(results, SearchResult{
			Title: result,
			Sum:   strings.TrimSpace(cutSum),
			URL:   "/wiki/" + lang + "/" + result,
		})
	}

	return c.JSON(results)
}
