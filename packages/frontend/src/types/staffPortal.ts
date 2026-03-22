/** スタッフポータル関連の型定義（BFF → Frontend） */

/** プロフィール下書き */
export interface ProfileDraft {
  id?: string
  staffId: string
  name: string
  role: string
  bio: string
  imageUrl: string
  imageCropPosition: string
  status: string
  adminComment: string
  submittedAt?: string
  reviewedAt?: string
  createdAt?: string
  updatedAt?: string
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
