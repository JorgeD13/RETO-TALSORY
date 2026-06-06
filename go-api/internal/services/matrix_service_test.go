// Tests de la capa de servicio: valida reglas de negocio, factorización QR
// y comportamiento ante fallos del cliente externo (Node API).
package services

import (
	"context"
	"errors"
	"math"
	"testing"

	"go-api/internal/models"
)

// ── Mock del NodeClient ───────────────────────────────────────────────────────

// mockNodeClient implementa clients.NodeClient sin realizar llamadas HTTP reales.
// Permite controlar la respuesta o el error en cada caso de prueba.
type mockNodeClient struct {
	response *models.NodeResponse
	err      error
}

func (m *mockNodeClient) SendQR(_ context.Context, _ models.NodePayload) (*models.NodeResponse, error) {
	return m.response, m.err
}

// ── Helpers matemáticos ───────────────────────────────────────────────────────

func approxEqual(a, b, eps float64) bool { return math.Abs(a-b) < eps }

// matrixApproxEqual verifica que dos matrices tengan los mismos valores
// dentro de una tolerancia eps (necesaria por errores de punto flotante).
func matrixApproxEqual(a, b [][]float64, eps float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if !approxEqual(a[i][j], b[i][j], eps) {
				return false
			}
		}
	}
	return true
}

