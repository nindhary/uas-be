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

	// AUTH REPO
	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)

	// ADMIN REPO DAN SERVICE
	adminRepo := repository.NewAdminUserRepository(database.DB)
	adminUserService := service.NewAdminUserService(adminRepo)

	// MIDDLEWARE
	jwt := middleware.NewJWTMiddleware(userRepo)
	rbac := middleware.NewRBACMiddleware(adminRepo)

	// FIBER
	app := fiber.New()

	// gunakan route dengan 4 parameter seperti yang diminta compiler
	route.RegisterRoutes(app, authService, adminUserService, jwt, rbac)

	log.Println("Running on: http://localhost:3000")
	app.Listen("0.0.0.0:3000")
}
