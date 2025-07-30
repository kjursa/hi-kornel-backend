package response

import "github.com/gofiber/fiber/v2"

func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

func BadRequestError(c *fiber.Ctx, err error) error {
	return Error(c, fiber.StatusBadRequest, err.Error())
}

func InvalidJsonError(c *fiber.Ctx) error {
	return Error(c, fiber.StatusBadRequest, "invalid Json")
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}
