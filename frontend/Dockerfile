# ── Stage 1: build ───────────────────────────────────────────────────────────
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

# Sin VITE_API_URL → la app usa URLs relativas (/api/...).
# Nginx intercepta /api/* y lo proxies a go-api:8080 internamente.
RUN npm run build

# ── Stage 2: nginx como servidor estático + reverse proxy ────────────────────
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# nginx.conf:
#   - Sirve la SPA (index.html para cualquier ruta desconocida).
#   - Proxies /api/* → go-api:8080 (resolución DNS interna de Docker).
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
