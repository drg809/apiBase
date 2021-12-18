package utils

import (
	"errors"
	"fmt"
	"math"
	"time"

	_ "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"github.com/nikola43/fibergormapitemplate/models"
)

func GenerateUserToken(address string, userID uint) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["address"] = address
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	fmt.Println(claims["id"])
	fmt.Println(claims["address"])
	fmt.Println(claims["exp"])

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(GetEnvVariable("JWT_USER_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserTokenClaims(context *fiber.Ctx) (*models.UserTokenClaims, error) {

	if context.Locals("user") == nil {
		return nil, errors.New("invalid claims nil token")
	}

	user := context.Locals("user").(*jwt.Token)
	if claims, ok := user.Claims.(jwt.MapClaims); ok && user.Valid {
		userTokenClaims := new(models.UserTokenClaims)

		if claims["id"] != nil {
			userTokenClaims.ID = uint(math.Round(claims["id"].(float64)))
		}

		if claims["address"] != nil {
			userTokenClaims.Address = claims["address"].(string)
		}

		if claims["exp"] != nil {
			userTokenClaims.Exp = uint(math.Round(claims["exp"].(float64)))
		}

		return userTokenClaims, nil
	} else {
		return nil, errors.New("invalid claims")
	}
}

func GeneratePasswordRecoveryToken(recoveryType string, id uint) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["type"] = recoveryType
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(GetEnvVariable("PASSWORD_RECOVERY_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
