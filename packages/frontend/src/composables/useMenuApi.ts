import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import type { MenuCategoryWithItems } from '@/types/menu'

/**
 * 公開メニュー API を呼び出す Composable。
 * PriceView でカテゴリ＋アイテム一覧を取得するために使用する。
 */
export function useMenuApi() {
  const menuList: Ref<MenuCategoryWithItems[]> = ref([])
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  async function fetchMenus(): Promise<void> {
    isLoading.value = true
    error.value = null
    try {
      const data = await apiFetch<MenuCategoryWithItems[]>('/api/menus')
      menuList.value = data ?? []
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'メニューの取得に失敗しました'
    } finally {
      isLoading.value = false
    }
  }

  return { menuList, isLoading, error, fetchMenus }
}
