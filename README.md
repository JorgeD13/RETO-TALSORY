# Reto Técnico — QR Matrix API

Sistema compuesto por dos APIs independientes que se comunican entre sí:

| Servicio   | Tecnología    | Puerto | Responsabilidad |
|------------|---------------|--------|-----------------|
| `go-api`   | Go + Fiber    | 8080   | Calcula la factorización QR de una matriz con Gonum y orquesta el flujo |
| `node-api` | Node + Express | 3000  | Recibe Q y R, calcula estadísticas y las devuelve a la API Go |

---

## Estructura del proyecto

```
reto/
├── docker-compose.yml          # Orquestador raíz
├── README.md
│
├── go-api/
│   ├── cmd/main.go             # Punto de entrada + DI
│   ├── internal/
│   │   ├── config/             # Variables de entorno
│   │   ├── models/             # DTOs
│   │   ├── handlers/           # Controladores HTTP
│   │   ├── services/           # Lógica de negocio + tests
│   │   ├── clients/            # Cliente HTTP → node-api
│   │   ├── routes/             # Registro de rutas
│   │   └── middleware/         # Recovery + ErrorHandler
│   ├── docs/                   # Swagger generado por swag
│   ├── Dockerfile
│   ├── .env.example
│   └── go.mod
│
└── node-api/
    ├── src/
    │   ├── app.js              # Entrada Express
    │   ├── routes/             # Definición de rutas + anotaciones OpenAPI
    │   ├── controllers/        # Controladores HTTP
    │   ├── services/           # Lógica de estadísticas + tests
    │   ├── middleware/         # ErrorHandler global
    │   ├── models/             # JSDoc types
    │   └── utils/              # Utilidades (matrix, swagger)
    ├── Dockerfile
    ├── .env.example
    └── package.json
```

---

## Flujo completo

```
Cliente
  │
  │  POST /api/v1/matrix/qr  { "matrix": [[...]] }
  ▼
go-api (Fiber :8080)
  │  Valida → QR con Gonum → obtiene Q, R
  │
  │  POST /api/v1/statistics  { "q": [...], "r": [...] }
  ▼
node-api (Express :3000)
  │  Calcula: max, min, average, sum, isDiagonalQ, isDiagonalR
  ▼
go-api
  │
  ▼
Cliente  ←  { "q": [...], "r": [...], "statistics": { ... } }
```

---

## Requisitos previos

- [Go 1.24+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker + Docker Compose](https://docs.docker.com/get-docker/)

---

## Ejecución local (sin Docker)

### 1. node-api

```bash
cd node-api
cp .env.example .env
npm install
npm start
# → http://localhost:3000
```

### 2. go-api

```bash
cd go-api
cp .env.example .env
go mod tidy
go run ./cmd/main.go
# → http://localhost:8080
```

> El `.env` de `go-api` apunta a `NODE_API_URL=http://localhost:3000/api/v1/statistics` por defecto.

---

## Ejecución con Docker Compose

```bash
# Desde la raíz del proyecto (reto/)
docker compose up --build

# En background
docker compose up -d --build

# Ver logs
docker compose logs -f

# Detener
docker compose down
```

Los servicios se comunican internamente mediante la red `app-network`. La go-api resuelve `http://node-api:3000` por nombre DNS de Docker.

---

## Tests

### Go

```bash
cd go-api
go test ./... -v
```

### Node.js

```bash
cd node-api
npm test
```

---

## Endpoints

### go-api — `POST /api/v1/matrix/qr`

**Request:**
```json
{
  "matrix": [
    [12, -51, 4],
    [6, 167, -68],
    [-4, 24, -41]
  ]
}
```

**Response `200 OK`:**
```json
{
  "q": [
    [-0.8571,  0.3943,  0.3314],
    [-0.4286, -0.9029, -0.0292],
    [ 0.2857, -0.1714,  0.9429]
  ],
  "r": [
    [-14,  -21,   14],
    [  0, -175,   70],
    [  0,    0,  -35]
  ],
  "statistics": {
    "max": 70,
    "min": -175,
    "average": -12.34,
    "sum": -222.1,
    "isDiagonalQ": false,
    "isDiagonalR": false
  }
}
```

**Error `422 Unprocessable Entity`:**
```json
{
  "error": "unprocessable_entity",
  "message": "row 1 has 2 columns, expected 3"
}
```

### node-api — `POST /api/v1/statistics`

**Request:**
```json
{
  "q": [[-0.857, 0.394], [-0.429, -0.903]],
  "r": [[-14, -21], [0, -175]]
}
```

**Response `200 OK`:**
```json
{
  "max": 0,
  "min": -175,
  "average": -48.9,
  "sum": -391.2,
  "isDiagonalQ": false,
  "isDiagonalR": false
}
```

---

## Swagger / Documentación

| API      | URL                                          |
|----------|----------------------------------------------|
| go-api   | http://localhost:8080/swagger/index.html     |
| node-api | http://localhost:3000/api-docs               |

---

## Health checks

```bash
curl http://localhost:8080/health   # go-api
curl http://localhost:3000/health   # node-api
```
