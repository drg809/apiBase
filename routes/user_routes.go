package routes

import (
	"github.com/drg809/apiBase/controllers"
	"github.com/drg809/apiBase/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func UserRoutes(router fiber.Router) {
	// /api/v1/user
	userRouter := router.Group("/user")

	// /api/v1/user/totalscore
	//userRouter.Get("/:id/totalscore", controllers.GetUserTotalScoreByID)

	// protected by jwt
	userRouter.Use(jwtware.New(jwtware.Config{SigningKey: []byte(utils.GetEnvVariable("JWT_USER_KEY"))}))

	// /api/v1/user/:user | GET
	userRouter.Get("/:id", controllers.GetUserByID)

	// /api/v1/user/:user | UPDATE
	userRouter.Patch("/", controllers.UpdateUserByID)

	// /api/v1/user/:user | DELETE
	userRouter.Delete("/:id", controllers.DeleteUserByID)

	// /api/v1/user/refresh
	userRouter.Get("/:id/refresh", controllers.RefreshUser)

}
