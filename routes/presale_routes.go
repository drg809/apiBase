package routes

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/nikola43/fibergormapitemplate/controllers"
	"github.com/nikola43/fibergormapitemplate/utils"
)

func PresaleRoutes(router fiber.Router) {
	// /api/v1/presale
	presaleRouter := router.Group("/presale")

	// protected by jwt
	presaleRouter.Use(jwtware.New(jwtware.Config{SigningKey: []byte(utils.GetEnvVariable("JWT_USER_KEY"))}))

	// /api/v1/presale/:presale | GET
	presaleRouter.Get("/:id", controllers.GetPresaleByUserID)

	// /api/v1/presale | INSERT
	presaleRouter.Post("/", controllers.InsertPresale)

	// /api/v1/presale/claim | UPDATE
	presaleRouter.Patch("/claim", controllers.SetUserClaim)

}
