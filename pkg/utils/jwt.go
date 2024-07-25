package utils

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenPair struct {
	AccessToken  string
	RefreshToken string
}

func CreateJWT(email string, name string, role string, expirationMinutes int) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"name":  name,
		"role":  role,
		"exp":   time.Now().Add(time.Minute * time.Duration(expirationMinutes)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func CreateAuthTokenPair(c *fiber.Ctx, email string, name string, role string) (AuthTokenPair, error) {
	accessToken, err := CreateJWT(email, name, role, 5) // 5 minutes
	if err != nil {
		return AuthTokenPair{}, errors.New("failed to create access token")
	}

	refreshToken, err := CreateJWT(email, name, role, 7*24*60) // 7 days
	if err != nil {
		return AuthTokenPair{}, errors.New("failed to create refresh token")
	}

	authTokenPair := AuthTokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authTokenPair, nil
}

func CreateSecureCookie(name string, value string, expiration time.Duration) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(expiration),
		HTTPOnly: true,
		SameSite: "Strict",
		Secure:   false,
	}
}

func InvalidateCookie(name string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Strict",
		Secure:   false,
	}
}
