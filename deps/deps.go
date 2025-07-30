package deps

import (
	"my-go-backend/handlers"
	"my-go-backend/internal"
	"my-go-backend/middleware"
	"my-go-backend/repos"
	"my-go-backend/services"
	"my-go-backend/utils"
	"os"
	"time"
)

type Dependencies struct {
	AuthHandler       handlers.AuthHandler
	SyncHandler       handlers.SyncHandler
	ContactHandler    handlers.ContactHandler
	ChatHandler       handlers.ChatHandler
	SessionMiddleware middleware.SessionMiddleware
	FirestoreClient   *internal.FirestoreClient
	Port              string
}

var secret = getenv("SECRET")
var tokenDuration = 60 * time.Minute
var firestorePath = getenv("FS_PATH")
var firestoreProjectId = getenv("FS_PROJECT_ID")
var firestoreDatabaseId = getenv("FS_DB_ID")
var aiApiKey = getenv("AI_API_KEY")
var aiModel = getenv("AI_MODEL")
var aiUrl = getenv("AI_URL")
var port = getenv("PORT")

func NewDependencies() *Dependencies {
	tokenManager := utils.NewTokenManager(secret, tokenDuration)
	firestoreConfig := internal.NewFirestoreConfig(firestorePath, firestoreProjectId, firestoreDatabaseId)
	firestoreClient := internal.NewFirestoreClient(*firestoreConfig)
	tokenRepo := repos.NewTokenRepository(firestoreClient)
	userRepository := repos.NewUserRepository(firestoreClient)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(userService, tokenRepo, tokenManager)
	authHandler := *handlers.NewAuthHandler(authService, userService, tokenManager)
	profileRepo := repos.NewProfileRepository(firestoreClient)
	experienceRepo := repos.NewExperienceRepository(firestoreClient)
	syncHandler := *handlers.NewSyncHandler(profileRepo, experienceRepo)
	contactRepo := repos.NewContactRepository(firestoreClient)
	contactHandler := *handlers.NewContactHandler(contactRepo)
	chatRepository := repos.NewChatRepository(firestoreClient)
	aiConfig := *internal.NewAiClientConfig(aiApiKey, aiModel, aiUrl)
	aiClient := *internal.NewAiClient(aiConfig)
	chatService := services.NewChatService(chatRepository, aiClient)
	chatHandler := *handlers.NewChatHandler(chatService)
	sessionMiddleware := *middleware.NewSessionMiddleware(tokenManager)

	return &Dependencies{
		AuthHandler:       authHandler,
		SyncHandler:       syncHandler,
		ContactHandler:    contactHandler,
		ChatHandler:       chatHandler,
		SessionMiddleware: sessionMiddleware,
		FirestoreClient:   firestoreClient,
		Port:              port,
	}
}

func (d *Dependencies) Init() {
	d.FirestoreClient.Init()
}

func (d *Dependencies) Close() {
	d.FirestoreClient.Close()
}

func getenv(name string) string {
	key := os.Getenv(name)
	if key == "" {
		panic("missing env key: " + name)
	}
	return key
}
