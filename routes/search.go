package routes

import (
	"strings"
	"sync"

	"github.com/abadojack/whatlanggo"
	"github.com/gofiber/fiber/v2"
	gowiki "github.com/trietmn/go-wiki"
)

var searchCache sync.Map

type cacheEntryy struct {
	results []SearchResult
}

func init() {
	Register(Route{
		Name:   "/search",
		Method: "GET",
		Run:    search,
	})
}

type SearchResult struct {
	Title string `json:"title"   example:"French Revolution"`
	Sum   string `json:"summary" example:"The French Revolution was a period of radical political and societal change in France..."`
	URL   string `json:"url"     example:"/wiki/en/French_Revolution"`
	Lang  string `json:"lang"    example:"en"`
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
func search(c *fiber.Ctx) error {
	query := c.Query("q")
	lang := c.Query("lang")

	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	var langMap = map[string]string{
		"eng": "en", "ara": "ar", "fra": "fr", "deu": "de", "per": "ar",
	}

	if lang == "" {
		info := whatlanggo.Detect(query)
		langCode := info.Lang.Iso6393()
		if val, ok := langMap[langCode]; ok {
			lang = val
		} else {
			lang = "en"
		}
	}

	cacheKey := query + "|" + lang
	if entry, found := searchCache.Load(cacheKey); found {
		return c.JSON(entry.(*cacheEntryy).results)
	}

	gowiki.SetLanguage(lang)
	searchResult, _, err := gowiki.Search(query, 3, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Wikipedia\n" + err.Error(),
		})
	}

	var results []SearchResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, result := range searchResult {
		if strings.Contains(result, "(disambiguation)") || strings.Contains(result, "(توضيح)") {
			continue
		}
		wg.Add(1)

		go func(res string) {
			defer wg.Done()
			sum, _ := gowiki.Summary(res, 2, -1, false, true)
			cutSum := strings.Split(sum, "==")[0]

			searchRes := SearchResult{
				Title: res,
				Sum:   strings.TrimSpace(cutSum),
				URL:   "/wiki/" + lang + "/" + res,
				Lang:  lang,
			}

			mu.Lock()
			results = append(results, searchRes)
			mu.Unlock()
		}(result)
	}

	wg.Wait()

	searchCache.Store(cacheKey, &cacheEntryy{results: results})
	return c.JSON(results)
}
