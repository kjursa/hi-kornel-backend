package services

import (
	"errors"
	"my-go-backend/models"
	"my-go-backend/repos"
	"my-go-backend/utils"
)

var ErrInvalidName = errors.New("invalid name")

type UserService interface {
	CreateUser(name string) (*models.User, error)
}

type userServiceImpl struct {
	userRepository repos.UserRepository
}

func (u *userServiceImpl) CreateUser(name string) (*models.User, error) {
	if !validateName(name) {
		return nil, ErrInvalidName
	}
	return u.userRepository.Create(utils.GenerateUserID(), name)
}

func validateName(name string) bool {
	return len(name) >= 1 && len(name) <= 20
}

func NewUserService(userRepository repos.UserRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}
