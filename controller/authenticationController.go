package controller

import (
	"antivirus/database"
	"antivirus/request"
	"antivirus/service"
	"antivirus/service/impl"
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
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
	flag, password := database.FindBy("email", authenticationRequest.Email)
	if !flag {
		return c.Status(fiber.StatusUnauthorized).SendString("Такого пользоателя не существует")
	}
	var userService service.UserService = &impl.UserServiceImpl{}
	if !userService.Authenticate(password, authenticationRequest.Password) {
		return c.Status(fiber.StatusUnauthorized).SendString("Неверный пароль")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": authenticationRequest.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Ошибка генерации токена")
	}
	return c.JSON(fiber.Map{
		"email": authenticationRequest.Email,
		"token": tokenStr,
	})
}
