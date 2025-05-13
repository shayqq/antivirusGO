package controller

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterDemoRoutes(app *fiber.App) {
	app.Get("/hello", hello)
}

func hello(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
