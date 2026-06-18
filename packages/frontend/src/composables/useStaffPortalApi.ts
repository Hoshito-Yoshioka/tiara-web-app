import { ref, type Ref } from 'vue'
import { apiFetch, apiUpload } from '@/lib/api'
import { useStaffAuthStore } from '@/stores/staffAuth'
import type { ProfileDraft, ScheduleDraft, ScheduleDraftItem } from '@/types/staffPortal'
import type { StaffImage } from '@/types/staff'

/**
 * スタッフポータル API を呼び出す Composable。
 * プロフィール・スケジュール下書きの CRUD + 承認申請を担う。
 */
export function useStaffPortalApi() {
  const profileDraft: Ref<ProfileDraft | null> = ref(null)
  const scheduleDraft: Ref<ScheduleDraft | null> = ref(null)
  const isLoading = ref(false)
  const error: Ref<string | null> = ref(null)
  const saveError: Ref<string | null> = ref(null)

  /** 認証ヘッダーを取得 */
  function authHeaders(): Record<string, string> {
    const store = useStaffAuthStore()
    return store.token ? { Authorization: `Bearer ${store.token}` } : {}
  }

  // --- Profile Draft ---

  /** 自分のプロフィール下書きを取得 */
  async function fetchMyProfileDraft(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      profileDraft.value = await apiFetch<ProfileDraft>('/api/v1/portal/profile', {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'プロフィール下書きの取得に失敗しました'
    } finally {
      isLoading.value = false
    }
  }

  /** プロフィール下書きを保存 */
  async function saveProfileDraft(data: {
    name: string
    role: string
    bio: string
    imageUrl: string
    externalScheduleUrl: string
    imageCropPosition: string
  }): Promise<ProfileDraft | null> {
    saveError.value = null

    try {
      const result = await apiFetch<ProfileDraft>('/api/v1/portal/profile', {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(data),
      })
      profileDraft.value = result
      return result
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : '保存に失敗しました'
      return null
    }
  }

  /** プロフィール下書きを承認申請 */
  async function submitProfileDraft(draftId: string): Promise<ProfileDraft | null> {
    saveError.value = null

    try {
      const result = await apiFetch<ProfileDraft>(`/api/v1/portal/profile/${draftId}/submit`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({}),
      })
      profileDraft.value = result
      return result
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : '承認申請に失敗しました'
      return null
    }
  }

  // --- Schedule Draft ---

  /** 自分のスケジュール下書きを取得 */
  async function fetchMyScheduleDraft(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      scheduleDraft.value = await apiFetch<ScheduleDraft>('/api/v1/portal/schedule', {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'スケジュール下書きの取得に失敗しました'
    } finally {
      isLoading.value = false
    }
  }

  /** スケジュール下書きを保存 */
  async function saveScheduleDraft(items: ScheduleDraftItem[]): Promise<ScheduleDraft | null> {
    saveError.value = null

    try {
      const result = await apiFetch<ScheduleDraft>('/api/v1/portal/schedule', {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ items }),
      })
      scheduleDraft.value = result
      return result
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : '保存に失敗しました'
      return null
    }
  }

  /** スケジュール下書きを承認申請 */
  async function submitScheduleDraft(draftId: string): Promise<ScheduleDraft | null> {
    saveError.value = null

    try {
      const result = await apiFetch<ScheduleDraft>(`/api/v1/portal/schedule/${draftId}/submit`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({}),
      })
      scheduleDraft.value = result
      return result
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : '承認申請に失敗しました'
      return null
    }
  }

  // --- Image Management ---

  /** 自分の画像一覧を取得 */
  async function fetchMyImages(): Promise<StaffImage[]> {
    try {
      return await apiFetch<StaffImage[]>('/api/v1/portal/images', {
        headers: authHeaders(),
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '画像一覧の取得に失敗しました'
      return []
    }
  }

  /** 画像をアップロード */
  async function uploadMyImage(file: File): Promise<StaffImage | null> {
    saveError.value = null
    try {
      const formData = new FormData()
      formData.append('image', file)
      return await apiUpload<StaffImage>('/api/v1/portal/images', formData, {
        headers: authHeaders(),
      })
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : '画像のアップロードに失敗しました'
      return null
    }
  }

  /** 画像を削除 */
  async function deleteMyImage(imageId: string): Promise<boolean> {
    saveError.value = null
    try {
      await apiFetch(`/api/v1/portal/images/${imageId}`, {
        method: 'DELETE',
        headers: authHeaders(),
      })
      return true
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : '画像の削除に失敗しました'
      return false
    }
  }

  /** メイン画像を設定 */
  async function setMyMainImage(imageId: string): Promise<boolean> {
    saveError.value = null
    try {
      await apiFetch('/api/v1/portal/images/main', {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ imageId }),
      })
      return true
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : 'メイン画像の設定に失敗しました'
      return false
    }
  }

  /** 画像のクロップ位置を更新 */
  async function updateMyImageCropPosition(
    imageId: string,
    cropPosition: string
  ): Promise<StaffImage | null> {
    saveError.value = null
    try {
      return await apiFetch<StaffImage>(`/api/v1/portal/images/${imageId}/crop`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ cropPosition }),
      })
    } catch (e) {
      saveError.value = e instanceof Error ? e.message : 'クロップ位置の更新に失敗しました'
      return null
    }
  }

  return {
    profileDraft,
    scheduleDraft,
    isLoading,
    error,
    saveError,
    fetchMyProfileDraft,
    saveProfileDraft,
    submitProfileDraft,
    fetchMyScheduleDraft,
    saveScheduleDraft,
    submitScheduleDraft,
    fetchMyImages,
    uploadMyImage,
    deleteMyImage,
    setMyMainImage,
    updateMyImageCropPosition,
  }
}
