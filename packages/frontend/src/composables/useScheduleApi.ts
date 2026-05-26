import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import type { StaffWithSchedules } from '@/types/staff'

/**
 * Schedule API を呼び出す Composable。
 * 全スタッフの出勤スケジュールを取得する。
 */
export function useScheduleApi() {
  const scheduleData: Ref<StaffWithSchedules[]> = ref([])
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  /** 全スタッフの出勤スケジュールを取得 */
  async function fetchSchedules(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      scheduleData.value = await apiFetch<StaffWithSchedules[]>('/api/v1/schedules')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch schedules'
    } finally {
      isLoading.value = false
    }
  }

  return {
    scheduleData,
    isLoading,
    error,
    fetchSchedules,
  }
}
