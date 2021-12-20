package routes

import (
	"github.com/drg809/apiBase/controllers"
	"github.com/gofiber/fiber/v2"
)

func SignUpRoutes(router fiber.Router) {
	// /api/v1/signup
	authRouter := router.Group("/signup")

	// /api/v1/signup/user
	authRouter.Post("/user", controllers.SignupUser)
}
