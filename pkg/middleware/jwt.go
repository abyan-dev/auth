package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/abyan-dev/auth/pkg/model"
	"github.com/abyan-dev/auth/pkg/response"
	"github.com/abyan-dev/auth/pkg/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func RequireAuthenticated() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler:   checkForRefresh,
		SuccessHandler: checkForRevocation,
		TokenLookup:    "cookie:access_token",
	})
}

func checkForRefresh(c *fiber.Ctx, err error) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Unauthorized(c, "Missing or malformed token")
	}

	token, parseErr := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if parseErr != nil || !token.Valid {
		return response.Unauthorized(c, "Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return response.Unauthorized(c, "Invalid refresh token claims")
	}

	email, emailOk := claims["email"].(string)
	name, nameOk := claims["name"].(string)
	role, roleOk := claims["role"].(string)

	if !emailOk || !nameOk || !roleOk {
		return response.Unauthorized(c, "Invalid refresh token claims")
	}

	accessToken, createErr := utils.CreateJWT(email, name, role, 5)
	if createErr != nil {
		return response.InternalServerError(c, "Something went wrong")
	}

	accessCookie := utils.CreateSecureCookie("access_token", accessToken, 5*time.Minute)
	c.Cookie(accessCookie)

	c.Locals("user", claims)
	return c.Next()
}

func checkForRevocation(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return response.Unauthorized(c, "Missing or malformed token")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !parsedToken.Valid {
		return checkForRefresh(c, err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return response.Unauthorized(c, "Invalid access token claims")
	}

	db := c.Locals("db").(*gorm.DB)

	var revokedToken model.RevokedToken
	err = db.Where("token = ?", token).First(&revokedToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Locals("user", claims)
			return c.Next()
		}
		return response.InternalServerError(c, "Database error")
	}

	return response.Unauthorized(c, "Token is revoked")
}
