import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
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
import { StaffWithSchedulesSchema, ErrorResponseSchema } from '../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const scheduleRoutes = new OpenAPIHono()

/** Backend の StaffResponse を Frontend 向けに変換 */
function toStaff(raw: StaffResponse): Staff {
  return {
    id: raw.ID,
    shopId: raw.ShopID,
    name: raw.Name,
    role: raw.Role,
    bio: raw.Bio,
    imageUrl: raw.ImageURL,
    externalScheduleUrl: raw.ExternalScheduleURL,
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

// --- Route definition ---

const listSchedulesRoute = createRoute({
  method: 'get',
  path: '/',
  tags: ['スケジュール'],
  summary: '全スタッフの出勤スケジュール取得',
  responses: {
    200: {
      content: { 'application/json': { schema: z.array(StaffWithSchedulesSchema) } },
      description: '全スタッフのスケジュール一覧',
    },
    502: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'Backend エラー',
    },
  },
})

// --- Handler ---

scheduleRoutes.openapi(listSchedulesRoute, async (c) => {
  const res = await fetch(`${BACKEND_URL}/api/v1/schedules`)

  if (!res.ok) {
    return c.json({ error: 'Failed to fetch schedules from backend' }, 502)
  }

  const data: StaffWithSchedulesResponse[] = await res.json()
  const result: StaffWithSchedules[] = data.map((item) => ({
    staff: toStaff(item.Staff),
    schedules: item.Schedules ? item.Schedules.map(toSchedule) : [],
    images: item.Images ? item.Images.map(toImage) : [],
  }))

  return c.json(result, 200)
})

export { scheduleRoutes }
