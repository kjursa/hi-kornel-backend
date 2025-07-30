package handlers

import (
	"errors"
	"my-go-backend/response"
	"my-go-backend/services"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	chatService services.ChatService
}

func NewChatHandler(chatService services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

type askQuestionBody struct {
	Question  string   `json:"question"`
	LastChats []string `json:"last_chats"`
}

func (h *ChatHandler) AskQuestion(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	if userId == "" {
		return response.BadRequest(c, "missing param")
	}
	chatId := c.Params("chatId")
	if chatId == "" {
		return response.BadRequest(c, "missing param")
	}
	var body askQuestionBody
	if err := c.BodyParser(&body); err != nil {
		return response.InvalidJsonError(c)
	}
	if err := validateBody(body); err != nil {
		return response.BadRequestError(c, err)
	}

	res, err := h.chatService.AnswerQuestion(userId, chatId, body.Question, body.LastChats)
	if err != nil {
		return response.BadRequestError(c, err)
	}
	return c.JSON(res)
}

func validateBody(body askQuestionBody) error {
	if len(body.Question) == 0 {
		return errors.New("question is empty")
	}
	if len(body.LastChats) >= 10 {
		return errors.New("too many params")
	}
	return nil
}