// multiplyMatrices calcula C = A·B para verificar que Q·R ≈ A original.
func multiplyMatrices(a, b [][]float64) [][]float64 {
	rows, cols, inner := len(a), len(b[0]), len(b)
	c := make([][]float64, rows)
	for i := range c {
		c[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			for k := 0; k < inner; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

// ── validateMatrix ────────────────────────────────────────────────────────────

// TestValidateMatrix cubre todas las restricciones de entrada que el negocio
// define: la matriz debe existir, tener al menos 1 fila y 1 columna, y ser
// rectangular (todas las filas con igual longitud).
func TestValidateMatrix(t *testing.T) {
	t.Run("nil — campo 'matrix' ausente en el JSON", func(t *testing.T) {
		// Caso: el cliente envía {} sin el campo matrix.
		// Debe retornar error para que el handler responda 422.
		if err := validateMatrix(nil); err == nil {
			t.Fatal("expected error: nil matrix should be rejected")
		}
	})

	t.Run("slice vacío — matrix=[]", func(t *testing.T) {
		// Caso: el cliente envía "matrix": [] sin ninguna fila.
		// Una matriz sin filas no puede factorizarse.
		if err := validateMatrix([][]float64{}); err == nil {
			t.Fatal("expected error: empty matrix should be rejected")
		}
	})

	t.Run("fila vacía — matrix=[[]]", func(t *testing.T) {
		// Caso: el cliente envía una fila sin columnas.
		// Una matriz sin columnas no tiene rango, no es factorizable.
		if err := validateMatrix([][]float64{{}}); err == nil {
			t.Fatal("expected error: row with no columns should be rejected")
		}
	})

	t.Run("matriz jagged — filas con longitudes distintas", func(t *testing.T) {
		// Caso: [[1,2,3],[4,5]] — Gonum requiere matrices rectangulares.
		// La factorización QR está indefinida para matrices no rectangulares.
		jagged := [][]float64{{1, 2, 3}, {4, 5}}
		if err := validateMatrix(jagged); err == nil {
			t.Fatal("expected error: jagged matrix should be rejected")
		}
	})

	t.Run("ancha (3 filas × 4 cols) — Gonum no soporta m < n", func(t *testing.T) {
		// Gonum requiere m >= n para la factorización QR.
		// Sin esta validación, qr.RTo() produce un panic en runtime.
		wide := [][]float64{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
		if err := validateMatrix(wide); err == nil {
			t.Fatal("expected error for wide matrix (cols > rows)")
		}
	})

	t.Run("1×1 — mínima matriz válida", func(t *testing.T) {
		// Caso límite: una sola celda es una matriz válida y factorizable.
		if err := validateMatrix([][]float64{{42}}); err != nil {
			t.Fatalf("expected no error for 1×1 matrix, got: %v", err)
		}
	})

	t.Run("cuadrada n×n (3×3) — válida explícitamente", func(t *testing.T) {
		// Caso estándar: matriz cuadrada, m == n, cumple m >= n.
		m := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
		if err := validateMatrix(m); err != nil {
			t.Fatalf("expected no error for square 3×3, got: %v", err)
		}
	})

	t.Run("columna única (4×1) — válida, m >> n", func(t *testing.T) {
		// Caso: vector columna. m=4, n=1. Cumple m >= n ampliamente.
		col := [][]float64{{1}, {2}, {3}, {4}}
		if err := validateMatrix(col); err != nil {
			t.Fatalf("expected no error for 4×1 column, got: %v", err)
		}
	})

	t.Run("rectangular alta (4×3) — caso feliz", func(t *testing.T) {
		// Caso: matriz 4×3 bien formada (más filas que columnas).
		// Gonum QR solo acepta m >= n, así que debe pasar.
		m := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}}
		if err := validateMatrix(m); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})
}

// ── ComputeQR — propiedad matemática ─────────────────────────────────────────

// TestComputeQR_MathProperty verifica que la factorización sea correcta
// comprobando la propiedad fundamental: A = Q·R (dentro de tolerancia float64).
func TestComputeQR_MathProperty(t *testing.T) {
	cases := []struct {
		name   string
		matrix [][]float64
	}{
		{
			// Caso estándar del enunciado del reto: 3×3 con valores mixtos.
			name: "3×3 — ejemplo del reto",
			matrix: [][]float64{
				{12, -51, 4},
				{6, 167, -68},
				{-4, 24, -41},
			},
		},
		{
			// Caso límite: la identidad debe factorizarse en Q=I, R=I.
			// Q·R = I·I = I, por lo que A ≈ Q·R debe cumplirse exactamente.
			name:   "identidad 2×2 — Q·R debe ser I",
			matrix: [][]float64{{1, 0}, {0, 1}},
		},
		{
			// Caso con valores negativos: verifica que el signo se preserva.
			name:   "valores negativos",
			matrix: [][]float64{{-3, 1}, {4, -2}},
		},
		{
			// Caso con decimal/flotante: verifica que la precisión se mantiene.
			name:   "valores decimales",
			matrix: [][]float64{{1.5, 2.5}, {3.5, 4.5}},
		},
		{
			// Caso 1×1: la factorización trivial. Q = [[±1]], R = [[valor]].
			// Q·R debe reproducir exactamente el único elemento.
			name:   "1×1 — caso trivial",
			matrix: [][]float64{{7}},
		},
		{
			// Caso rectangular alta 4×3: Q será 4×4, R será 4×3.
			// multiplyMatrices(Q, R) debe devolver la matriz original 4×3.
			name: "4×3 rectangular alta — full QR",
			matrix: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
				{10, 11, 12},
			},
		},
		{
			// Caso singular: rango 1, la segunda fila es múltiplo de la primera.
			// QR siempre existe (aunque R tenga cero en la diagonal).
			// La propiedad A = Q·R se cumple incluso para matrices singulares.
			name:   "singular (rango 1) — [[1,2],[2,4]]",
			matrix: [][]float64{{1, 2}, {2, 4}},
		},
		{
			// Caso todo ceros: la factorización degenera en Q=I, R=0.
			// Q·R = I·0 = 0 = A. La propiedad debe cumplirse.
			name:   "todo ceros 2×2",
			matrix: [][]float64{{0, 0}, {0, 0}},
		},
	}

	nodeResp := &models.NodeResponse{Max: 1, Min: -1}
	svc := NewMatrixService(&mockNodeClient{response: nodeResp})
	const eps = 1e-9

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.ComputeQR(context.Background(), models.MatrixRequest{Matrix: tc.matrix})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			reconstructed := multiplyMatrices(result.Q, result.R)
			if !matrixApproxEqual(tc.matrix, reconstructed, eps) {
				t.Errorf("A ≠ Q·R\noriginal:      %v\nreconstructed: %v", tc.matrix, reconstructed)
			}
		})
	}
}

// ── ComputeQR — dimensiones de salida ────────────────────────────────────────

