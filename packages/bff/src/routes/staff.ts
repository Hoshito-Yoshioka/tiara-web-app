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

const staffRoutes = new Hono()

/** Backend の StaffResponse（PascalCase）を Frontend 向け（camelCase）に変換 */
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

/** GET /api/staffs — スタッフ一覧を取得 */
staffRoutes.get('/', async (c) => {
  const res = await fetch(`${BACKEND_URL}/staffs`)

  if (!res.ok) {
    return c.json({ error: 'Failed to fetch staffs from backend' }, 502)
  }

  const data: StaffResponse[] = await res.json()
  const staffs: Staff[] = data.map(toStaff)

  return c.json(staffs)
})

/** GET /api/staffs/:id — スタッフ詳細（スケジュール付き）を取得 */
staffRoutes.get('/:id', async (c) => {
  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/staffs/${id}`)

  if (!res.ok) {
    return c.json({ error: 'Staff not found' }, 404)
  }

  const data: StaffWithSchedulesResponse = await res.json()
  const result: StaffWithSchedules = {
    staff: toStaff(data.Staff),
    schedules: data.Schedules ? data.Schedules.map(toSchedule) : [],
    images: data.Images ? data.Images.map(toImage) : [],
  }

  return c.json(result)
})

export { staffRoutes }
