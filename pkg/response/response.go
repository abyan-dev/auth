package response

import (
	"github.com/gofiber/fiber/v2"
)

type Payload struct {
	Status  uint        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusOK, message, data...)
}

func Created(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusCreated, message, data...)
}

func Accepted(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusAccepted, message, data...)
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func BadRequest(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusBadRequest, message, data...)
}

func Unauthorized(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusUnauthorized, message, data...)
}

func Forbidden(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusForbidden, message, data...)
}

func NotFound(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusNotFound, message, data...)
}

func InternalServerError(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusInternalServerError, message, data...)
}

func ServiceUnavailable(c *fiber.Ctx, message string, data ...interface{}) error {
	return jsonResponse(c, fiber.StatusServiceUnavailable, message, data...)
}

func jsonResponse(c *fiber.Ctx, status uint, message string, data ...interface{}) error {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = nil
	}

	payload := Payload{
		Status:  status,
		Message: message,
		Data:    responseData,
	}

	return c.Status(int(status)).JSON(payload)
}
