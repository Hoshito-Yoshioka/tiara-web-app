import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { ProfileDraft, ScheduleDraft } from '@/types/staffPortal'

/**
 * 管理者用レビュー API を呼び出す Composable。
 * 承認待ち下書きの一覧取得・承認・却下を担う。
 */
export function useAdminReviewApi() {
  const pendingProfiles: Ref<ProfileDraft[]> = ref([])
  const pendingSchedules: Ref<ScheduleDraft[]> = ref([])
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)

  /** 管理者認証ヘッダーを取得 */
  function authHeaders(): Record<string, string> {
    const store = useAuthStore()
    return store.token ? { Authorization: `Bearer ${store.token}` } : {}
  }

  /** 承認待ちプロフィール下書き一覧を取得 */
  async function fetchPendingProfiles(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      pendingProfiles.value = await apiFetch<ProfileDraft[]>('/api/admin/reviews/profiles', {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '取得に失敗しました'
    } finally {
      isLoading.value = false
    }
  }

  /** 承認待ちスケジュール下書き一覧を取得 */
  async function fetchPendingSchedules(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      pendingSchedules.value = await apiFetch<ScheduleDraft[]>('/api/admin/reviews/schedules', {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '取得に失敗しました'
    } finally {
      isLoading.value = false
    }
  }

  /** プロフィール下書きをレビュー（承認 or 却下） */
  async function reviewProfileDraft(
    draftId: string,
    status: 'approved' | 'rejected',
    adminComment: string
  ): Promise<ProfileDraft | null> {
    try {
      const result = await apiFetch<ProfileDraft>(`/api/admin/reviews/profiles/${draftId}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ status, adminComment }),
      })
      // 一覧から除去
      pendingProfiles.value = pendingProfiles.value.filter((d) => d.id !== draftId)
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'レビューに失敗しました'
      return null
    }
  }

  /** スケジュール下書きをレビュー（承認 or 却下） */
  async function reviewScheduleDraft(
    draftId: string,
    status: 'approved' | 'rejected',
    adminComment: string
  ): Promise<ScheduleDraft | null> {
    try {
      const result = await apiFetch<ScheduleDraft>(`/api/admin/reviews/schedules/${draftId}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ status, adminComment }),
      })
      pendingSchedules.value = pendingSchedules.value.filter((d) => d.id !== draftId)
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'レビューに失敗しました'
      return null
    }
  }

  return {
    pendingProfiles,
    pendingSchedules,
    isLoading,
    error,
    fetchPendingProfiles,
    fetchPendingSchedules,
    reviewProfileDraft,
    reviewScheduleDraft,
  }
}
