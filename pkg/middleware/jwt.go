package middleware

import (
	"os"

	"github.com/abyan.dev/auth/pkg/response"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func RequireAuthenticated() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: jwtErrorHandler,
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return response.Unauthorized(c, err.Error())
}
