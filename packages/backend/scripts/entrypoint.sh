#!/bin/sh
# ============================================================
# Docker エントリーポイント
# 1. マイグレーション実行（べき等）
# 2. シード実行（APP_ENV=production の場合はスキップ）
# 3. Go サーバー起動
# ============================================================
set -e

echo "Running migrations..."
sh /app/scripts/migrate.sh

if [ "${APP_ENV}" = "production" ]; then
  echo "Skipping seed data (APP_ENV=production)"
else
  echo "Applying seed data (development mode)..."
  psql "${DATABASE_URL}" -f /app/migrations/seed.sql
fi

echo "Starting server..."
exec ./server
