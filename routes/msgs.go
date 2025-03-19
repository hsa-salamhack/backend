package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hsa-salamhack/backend/database"
)

func init() {
	Register(Route{
		Name:   "/messages/:chat_id",
		Method: "GET",
		Run:    handler,
	})
}

type MessageResponse struct {
	ID      uint   `json:"id"      example:"2"`
	ChatID  string `json:"chat_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Content string `json:"content" example:"What is the capital of France?"`
	Role    string `json:"role"    example:"User"`
	Topic   string `json:"topic"   example:"French_Revolution"`
}

type MessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
}

// @Summary Get chat messages
// @Description Retrieves all messages for a specific chat ID
// @Tags Messages
// @Accept json
// @Produce json
// @Param chat_id path string true "Chat ID"
// @Success 200 {object} MessagesResponse "List of messages"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /messages/{chat_id} [get]
func handler(c *fiber.Ctx) error {
	database.Connect()
	chatID := c.Params("chat_id")

	var messages []database.Message
	database.DB.Where("chat_id = ?", chatID).Find(&messages)

	if len(messages) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Chat not found",
		})
	}

	return c.JSON(fiber.Map{"messages": messages})
}
