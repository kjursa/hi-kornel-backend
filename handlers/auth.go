package handlers

import (
	"my-go-backend/response"
	"my-go-backend/services"
	"my-go-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService  services.AuthService
	tokenManager utils.TokenManager
}

func NewAuthHandler(
	authService services.AuthService,
	userService services.UserService,
	tokenManager utils.TokenManager) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenManager: tokenManager,
	}
}

type createUserRequest struct {
	Name string `json:"name"`
}

type RegisterResponse struct {
	Name         string `json:"name"`
	UserId       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

// POST /users
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.InvalidJsonError(c)
	}

	userToken, err := h.authService.Register(req.Name)
	if err != nil {
		return response.BadRequestError(c, err)
	}

	res := RegisterResponse{
		Name:         userToken.User.Name,
		UserId:       userToken.User.Uid,
		RefreshToken: userToken.Token.Refresh,
		AccessToken:  userToken.Token.Access,
	}

	return c.JSON(res)
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req refreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.InvalidJsonError(c)
	}

	token, err := h.authService.RefreshToken(req.AccessToken, req.RefreshToken)
	if err != nil {
		return response.BadRequestError(c, err)
	}

	res := RefreshTokenResponse{
		RefreshToken: token.Refresh,
		AccessToken:  token.Access,
	}

	return c.JSON(res)
}
