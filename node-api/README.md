# Statistics API — Node.js + Express

API secundaria encargada de calcular estadísticas sobre las matrices Q y R
recibidas desde la Go API tras la factorización QR.

---

## Stack

| Tecnología | Uso |
|---|---|
| Node.js 18+ | Runtime |
| Express 4 | Framework HTTP |
| Jest + Supertest | Testing |
| swagger-jsdoc + swagger-ui-express | Documentación OpenAPI |
| dotenv | Variables de entorno |

---

## Estructura

```
node-api/
├── src/
│   ├── server.js              # Arranca el servidor HTTP (separado de app.js)
│   ├── app.js                 # Configura Express, middlewares y rutas
│   ├── routes/
│   │   └── statistics.routes.js    # Registro de rutas + anotaciones OpenAPI
│   ├── controllers/
│   │   └── statistics.controller.js  # Valida entrada, llama al service
│   ├── services/
│   │   └── statistics.service.js     # Lógica de cálculo pura
│   ├── middleware/
│   │   └── errorHandler.js    # Manejador global de errores Express
│   ├── models/
│   │   └── statistics.js      # JSDoc types (contratos documentados)
│   ├── utils/
│   │   ├── matrix.utils.js    # flatten, isDiagonal, validateMatrix2D
│   │   └── swagger.js         # Configuración OpenAPI
│   └── __tests__/
│       ├── statistics.controller.test.js
│       ├── statistics.service.test.js
│       └── matrix.utils.test.js
├── Dockerfile
├── .env.example
└── package.json
```

---

## Variables de entorno

```bash
cp .env.example .env
```

| Variable | Default | Descripción |
|---|---|---|
| `PORT` | `3000` | Puerto del servidor |
| `NODE_ENV` | `development` | Entorno de ejecución |

---

## Ejecutar localmente

```bash
npm install
npm start
# → http://localhost:3000
```

### Modo desarrollo (recarga automática)

```bash
npm run dev
```

---

## Ejecutar con Docker

El servicio se levanta junto al resto del sistema desde la raíz del proyecto:

```bash
# Desde reto/
docker compose up --build
```

---

## Endpoint

### `POST /api/v1/statistics`

Recibe las matrices Q y R y devuelve estadísticas calculadas sobre el conjunto de todos sus elementos.

#### Request

```json
{
  "q": [
    [-0.8571,  0.3943,  0.3314],
    [-0.4286, -0.9029, -0.0292],
    [ 0.2857, -0.1714,  0.9429]
  ],
  "r": [
    [-14,    -21,  14],
    [  0,  -175,  70],
    [  0,     0, -35]
  ]
}
```

#### Response `200 OK`

```json
{
  "max": 70,
  "min": -175,
  "average": -8.97,
  "sum": -161.4,
  "isDiagonalQ": false,
  "isDiagonalR": false
}
```

#### Error `400 Bad Request`

```json
{
  "error": "bad_request",
  "message": "q: must not be empty"
}
```

---

## Validaciones

El endpoint valida que `q` y `r`:

- Sean arrays (no null, no string, no número)
- No estén vacíos
- No tengan filas vacías
- Sean rectangulares (todas las filas con igual longitud)
- Contengan solo números finitos (rechaza `NaN`, `Infinity`, strings)

---

## Tests

```bash
npm test
```

Cobertura actual: **97%+** sobre statements y funciones.

```
Test Suites: 3 passed
Tests:       53 passed
```

| Suite | Qué prueba |
|---|---|
| `statistics.controller.test.js` | Capa HTTP con Supertest (validaciones, status codes, esquema de error) |
| `statistics.service.test.js` | Cálculo de max, min, average, sum, isDiagonal |
| `matrix.utils.test.js` | flatten, isSquare, isDiagonal, validateMatrix2D |

---

## Documentación

Swagger UI disponible con el servidor corriendo:

```
http://localhost:3000/api-docs
```

---

## Health check

```bash
curl http://localhost:3000/health
# {"status":"ok","env":"development"}
```
