package middleware

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(c fiber.Ctx) error {
	auth := c.Get("Authorization")
	if len(auth) < 7 || auth[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).SendString("Токен отсутствует или неверный")
	}
	tokenStr := auth[7:]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Недопустимый метод подписи")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Недействительный токен")
	}
	return c.Next()
}

func GenerateToken() {

}
