package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hsa-salamhack/backend/database"
)

func init() {
	Register(Route{
		Name:   "/messages/:chat_id",
		Method: "GET",
		Run: func(c *fiber.Ctx) error {
			database.Connect()
			chatID := c.Params("chat_id")

			var messages []database.Message
			database.DB.Where("chat_id = ?", chatID).Find(&messages)

			return c.JSON(fiber.Map{"messages": messages})
		},
	})
}
