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

	// db connection
	database.ConnectDB()
	database.ConnectMongo()

	// auth
	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)

	// admin
	adminRepo := repository.NewAdminUserRepository(database.DB)
	adminUserService := service.NewAdminUserService(adminRepo)

	// student dan lecturer
	studentRepo := repository.NewStudentRepository(database.DB)
	lecturerRepo := repository.NewLecturerRepository(database.DB)

	// achievements
	achievementPGRepo := repository.NewAchievementRepository(database.DB)
	achievementMongoRepo := repository.NewMongoAchievementRepository(database.Mongo)

	studentAch := service.NewStudentAchievementService(
		achievementPGRepo,
		studentRepo,
		achievementMongoRepo,
	)

	lecturerAch := service.NewLecturerAchievementService(
		achievementPGRepo,
		studentRepo,
		lecturerRepo,
		achievementMongoRepo,
	)

	jwt := middleware.NewJWTMiddleware(userRepo)
	rbac := middleware.NewRBACMiddleware(adminRepo)

	app := fiber.New()
	route.RegisterRoutes(app, authService, adminUserService, jwt, rbac, studentAch, lecturerAch)

	app.Static("/uploads", "./uploads")

	log.Println("Running on: http://localhost:3000")
	app.Listen("0.0.0.0:3000")
}
