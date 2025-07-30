package middleware

import (
	"my-go-backend/utils"

	"github.com/gofiber/fiber/v2"
)

var TokenHeaderKey = "token"
var UserIdKey = "userId"
var AdminUserId = "admin"

type SessionMiddleware struct {
	tokenManager utils.TokenManager
}

func NewSessionMiddleware(tokenManager utils.TokenManager) *SessionMiddleware {
	return &SessionMiddleware{
		tokenManager: tokenManager,
	}
}

func (s *SessionMiddleware) HandleToken(c *fiber.Ctx) error {
	token := c.Get(TokenHeaderKey)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing token",
		})
	}

	userId, err := s.tokenManager.VerifyAccessToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Locals(UserIdKey, userId)

	return c.Next()
}

func (s *SessionMiddleware) HandleAdminToken(c *fiber.Ctx) error {
	token := c.Get(TokenHeaderKey)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing token",
		})
	}

	userId, err := s.tokenManager.VerifyAccessToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if userId != AdminUserId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing permission",
		})
	}
	c.Locals(UserIdKey, userId)

	return c.Next()
}
