import { ref, type Ref } from 'vue'
import { apiFetch } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { ProfileDraft, ScheduleDraft, ScheduleDraftItem } from '@/types/staffPortal'

/**
 * 管理者用レビュー API を呼び出す Composable。
 * 承認待ち下書きの一覧取得・承認・却下・内容修正を担う。
 */
export function useAdminReviewApi() {
  const pendingProfiles: Ref<ProfileDraft[]> = ref([])
  const pendingSchedules: Ref<ScheduleDraft[]> = ref([])
  const approvedSchedules: Ref<ScheduleDraft[]> = ref([])
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
      pendingProfiles.value = await apiFetch<ProfileDraft[]>('/api/v1/admin/reviews/profiles', {
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
      pendingSchedules.value = await apiFetch<ScheduleDraft[]>('/api/v1/admin/reviews/schedules', {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '取得に失敗しました'
    } finally {
      isLoading.value = false
    }
  }

  /** プロフィール下書き単体取得 */
  async function fetchProfileDraft(draftId: string): Promise<ProfileDraft | null> {
    try {
      return await apiFetch<ProfileDraft>(`/api/v1/admin/reviews/profiles/${draftId}`, {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '取得に失敗しました'
      return null
    }
  }

  /** スケジュール下書き単体取得 */
  async function fetchScheduleDraft(draftId: string): Promise<ScheduleDraft | null> {
    try {
      return await apiFetch<ScheduleDraft>(`/api/v1/admin/reviews/schedules/${draftId}`, {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '取得に失敗しました'
      return null
    }
  }

  /** プロフィール下書きの内容を修正（ステータス変更なし） */
  async function updateProfileDraftContent(
    draftId: string,
    data: { name: string; role: string; bio: string; imageUrl: string; imageCropPosition: string }
  ): Promise<ProfileDraft | null> {
    try {
      return await apiFetch<ProfileDraft>(`/api/v1/admin/reviews/profiles/${draftId}/content`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(data),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '修正に失敗しました'
      return null
    }
  }

  /** スケジュール下書きの内容を修正（ステータス変更なし） */
  async function updateScheduleDraftContent(
    draftId: string,
    items: ScheduleDraftItem[]
  ): Promise<ScheduleDraft | null> {
    try {
      return await apiFetch<ScheduleDraft>(`/api/v1/admin/reviews/schedules/${draftId}/content`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ items }),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '修正に失敗しました'
      return null
    }
  }

  /** プロフィール下書きをレビュー（承認 or 却下） */
  async function reviewProfileDraft(
    draftId: string,
    status: 'approved' | 'rejected',
    adminComment: string
  ): Promise<ProfileDraft | null> {
    try {
      const result = await apiFetch<ProfileDraft>(`/api/v1/admin/reviews/profiles/${draftId}`, {
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
      const result = await apiFetch<ScheduleDraft>(`/api/v1/admin/reviews/schedules/${draftId}`, {
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

  /** 承認済み（未反映）スケジュール下書き一覧を取得 */
  async function fetchApprovedSchedules(): Promise<void> {
    try {
      approvedSchedules.value = await apiFetch<ScheduleDraft[]>(
        '/api/v1/admin/reviews/schedules/approved',
        {
          headers: authHeaders(),
        }
      )
    } catch (e) {
      error.value = e instanceof Error ? e.message : '取得に失敗しました'
    }
  }

  /** 承認済みスケジュールを店舗ページに反映 */
  async function publishScheduleDraft(draftId: string): Promise<boolean> {
    try {
      await apiFetch(`/api/v1/admin/reviews/schedules/${draftId}/publish`, {
        method: 'POST',
        headers: authHeaders(),
      })
      approvedSchedules.value = approvedSchedules.value.filter((d) => d.id !== draftId)
      return true
    } catch (e) {
      error.value = e instanceof Error ? e.message : '反映に失敗しました'
      return false
    }
  }

  return {
    pendingProfiles,
    pendingSchedules,
    approvedSchedules,
    isLoading,
    error,
    fetchPendingProfiles,
    fetchPendingSchedules,
    fetchApprovedSchedules,
    fetchProfileDraft,
    fetchScheduleDraft,
    updateProfileDraftContent,
    updateScheduleDraftContent,
    reviewProfileDraft,
    reviewScheduleDraft,
    publishScheduleDraft,
  }
}
