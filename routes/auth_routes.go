package routes

import (
	"fibergo/responses"
	"fibergo/security"
	"fibergo/utils"
	"github.com/gofiber/fiber/v2"
	jwtaware "github.com/gofiber/jwt/v2"
)

func Authenticate(c *fiber.Ctx) error {
	return jwtaware.New(jwtaware.Config{
		SigningKey:    security.SecretKey,
		SigningMethod: security.SigningMethod,
		TokenLookup:   "header:x-access-token",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(responses.UserResponse{
					Status:  fiber.StatusUnauthorized,
					Message: "error",
					Data:    utils.NewError(err)})
		},
	})(c)
}
