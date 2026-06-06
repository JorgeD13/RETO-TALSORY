# Frontend — Vite + React + Tailwind CSS

SPA que permite ingresar una matriz numérica, enviarla a la Go API
y visualizar las matrices Q, R y las estadísticas resultantes.

---

## Stack

| Tecnología | Uso |
|---|---|
| Vite 6 | Bundler y dev server |
| React 19 | UI |
| Tailwind CSS 4 | Estilos |

---

## Estructura

```
frontend/
├── src/
│   ├── api/
│   │   └── matrix.js          # Llamada HTTP a la Go API
│   ├── hooks/
│   │   └── useQR.js           # Estado (loading, error, result)
│   ├── components/
│   │   ├── MatrixInput.jsx    # Grilla editable dinámica
│   │   ├── MatrixDisplay.jsx  # Visualización de Q y R
│   │   └── StatisticsDisplay.jsx  # Tarjetas de métricas
│   ├── App.jsx                # Orquestador principal
│   └── main.jsx               # Entry point React
├── nginx.conf                 # Proxy /api → go-api (Docker)
├── Dockerfile
└── .env.example
```

---

## Variables de entorno

Solo necesaria en producción/Docker. En desarrollo el proxy de Vite
redirige `/api/*` automáticamente a `localhost:8080`.

```bash
cp .env.example .env
```

| Variable | Descripción |
|---|---|
| `VITE_API_URL` | URL base de la Go API (vacío = URLs relativas) |

---

## Ejecutar localmente

Requiere que la Go API esté corriendo en `localhost:8080`.

```bash
npm install
npm run dev
# → http://localhost:5173
```

El proxy de Vite está configurado en `vite.config.js`:
cualquier request a `/api/*` se redirige a `http://localhost:8080`.

---

## Build de producción

```bash
npm run build
# Genera dist/ listo para servir con cualquier servidor estático
```

---

## Ejecutar con Docker

Desde la raíz del proyecto (`reto/`):

```bash
docker compose up --build
# → http://localhost
```

En Docker, nginx sirve el build estático y actúa como reverse proxy:
`/api/*` se reenvía internamente a `go-api:8080` sin exponer esa URL al browser.
