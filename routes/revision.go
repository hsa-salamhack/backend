package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type WikiRaw struct {
	Query struct {
		Pages map[string]struct {
			Revisions []struct {
				User      string `json:"user"`
				Timestamp string `json:"timestamp"`
				Comment   string `json:"comment"`
				Anon      string `json:"anon,omitempty"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

type Revision struct {
	User      string `json:"user"`
	Timestamp int64  `json:"timestamp"`
	Comment   string `json:"comment"`
	Action    string `json:"action"`
}

func CleanUp(comment string) string {
	re := regexp.MustCompile(`\[\[.*?\]\]|\[\w+:\w+\||\(|\)|\[\w+:[\w/]+\|`)
	comment = re.ReplaceAllString(comment, "")

	comment = strings.ReplaceAll(comment, "[", "")
	comment = strings.ReplaceAll(comment, "]", "")

	return strings.TrimSpace(comment)
}

func DetectAction(comment string) string {
	comment = strings.ToLower(comment)

	actionPatterns := map[string][]string{
		"revert":   {"revert", "reverted", "reverting", "rv", "rvv"},
		"remove":   {"remove", "removed", "removing", "deletion", "deleted", "deleting"},
		"add":      {"add", "added", "adding", "insertion", "inserted", "inserting"},
		"update":   {"update", "updated", "updating", "change", "changed", "changing"},
		"format":   {"format", "formatted", "formatting", "style", "styled", "styling"},
		"fix":      {"fix", "fixed", "fixing", "correct", "corrected", "correcting", "typo"},
		"merge":    {"merge", "merged", "merging"},
		"redirect": {"redirect", "redirected", "redirecting"},
		"cleanup":  {"cleanup", "cleaned", "cleaning"},
	}

	for action, patterns := range actionPatterns {
		for _, pattern := range patterns {
			if strings.Contains(comment, pattern) {
				return action
			}
		}
	}

	return "edit"
}

func init() {
	Register(Route{
		Name:   "/revision",
		Method: "GET",
		Run: func(c *fiber.Ctx) error {
			query := c.Query("q")
			lang := c.Query("lang")

			if query == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Query parameter 'q' is required",
				})
			}

			if lang == "" {
				lang = "en"
			}

			resp, err := http.Get(
				"https://" + lang + ".wikipedia.org/w/api.php?action=query&format=json&prop=revisions&rvlimit=100&titles=" + query,
			)
			if err != nil {
				fmt.Println("Error fetching data:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to fetch data from Wikipedia",
				})
			}
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to read response data",
				})
			}

			var input WikiRaw
			if err := json.Unmarshal(data, &input); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to parse Wikipedia API response",
				})
			}

			result := []Revision{}
			for _, page := range input.Query.Pages {
				for _, rev := range page.Revisions {
					user := rev.User
					if ipRegex := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`); ipRegex.MatchString(
						user,
					) {
						parts := strings.Split(user, ".")
						if len(parts) == 4 {
							user = fmt.Sprintf("%s...%s", parts[0], parts[3])
						}
					}

					comment := CleanUp(rev.Comment)
					timestamp, _ := time.Parse(time.RFC3339, rev.Timestamp)
					action := DetectAction(comment)

					result = append(result, Revision{
						User:      user,
						Timestamp: timestamp.Unix(),
						Comment:   comment,
						Action:    action,
					})
				}
			}

			return c.JSON(result)
		},
	})
}
