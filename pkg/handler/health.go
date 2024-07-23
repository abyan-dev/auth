package handler

import (
	r "github.com/abyan.dev/auth/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	return r.Ok(c, "Hello, World from an UNPROTECTED route!")
}

func HealthProtected(c *fiber.Ctx) error {
	return r.Ok(c, "Hello, World from a PROTECTED route!")
}
