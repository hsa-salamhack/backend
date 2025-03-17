package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hsa-salamhack/backend/database"
)

func init() {
	Register(Route{
		Name:   "/messages/:chat_id",
		Method: "GET",
		Run:    getMessagesHandler,
	})
}

type MessageResponse struct {
	ID        uint   `json:"id"         example:"2"`
	ChatID    string `json:"chat_id"    example:"123e4567-e89b-12d3-a456-426614174000"`
	Content   string `json:"content"    example:"What is the capital of France?"`
	Role      string `json:"role"       example:"User"`
	CreatedAt string `json:"created_at" example:"2025-03-17T23:23:46Z"`
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
func getMessagesHandler(c *fiber.Ctx) error {
	database.Connect()
	chatID := c.Params("chat_id")

	var messages []database.Message
	database.DB.Where("chat_id = ?", chatID).Find(&messages)

	return c.JSON(fiber.Map{"messages": messages})
}
