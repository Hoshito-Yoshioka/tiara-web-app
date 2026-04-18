import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'

/** スタッフアカウント情報 */
export interface StaffAccount {
  id: string
  staffId: string
  username: string
  createdAt: string
  updatedAt: string
}

/**
 * 管理者用スタッフアカウント API を呼び出す Composable。
 * アカウント CRUD を管理者のみが操作する。
 */
export function useAdminAccountApi() {
  const account: Ref<StaffAccount | null> = ref(null)
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  /** 管理者認証ヘッダーを取得 */
  function authHeaders(): Record<string, string> {
    const store = useAuthStore()
    return store.token ? { Authorization: `Bearer ${store.token}` } : {}
  }

  /** 指定スタッフのアカウント情報を取得 */
  async function fetchAccountByStaffId(staffId: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const data = await apiFetch<StaffAccount | null>(
        `/api/admin/staff-accounts/staff/${staffId}`,
        {
          headers: authHeaders(),
        }
      )
      account.value = data
    } catch (e) {
      // 404 の場合はアカウント未作成
      account.value = null
      error.value = null
    } finally {
      isLoading.value = false
    }
  }

  /** スタッフアカウントを作成 */
  async function createAccount(
    staffId: string,
    username: string,
    password: string
  ): Promise<StaffAccount | null> {
    error.value = null

    try {
      const data = await apiFetch<StaffAccount>('/api/admin/staff-accounts', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({ staffId, username, password }),
      })
      account.value = data
      return data
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'アカウント作成に失敗しました'
      return null
    }
  }

  /** スタッフアカウントを更新 */
  async function updateAccount(
    accountId: string,
    username: string,
    password: string
  ): Promise<StaffAccount | null> {
    error.value = null

    try {
      const data = await apiFetch<StaffAccount>(`/api/admin/staff-accounts/${accountId}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ username, password }),
      })
      account.value = data
      return data
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'アカウント更新に失敗しました'
      return null
    }
  }

  /** スタッフアカウントを削除 */
  async function deleteAccount(accountId: string): Promise<boolean> {
    error.value = null

    try {
      await apiFetch(`/api/admin/staff-accounts/${accountId}`, {
        method: 'DELETE',
        headers: authHeaders(),
      })
      account.value = null
      return true
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'アカウント削除に失敗しました'
      return false
    }
  }

  return {
    account,
    isLoading,
    error,
    fetchAccountByStaffId,
    createAccount,
    updateAccount,
    deleteAccount,
  }
}
