import { Hono } from 'hono'
import type {
  Staff,
  StaffResponse,
  StaffImage,
  StaffImageResponse,
  StaffSchedule,
  StaffScheduleResponse,
  StaffWithSchedules,
  StaffWithSchedulesResponse,
} from '../types/staff'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const scheduleRoutes = new Hono()

/** Backend の StaffResponse を Frontend 向けに変換 */
function toStaff(raw: StaffResponse): Staff {
  return {
    id: raw.ID,
    shopId: raw.ShopID,
    name: raw.Name,
    role: raw.Role,
    bio: raw.Bio,
    imageUrl: raw.ImageURL,
    imageCropPosition: raw.ImageCropPosition,
    sortOrder: raw.SortOrder,
  }
}

/** Backend の StaffScheduleResponse を Frontend 向けに変換 */
function toSchedule(raw: StaffScheduleResponse): StaffSchedule {
  return {
    id: raw.ID,
    staffId: raw.StaffID,
    dayOfWeek: raw.DayOfWeek,
    startTime: raw.StartTime,
    endTime: raw.EndTime,
  }
}

/** Backend の StaffImageResponse を Frontend 向けに変換 */
function toImage(raw: StaffImageResponse): StaffImage {
  return {
    id: raw.ID,
    staffId: raw.StaffID,
    imageUrl: raw.ImageURL,
    isMain: raw.IsMain,
    sortOrder: raw.SortOrder,
  }
}

/** GET /api/schedules — 全スタッフの出勤スケジュールを取得 */
scheduleRoutes.get('/', async (c) => {
  const res = await fetch(`${BACKEND_URL}/schedules`)

  if (!res.ok) {
    return c.json({ error: 'Failed to fetch schedules from backend' }, 502)
  }

  const data: StaffWithSchedulesResponse[] = await res.json()
  const result: StaffWithSchedules[] = data.map((item) => ({
    staff: toStaff(item.Staff),
    schedules: item.Schedules ? item.Schedules.map(toSchedule) : [],
    images: item.Images ? item.Images.map(toImage) : [],
  }))

  return c.json(result)
})

export { scheduleRoutes }
