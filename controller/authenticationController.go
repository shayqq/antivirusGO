package controller

import (
	"antivirus/middleware"
	"antivirus/repository/user"
	"antivirus/request"
	"antivirus/service"
	"antivirus/service/impl"
	"database/sql"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthenticationRoute(app *fiber.App) {
	app.Post("/auth/login", login)
}

func login(c fiber.Ctx) error {
	var userService service.UserService = &impl.UserServiceImpl{}
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
	applicationUser, err := user.FindByEmail(authenticationRequest.Email)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка получения данных")
	}
	if applicationUser == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Такого пользоателя не существует")
	}
	if !userService.Authenticate(applicationUser.Password, authenticationRequest.Password) {
		return c.Status(fiber.StatusUnauthorized).SendString("Неверный пароль")
	}
	tokenStr, err := middleware.GenerateToken(applicationUser.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка генерации токена")
	}
	return c.JSON(fiber.Map{
		"email": authenticationRequest.Email,
		"token": tokenStr,
	})
}
