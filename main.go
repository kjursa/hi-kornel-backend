package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// deps := deps.NewDependencies()
	// deps.Init()
	// defer deps.Close()

	// api := app.Group("/api")
	// route.AuthRoute(api, deps)
	// route.SyncRoute(api, deps)
	// route.ChatRoute(api, deps)
	// route.ContactRoute(api, deps)

	// fmt.Print("elapsed:" + deps.Port)

	port := os.Getenv("PORT")

	app.Get("/env", func(c *fiber.Ctx) error {
		key := c.Query("key")
		if key == "" {
			return c.SendString("No KEY!")
		}

		value := os.Getenv(key)

		return c.SendString(key + " = " + value)
	})

	app.Listen(":" + port)
}
