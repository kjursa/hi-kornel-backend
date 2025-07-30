package handlers

import (
	"errors"
	"my-go-backend/repos"
	"my-go-backend/response"

	"github.com/gofiber/fiber/v2"
)

type ContactHandler struct {
	ContactRepo repos.ContactRepository
}

func NewContactHandler(repo repos.ContactRepository) *ContactHandler {
	return &ContactHandler{
		ContactRepo: repo,
	}
}

type sendMessageBody struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Project string `json:"project"`
}

func (ch *ContactHandler) SendMessage(c *fiber.Ctx) error {
	var body sendMessageBody

	if err := c.BodyParser(&body); err != nil {
		return response.InvalidJsonError(c)
	}

	if err := validate(body); err != nil {
		return response.BadRequestError(c, err)
	}

	message, err := ch.ContactRepo.Save(body.Name, body.Email, body.Project)
	if err != nil {
		return response.BadRequestError(c, err)
	}

	return c.JSON(message)
}

func (ch *ContactHandler) GetMessages(c *fiber.Ctx) error {
	messages, err := ch.ContactRepo.Load()
	if err != nil {
		return response.BadRequestError(c, err)
	}

	return c.JSON(messages)
}

func validate(body sendMessageBody) error {
	if len(body.Email) == 0 {
		return errors.New("invalid imput")
	}
	if len(body.Name) == 0 {
		return errors.New("invalid imput")
	}
	if len(body.Project) == 0 {
		return errors.New("invalid imput")
	}
	return nil
}
