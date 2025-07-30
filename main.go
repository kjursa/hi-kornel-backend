package main

import (
	"fmt"
	"my-go-backend/deps"
	"my-go-backend/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	deps := deps.NewDependencies()
	deps.Init()
	defer deps.Close()

	api := app.Group("/api")
	route.AuthRoute(api, deps)
	route.SyncRoute(api, deps)
	route.ChatRoute(api, deps)
	route.ContactRoute(api, deps)

	fmt.Print("elapsed:" + deps.Port)

	app.Listen(":" + deps.Port)
}
