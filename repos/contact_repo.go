package repos

import (
	"context"
	"my-go-backend/internal"
	"my-go-backend/models"
	"time"

	"google.golang.org/api/iterator"
)

type ContactRepository interface {
	Save(name string, email string, project string) (*models.Message, error)
	Load() ([]models.Message, error)
}

type contactRepositoryImpl struct {
	firestore *internal.FirestoreClient
}

func NewContactRepository(firestore *internal.FirestoreClient) *contactRepositoryImpl {
	return &contactRepositoryImpl{
		firestore: firestore,
	}
}

func (c *contactRepositoryImpl) Save(name string, email string, project string) (*models.Message, error) {
	ctx := context.Background()
	_, _, err := c.firestore.Get().Collection("messages").Add(ctx, map[string]interface{}{
		"name":      name,
		"email":     email,
		"project":   project,
		"createdAt": time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &models.Message{
		Name:    name,
		Email:   email,
		Project: project,
	}, nil
}

func (c *contactRepositoryImpl) Load() ([]models.Message, error) {
	ctx := context.Background()
	iter := c.firestore.Get().Collection("messages").Documents(ctx)
	var messages []models.Message
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
		}
		message := models.Message{
			Name:    doc.Data()["name"].(string),
			Email:   doc.Data()["email"].(string),
			Project: doc.Data()["project"].(string),
		}
		messages = append(messages, message)

	}

	return messages, nil
}
