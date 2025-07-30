package repos

import (
	"context"
	"my-go-backend/internal"
	"my-go-backend/models"
	"my-go-backend/utils"
	"sort"
	"time"

	"google.golang.org/api/iterator"
)

type ChatRepository interface {
	Find(userId string, chatId string, lastMessages []string) ([]models.Chat, error)
	Create(userId string, chatId string, question string, answer string) (*models.Chat, error)
}

type chatRepositoryImpl struct {
	firestore *internal.FirestoreClient
}

func (ch *chatRepositoryImpl) Find(userId string, chatId string, lastMessages []string) ([]models.Chat, error) {
	query := ch.firestore.Get().
		Collection("user-chats").
		Doc(userId).
		Collection("chats").
		Doc(chatId).
		Collection("messages").
		Where("messageId", "in", lastMessages)

	ctx := context.Background()
	iter := query.Documents(ctx)
	defer iter.Stop()

	var chats []models.Chat
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []models.Chat{}, nil
		}

		chat := models.Chat{
			Question:  doc.Data()["question"].(string),
			Answer:    doc.Data()["answer"].(string),
			MessageId: doc.Data()["messageId"].(string),
			ChatId:    doc.Data()["chatId"].(string),
			Timestamp: doc.Data()["createdAt"].(int64),
		}
		chats = append(chats, chat)
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[i].Timestamp >= chats[j].Timestamp
	})

	return chats, nil
}

func (ch *chatRepositoryImpl) Create(userId string, chatId string, question string, answer string) (*models.Chat, error) {
	messageId := utils.GenerateID()
	now := time.Now().Unix()
	ctx := context.Background()
	_, err := ch.firestore.Get().
		Collection("user-chats").
		Doc(userId).
		Collection("chats").
		Doc(chatId).
		Collection("messages").
		Doc(messageId).
		Set(ctx, map[string]any{
			"ownerId":   userId,
			"chatId":    chatId,
			"messageId": messageId,
			"question":  question,
			"answer":    answer,
			"createdAt": now,
		})

	if err != nil {
		return nil, err
	}

	chat := &models.Chat{
		ChatId:    chatId,
		MessageId: messageId,
		Question:  question,
		Answer:    answer,
		Timestamp: now,
	}
	return chat, nil
}

func NewChatRepository(firestore *internal.FirestoreClient) ChatRepository {
	return &chatRepositoryImpl{
		firestore: firestore,
	}
}
