#!/bin/sh
# ============================================================
# DB マイグレーションスクリプト（スキーマのみ）
# DATABASE_URL 環境変数を使用して PostgreSQL に接続する。
# べき等: テーブルが既に存在する場合はスキップ。
#
# シードデータは entrypoint.sh が APP_ENV に応じて制御するため、
# このスクリプトでは実行しない。
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

echo "Migration complete."
