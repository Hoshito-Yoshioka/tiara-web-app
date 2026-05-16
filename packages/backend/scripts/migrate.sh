#!/bin/sh
# ============================================================
# DB マイグレーション & シードスクリプト（Render 用）
# DATABASE_URL 環境変数を使用して PostgreSQL に接続する。
# べき等: テーブルが既に存在する場合はスキップ、seed はコンフリクト時に何もしない。
# ============================================================
set -e

if [ -z "$DATABASE_URL" ]; then
  echo "ERROR: DATABASE_URL is not set"
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MIGRATIONS_DIR="$SCRIPT_DIR/../migrations"

echo "Running schema migration..."
psql "$DATABASE_URL" -f "$MIGRATIONS_DIR/schema.sql"

echo "Running seed data..."
psql "$DATABASE_URL" -f "$MIGRATIONS_DIR/seed.sql"

echo "Migration complete."
