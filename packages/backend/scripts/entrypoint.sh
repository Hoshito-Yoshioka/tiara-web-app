#!/bin/sh
# ============================================================
# Docker エントリーポイント
# 1. マイグレーション & シード実行（べき等）
# 2. Go サーバー起動
# ============================================================
set -e

echo "Running migrations..."
sh /app/scripts/migrate.sh

echo "Starting server..."
exec ./server
