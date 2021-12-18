package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nikola43/fibergormapitemplate/controllers"
)

func SignUpRoutes(router fiber.Router) {
	// /api/v1/signup
	authRouter := router.Group("/signup")

	// /api/v1/signup/user
	authRouter.Post("/user", controllers.SignupUser)
}
