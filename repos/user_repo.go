package repos

import (
	"context"
	"my-go-backend/internal"
	"my-go-backend/models"
	"time"
)

type UserRepository interface {
	Find(userId string) (*models.User, error)
	Create(userId string, name string) (*models.User, error)
}

type userRepositoryImpl struct {
	firestore *internal.FirestoreClient
}

func (u *userRepositoryImpl) Find(userId string) (*models.User, error) {
	return nil, nil
}

func (u *userRepositoryImpl) Create(userId string, name string) (*models.User, error) {
	ctx := context.Background()
	_, err := u.firestore.Get().
		Collection("users").
		Doc(userId).
		Set(ctx, map[string]any{
			"name":      name,
			"uid":       userId,
			"createdAt": time.Now(),
		})

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name: name,
		Uid:  userId,
	}
	return user, nil
}

func NewUserRepository(firestore *internal.FirestoreClient) UserRepository {
	return &userRepositoryImpl{
		firestore: firestore,
	}
}
