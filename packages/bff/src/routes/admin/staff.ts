import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type {
  Staff,
  StaffResponse,
  StaffImage,
  StaffImageResponse,
  StaffSchedule,
  StaffScheduleResponse,
  StaffWithSchedules,
  StaffWithSchedulesResponse,
} from '../../types/staff'
import {
  createStaffSchema,
  updateStaffSchema,
  setMainImageSchema,
  updateCropPositionSchema,
} from '../../schemas'
import {
  StaffWithSchedulesSchema,
  StaffImageSchema,
  ErrorResponseSchema,
  MessageResponseSchema,
} from '../../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminStaffRoutes = new OpenAPIHono<AuthEnv>()

// 全ルートに認証ミドルウェアを適用
adminStaffRoutes.use('/*', authMiddleware)

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
    cropPosition: raw.CropPosition ?? '50 50',
  }
}

/** Backend の StaffWithSchedulesResponse を Frontend 向けに変換 */
function toStaffWithSchedules(data: StaffWithSchedulesResponse): StaffWithSchedules {
  return {
    staff: toStaff(data.Staff),
    schedules: data.Schedules ? data.Schedules.map(toSchedule) : [],
    images: data.Images ? data.Images.map(toImage) : [],
  }
}

// ============================================================
// Route definitions
// ============================================================

const createStaffRoute = createRoute({
  method: 'post',
  path: '/',
  tags: ['管理者 - スタッフ'],
  summary: 'スタッフ新規作成',
  security: [{ Bearer: [] }],
  request: {
    body: { content: { 'application/json': { schema: createStaffSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffWithSchedulesSchema } },
      description: '作成成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const updateStaffRoute = createRoute({
  method: 'put',
  path: '/{id}',
  tags: ['管理者 - スタッフ'],
  summary: 'スタッフ情報更新',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'スタッフ ID' }) }),
    body: { content: { 'application/json': { schema: updateStaffSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffWithSchedulesSchema } },
      description: '更新成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const deleteStaffRoute = createRoute({
  method: 'delete',
  path: '/{id}',
  tags: ['管理者 - スタッフ'],
  summary: 'スタッフ削除',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'スタッフ ID' }) }),
  },
  responses: {
    204: { description: '削除成功' },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const uploadStaffImageRoute = createRoute({
  method: 'post',
  path: '/{id}/images',
  tags: ['管理者 - スタッフ画像'],
  summary: 'スタッフ画像アップロード',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'スタッフ ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffImageSchema } },
      description: 'アップロード成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const deleteStaffImageRoute = createRoute({
  method: 'delete',
  path: '/{id}/images/{imageId}',
  tags: ['管理者 - スタッフ画像'],
  summary: 'スタッフ画像削除',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({
      id: z.string().openapi({ description: 'スタッフ ID' }),
      imageId: z.string().openapi({ description: '画像 ID' }),
    }),
  },
  responses: {
    204: { description: '削除成功' },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const setStaffMainImageRoute = createRoute({
  method: 'put',
  path: '/{id}/images/main',
  tags: ['管理者 - スタッフ画像'],
  summary: 'メイン画像設定',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'スタッフ ID' }) }),
    body: { content: { 'application/json': { schema: setMainImageSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: MessageResponseSchema } },
      description: '設定成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const updateStaffCropRoute = createRoute({
  method: 'put',
  path: '/{id}/images/{imageId}/crop',
  tags: ['管理者 - スタッフ画像'],
  summary: '画像クロップ位置更新',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({
      id: z.string().openapi({ description: 'スタッフ ID' }),
      imageId: z.string().openapi({ description: '画像 ID' }),
    }),
    body: {
      content: { 'application/json': { schema: updateCropPositionSchema } },
      required: true,
    },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffImageSchema } },
      description: '更新成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

// ============================================================
// Handlers
// ============================================================

adminStaffRoutes.openapi(createStaffRoute, async (c) => {
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staffs`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 500)
  }

  const data: StaffWithSchedulesResponse = await res.json()
  return c.json(toStaffWithSchedules(data), 200)
})

adminStaffRoutes.openapi(updateStaffRoute, async (c) => {
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staffs/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 500)
  }

  const data: StaffWithSchedulesResponse = await res.json()
  return c.json(toStaffWithSchedules(data), 200)
})

adminStaffRoutes.openapi(deleteStaffRoute, async (c) => {
  const { id } = c.req.valid('param')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staffs/${id}`, {
    method: 'DELETE',
    headers: {
      Authorization: authHeader,
    },
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete staff' }))
    return c.json(error as { error: string }, 500)
  }

  return c.body(null, 204)
})

// ============================================================
// Staff Image Management
// ============================================================

adminStaffRoutes.openapi(uploadStaffImageRoute, async (c) => {
  const { id } = c.req.valid('param')
  const authHeader = c.get('authHeader') as string

  // multipart body をそのまま Backend へ転送
  const body = await c.req.raw.arrayBuffer()
  const contentType = c.req.header('content-type') || ''

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staffs/${id}/images`, {
    method: 'POST',
    headers: {
      'Content-Type': contentType,
      Authorization: authHeader,
    },
    body,
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to upload image' }))
    return c.json(error as { error: string }, 500)
  }

  const data: StaffImageResponse = await res.json()
  return c.json(toImage(data), 200)
})

adminStaffRoutes.openapi(deleteStaffImageRoute, async (c) => {
  const { id, imageId } = c.req.valid('param')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staffs/${id}/images/${imageId}`, {
    method: 'DELETE',
    headers: {
      Authorization: authHeader,
    },
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete image' }))
    return c.json(error as { error: string }, 500)
  }

  return c.body(null, 204)
})

adminStaffRoutes.openapi(setStaffMainImageRoute, async (c) => {
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staffs/${id}/images/main`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to set main image' }))
    return c.json(error as { error: string }, 500)
  }

  return c.json({ message: 'ok' }, 200)
})

adminStaffRoutes.openapi(updateStaffCropRoute, async (c) => {
  const { id, imageId } = c.req.valid('param')
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staffs/${id}/images/${imageId}/crop`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to update crop position' }))
    return c.json(error as { error: string }, 500)
  }

  const data: StaffImageResponse = await res.json()
  return c.json(toImage(data), 200)
})

export { adminStaffRoutes }
