package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func HasRole(requiredRole int8) fiber.Handler {
	return func(c fiber.Ctx) error {
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
		if requiredRole == 1 {
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).SendString("Ошибка чтения токена")
			}
			roleFloat, ok := claims["role"].(float64)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).SendString("Роль отсутствует или неверного формата")
			}
			userRole := int8(roleFloat)
			if userRole != requiredRole {
				return c.Status(fiber.StatusUnauthorized).SendString("Недостаточно прав")
			}
		}
		return c.Next()
	}
}

func GenerateToken(role int8) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenStr, err
}
