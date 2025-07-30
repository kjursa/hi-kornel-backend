package repos

import (
	"context"
	"fmt"
	"my-go-backend/internal"
	"sync"
	"time"
)

type TokenRepository interface {
	Save(userID, token string) error
	Get(userID string) (string, error)
}

type tokenRepositoryImpl struct {
	mu        sync.RWMutex
	firestore *internal.FirestoreClient
}

func NewTokenRepository(firestore *internal.FirestoreClient) *tokenRepositoryImpl {
	return &tokenRepositoryImpl{
		firestore: firestore,
	}
}

func (t *tokenRepositoryImpl) Save(userId, token string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().Unix()
	ctx := context.Background()
	_, err := t.firestore.Get().
		Collection("tokens").
		Doc(userId).
		Set(ctx, map[string]any{
			"ownerId":      userId,
			"refreshToken": token,
			"updatedAt":    now,
		})

	return err
}

func (t *tokenRepositoryImpl) Get(userId string) (string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	ctx := context.Background()
	doc, err := t.firestore.Get().
		Collection("tokens").
		Doc(userId).
		Get(ctx)

	if err != nil {
		return "", err
	}

	token := doc.Data()["refreshToken"].(string)
	if token == "" {
		return "", ErrTokenNotFound
	}

	return token, nil
}

var ErrTokenNotFound = fmt.Errorf("refresh token not found")
