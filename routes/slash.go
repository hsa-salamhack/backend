package routes

import "github.com/gofiber/fiber/v2"

func init() {
	Register(Route{
		Name:   "/chat",
		Method: "GET",
		Run: func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "Test route"})
		},
	})
}
