package routes

import (
	"github.com/drg809/apiBase/controllers"
	"github.com/drg809/apiBase/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func PresaleRoutes(router fiber.Router) {
	// /api/v1/presale
	presaleRouter := router.Group("/presale")

	// protected by jwt
	presaleRouter.Use(jwtware.New(jwtware.Config{SigningKey: []byte(utils.GetEnvVariable("JWT_USER_KEY"))}))

	// /api/v1/presale/oracle | GET
	presaleRouter.Get("/oracle", controllers.CallBSC)

	// /api/v1/presale/:presale | GET
	presaleRouter.Get("/:id", controllers.GetPresaleByUserID)

	// /api/v1/presale | INSERT
	presaleRouter.Post("/", controllers.InsertPresale)

	// /api/v1/presale/claim | UPDATE
	presaleRouter.Patch("/claim", controllers.SetUserClaim)

	// /api/v1/presale/claim | UPDATE
	presaleRouter.Patch("/vesting", controllers.SetUserVesting)

}
