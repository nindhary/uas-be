package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"uas/app/repository"
	"uas/app/service"
	"uas/database"
	"uas/middleware"
	"uas/route"
)

func main() {
	godotenv.Load()

	database.ConnectDB()

	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)
	jwt := middleware.NewJWTMiddleware(userRepo)

	app := fiber.New()
	route.RegisterRoutes(app, authService, jwt)

	log.Println("Running on: http://localhost:3000")
	app.Listen("0.0.0.0:3000")
}
