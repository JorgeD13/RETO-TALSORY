# QR Matrix API — Go + Fiber

API REST que recibe una matriz numérica, calcula su factorización QR con **Gonum** y obtiene estadísticas desde una API Node.js.

---

## Stack

| Tecnología | Uso |
|------------|-----|
| Go 1.24+   | Lenguaje principal |
| Fiber v2   | Framework HTTP |
| Gonum      | Factorización QR |
| godotenv   | Variables de entorno |
| Docker     | Contenedorización |
| Swagger    | Documentación |

---

## Estructura

```
go-api/
├── cmd/main.go              # Punto de entrada + DI
├── internal/
│   ├── config/              # Carga de variables de entorno
│   ├── models/              # DTOs (request / response)
│   ├── handlers/            # Controladores HTTP
│   ├── services/            # Lógica de negocio + validación
│   ├── clients/             # Cliente HTTP → Node.js API
│   ├── routes/              # Registro de rutas
│   └── middleware/          # Recovery + error handler
├── docs/                    # Generado por swag init
├── .env.example
├── Dockerfile
└── docker-compose.yml
```

---

## Variables de entorno

Copia `.env.example` a `.env` y ajusta los valores:

```bash
cp .env.example .env
```

| Variable      | Default                                      | Descripción |
|---------------|----------------------------------------------|-------------|
| `APP_PORT`    | `8080`                                       | Puerto del servidor |
| `APP_ENV`     | `development`                                | Entorno de ejecución |
| `NODE_API_URL`| `http://localhost:3000/api/v1/statistics`    | URL de la API Node.js |

---

## Ejecutar localmente

### 1. Instalar dependencias

```bash
go mod tidy
```

### 2. Arrancar la API

```bash
go run ./cmd/main.go
```

El servidor queda disponible en `http://localhost:8080`.

### 3. Regenerar Swagger (opcional)

```bash
# Instalar swag CLI si no lo tienes
go install github.com/swaggo/swag/cmd/swag@latest

# Generar docs/
swag init -g cmd/main.go -o docs
```

Swagger UI: `http://localhost:8080/swagger/index.html`

---

## Ejecutar con Docker

```bash
# Construir y levantar
docker-compose up --build

# En background
docker-compose up -d --build

# Detener
docker-compose down
```

---

## Endpoint

### `POST /api/v1/matrix/qr`

#### Request

```json
{
  "matrix": [
    [12, -51, 4],
    [6, 167, -68],
    [-4, 24, -41]
  ]
}
```

#### Response exitosa `200 OK`

```json
{
  "q": [
    [-0.8571,  0.3943,  0.3314],
    [-0.4286, -0.9029, -0.0292],
    [ 0.2857, -0.1714,  0.9429]
  ],
  "r": [
    [-14.0000, -21.0000,  14.0000],
    [  0.0000, -175.0000,  70.0000],
    [  0.0000,    0.0000, -35.0000]
  ],
  "statistics": {
    "max": 0,
    "min": 0,
    "average": 0,
    "sum": 0,
    "isDiagonalQ": false,
    "isDiagonalR": false
  }
}
```

> **Nota:** si la API Node.js no está disponible, `statistics` devuelve ceros y la API Go responde igualmente con `200 OK`.

#### Errores de validación `422`

```json
{
  "error": "unprocessable_entity",
  "message": "row 1 has 2 columns, expected 3"
}
```

#### Body inválido `400`

```json
{
  "error": "bad_request",
  "message": "invalid request body: ..."
}
```

---

## Health check

```bash
curl http://localhost:8080/health
# {"status":"ok","env":"development"}
```
