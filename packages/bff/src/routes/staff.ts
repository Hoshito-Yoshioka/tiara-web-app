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
  PaginatedStaffs,
  PaginatedStaffsResponse,
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
    cropPosition: raw.CropPosition ?? '50 50',
  }
}

/** GET /api/staffs — スタッフ一覧を取得（page パラメータ指定時はページネーション付き） */
staffRoutes.get('/', async (c) => {
  const page = c.req.query('page')
  const url = page ? `${BACKEND_URL}/staffs?page=${page}` : `${BACKEND_URL}/staffs`

  const res = await fetch(url)

  if (!res.ok) {
    const error =
      typeof res.json === 'function'
        ? await res.json().catch(() => ({ error: 'Failed to fetch staffs from backend' }))
        : { error: 'Failed to fetch staffs from backend' }
    return c.json(error, res.status >= 500 ? 502 : (res.status as 400))
  }

  if (page) {
    const data: PaginatedStaffsResponse = await res.json()
    const result: PaginatedStaffs = {
      data: data.data.map(toStaff),
      pagination: data.pagination,
    }
    return c.json(result)
  }

  const data: StaffResponse[] = await res.json()
  return c.json(data.map(toStaff))
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
