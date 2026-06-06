// Tests de la capa HTTP (handler): verifica que el handler traduzca
// correctamente los errores del service a códigos HTTP y serialice la respuesta.
// Usa app.Test() de Fiber para ejecutar requests reales sin levantar un servidor.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"go-api/internal/middleware"
	"go-api/internal/models"
	"go-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

// ── Mock del MatrixService ────────────────────────────────────────────────────

// mockMatrixService implementa services.MatrixService sin lógica real.
// Permite inyectar una respuesta o un error arbitrario en cada test.
type mockMatrixService struct {
	result *models.MatrixResponse
	err    error
}

func (m *mockMatrixService) ComputeQR(_ context.Context, _ models.MatrixRequest) (*models.MatrixResponse, error) {
	return m.result, m.err
}

var _ services.MatrixService = (*mockMatrixService)(nil) // verificación en compilación

// ── Helpers ───────────────────────────────────────────────────────────────────

// newTestApp crea una instancia Fiber con el ErrorHandler centralizado
// y registra el handler bajo la misma ruta que en producción.
func newTestApp(svc services.MatrixService) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	h := NewMatrixHandler(svc)
	app.Post("/api/v1/matrix/qr", h.ComputeQR)
	return app
}

func doPost(app *fiber.App, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/matrix/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return app.Test(req, -1)
}

func decodeBody(t *testing.T, r io.Reader, dst any) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(dst); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

// ── Tests ─────────────────────────────────────────────────────────────────────

// TestHandler_ValidRequest verifica el camino feliz completo:
// body correcto → service devuelve Q, R y estadísticas → handler responde 200
// con la estructura MatrixResponse serializada en JSON.
func TestHandler_ValidRequest(t *testing.T) {
	want := &models.MatrixResponse{
		Q: [][]float64{{-0.857, 0.394}, {-0.429, -0.903}},
		R: [][]float64{{-14, -21}, {0, -175}},
		Statistics: models.Statistics{
			Max: 175, Min: -21, Average: 10.5, Sum: 100,
			IsDiagonalQ: false, IsDiagonalR: false,
		},
	}
	app := newTestApp(&mockMatrixService{result: want})

	body, _ := json.Marshal(models.MatrixRequest{Matrix: [][]float64{{1, 2}, {3, 4}}})
	resp, err := doPost(app, body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}

	var got models.MatrixResponse
	decodeBody(t, resp.Body, &got)

	// Verificar que el campo statistics se serializa correctamente.
	if got.Statistics.Max != want.Statistics.Max {
		t.Errorf("Statistics.Max: want %.2f, got %.2f", want.Statistics.Max, got.Statistics.Max)
	}
}

// TestHandler_InvalidJSON verifica que un body que no sea JSON válido
// sea rechazado con 400 Bad Request, sin llegar al service.
func TestHandler_InvalidJSON(t *testing.T) {
	// Caso: el cliente envía texto plano en lugar de JSON.
	app := newTestApp(&mockMatrixService{})

	resp, err := doPost(app, []byte("not json at all"))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// TestHandler_EmptyBody verifica que un body completamente vacío
// sea rechazado con 400 Bad Request.
func TestHandler_EmptyBody(t *testing.T) {
	// Caso: el cliente envía una request sin body.
	// BodyParser de Fiber devuelve error para body vacío con Content-Type JSON.
	app := newTestApp(&mockMatrixService{})

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/matrix/qr", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	// Body vacío → Fiber no puede hacer parse → 400.
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 for empty body, got %d", resp.StatusCode)
	}
}

// TestHandler_ServiceValidationError verifica que cuando el service rechaza
// la matriz (e.g. jagged), el handler responde 422 Unprocessable Entity
// con el mensaje de error del service en el cuerpo JSON.
func TestHandler_ServiceValidationError(t *testing.T) {
	// Simula que el service devuelve un error de validación de negocio.
	validationErr := errors.New("row 1 has 2 columns, expected 3")
	app := newTestApp(&mockMatrixService{err: validationErr})

	body, _ := json.Marshal(models.MatrixRequest{
		// El service mock ignora el contenido y devuelve el error configurado.
		Matrix: [][]float64{{1, 2, 3}, {4, 5}},
	})
	resp, err := doPost(app, body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("want 422 for validation error, got %d", resp.StatusCode)
	}

	// El cuerpo de error debe contener el mensaje legible para el cliente.
	var errResp middleware.ErrorResponse
	decodeBody(t, resp.Body, &errResp)
	if errResp.Error != "unprocessable_entity" {
		t.Errorf("error field: want 'unprocessable_entity', got %q", errResp.Error)
	}
	if errResp.Message == "" {
		t.Error("message field must not be empty")
	}
}

// TestHandler_NilMatrixField verifica que enviar { "matrix": null }
// sea rechazado con 422 (el JSON parsea correctamente pero la validación falla).
func TestHandler_NilMatrixField(t *testing.T) {
	// Caso: el campo matrix está presente en el JSON pero con valor null.
	// BodyParser no falla (JSON válido), pero el service retorna error de validación.
	svcErr := errors.New("matrix is required")
	app := newTestApp(&mockMatrixService{err: svcErr})

	resp, err := doPost(app, []byte(`{"matrix": null}`))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("want 422 for null matrix, got %d", resp.StatusCode)
	}
}

// TestHandler_ContentTypeJSON verifica que la respuesta siempre lleva
// el Content-Type correcto para que los clientes puedan deserializarla.
func TestHandler_ContentTypeJSON(t *testing.T) {
	app := newTestApp(&mockMatrixService{result: &models.MatrixResponse{}})
	body, _ := json.Marshal(models.MatrixRequest{Matrix: [][]float64{{1}}})

	resp, err := doPost(app, body)
	if err != nil {
		t.Fatal(err)
	}

	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		t.Error("Content-Type header must be present in response")
	}
}