// TestComputeQR_OutputDimensions verifica que Q y R tengan las dimensiones
// correctas según la definición de la factorización QR completa:
//   - Q siempre es cuadrada: m×m
//   - R tiene las mismas dimensiones que A: m×n
func TestComputeQR_OutputDimensions(t *testing.T) {
	cases := []struct {
		name     string
		matrix   [][]float64
		wantQRow int
		wantQCol int
		wantRRow int
		wantRCol int
	}{
		{
			// Matriz cuadrada 3×3: Q→3×3, R→3×3.
			name:     "3×3 cuadrada",
			matrix:   [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			wantQRow: 3, wantQCol: 3,
			wantRRow: 3, wantRCol: 3,
		},
		{
			// Matriz rectangular tall 3×2: Q→3×3, R→3×2.
			// Gonum expande Q a la forma completa (full QR).
			name:     "3×2 rectangular",
			matrix:   [][]float64{{1, 2}, {3, 4}, {5, 6}},
			wantQRow: 3, wantQCol: 3,
			wantRRow: 3, wantRCol: 2,
		},
		{
			// Matriz 1×1: ambas matrices son escalares, Q→1×1, R→1×1.
			name:     "1×1 trivial",
			matrix:   [][]float64{{5}},
			wantQRow: 1, wantQCol: 1,
			wantRRow: 1, wantRCol: 1,
		},
		{
			// Matriz rectangular alta 4×3: Q→4×4 (cuadrada), R→4×3.
			// Las dos últimas filas de R son cero (full QR).
			name: "4×3 tall",
			matrix: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
				{10, 11, 12},
			},
			wantQRow: 4, wantQCol: 4,
			wantRRow: 4, wantRCol: 3,
		},
		{
			// Vector columna 4×1: Q→4×4, R→4×1.
			// Caso extremo donde n=1.
			name:     "4×1 vector columna",
			matrix:   [][]float64{{1}, {2}, {3}, {4}},
			wantQRow: 4, wantQCol: 4,
			wantRRow: 4, wantRCol: 1,
		},
	}

	nodeResp := &models.NodeResponse{}
	svc := NewMatrixService(&mockNodeClient{response: nodeResp})

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.ComputeQR(context.Background(), models.MatrixRequest{Matrix: tc.matrix})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result.Q) != tc.wantQRow || len(result.Q[0]) != tc.wantQCol {
				t.Errorf("Q: want %dx%d, got %dx%d", tc.wantQRow, tc.wantQCol, len(result.Q), len(result.Q[0]))
			}
			if len(result.R) != tc.wantRRow || len(result.R[0]) != tc.wantRCol {
				t.Errorf("R: want %dx%d, got %dx%d", tc.wantRRow, tc.wantRCol, len(result.R), len(result.R[0]))
			}
		})
	}
}

// ── ComputeQR — resiliencia ante fallo de Node ───────────────────────────────

// TestComputeQR_NodeAPIFailure verifica que el servicio responde con Q y R
// válidas aunque la API Node.js no esté disponible, en lugar de propagar
// el error al cliente. Las estadísticas quedan como zero values.
func TestComputeQR_NodeAPIFailure(t *testing.T) {
	// Simula que Node está caído (timeout, connection refused, etc.).
	svc := NewMatrixService(&mockNodeClient{err: errors.New("connection refused")})

	result, err := svc.ComputeQR(context.Background(), models.MatrixRequest{
		Matrix: [][]float64{{1, 2}, {3, 4}},
	})

	// El servicio NO debe propagar el error: la factorización es útil
	// aunque no tengamos estadísticas.
	if err != nil {
		t.Fatalf("service must not propagate node client error, got: %v", err)
	}

	// Q y R deben estar presentes con los datos de la factorización.
	if len(result.Q) == 0 || len(result.R) == 0 {
		t.Error("Q and R must be populated even when node API is unavailable")
	}

	// Las estadísticas deben ser zero values, no datos corruptos.
	if result.Statistics != (models.Statistics{}) {
		t.Errorf("statistics must be zero values when node API fails, got: %+v", result.Statistics)
	}
}

// ── ComputeQR — propagación de estadísticas ───────────────────────────────────

// TestComputeQR_StatisticsPropagation verifica que los valores devueltos por
// el NodeClient se mapeen correctamente a la respuesta final sin pérdida de datos.
func TestComputeQR_StatisticsPropagation(t *testing.T) {
	want := &models.NodeResponse{
		Max: 175, Min: -68, Average: 10.5, Sum: 100,
		IsDiagonalQ: false, IsDiagonalR: true,
	}
	svc := NewMatrixService(&mockNodeClient{response: want})

	result, err := svc.ComputeQR(context.Background(), models.MatrixRequest{
		Matrix: [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}},
	})
	if err != nil {
		t.Fatal(err)
	}

	got := result.Statistics
	// Verificar cada campo individualmente para mensajes de error claros.
	if got.Max != want.Max {
		t.Errorf("Max: want %.2f, got %.2f", want.Max, got.Max)
	}
	if got.Min != want.Min {
		t.Errorf("Min: want %.2f, got %.2f", want.Min, got.Min)
	}
	if got.Average != want.Average {
		t.Errorf("Average: want %.2f, got %.2f", want.Average, got.Average)
	}
	if got.Sum != want.Sum {
		t.Errorf("Sum: want %.2f, got %.2f", want.Sum, got.Sum)
	}
	if got.IsDiagonalQ != want.IsDiagonalQ {
		t.Errorf("IsDiagonalQ: want %v, got %v", want.IsDiagonalQ, got.IsDiagonalQ)
	}
	if got.IsDiagonalR != want.IsDiagonalR {
		t.Errorf("IsDiagonalR: want %v, got %v", want.IsDiagonalR, got.IsDiagonalR)
	}
}
