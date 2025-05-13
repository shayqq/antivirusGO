package controller

import (
	"antivirus/request"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthenticationRoute(app *fiber.App) {
	app.Post("/auth/login", login)
}

func login(c fiber.Ctx) error {
	var authenticationRequest request.AuthenticationRequest
	body := c.Body()
	err := json.Unmarshal(body, &authenticationRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Невалидный формат запроса")
	}
	if authenticationRequest.Email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Введите email!")
	}
	if authenticationRequest.Email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Введите пароль!")
	}

	return nil
}
