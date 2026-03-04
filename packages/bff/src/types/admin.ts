/** Login リクエスト型 */
export interface LoginRequest {
  username: string
  password: string
}

/** Login レスポンス型 (Backend → BFF) */
export interface LoginResponse {
  token: string
}

/** Shop 更新リクエスト型 */
export interface UpdateShopRequest {
  name: string
  address: string
  openingTime: string
  closingTime: string
}

/** Staff 作成リクエスト型 */
export interface CreateStaffRequest {
  shopId: string
  name: string
  role: string
  bio: string
  imageUrl: string
  sortOrder: number
  schedules: ScheduleInput[]
}

/** Staff 更新リクエスト型 */
export interface UpdateStaffRequest {
  name: string
  role: string
  bio: string
  imageUrl: string
  sortOrder: number
  schedules: ScheduleInput[]
}

/** Schedule 入力型 */
export interface ScheduleInput {
  dayOfWeek: number
  startTime: string
  endTime: string
}
