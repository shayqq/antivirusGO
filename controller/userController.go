package controller

import (
	"antivirus/database"
	"antivirus/middleware"
	"antivirus/model"
	"antivirus/repository/user"
	"antivirus/request"
	"antivirus/service"
	"antivirus/service/impl"
	"database/sql"
	"encoding/json"
	"net/mail"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(app *fiber.App) {
	app.Get("/user/showAll", showAll, middleware.HasRole(0))
	app.Post("/user/createadm", createadm, middleware.HasRole(1))
	app.Put("/user/updateadm", updateadm, middleware.HasRole(1))
	app.Delete("/user/deleteadm", deleteadm, middleware.HasRole(1))
}

func showAll(c fiber.Ctx) error {
	var userService service.UserService = &impl.UserServiceImpl{}
	applicationUsers, err := userService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка получения данных")
	}
	data := make([]map[string]any, 0, len(applicationUsers))
	for _, value := range applicationUsers {
		data = append(data, map[string]any{
			"email": value.Email,
			"role":  value.Role,
		})
	}
	return c.JSON(data)
}

func createadm(c fiber.Ctx) error {
	var userService service.UserService = &impl.UserServiceImpl{}
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
	hashedPassword, err := userService.HashPassword(registrationRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при хэшировании пароля")
	}
	applicationUser := model.ApplicationUser{
		Username: registrationRequest.Username,
		Email:    registrationRequest.Email,
		Password: hashedPassword,
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

func updateadm(c fiber.Ctx) error {
	var userService service.UserService = &impl.UserServiceImpl{}
	var data map[string]any
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Невалидный JSON")
	}
	id := data["id"]
	if id == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Введите id пользователя")
	}
	var applicationUser model.ApplicationUser
	body := c.Body()
	err = json.Unmarshal(body, &applicationUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Невалидный формат запроса")
	}
	applicationUser1, err := user.FindById(int(id.(float64)))
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка получения данных")
	}
	if applicationUser1 == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Такого пользователя не существует")
	}
	if applicationUser.Username != "" {
		applicationUser1.Username = applicationUser.Username
	}
	if applicationUser.Email != "" {
		applicationUser1.Email = applicationUser.Email
	}
	if applicationUser.Password != "" {
		hashedPassword, err := userService.HashPassword(applicationUser.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при хэшировании пароля")
		}
		applicationUser1.Password = hashedPassword
	}
	if data["role"] != nil {
		applicationUser1.Role = applicationUser.Role
	}
	result := database.Update(*applicationUser1, id)
	if result != "" {
		return c.Status(fiber.StatusInternalServerError).SendString(result)
	}
	return c.Status(fiber.StatusOK).SendString("Пользователь успешно обновлен!")
}

func deleteadm(c fiber.Ctx) error {
	var deleteUserRequest request.DeleteUserRequest
	err := json.Unmarshal(c.Body(), &deleteUserRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Невалидный JSON")
	}
	applicationUser, err := user.FindById(deleteUserRequest.Id)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка получения данных")
	}
	if applicationUser == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Такого пользователя не существует")
	}
	result := database.Delete(*applicationUser, deleteUserRequest.Id)
	if result != "" {
		return c.Status(fiber.StatusInternalServerError).SendString(result)
	}
	return c.Status(fiber.StatusOK).SendString("Пользователь успешно удален!")
}
