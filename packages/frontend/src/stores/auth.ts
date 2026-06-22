import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiFetch } from '@/lib/api'
import type { LoginResponse } from '@/types/admin'

/**
 * 認証ストア。
 * JWT トークンの管理（取得・保存・削除・検証）を担う。
 * Pinia の Setup Store 構文を使用（Composition API スタイル）。
 *
 * トークンは localStorage に永続化し、ページリロード後も認証状態を維持する。
 */
export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('tiara_admin_token'))
  const isAuthenticated = computed(() => !!token.value)

  /** ユーザー名とパスワードでログインし、JWT トークンを取得・保存する */
  async function login(username: string, password: string): Promise<void> {
    const data = await apiFetch<LoginResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
    token.value = data.token
    localStorage.setItem('tiara_admin_token', data.token)
  }

  /** ログアウトし、トークンを削除する */
  function logout(): void {
    token.value = null
    localStorage.removeItem('tiara_admin_token')
  }

  /** 保持しているトークンの有効性を Backend に問い合わせる */
  async function verify(): Promise<boolean> {
    if (!token.value) return false

    try {
      await apiFetch('/api/v1/auth/verify', {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
      return true
    } catch {
      // verify失敗時: 他タブでログアウトされている可能性があるため、
      // localStorage を再チェック。あれば同期、なければ logout する。
      const latestToken = localStorage.getItem('tiara_admin_token')

      if (latestToken && latestToken !== token.value) {
        // 他タブで更新されたトークンを検出 → Pinia に同期
        token.value = latestToken
        return true
      }

      // トークンがない、またはverifyに本当に失敗 → ログアウト
      logout()
      return false
    }
  }

  return {
    token,
    isAuthenticated,
    login,
    logout,
    verify,
  }
})
