package services

import (
	"errors"
	"my-go-backend/models"
	"my-go-backend/repos"
	"my-go-backend/utils"
)

type AuthService interface {
	Register(name string) (*models.UserToken, error)
	RefreshToken(access string, refresh string) (*models.Token, error)
}

type authServiceImpl struct {
	userService  UserService
	tokenManager utils.TokenManager
	tokenRepo    repos.TokenRepository
}

func (t *authServiceImpl) Register(name string) (*models.UserToken, error) {
	user, err := t.userService.CreateUser(name)
	if err != nil {
		return nil, err
	}

	accessToken := t.tokenManager.GenerateAccessToken(user.Uid)
	refreshToken := t.tokenManager.GenerateRefreshToken()

	t.tokenRepo.Save(user.Uid, refreshToken)

	userToken := models.UserToken{
		User: *user,
		Token: models.Token{
			Access:  accessToken,
			Refresh: refreshToken,
		},
	}

	return &userToken, nil
}

func (t *authServiceImpl) RefreshToken(access string, refresh string) (*models.Token, error) {
	userId, err := t.tokenManager.VerifyAccessTokenSignature(access)
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	currentRefresh, err := t.tokenRepo.Get(userId)
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	if refresh != currentRefresh {
		return nil, errors.New("invalid refresh token")
	}

	accessToken := t.tokenManager.GenerateAccessToken(userId)
	refreshToken := t.tokenManager.GenerateRefreshToken()

	t.tokenRepo.Save(userId, refreshToken)

	userToken := &models.Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}

	return userToken, nil
}

func NewAuthService(userService UserService, tokenRepo repos.TokenRepository, tokenManager utils.TokenManager) AuthService {
	return &authServiceImpl{
		userService:  userService,
		tokenManager: tokenManager,
		tokenRepo:    tokenRepo,
	}
}
