# TIARA Web App

「TIARA」の Web サイト。店舗・スタッフ情報を顧客に提供するモノレポ構成のフルスタックアプリケーション。

## Tech Stack

| レイヤー | 技術                                                  |
| -------- | ----------------------------------------------------- |
| Frontend | Vue 3 / TypeScript / Vite / Tailwind CSS / shadcn-vue |
| BFF      | Hono / Node.js / TypeScript                           |
| Backend  | Go / Echo / PostgreSQL / sqlc / pgx                   |
| Infra    | Docker / docker-compose / pnpm workspaces             |

---

## 必要なツール

以下をあらかじめインストールしてください。

| ツール                                                            | 推奨バージョン | インストール先          |
| ----------------------------------------------------------------- | -------------- | ----------------------- |
| [Node.js](https://nodejs.org/)                                    | v20 以上       | https://nodejs.org/     |
| [pnpm](https://pnpm.io/)                                          | v9 以上        | `npm install -g pnpm`   |
| [Go](https://go.dev/)                                             | v1.22 以上     | https://go.dev/dl/      |
| [Docker Desktop](https://www.docker.com/products/docker-desktop/) | 最新安定版     | https://www.docker.com/ |

バージョン確認コマンド：

```bash
node -v
pnpm -v
go version
docker -v
```

---

## セットアップ手順

### 1. リポジトリをクローン

```bash
git clone <repository-url>
cd tiara-web-app
```

### 2. 環境変数を設定

リポジトリルートにある `.env.example` をコピーして `.env` を作成し、値を編集します。

```bash
cp .env.example .env
```

```dotenv
# .env
POSTGRES_DB=tiara_db
POSTGRES_USER=your_db_user          # 任意のユーザー名
POSTGRES_PASSWORD=your_db_password  # 任意のパスワード
DATABASE_URL=postgres://your_db_user:your_db_password@localhost:5432/tiara_db
```

### 3. npm パッケージをインストール

```bash
pnpm install
```

> ルートで実行 >> `frontend` / `bff` すべてのパッケージが一括インストールされます。

---

## 起動手順

以下の **3 つのプロセス** をそれぞれ別ターミナルで起動します。

### ① データベース（PostgreSQL）を起動

```bash
docker-compose up -d
```

起動確認：

```bash
docker-compose ps
# db コンテナが running であればOK
```

### ② バックエンド（Go / Echo）を起動

```bash
cd packages/backend
source ../../.env && go run ./cmd/server/main.go
```

> **Note:** Go は `.env` ファイルを自動で読み込みません。`source` で環境変数をシェルにロードしてから起動します。

起動確認：ターミナルに `Successfully connected to the database!` が表示されれば OK。

> `http://localhost:1323` でAPI が起動します。

### ③ BFF（Hono）を起動

```bash
pnpm --filter bff dev
```

> `http://localhost:3000`（または `src/index.ts` に記載のポート）で BFF が起動します。

### ④ フロントエンド（Vue / Vite）を起動

```bash
pnpm --filter frontend dev
# または
pnpm dev:frontend
```

> `http://localhost:5173` でページが確認できます。

---

## ページ一覧

| パス         | ページ           |
| ------------ | ---------------- |
| `/`          | ホーム           |
| `/shop`      | 店舗紹介         |
| `/staff`     | スタッフ一覧     |
| `/staff/:id` | スタッフ詳細     |
| `/schedule`  | 出勤スケジュール |
| `/price`     | 料金システム     |
| `/access`    | アクセス         |

---

## プロジェクト構成

```
tiara-web-app/
├── packages/
│   ├── frontend/       # Vue 3 (Vite)
│   ├── bff/            # Hono (Node.js)
│   └── backend/        # Go (Echo) / Clean Architecture
│       ├── cmd/server/ # エントリーポイント
│       └── internal/
│           ├── domain/
│           ├── usecase/
│           ├── interface/
│           └── infrastructure/
├── docker-compose.yml
├── .env.example
├── pnpm-workspace.yaml
└── package.json
```

---

## よくあるエラー

**`DATABASE_URL environment variable is not set`**
→ Go は `.env` を自動で読み込みません。起動コマンドに `source ../../.env &&` を付けて実行してください（[② バックエンドを起動](#②-バックエンドgo--echoを起動) 参照）。`.env` 自体が存在しない場合は [環境変数の設定](#2-環境変数を設定) を確認してください。

**`Unable to connect to database`**
→ Docker コンテナが起動していません。`docker-compose up -d` を実行してから再試行してください。

**`command not found: pnpm`**
→ `npm install -g pnpm` でインストールしてください。
