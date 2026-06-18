import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiFetch } from '@/lib/api'
import type { StaffLoginResponse, RefreshTokenResponse } from '@/types/staffPortal'

/**
 * スタッフ認証ストア。
 * 管理者用 useAuthStore とは独立して動作する。
 * localStorage キーも別（tiara_staff_token / tiara_staff_refresh_token）。
 */
export const useStaffAuthStore = defineStore('staffAuth', () => {
  const token = ref<string | null>(localStorage.getItem('tiara_staff_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('tiara_staff_refresh_token'))
  const staffId = ref<string | null>(localStorage.getItem('tiara_staff_id'))
  const isAuthenticated = computed(() => !!token.value)

  /** スタッフのユーザー名とパスワードでログインし、JWT トークンを取得・保存する */
  async function login(username: string, password: string): Promise<void> {
    const data = await apiFetch<StaffLoginResponse>('/api/v1/staff-auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
    token.value = data.token
    refreshToken.value = data.refreshToken
    staffId.value = data.staffId
    localStorage.setItem('tiara_staff_token', data.token)
    localStorage.setItem('tiara_staff_refresh_token', data.refreshToken)
    localStorage.setItem('tiara_staff_id', data.staffId)
  }

  /** ログアウトし、トークンを削除する */
  async function logout(): Promise<void> {
    if (refreshToken.value) {
      try {
        await apiFetch('/api/v1/staff-auth/logout', {
          method: 'POST',
          body: JSON.stringify({ refreshToken: refreshToken.value }),
        })
      } catch {
        // サーバーエラーでもクライアント側のトークンは削除する
      }
    }
    token.value = null
    refreshToken.value = null
    staffId.value = null
    localStorage.removeItem('tiara_staff_token')
    localStorage.removeItem('tiara_staff_refresh_token')
    localStorage.removeItem('tiara_staff_id')
  }

  /** 保持しているトークンの有効性を Backend に問い合わせる */
  async function verify(): Promise<boolean> {
    if (!token.value) return false

    try {
      await apiFetch('/api/v1/staff-auth/verify', {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
      return true
    } catch {
      // アクセストークン失効 → リフレッシュ試行
      return await refresh()
    }
  }

  /** リフレッシュトークンを使ってアクセストークンを再発行する */
  async function refresh(): Promise<boolean> {
    if (!refreshToken.value) {
      await logout()
      return false
    }
    try {
      const data = await apiFetch<RefreshTokenResponse>('/api/v1/staff-auth/refresh', {
        method: 'POST',
        body: JSON.stringify({ refreshToken: refreshToken.value }),
      })
      token.value = data.token
      refreshToken.value = data.refreshToken
      localStorage.setItem('tiara_staff_token', data.token)
      localStorage.setItem('tiara_staff_refresh_token', data.refreshToken)
      return true
    } catch {
      await logout()
      return false
    }
  }

  return {
    token,
    refreshToken,
    staffId,
    isAuthenticated,
    login,
    logout,
    verify,
    refresh,
  }
})
