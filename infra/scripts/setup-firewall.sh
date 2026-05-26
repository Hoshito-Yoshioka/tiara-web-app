#!/bin/bash
# ============================================================
# ファイアウォール設定 — ufw
#
# 許可ポート: 22 (SSH), 80 (HTTP), 443 (HTTPS)
# それ以外はすべて拒否（Backend 1323 / BFF 3001 / PostgreSQL 5432 は外部遮断）
#
# 実行: sudo bash infra/scripts/setup-firewall.sh
# ============================================================
set -euo pipefail

echo "=== ファイアウォール設定 ==="

# ufw がインストールされていない場合
if ! command -v ufw &> /dev/null; then
    echo "ufw をインストール中..."
    apt-get update && apt-get install -y ufw
fi

# デフォルトポリシー: 受信拒否・送信許可
ufw default deny incoming
ufw default allow outgoing

# 許可ポート
ufw allow 22/tcp    comment 'SSH'
ufw allow 80/tcp    comment 'HTTP'
ufw allow 443/tcp   comment 'HTTPS'

# ufw を有効化（既に有効な場合はリロード）
if ufw status | grep -q "Status: active"; then
    echo "ufw をリロード中..."
    ufw reload
else
    echo "ufw を有効化中..."
    echo "y" | ufw enable
fi

echo ""
echo "=== 現在のファイアウォール状態 ==="
ufw status verbose

echo ""
echo "✅ ファイアウォール設定完了"
echo "   許可: SSH(22), HTTP(80), HTTPS(443)"
echo "   拒否: その他すべて（Backend 1323 / BFF 3001 / PostgreSQL 5432 含む）"
