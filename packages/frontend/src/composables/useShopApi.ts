import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import type { Shop } from '@/types/shop'

/**
 * Shop API を呼び出す Composable。
 * ビジネスロジック（データ取得・ローディング・エラー管理）を
 * コンポーネントから分離し、再利用可能にする。
 */
export function useShopApi() {
  const shops: Ref<Shop[]> = ref([])
  const shop: Ref<Shop | null> = ref(null)
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  /** 店舗一覧を取得 */
  async function fetchShops(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      shops.value = await apiFetch<Shop[]>('/api/shops')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch shops'
    } finally {
      isLoading.value = false
    }
  }

  /** 店舗詳細を取得 */
  async function fetchShopById(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      shop.value = await apiFetch<Shop>(`/api/shops/${id}`)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch shop'
    } finally {
      isLoading.value = false
    }
  }

  return {
    shops,
    shop,
    isLoading,
    error,
    fetchShops,
    fetchShopById,
  }
}
