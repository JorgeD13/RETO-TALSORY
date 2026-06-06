// Package services contiene la lógica de negocio de la aplicación.
package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go-api/internal/clients"
	"go-api/internal/models"

	"gonum.org/v1/gonum/mat"
)

// MatrixService define el contrato de la capa de servicio para operaciones de matriz.
type MatrixService interface {
	ComputeQR(ctx context.Context, req models.MatrixRequest) (*models.MatrixResponse, error)
}

// matrixService es la implementación concreta de MatrixService.
type matrixService struct {
	nodeClient clients.NodeClient
}

// NewMatrixService construye un matrixService con las dependencias inyectadas.
func NewMatrixService(nodeClient clients.NodeClient) MatrixService {
	return &matrixService{nodeClient: nodeClient}
}

// ComputeQR valida la matriz de entrada, calcula su factorización QR y
// devuelve Q, R y las estadísticas provistas por la API Node.js.
// El contexto se propaga al cliente HTTP para permitir cancelación.
func (s *matrixService) ComputeQR(ctx context.Context, req models.MatrixRequest) (*models.MatrixResponse, error) {
	if err := validateMatrix(req.Matrix); err != nil {
		return nil, err
	}

	rows := len(req.Matrix)
	cols := len(req.Matrix[0])

	// Construir la matriz densa de Gonum en orden row-major.
	data := make([]float64, 0, rows*cols)
	for _, row := range req.Matrix {
		data = append(data, row...)
	}
	A := mat.NewDense(rows, cols, data)

	// Ejecutar la factorización QR.
	var qr mat.QR
	qr.Factorize(A)

	var Q mat.Dense
	var R mat.Dense
	qr.QTo(&Q)
	qr.RTo(&R)

	qSlice := denseToSlice(&Q)
	rSlice := denseToSlice(&R)

	log.Printf("[matrix_service] QR computed: Q%dx%d R%dx%d", rows, rows, rows, cols)

	// Enviar Q y R a la API Node.js propagando el contexto de la request.
	nodeResp, err := s.nodeClient.SendQR(ctx, models.NodePayload{Q: qSlice, R: rSlice})
	if err != nil {
		// Si Node no está disponible devolvemos Q y R con estadísticas vacías
		// en lugar de abortar la respuesta completa.
		log.Printf("[matrix_service] node API unavailable: %v — returning empty statistics", err)
		return &models.MatrixResponse{
			Q:          qSlice,
			R:          rSlice,
			Statistics: models.Statistics{},
		}, nil
	}

	return &models.MatrixResponse{
		Q:          qSlice,
		R:          rSlice,
		Statistics: *nodeResp,
	}, nil
}

// validateMatrix verifica que la matriz cumpla todas las restricciones de negocio.
func validateMatrix(matrix [][]float64) error {
	if matrix == nil {
		return errors.New("matrix is required")
	}
	if len(matrix) == 0 {
		return errors.New("matrix must not be empty")
	}
	cols := len(matrix[0])
	if cols == 0 {
		return errors.New("matrix must have at least one column")
	}
	for i, row := range matrix {
		if len(row) != cols {
			return fmt.Errorf("row %d has %d columns, expected %d", i, len(row), cols)
		}
	}
	// Gonum QR requiere m >= n (filas >= columnas).
	// Una matriz ancha produce panic en qr.RTo().
	if len(matrix) < cols {
		return fmt.Errorf(
			"matrix must have at least as many rows as columns for QR factorization (got %d×%d)",
			len(matrix), cols,
		)
	}
	return nil
}

// denseToSlice convierte una mat.Dense en [][]float64.
func denseToSlice(m *mat.Dense) [][]float64 {
	rows, cols := m.Dims()
	result := make([][]float64, rows)
	for i := range result {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = m.At(i, j)
		}
	}
	return result
}
