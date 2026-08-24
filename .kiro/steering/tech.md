# Technology Stack

## Architecture

3 層構成の pnpm workspace モノレポ。Frontend（Vue 3 / vite-ssg）→ BFF（Hono / Node.js）→ Backend（Go / Echo）→ PostgreSQL の一方向にリクエストが流れる。

- **BFF パターン**: Frontend は Backend に直接アクセスせず、必ず BFF（`http://localhost:3001`、`/api/v1/*`）を経由する。BFF は zod スキーマによる入出力バリデーションとクライアント向けデータ整形、OpenAPI ドキュメント（Swagger UI / Redoc）の提供を担う
- **クリーンアーキテクチャ（Backend）**: 依存は常に内側（domain ← usecase ← interface / infrastructure）へ向ける。Echo・sqlc/pgx の詳細は infrastructure / interface 層に閉じ込める（docs/ADR.md 001 参照）
- **SSG**: 公開ページは `vite-ssg build` でプリレンダリング。管理画面・ポータルは SPA のまま（プリレンダリング対象外）

## Core Technologies

- **Frontend**: Vue 3（Composition API）/ TypeScript / Vite 5 / vite-ssg / Tailwind CSS / shadcn-vue（radix-vue）/ Pinia / vue-router / @unhead/vue
- **BFF**: Hono / @hono/zod-openapi / zod / Node.js（tsx 実行）
- **Backend**: Go 1.25 / Echo v4 / pgx v5 / sqlc / viper（設定）/ golang-jwt（認証）
- **DB**: PostgreSQL 16（Docker）
- **ツールチェーン**: pnpm 9 workspace / Node.js 22.11（.tool-versions で固定）

## Development Standards

### Code Quality
- ルートの flat config（eslint.config.js）で ESLint（typescript-eslint + eslint-plugin-vue）、Prettier 併用
- husky + lint-staged でコミット時に `packages/frontend/src`・`packages/bff/src` の ts/vue を自動 lint/format
- Go は golangci-lint（packages/backend/.golangci.yml）

### Testing
- Frontend: Vitest + @vue/test-utils（environment: happy-dom）、`src/__tests__/`
- BFF: Vitest、`src/__tests__/`
- Backend: Go 標準 testing + testify。`*_test.go` を実装と同居させる。モックは testutil パッケージ
- ルートから `pnpm test`（frontend + bff）で一括実行

### API 契約
- BFF のルートは @hono/zod-openapi の `createRoute` + zod スキーマ（`src/schemas/`）で定義し、OpenAPI 仕様を自動生成する
- Backend の DB アクセスは `queries/*.sql` から sqlc で生成（手書き SQL をコードに埋め込まない）

## Development Environment

### Common Commands
```bash
docker-compose up -d              # DB 起動
cd packages/backend && go run ./cmd/server/   # Backend (:1323)
pnpm --filter bff dev             # BFF (:3001)
pnpm --filter frontend dev        # Frontend (:5173)
pnpm test / pnpm lint             # ルートから一括
pnpm --filter frontend build      # sitemap 生成 + vue-tsc + vite-ssg build
```

- 環境変数はルート `.env`（.env.example 参照）。Backend は viper が自動読込。必須: DATABASE_URL / JWT_SECRET
- `VITE_ADMIN_BASE_PATH` で管理画面のベースパスを変更可能（本番は推測困難な値を推奨）

## Key Technical Decisions

- **SSG ビルド制約**: `vite-ssg build` と sitemap 生成は `SITEMAP_API_BASE_URL` からスタッフ一覧を fetch してスタッフ詳細ページをプリレンダリングする。未設定だと localhost:3001 フォールバックとなり本番ビルドでスタッフページが漏れる（CI では本番 URL を設定済み）
- **バージョン固定**: Node 22.11 のため jsdom はルート pnpm.overrides で `^26.1.0` に固定（jsdom 28 は Node 22.12+ 必須）。unhead は 2.1.17 に統一（vite-ssg 28 が unhead v2 前提）
- **公開ビューのデータ取得**: async setup + Suspense（App.vue）で行い、SSG HTML に本文を含める
- **本番構成**: Xserver VPS 上で docker-compose.prod.yml（db / backend / bff）を稼働。Frontend は Nginx が静的配信（infra/nginx/tiara.conf、VPS への反映は手動。CI は dist の rsync のみ）

---
_更新日: 2026-08-24（bootstrap 生成）_
