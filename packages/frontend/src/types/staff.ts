/** BFF から受け取る Staff 型 */
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

/** BFF から受け取る StaffSchedule 型 */
export interface StaffSchedule {
  id: string
  staffId: string
  dayOfWeek: number
  startTime: string
  endTime: string
}

/** BFF から受け取る StaffImage 型 */
export interface StaffImage {
  id: string
  staffId: string
  imageUrl: string
  isMain: boolean
  sortOrder: number
}

/** BFF から受け取る StaffWithSchedules 型 */
export interface StaffWithSchedules {
  staff: Staff
  schedules: StaffSchedule[]
  images: StaffImage[]
}
