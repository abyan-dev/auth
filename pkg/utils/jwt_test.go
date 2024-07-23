package utils

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

func TestCreateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		email             string
		name              string
		role              string
		expirationMinutes int
	}{
		{"test@example.com", "Test User", "user", 60},
		{"admin@example.com", "Admin User", "admin", 120},
	}

	for _, test := range tests {
		token, err := CreateJWT(test.email, test.name, test.role, test.expirationMinutes)
		if err != nil {
			t.Fatalf("Failed to create JWT: %v", err)
		}

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			t.Fatalf("Failed to parse JWT: %v", err)
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			if claims["email"] != test.email {
				t.Errorf("Email claim mismatch: got %v, want %v", claims["email"], test.email)
			}
			if claims["name"] != test.name {
				t.Errorf("Name claim mismatch: got %v, want %v", claims["name"], test.name)
			}
			if claims["role"] != test.role {
				t.Errorf("Role claim mismatch: got %v, want %v", claims["role"], test.role)
			}
			exp := time.Unix(int64(claims["exp"].(float64)), 0)
			expectedExp := time.Now().Add(time.Minute * time.Duration(test.expirationMinutes))
			if exp.Sub(expectedExp).Seconds() > 1 {
				t.Errorf("Expiration claim mismatch: got %v, want %v", exp, expectedExp)
			}
		} else {
			t.Errorf("Invalid JWT claims or token")
		}
	}
}

func TestCreateAuthTokenPair(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	defer os.Unsetenv("JWT_SECRET")

	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	tests := []struct {
		email    string
		name     string
		role     string
		expected AuthTokenPair
	}{
		{"test@example.com", "Test User", "user", AuthTokenPair{}},
		{"admin@example.com", "Admin User", "admin", AuthTokenPair{}},
	}

	for _, test := range tests {
		authTokenPair, err := CreateAuthTokenPair(c, test.email, test.name, test.role)
		if err != nil {
			t.Fatalf("Failed to create auth token pair: %v", err)
		}

		accessToken, err := jwt.Parse(authTokenPair.AccessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			t.Fatalf("Failed to parse access token: %v", err)
		}

		if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
			if claims["email"] != test.email {
				t.Errorf("Email claim mismatch in access token: got %v, want %v", claims["email"], test.email)
			}
			if claims["name"] != test.name {
				t.Errorf("Name claim mismatch in access token: got %v, want %v", claims["name"], test.name)
			}
			if claims["role"] != test.role {
				t.Errorf("Role claim mismatch in access token: got %v, want %v", claims["role"], test.role)
			}
		} else {
			t.Errorf("Invalid access token claims or token")
		}

		refreshToken, err := jwt.Parse(authTokenPair.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			t.Fatalf("Failed to parse refresh token: %v", err)
		}

		if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
			if claims["email"] != test.email {
				t.Errorf("Email claim mismatch in refresh token: got %v, want %v", claims["email"], test.email)
			}
			if claims["name"] != test.name {
				t.Errorf("Name claim mismatch in refresh token: got %v, want %v", claims["name"], test.name)
			}
			if claims["role"] != test.role {
				t.Errorf("Role claim mismatch in refresh token: got %v, want %v", claims["role"], test.role)
			}
		} else {
			t.Errorf("Invalid refresh token claims or token")
		}
	}

	app.ReleaseCtx(c)
}

func TestCreateSecureCookie(t *testing.T) {
	name := "testCookie"
	value := "testValue"
	expiration := time.Hour

	cookie := CreateSecureCookie(name, value, expiration)

	if cookie.Name != name {
		t.Errorf("Name mismatch: got %v, want %v", cookie.Name, name)
	}
	if cookie.Value != value {
		t.Errorf("Value mismatch: got %v, want %v", cookie.Value, value)
	}
	expectedExpiration := time.Now().Add(expiration)
	if !cookie.Expires.After(expectedExpiration.Add(-time.Second)) || !cookie.Expires.Before(expectedExpiration.Add(time.Second)) {
		t.Errorf("Expiration mismatch: got %v, want around %v", cookie.Expires, expectedExpiration)
	}
	if !cookie.HTTPOnly {
		t.Errorf("HTTPOnly mismatch: got %v, want true", cookie.HTTPOnly)
	}
	if cookie.SameSite != "Strict" {
		t.Errorf("SameSite mismatch: got %v, want Strict", cookie.SameSite)
	}
}

func TestInvalidateCookie(t *testing.T) {
	name := "testCookie"

	cookie := InvalidateCookie(name)

	if cookie.Name != name {
		t.Errorf("Name mismatch: got %v, want %v", cookie.Name, name)
	}
	if cookie.Value != "" {
		t.Errorf("Value mismatch: got %v, want empty string", cookie.Value)
	}
	if !cookie.Expires.Before(time.Now()) {
		t.Errorf("Expiration mismatch: got %v, want before %v", cookie.Expires, time.Now())
	}
	if !cookie.HTTPOnly {
		t.Errorf("HTTPOnly mismatch: got %v, want true", cookie.HTTPOnly)
	}
	if cookie.SameSite != "Strict" {
		t.Errorf("SameSite mismatch: got %v, want Strict", cookie.SameSite)
	}
}
