import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Shop } from '@/types/shop'
import type { StaffWithSchedules } from '@/types/staff'
import type { MenuCategory, MenuItem } from '@/types/menu'
import type {
  UpdateShopInput,
  CreateStaffInput,
  UpdateStaffInput,
  CreateMenuCategoryInput,
  UpdateMenuCategoryInput,
  CreateMenuItemInput,
  UpdateMenuItemInput,
} from '@/types/admin'

/**
 * Admin API を呼び出す Composable。
 * 認証が必要な管理者向け API（店舗編集・スタッフ CRUD・メニュー CRUD）のロジックを
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
      return await apiFetch<Shop>(`/api/admin/shops/${id}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
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
      return await apiFetch<StaffWithSchedules>('/api/admin/staffs', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
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
      return await apiFetch<StaffWithSchedules>(`/api/admin/staffs/${id}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
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

  // ============================================================
  // Menu Category CRUD
  // ============================================================

  /** メニューカテゴリを新規作成 */
  async function createMenuCategory(input: CreateMenuCategoryInput): Promise<MenuCategory | null> {
    isLoading.value = true
    error.value = null
    try {
      return await apiFetch<MenuCategory>('/api/admin/menu/categories', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'カテゴリの作成に失敗しました'
      return null
    } finally {
      isLoading.value = false
    }
  }

  /** メニューカテゴリを更新 */
  async function updateMenuCategory(
    id: string,
    input: UpdateMenuCategoryInput
  ): Promise<MenuCategory | null> {
    isLoading.value = true
    error.value = null
    try {
      return await apiFetch<MenuCategory>(`/api/admin/menu/categories/${id}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'カテゴリの更新に失敗しました'
      return null
    } finally {
      isLoading.value = false
    }
  }

  /** メニューカテゴリを削除（配下アイテムも DB の CASCADE で削除される） */
  async function deleteMenuCategory(id: string): Promise<boolean> {
    isLoading.value = true
    error.value = null
    try {
      await apiFetch(`/api/admin/menu/categories/${id}`, {
        method: 'DELETE',
        headers: authHeaders(),
      })
      return true
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'カテゴリの削除に失敗しました'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // ============================================================
  // Menu Item CRUD
  // ============================================================

  /** メニューアイテムを新規作成 */
  async function createMenuItem(input: CreateMenuItemInput): Promise<MenuItem | null> {
    isLoading.value = true
    error.value = null
    try {
      return await apiFetch<MenuItem>('/api/admin/menu/items', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'メニューアイテムの作成に失敗しました'
      return null
    } finally {
      isLoading.value = false
    }
  }

  /** メニューアイテムを更新 */
  async function updateMenuItem(id: string, input: UpdateMenuItemInput): Promise<MenuItem | null> {
    isLoading.value = true
    error.value = null
    try {
      return await apiFetch<MenuItem>(`/api/admin/menu/items/${id}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(input),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'メニューアイテムの更新に失敗しました'
      return null
    } finally {
      isLoading.value = false
    }
  }

  /** メニューアイテムを削除 */
  async function deleteMenuItem(id: string): Promise<boolean> {
    isLoading.value = true
    error.value = null
    try {
      await apiFetch(`/api/admin/menu/items/${id}`, {
        method: 'DELETE',
        headers: authHeaders(),
      })
      return true
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'メニューアイテムの削除に失敗しました'
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
    createMenuCategory,
    updateMenuCategory,
    deleteMenuCategory,
    createMenuItem,
    updateMenuItem,
    deleteMenuItem,
  }
}
