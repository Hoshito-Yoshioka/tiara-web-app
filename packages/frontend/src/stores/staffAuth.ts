import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiFetch } from '@/lib/api'
import type { StaffLoginResponse } from '@/types/staffPortal'

/**
 * スタッフ認証ストア。
 * 管理者用 useAuthStore とは独立して動作する。
 * localStorage キーも別（tiara_staff_token）。
 */
export const useStaffAuthStore = defineStore('staffAuth', () => {
  const token = ref<string | null>(localStorage.getItem('tiara_staff_token'))
  const staffId = ref<string | null>(localStorage.getItem('tiara_staff_id'))
  const isAuthenticated = computed(() => !!token.value)

  /** スタッフのユーザー名とパスワードでログインし、JWT トークンを取得・保存する */
  async function login(username: string, password: string): Promise<void> {
    const data = await apiFetch<StaffLoginResponse>('/api/staff-auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
    token.value = data.token
    staffId.value = data.staffId
    localStorage.setItem('tiara_staff_token', data.token)
    localStorage.setItem('tiara_staff_id', data.staffId)
  }

  /** ログアウトし、トークンを削除する */
  function logout(): void {
    token.value = null
    staffId.value = null
    localStorage.removeItem('tiara_staff_token')
    localStorage.removeItem('tiara_staff_id')
  }

  /** 保持しているトークンの有効性を Backend に問い合わせる */
  async function verify(): Promise<boolean> {
    if (!token.value) return false

    try {
      await apiFetch('/api/staff-auth/verify', {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
      return true
    } catch {
      logout()
      return false
    }
  }

  return {
    token,
    staffId,
    isAuthenticated,
    login,
    logout,
    verify,
  }
})
