package main

import (
	"antivirus/controller"
	"antivirus/database"
	"antivirus/textcolor"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(textcolor.RedErrorText("Не удалось найти env файл "), err)
	}
	database.Connect()
	defer database.DB.Close()
	app := fiber.New()
	controller.RegisterRegistrationRoute(app)
	controller.RegisterAuthenticationRoute(app)
	controller.RegisterUserRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
