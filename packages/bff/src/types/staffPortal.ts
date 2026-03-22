/** スタッフポータル関連の型定義 */

// --- Staff Auth ---

/** スタッフログインリクエスト */
export interface StaffLoginRequest {
  username: string
  password: string
}

/** スタッフログインレスポンス（Backend がそのまま camelCase で返す） */
export interface StaffLoginResponse {
  token: string
  staffId: string
}

// --- Profile Draft ---

/** プロフィール下書き保存リクエスト */
export interface SaveProfileDraftRequest {
  name: string
  role: string
  bio: string
  imageUrl: string
  imageCropPosition: string
}

/** プロフィール下書きレスポンス */
export interface ProfileDraftResponse {
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

// --- Schedule Draft ---

/** スケジュール下書きアイテム */
export interface ScheduleDraftItem {
  id?: string
  dayOfWeek: number
  startTime: string
  endTime: string
}

/** スケジュール下書き保存リクエスト */
export interface SaveScheduleDraftRequest {
  items: ScheduleDraftItem[]
}

/** スケジュール下書きレスポンス */
export interface ScheduleDraftResponse {
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

// --- Admin Review ---

/** レビューリクエスト */
export interface ReviewDraftRequest {
  status: 'approved' | 'rejected'
  adminComment: string
}
