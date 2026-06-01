#!/bin/bash
# ============================================================
# SSH ハードニング
#
# - パスワード認証を無効化（鍵認証のみ）
# - root ログインを無効化
# - 設定変更前に必ず鍵認証でログインできることを確認
#
# 実行: sudo bash infra/scripts/setup-ssh.sh
#
# ⚠️ 注意: このスクリプトを実行する前に、必ず SSH 公開鍵が
#    ~/.ssh/authorized_keys に登録済みであることを確認すること。
#    鍵なしでパスワード認証を無効化するとログインできなくなる。
# ============================================================
set -euo pipefail

SSHD_CONFIG="/etc/ssh/sshd_config"
BACKUP="${SSHD_CONFIG}.bak.$(date +%Y%m%d%H%M%S)"

echo "=== SSH ハードニング ==="

# 事前チェック: authorized_keys が存在するか
CURRENT_USER="${SUDO_USER:-$USER}"
AUTH_KEYS="/home/${CURRENT_USER}/.ssh/authorized_keys"

if [ ! -f "$AUTH_KEYS" ] || [ ! -s "$AUTH_KEYS" ]; then
    echo "❌ エラー: ${AUTH_KEYS} が見つからないか空です"
    echo "   先に SSH 公開鍵を登録してください:"
    echo "   ssh-copy-id ${CURRENT_USER}@<VPS_IP>"
    exit 1
fi

echo "✅ SSH 公開鍵を確認: ${AUTH_KEYS}"
echo "   $(wc -l < "$AUTH_KEYS") 個の鍵が登録済み"

# バックアップ
echo ""
echo "📋 sshd_config をバックアップ: ${BACKUP}"
cp "$SSHD_CONFIG" "$BACKUP"

# パスワード認証を無効化
echo "🔒 パスワード認証を無効化..."
sed -i 's/^#\?PasswordAuthentication.*/PasswordAuthentication no/' "$SSHD_CONFIG"
sed -i 's/^#\?ChallengeResponseAuthentication.*/ChallengeResponseAuthentication no/' "$SSHD_CONFIG"

# root ログインを無効化
echo "🔒 root ログインを無効化..."
sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin no/' "$SSHD_CONFIG"

# 公開鍵認証を明示的に有効化
sed -i 's/^#\?PubkeyAuthentication.*/PubkeyAuthentication yes/' "$SSHD_CONFIG"

# 設定をテスト
echo ""
echo "🔍 sshd 設定をテスト中..."
if sshd -t; then
    echo "✅ 設定テスト OK"
else
    echo "❌ 設定にエラーがあります。バックアップから復元します..."
    cp "$BACKUP" "$SSHD_CONFIG"
    exit 1
fi

# sshd を再起動（Ubuntu 24.04 では ssh.service、旧バージョンでは sshd.service）
echo ""
echo "🔄 sshd を再起動中..."
if systemctl list-units --type=service --all | grep -q 'ssh\.service'; then
    systemctl restart ssh
elif systemctl list-units --type=service --all | grep -q 'sshd\.service'; then
    systemctl restart sshd
else
    echo "❌ SSH サービスが見つかりません"
    exit 1
fi

echo ""
echo "=== 設定完了 ==="
echo "   PasswordAuthentication: no"
echo "   PermitRootLogin: no"
echo "   PubkeyAuthentication: yes"
echo ""
echo "⚠️ このターミナルは閉じずに、別ターミナルから SSH 接続できることを確認してください"
echo "   ssh ${CURRENT_USER}@<VPS_IP>"
