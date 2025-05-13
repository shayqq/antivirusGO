package controller

import (
	"antivirus/database"
	"antivirus/model"
	"antivirus/request"
	"antivirus/service"
	"antivirus/service/impl"
	"encoding/json"
	"net/mail"

	"github.com/gofiber/fiber/v3"
)

func RegisterRegistrationRoute(app *fiber.App) {
	app.Post("/registration", registration)
}

func registration(c fiber.Ctx) error {
	var registrationRequest request.RegistrationRequest
	body := c.Body()
	err := json.Unmarshal(body, &registrationRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Невалидный формат запроса")
	}
	if registrationRequest.Username == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Введите логин!")
	}
	if registrationRequest.Email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Введите email!")
	}
	if registrationRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Введите пароль!")
	}
	_, err = mail.ParseAddress(registrationRequest.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Неверный формат email!")
	}
	var userService service.UserService = &impl.UserServiceImpl{}
	hashedPassword, err := userService.HashPassword(registrationRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при хэшировании пароля")
	}
	applicationUser := model.ApplicationUser{
		Username: registrationRequest.Username,
		Email:    registrationRequest.Email,
		Password: string(hashedPassword),
		Role:     0,
	}
	result := database.Save(applicationUser)
	if result != "" {
		if result != "Ошибка сервера" {
			return c.Status(fiber.StatusConflict).SendString(result)
		}
		return c.Status(fiber.StatusInternalServerError).SendString(result)
	}
	return c.Status(fiber.StatusOK).SendString("Пользователь успешно создан!")
}
