# ── Stage 1: dependencias ────────────────────────────────────────────────────
FROM node:20-alpine AS deps

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

# ── Stage 2: runtime ─────────────────────────────────────────────────────────
FROM node:20-alpine

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY src ./src
COPY package.json ./
COPY .env.example .env

EXPOSE 3000

CMD ["node", "src/server.js"]
