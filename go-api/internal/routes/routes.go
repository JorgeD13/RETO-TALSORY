// Package routes registra todas las rutas HTTP de la aplicación.
package routes

import (
	"go-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// Register monta todos los grupos de rutas sobre la instancia Fiber.
func Register(app *fiber.App, matrixHandler *handlers.MatrixHandler) {
	api := app.Group("/api/v1")

	matrix := api.Group("/matrix")
	matrix.Post("/qr", matrixHandler.ComputeQR)
}
