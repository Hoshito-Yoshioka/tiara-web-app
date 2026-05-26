import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import type { Staff, StaffWithSchedules, PaginatedStaffs, Pagination } from '@/types/staff'

/**
 * Staff API を呼び出す Composable。
 * スタッフ一覧・詳細取得のロジックをコンポーネントから分離する。
 */
export function useStaffApi() {
  const staffList: Ref<Staff[]> = ref([])
  const staffDetail: Ref<StaffWithSchedules | null> = ref(null)
  const pagination: Ref<Pagination | null> = ref(null)
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  /** スタッフ一覧を取得（全件） */
  async function fetchStaffs(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      staffList.value = await apiFetch<Staff[]>('/api/v1/staffs')
      pagination.value = null
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch staffs'
    } finally {
      isLoading.value = false
    }
  }

  /** スタッフ一覧をページネーション付きで取得 */
  async function fetchStaffsPaginated(page: number): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const result = await apiFetch<PaginatedStaffs>(`/api/v1/staffs?page=${page}`)
      staffList.value = result.data
      pagination.value = result.pagination
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch staffs'
    } finally {
      isLoading.value = false
    }
  }

  /** スタッフ詳細（スケジュール付き）を取得 */
  async function fetchStaffById(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      staffDetail.value = await apiFetch<StaffWithSchedules>(`/api/v1/staffs/${id}`)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch staff detail'
    } finally {
      isLoading.value = false
    }
  }

  return {
    staffList,
    staffDetail,
    pagination,
    isLoading,
    error,
    fetchStaffs,
    fetchStaffsPaginated,
    fetchStaffById,
  }
}
