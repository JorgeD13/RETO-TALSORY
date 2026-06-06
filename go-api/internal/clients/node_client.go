// Package clients encapsula los clientes HTTP hacia servicios externos.
package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go-api/internal/models"
)

// NodeClient define el contrato para comunicarse con la API Node.js.
type NodeClient interface {
	SendQR(ctx context.Context, payload models.NodePayload) (*models.NodeResponse, error)
}

// httpNodeClient es la implementación real que usa net/http.
type httpNodeClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewNodeClient construye un NodeClient con timeout configurado.
func NewNodeClient(baseURL string) NodeClient {
	return &httpNodeClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendQR serializa el payload Q/R y lo envía a la API Node.js.
// Acepta un contexto para permitir cancelación y timeouts externos.
func (c *httpNodeClient) SendQR(ctx context.Context, payload models.NodePayload) (*models.NodeResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("node_client: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("node_client: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[node_client] POST %s", c.baseURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("node_client: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("node_client: failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("node_client: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var nodeResp models.NodeResponse
	if err := json.Unmarshal(respBody, &nodeResp); err != nil {
		return nil, fmt.Errorf("node_client: failed to decode response: %w", err)
	}

	log.Printf("[node_client] response received: max=%.4f min=%.4f avg=%.4f", nodeResp.Max, nodeResp.Min, nodeResp.Average)
	return &nodeResp, nil
}
