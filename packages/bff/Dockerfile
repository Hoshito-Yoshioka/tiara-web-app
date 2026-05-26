# ==============================
# BFF (Hono) - Multi-stage build
# Docker context はリポジトリルートを想定
# ==============================
FROM node:20-alpine AS builder

RUN corepack enable && corepack prepare pnpm@9 --activate

WORKDIR /app
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages/bff/package.json ./packages/bff/
RUN pnpm install --frozen-lockfile --filter bff

COPY packages/bff/ ./packages/bff/
RUN pnpm --filter bff build

# ==============================
# Runtime
# ==============================
FROM node:20-alpine

WORKDIR /app
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages/bff/package.json ./packages/bff/

RUN corepack enable && corepack prepare pnpm@9 --activate \
    && pnpm install --frozen-lockfile --filter bff --prod --ignore-scripts

COPY --from=builder /app/packages/bff/dist ./packages/bff/dist

WORKDIR /app/packages/bff

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:3001/ || exit 1

CMD ["node", "dist/index.js"]
