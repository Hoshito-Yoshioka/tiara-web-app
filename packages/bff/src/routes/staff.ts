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
  PaginatedStaffs,
  PaginatedStaffsResponse,
} from '../types/staff'
import {
  StaffSchema,
  StaffWithSchedulesSchema,
  PaginatedStaffsSchema,
  ErrorResponseSchema,
} from '../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const staffRoutes = new OpenAPIHono()

/** Backend の StaffResponse（PascalCase）を Frontend 向け（camelCase）に変換 */
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

// --- Route definitions ---

const listStaffsRoute = createRoute({
  method: 'get',
  path: '/',
  tags: ['スタッフ'],
  summary: 'スタッフ一覧取得',
  description: 'page クエリパラメータ指定時はページネーション付き',
  request: {
    query: z.object({
      page: z.string().optional().openapi({ description: 'ページ番号' }),
    }),
  },
  responses: {
    200: {
      content: {
        'application/json': {
          schema: z.union([z.array(StaffSchema), PaginatedStaffsSchema]),
        },
      },
      description: 'スタッフ一覧（page 指定時はページネーション付き）',
    },
    400: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'リクエストエラー',
    },
    502: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'Backend エラー',
    },
  },
})

const getStaffRoute = createRoute({
  method: 'get',
  path: '/{id}',
  tags: ['スタッフ'],
  summary: 'スタッフ詳細取得（スケジュール・画像付き）',
  request: {
    params: z.object({ id: z.string().openapi({ description: 'スタッフ ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffWithSchedulesSchema } },
      description: 'スタッフ詳細',
    },
    404: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'スタッフが見つからない',
    },
  },
})

// --- Handlers ---

staffRoutes.openapi(listStaffsRoute, async (c) => {
  const { page } = c.req.valid('query')
  const url = page ? `${BACKEND_URL}/api/v1/staffs?page=${page}` : `${BACKEND_URL}/api/v1/staffs`

  const res = await fetch(url)

  if (!res.ok) {
    const error =
      typeof res.json === 'function'
        ? await res.json().catch(() => ({ error: 'Failed to fetch staffs from backend' }))
        : { error: 'Failed to fetch staffs from backend' }
    if (res.status >= 500) {
      return c.json(error as { error: string }, 502)
    }
    return c.json(error as { error: string }, 400)
  }

  if (page) {
    const data: PaginatedStaffsResponse = await res.json()
    const result: PaginatedStaffs = {
      data: data.data.map(toStaff),
      pagination: data.pagination,
    }
    return c.json(result, 200)
  }

  const data: StaffResponse[] = await res.json()
  return c.json(data.map(toStaff), 200)
})

staffRoutes.openapi(getStaffRoute, async (c) => {
  const { id } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/staffs/${id}`)

  if (!res.ok) {
    return c.json({ error: 'Staff not found' }, 404)
  }

  const data: StaffWithSchedulesResponse = await res.json()
  const result: StaffWithSchedules = {
    staff: toStaff(data.Staff),
    schedules: data.Schedules ? data.Schedules.map(toSchedule) : [],
    images: data.Images ? data.Images.map(toImage) : [],
  }

  return c.json(result, 200)
})

export { staffRoutes }
