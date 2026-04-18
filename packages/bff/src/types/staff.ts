/** Backend の Staff レスポンス型（PascalCase） */
export interface StaffResponse {
  ID: string
  ShopID: string
  Name: string
  Role: string
  Bio: string
  ImageURL: string
  ImageCropPosition: string
  SortOrder: number
  CreatedAt: string
  UpdatedAt: string
}

/** Backend の StaffSchedule レスポンス型 */
export interface StaffScheduleResponse {
  ID: string
  StaffID: string
  DayOfWeek: number
  StartTime: string
  EndTime: string
  CreatedAt: string
  UpdatedAt: string
}

/** Backend の StaffImage レスポンス型（PascalCase） */
export interface StaffImageResponse {
  ID: string
  StaffID: string
  ImageURL: string
  IsMain: boolean
  SortOrder: number
  CropPosition: string
  CreatedAt: string
  UpdatedAt: string
}

/** Backend の StaffWithSchedules レスポンス型 */
export interface StaffWithSchedulesResponse {
  Staff: StaffResponse
  Schedules: StaffScheduleResponse[]
  Images: StaffImageResponse[] | null
}

/** BFF → Frontend 向け Staff 型（camelCase） */
export interface Staff {
  id: string
  shopId: string
  name: string
  role: string
  bio: string
  imageUrl: string
  imageCropPosition: string
  sortOrder: number
}

/** BFF → Frontend 向け StaffSchedule 型 */
export interface StaffSchedule {
  id: string
  staffId: string
  dayOfWeek: number
  startTime: string
  endTime: string
}

/** BFF → Frontend 向け StaffImage 型 */
export interface StaffImage {
  id: string
  staffId: string
  imageUrl: string
  isMain: boolean
  sortOrder: number
  cropPosition: string
}

/** BFF → Frontend 向け StaffWithSchedules 型 */
export interface StaffWithSchedules {
  staff: Staff
  schedules: StaffSchedule[]
  images: StaffImage[]
}
