#!/bin/bash
# ============================================================
# Let's Encrypt SSL 証明書セットアップ — Nginx + certbot
#
# 前提:
#   - Nginx がインストール済みで tiara.conf が sites-enabled にある
#   - ドメインの DNS が VPS の IP を指している
#   - ポート 80 がファイアウォールで許可されている
#
# 実行: sudo bash infra/scripts/setup-certbot.sh your-domain.com
# ============================================================
set -euo pipefail

DOMAIN="${1:-}"

if [ -z "$DOMAIN" ]; then
    echo "使い方: sudo bash infra/scripts/setup-certbot.sh your-domain.com"
    exit 1
fi

echo "=== Let's Encrypt SSL 証明書セットアップ ==="
echo "ドメイン: ${DOMAIN}"
echo ""

# certbot インストール
if ! command -v certbot &> /dev/null; then
    echo "📦 certbot をインストール中..."
    apt-get update
    apt-get install -y certbot python3-certbot-nginx
else
    echo "✅ certbot は既にインストール済み"
fi

# Nginx 設定の構文チェック
echo ""
echo "🔍 Nginx 設定をテスト中..."
nginx -t

# certbot で証明書取得 + Nginx 設定を自動更新
echo ""
echo "🔒 SSL 証明書を取得中..."

CERTBOT_ARGS=(
    --nginx
    -d "$DOMAIN"
    --non-interactive
    --agree-tos
    --redirect
)

if [ -n "${CERTBOT_EMAIL:-}" ]; then
    CERTBOT_ARGS+=(--email "$CERTBOT_EMAIL")
else
    CERTBOT_ARGS+=(--register-unsafely-without-email)
fi

certbot "${CERTBOT_ARGS[@]}"

# 自動更新の確認
echo ""
echo "🔄 自動更新のテスト..."
certbot renew --dry-run

# 自動更新の cron/timer 確認
echo ""
if systemctl is-active --quiet certbot.timer; then
    echo "✅ certbot.timer が有効です（自動更新）"
else
    echo "⚠️ certbot.timer が無効です。手動で cron を設定します..."
    # cron による自動更新（毎日2回チェック — Let's Encrypt 推奨）
    CRON_JOB="0 0,12 * * * certbot renew --quiet --post-hook 'systemctl reload nginx'"
    if ! crontab -l 2>/dev/null | grep -q "certbot renew"; then
        (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
        echo "✅ cron ジョブを追加しました"
    else
        echo "✅ cron ジョブは既に存在します"
    fi
fi

echo ""
echo "=== SSL 設定完了 ==="
echo "   https://${DOMAIN} でアクセス可能です"
echo "   証明書は自動で更新されます（有効期限: 90日 / 自動更新: 60日ごと）"
