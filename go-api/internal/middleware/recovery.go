// Package middleware provee middlewares reutilizables para Fiber.
package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse es la estructura estándar para respuestas de error.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Recovery captura cualquier panic durante el manejo de una request y
// responde con HTTP 500 sin interrumpir el servidor.
func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[recovery] panic caught: %v", r)
				_ = c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
					Error:   "internal_server_error",
					Message: "An unexpected error occurred",
				})
			}
		}()
		return c.Next()
	}
}

// ErrorHandler es el manejador centralizado de errores de Fiber.
// Se registra en fiber.Config.ErrorHandler.
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("[error] %d - %s", code, message)

	return c.Status(code).JSON(ErrorResponse{
		Error:   httpStatusText(code),
		Message: message,
	})
}

func httpStatusText(code int) string {
	switch code {
	case 400:
		return "bad_request"
	case 404:
		return "not_found"
	case 422:
		return "unprocessable_entity"
	case 500:
		return "internal_server_error"
	default:
		return "error"
	}
}
