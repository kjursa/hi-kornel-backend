package services

import (
	"my-go-backend/internal"
	"my-go-backend/models"
	"my-go-backend/repos"
)

type ChatService interface {
	AnswerQuestion(userId string, chatId string, question string, lastChatIds []string) (*models.Chat, error)
}

type chatServiceImpl struct {
	chatRepository repos.ChatRepository
	aiClient       internal.AiClient
}

var systemContent = `
Jesteś profesjonalnym asystentem AI. 
W imieniu Kornela Jursy odpowiadasz na pytania o jego karierze zawodowej, doświadczeniu, umiejętnościach technicznych i osiągnięciach. 
Zawsze odpowiadaj konkretnie, szczerze i profesjonalnie. 
Jeśli to pasuje, możesz podawać przykłady projektów, technologii i wyzwań, które Kornel rozwiązywał. 
Nie zmyślaj faktów.
`

var userContext = `
Kornel Jursa to Senior Android Developer z ponad 10-letnim doświadczeniem. Urodzony w 1988.
Specjalizuje się w Kotlinie, Jetpack Compose, Coroutines i nowoczesnej, modularnej architekturze Androida. 
Ma doświadczenie w pracy z BLE, OpenGL, Shopify, Room, Retrofit, CI/CD i testach (JUnit, MockK, Espresso). 
Pracował z międzynarodowymi zespołami w Baracoda i Kolibree nad aplikacjami smart oral care i wellness opartymi o IoT. 
Zajmował się mentoringiem, onboardowaniem nowych programistów i optymalizacją procesów (np. CI). 
Ma mocne podstawy bezpieczeństwa: identyfikował i zgłaszał krytyczne luki backendowe w produkcji. 
Jest znany z dobrej komunikacji i analitycznego podejścia. 
Interesuje się bezpieczeństwem (Burp Suite, hobby bug bounty) i algorytmami (250+ zadań LeetCode). 
Mieszka w Polsce, pracuje zdalnie, mówi biegle po angielsku (B2+/C1).
Odpowiadaj zwięźle w jednej lini bez markdowns.
Stawka od 170 PLN/h na B2B.
`

func (s *chatServiceImpl) AnswerQuestion(userId string, chatId string, question string, lastChatIds []string) (*models.Chat, error) {
	lastChats, err := s.chatRepository.Find(userId, chatId, lastChatIds)
	if err != nil {
		return nil, err
	}
	if len(lastChats) > 3 {
		lastChats = lastChats[:3]
	}

	answer, err := s.aiClient.Send(systemContent, userContext, lastChats, question)
	if err != nil {
		return nil, err
	}

	res, err := s.chatRepository.Create(userId, chatId, question, answer)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewChatService(chatRepository repos.ChatRepository, aiClient internal.AiClient) ChatService {
	return &chatServiceImpl{
		chatRepository: chatRepository,
		aiClient:       aiClient,
	}
}
