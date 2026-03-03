import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import type { Staff, StaffWithSchedules } from '@/types/staff'

/**
 * Staff API を呼び出す Composable。
 * スタッフ一覧・詳細取得のロジックをコンポーネントから分離する。
 */
export function useStaffApi() {
  const staffList: Ref<Staff[]> = ref([])
  const staffDetail: Ref<StaffWithSchedules | null> = ref(null)
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  /** スタッフ一覧を取得 */
  async function fetchStaffs(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      staffList.value = await apiFetch<Staff[]>('/api/staffs')
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
      staffDetail.value = await apiFetch<StaffWithSchedules>(`/api/staffs/${id}`)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch staff detail'
    } finally {
      isLoading.value = false
    }
  }

  return {
    staffList,
    staffDetail,
    isLoading,
    error,
    fetchStaffs,
    fetchStaffById,
  }
}
