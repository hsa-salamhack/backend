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
		Run: func(c *fiber.Ctx) error {
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

			type Result struct {
				Title string `json:"title"`
				Sum   string `json:"summary"`
				URL   string `json:"url"`
			}

			var results []Result
			for _, result := range search_result {
				if strings.Contains(result, "(disambiguation)") {
					continue
				}
				sum, _ := gowiki.Summary(result, 2, -1, false, true)
				cutSum := strings.Split(sum, "==")[0]

				results = append(results, Result{
					Title: result,
					Sum:   strings.TrimSpace(cutSum),
					URL:   "/wiki/" + lang + "/" + result,
				})
			}

			return c.JSON(results)
		},
	})
}
