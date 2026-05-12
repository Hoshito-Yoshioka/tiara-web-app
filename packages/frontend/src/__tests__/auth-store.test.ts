import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('初期状態では未認証', () => {
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBeNull()
  })

  it('logout でトークンが削除される', () => {
    localStorage.setItem('tiara_admin_token', 'test-token')
    // Pinia を再作成して localStorage の値を読み込む
    setActivePinia(createPinia())
    const store = useAuthStore()

    expect(store.isAuthenticated).toBe(true)
    expect(store.token).toBe('test-token')

    store.logout()

    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBeNull()
    expect(localStorage.getItem('tiara_admin_token')).toBeNull()
  })

  it('verify はトークンがない場合 false を返す', async () => {
    const store = useAuthStore()
    const result = await store.verify()
    expect(result).toBe(false)
  })
})
