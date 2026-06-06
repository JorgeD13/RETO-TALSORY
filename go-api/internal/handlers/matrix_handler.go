// Package handlers contiene los controladores HTTP de la API.
package handlers

import (
	"log"

	"go-api/internal/models"
	"go-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

// MatrixHandler agrupa los endpoints relacionados con operaciones de matriz.
type MatrixHandler struct {
	service services.MatrixService
}

// NewMatrixHandler construye un MatrixHandler con la dependencia inyectada.
func NewMatrixHandler(service services.MatrixService) *MatrixHandler {
	return &MatrixHandler{service: service}
}

// ComputeQR godoc
// @Summary      Calcula la factorización QR de una matriz
// @Description  Recibe una matriz numérica, calcula Q y R con Gonum y obtiene estadísticas de la API Node.js
// @Tags         matrix
// @Accept       json
// @Produce      json
// @Param        body  body      models.MatrixRequest   true  "Matriz de entrada"
// @Success      200   {object}  models.MatrixResponse
// @Failure      400   {object}  middleware.ErrorResponse
// @Failure      422   {object}  middleware.ErrorResponse
// @Failure      500   {object}  middleware.ErrorResponse
// @Router       /api/v1/matrix/qr [post]
func (h *MatrixHandler) ComputeQR(c *fiber.Ctx) error {
	var req models.MatrixRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[handler] invalid body: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body: "+err.Error())
	}

	// Propagar el contexto de la request HTTP al service y al cliente HTTP.
	result, err := h.service.ComputeQR(c.UserContext(), req)
	if err != nil {
		log.Printf("[handler] service error: %v", err)
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
