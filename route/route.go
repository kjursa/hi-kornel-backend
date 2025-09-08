package route

import (
	"my-go-backend/deps"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app fiber.Router, deps *deps.Dependencies) {
	auth := app.Group("/auth")
	auth.Post("/register", deps.AuthHandler.Register)
	auth.Post("/refresh", deps.AuthHandler.RefreshToken)
	auth.Post("/ping", deps.AuthHandler.Ping)
}

func SyncRoute(app fiber.Router, deps *deps.Dependencies) {
	sync := app.Group("/sync", deps.SessionMiddleware.HandleToken)
	sync.Get("/", deps.SyncHandler.Sync)

	admin := sync.Group("/admin", deps.SessionMiddleware.HandleAdminToken)
	admin.Put("/profile", deps.SyncHandler.UpdateProfile)
	admin.Put("/experience", deps.SyncHandler.UpdateExperience)
	admin.Put("/projects", deps.SyncHandler.UpdateProjects)
}

func ContactRoute(app fiber.Router, deps *deps.Dependencies) {
	contact := app.Group("/contact", deps.SessionMiddleware.HandleToken)
	contact.Post("/message", deps.ContactHandler.SendMessage)
	contact.Get("/message", deps.ContactHandler.GetMessages)
}

func ChatRoute(app fiber.Router, deps *deps.Dependencies) {
	contact := app.Group("/chat", deps.SessionMiddleware.HandleToken)
	contact.Post("/:chatId/question", deps.ChatHandler.AskQuestion)
}
