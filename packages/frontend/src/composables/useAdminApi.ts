import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Shop } from '@/types/shop'
import type { StaffWithSchedules } from '@/types/staff'
import type { UpdateShopInput, CreateStaffInput, UpdateStaffInput } from '@/types/admin'

/**
 * Admin API を呼び出す Composable。
 * 認証が必要な管理者向け API（店舗編集・スタッフ CRUD）のロジックを
 * コンポーネントから分離する。
 */
export function useAdminApi() {
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  /** 認証ヘッダーを生成するヘルパー */
  function authHeaders(): Record<string, string> {
    const authStore = useAuthStore()
    return {
      Authorization: `Bearer ${authStore.token}`,
    }
  }

  /** 店舗情報を更新 */
  async function updateShop(id: string, input: UpdateShopInput): Promise<Shop | null> {
    isLoading.value = true
    error.value = null

    try {
      const result = await apiFetch<Shop>(`/api/admin/shops/${id}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : '店舗情報の更新に失敗しました'
      return null
    } finally {
      isLoading.value = false
    }
  }

  /** スタッフを新規作成 */
  async function createStaff(input: CreateStaffInput): Promise<StaffWithSchedules | null> {
    isLoading.value = true
    error.value = null

    try {
      const result = await apiFetch<StaffWithSchedules>('/api/admin/staffs', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'スタッフの作成に失敗しました'
      return null
    } finally {
      isLoading.value = false
    }
  }

  /** スタッフ情報を更新 */
  async function updateStaff(
    id: string,
    input: UpdateStaffInput
  ): Promise<StaffWithSchedules | null> {
    isLoading.value = true
    error.value = null

    try {
      const result = await apiFetch<StaffWithSchedules>(`/api/admin/staffs/${id}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'スタッフ情報の更新に失敗しました'
      return null
    } finally {
      isLoading.value = false
    }
  }

  /** スタッフを削除 */
  async function deleteStaff(id: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      await apiFetch(`/api/admin/staffs/${id}`, {
        method: 'DELETE',
        headers: authHeaders(),
      })
      return true
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'スタッフの削除に失敗しました'
      return false
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    error,
    updateShop,
    createStaff,
    updateStaff,
    deleteStaff,
  }
}
