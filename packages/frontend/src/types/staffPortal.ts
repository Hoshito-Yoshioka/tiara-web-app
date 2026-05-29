/** スタッフポータル関連の型定義（BFF → Frontend） */

/** プロフィール下書き */
export interface ProfileDraft {
  id?: string
  staffId: string
  staffName?: string
  name: string
  role: string
  bio: string
  imageUrl: string
  externalScheduleUrl: string
  imageCropPosition: string
  status: string
  adminComment: string
  submittedAt?: string
  reviewedAt?: string
  createdAt?: string
  updatedAt?: string
  /** スタッフに紐づく画像一覧（staff_images） */
  images?: StaffImageForDraft[]
}

/** ドラフトレスポンスに含まれるスタッフ画像 */
export interface StaffImageForDraft {
  id: string
  staffId: string
  imageUrl: string
  isMain: boolean
  sortOrder: number
  cropPosition: string
}

/** スケジュール下書きアイテム */
export interface ScheduleDraftItem {
  id?: string
  dayOfWeek: number
  startTime: string
  endTime: string
}

/** スケジュール下書き */
export interface ScheduleDraft {
  id?: string
  staffId: string
  staffName?: string
  status: string
  adminComment: string
  submittedAt?: string
  reviewedAt?: string
  createdAt?: string
  updatedAt?: string
  items: ScheduleDraftItem[]
}

/** スタッフログインレスポンス */
export interface StaffLoginResponse {
  token: string
  staffId: string
}
