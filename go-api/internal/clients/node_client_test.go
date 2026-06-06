// Tests del cliente HTTP que se comunica con la API Node.js.
// Usa httptest.NewServer para levantar un servidor HTTP falso y controlar
// exactamente qué responde Node en cada escenario, sin dependencias de red.
package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-api/internal/models"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

// samplePayload es el payload mínimo válido para todos los casos.
var samplePayload = models.NodePayload{
	Q: [][]float64{{-0.857, 0.394}, {-0.429, -0.903}},
	R: [][]float64{{-14, -21}, {0, -175}},
}

// newFakeServer crea un servidor HTTP de test que responde con el statusCode
// y el body dados. El caller debe invocar ts.Close() al terminar.
func newFakeServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))
}

// ── Tests ─────────────────────────────────────────────────────────────────────

// TestSendQR_Success verifica el camino feliz:
// Node responde 200 con un JSON de estadísticas válido → el cliente lo parsea
// correctamente y devuelve un *NodeResponse sin error.
func TestSendQR_Success(t *testing.T) {
	want := models.NodeResponse{
		Max: 70, Min: -175, Average: -12.3, Sum: -110,
		IsDiagonalQ: false, IsDiagonalR: false,
	}
	body, _ := json.Marshal(want)

	ts := newFakeServer(http.StatusOK, string(body))
	defer ts.Close()

	client := NewNodeClient(ts.URL)
	got, err := client.SendQR(context.Background(), samplePayload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Max != want.Max {
		t.Errorf("Max: want %.4f, got %.4f", want.Max, got.Max)
	}
	if got.Min != want.Min {
		t.Errorf("Min: want %.4f, got %.4f", want.Min, got.Min)
	}
	if got.IsDiagonalR != want.IsDiagonalR {
		t.Errorf("IsDiagonalR: want %v, got %v", want.IsDiagonalR, got.IsDiagonalR)
	}
}

// TestSendQR_NodeReturns500 verifica que un status 5xx de Node se traduce
// en un error descriptivo. El cliente nunca debe silenciar errores HTTP.
func TestSendQR_NodeReturns500(t *testing.T) {
	// Caso: Node está degradado y responde Internal Server Error.
	ts := newFakeServer(http.StatusInternalServerError, `{"error":"internal_server_error"}`)
	defer ts.Close()

	client := NewNodeClient(ts.URL)
	_, err := client.SendQR(context.Background(), samplePayload)

	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
	// El mensaje de error debe incluir el status code para facilitar el debug.
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error message should mention status code 500, got: %q", err.Error())
	}
}

// TestSendQR_NodeReturns4xx verifica que respuestas 4xx (e.g. 400 Bad Request)
// también se traten como errores. El cliente no distingue el tipo de no-2xx.
func TestSendQR_NodeReturns4xx(t *testing.T) {
	// Caso: Node rechaza el payload con 400.
	ts := newFakeServer(http.StatusBadRequest, `{"error":"bad_request","message":"q must be array"}`)
	defer ts.Close()

	client := NewNodeClient(ts.URL)
	_, err := client.SendQR(context.Background(), samplePayload)

	if err == nil {
		t.Fatal("expected error for 400 response, got nil")
	}
}

// TestSendQR_MalformedJSONResponse verifica que si Node devuelve 200 pero con
// un body que no es JSON válido, el cliente retorna error de decodificación.
// Previene que datos corruptos lleguen al service.
func TestSendQR_MalformedJSONResponse(t *testing.T) {
	// Caso: Node devuelve HTML de error en lugar de JSON (e.g. proxy intermedio).
	ts := newFakeServer(http.StatusOK, "<html>Service Unavailable</html>")
	defer ts.Close()

	client := NewNodeClient(ts.URL)
	_, err := client.SendQR(context.Background(), samplePayload)

	if err == nil {
		t.Fatal("expected decode error for non-JSON response, got nil")
	}
}

// TestSendQR_ContextCancelled verifica que el cliente respeta la cancelación
// del contexto y no espera la respuesta de Node cuando la request original
// ya fue cancelada (e.g. el cliente HTTP desconectó).
func TestSendQR_ContextCancelled(t *testing.T) {
	// El servidor artificial tarda 'para siempre' — nunca responde.
	blockForever := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-blockForever // bloquea hasta que el canal se cierre
	}))
	defer ts.Close()
	defer close(blockForever)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancelar inmediatamente

	client := NewNodeClient(ts.URL)
	_, err := client.SendQR(ctx, samplePayload)

	// Con el contexto cancelado, Do() debe retornar error antes de completar.
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

// TestSendQR_SendsCorrectPayload verifica que el cliente serializa Q y R
// en el body de la request y los envía a Node con el Content-Type correcto.
func TestSendQR_SendsCorrectPayload(t *testing.T) {
	var receivedPayload models.NodePayload
	var receivedContentType string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capturar el Content-Type que el cliente envía.
		receivedContentType = r.Header.Get("Content-Type")

		// Capturar y validar que el body tenga Q y R.
		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			http.Error(w, "bad payload", http.StatusBadRequest)
			return
		}

		// Respuesta válida para que el cliente no falle en decode.
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.NodeResponse{})
	}))
	defer ts.Close()

	client := NewNodeClient(ts.URL)
	_, err := client.SendQR(context.Background(), samplePayload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verificar Content-Type enviado al servidor.
	if !strings.Contains(receivedContentType, "application/json") {
		t.Errorf("Content-Type: want 'application/json', got %q", receivedContentType)
	}

	// Verificar que las dimensiones de Q y R llegaron íntegras.
	if len(receivedPayload.Q) != len(samplePayload.Q) {
		t.Errorf("Q rows: want %d, got %d", len(samplePayload.Q), len(receivedPayload.Q))
	}
	if len(receivedPayload.R) != len(samplePayload.R) {
		t.Errorf("R rows: want %d, got %d", len(samplePayload.R), len(receivedPayload.R))
	}
}
