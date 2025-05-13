package controller

import (
	"antivirus/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterDemoRoutes(app *fiber.App) {
	app.Get("/hello", hello, middleware.Protected)
}

func hello(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
