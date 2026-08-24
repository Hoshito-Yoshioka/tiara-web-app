# Project Structure

## Organization Philosophy

pnpm workspace によるモノレポで、レイヤーごとに独立したパッケージ（frontend / bff / backend）に分割する。各パッケージ内はそれぞれの技術スタックの標準構成に従う（Frontend: 役割別ディレクトリ、Backend: クリーンアーキテクチャの層別ディレクトリ）。

## Directory Patterns

### Workspace ルート
**Location**: `/`
**Purpose**: 共通ツール設定（eslint.config.js / .prettierrc / husky / lint-staged）、docker-compose 各種、`.env`。パッケージ横断のスクリプトは `pnpm --filter <pkg>` 形式で定義

### Frontend（`packages/frontend/src/`）
- `views/`: ページコンポーネント。公開ページは直下（`HomeView.vue` 等）、管理画面は `views/admin/`、スタッフポータルは `views/portal/` に分離
- `components/ui/`: shadcn-vue ベースの再利用 UI プリミティブ、`components/layout/`: ヘッダー等のレイアウト
- `composables/`: API アクセスは `useXxxApi.ts` 形式の composable に集約（例: `useStaffApi.ts`、`useAdminApi.ts`）。ページメタは `usePageMeta.ts`
- `lib/`: 共通ユーティリティ（`api.ts` = fetch ラッパー、`seo.ts`、`utils.ts`）
- `stores/`: Pinia ストア。認証状態のみ（`auth.ts` = 管理者、`staffAuth.ts` = スタッフ）
- `router/`: vue-router 定義、`types/`: 型定義、`__tests__/`: Vitest テスト

### BFF（`packages/bff/src/`）
- `routes/`: リソース単位のルート定義（`staff.ts`、`shop.ts` 等）。管理系は `routes/admin/` に分離
- `schemas/`: zod スキーマ（OpenAPI 契約の源泉）
- `middleware/`: 認証等のミドルウェア、`types/`: 型定義、`app.ts`: ルート集約、`index.ts`: エントリ

### Backend（`packages/backend/`）— クリーンアーキテクチャ
- `cmd/server/`: エントリーポイント（main.go / routes.go）
- `internal/domain/`: エンティティとドメインエラー（フレームワーク非依存）
- `internal/usecase/`: ビジネスロジック。リソース単位のファイル（staff.go、menu.go 等）
- `internal/interface/`: handler/（HTTP ハンドラ）と middleware/
- `internal/infrastructure/db/`: sqlc 生成コードと DB 実装
- `internal/config/`: viper 設定、`internal/testutil/`: テスト用モック
- `queries/*.sql` + `sqlc.yaml`: sqlc の入力、`migrations/`: DB マイグレーション

### インフラ（`infra/`）
- `nginx/tiara.conf`: 本番 Nginx 設定（リポジトリ管理・VPS 反映は手動）、`scripts/`: VPS セットアップスクリプト

### ドキュメント（`docs/`）
- ADR（`ADR.md`）や技術評価メモを配置

## Naming Conventions

- **Vue コンポーネント**: PascalCase。ページは `XxxView.vue` サフィックス
- **composables**: `useXxx.ts`（API アクセスは `useXxxApi.ts`）
- **BFF ルート/スキーマ**: リソース名の kebab-case / 単数形 ts ファイル（`staff-auth.ts` 等）
- **Go**: パッケージは層名、ファイルはリソース単位の snake_case（`staff_portal.go`）、テストは `*_test.go` 同居

## Import Organization

```typescript
// Frontend: '@' は packages/frontend/src へのエイリアス（vite.config.ts / tsconfig）
import { Button } from '@/components/ui/button'
import { useStaffApi } from '@/composables/useStaffApi'
```

- Go の import は `tiara-web-app/backend/internal/...`。依存方向は常に内側（domain）へ。domain は他層を import しない

## Code Organization Principles

- Frontend からのデータ取得は composable（`useXxxApi`）→ `lib/api.ts` 経由に統一し、コンポーネントに fetch を直書きしない
- BFF がクライアント契約（zod / OpenAPI）の唯一の定義場所。Frontend と Backend の型のズレは BFF で吸収する
- Backend のビジネスロジックは usecase 層に置き、handler は入出力変換に徹する。DB 操作は sqlc 生成コードのみ使用
- 公開ページ / 管理画面（admin）/ スタッフポータル（portal, mypage）の 3 区画を views・routes ともにディレクトリで分離する

---
_更新日: 2026-08-24（bootstrap 生成）_
