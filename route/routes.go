package route

import (
	"uas/app/service"
	"uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, auth service.AuthService, jwt *middleware.JWTMiddleware) {

	api := app.Group("/app")

	// auth
	authRoute := api.Group("/auth")
	authRoute.Post("/login", auth.LoginHandler)
	authRoute.Post("/refresh", auth.RefreshHandler)
	authRoute.Get("/profile", jwt.RequireAuth, auth.ProfileHandler)

}
