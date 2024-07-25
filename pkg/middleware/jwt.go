package middleware

import (
	"os"

	"github.com/abyan-dev/auth/pkg/model"
	"github.com/abyan-dev/auth/pkg/response"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RequireAuthenticated() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: jwtErrorHandler,
		// SuccessHandler: checkForRevocation,
		TokenLookup: "cookie:access_token",
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return response.Unauthorized(c, err.Error())
}

func checkForRevocation(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return response.Unauthorized(c, "Missing or malformed token")
	}

	db := c.Locals("db").(*gorm.DB)

	var revokedToken model.RevokedToken
	err := db.Where("token = ?", token).First(&revokedToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Next()
		}

		return response.InternalServerError(c, "Database error")
	}

	return response.Unauthorized(c, "Token is revoked")
}
