// Package models define los DTOs (Data Transfer Objects) de la API.
package models

// MatrixRequest es el body del endpoint POST /api/v1/matrix/qr.
type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// Statistics contiene las métricas calculadas sobre los valores de Q y R.
// Es también el tipo que devuelve la API Node.js (NodeResponse).
type Statistics struct {
	Max         float64 `json:"max"`
	Min         float64 `json:"min"`
	Average     float64 `json:"average"`
	Sum         float64 `json:"sum"`
	IsDiagonalQ bool    `json:"isDiagonalQ"`
	IsDiagonalR bool    `json:"isDiagonalR"`
}

// MatrixResponse es la respuesta final del endpoint QR.
type MatrixResponse struct {
	Q          [][]float64 `json:"q"`
	R          [][]float64 `json:"r"`
	Statistics Statistics  `json:"statistics"`
}

// NodePayload es el payload que se envía a la API Node.js.
type NodePayload struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

// NodeResponse es un alias de Statistics: la API Node.js devuelve
// exactamente el mismo contrato que se expone al cliente final.
type NodeResponse = Statistics
