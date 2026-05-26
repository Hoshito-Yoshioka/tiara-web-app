import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import type { ProfileDraftResponse, ScheduleDraftResponse } from '../types/staffPortal'
import type { StaffImageResponse, StaffImage } from '../types/staff'
import {
  saveProfileDraftSchema,
  saveScheduleDraftSchema,
  submitDraftSchema,
  setMainImageSchema,
  updateCropPositionSchema,
} from '../schemas'
import {
  ProfileDraftResponseSchema,
  ScheduleDraftResponseSchema,
  StaffImageSchema,
  ErrorResponseSchema,
  MessageResponseSchema,
} from '../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const portalRoutes = new OpenAPIHono()

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

// ============================================================
// Route definitions
// ============================================================

// --- Profile Draft ---

const getProfileRoute = createRoute({
  method: 'get',
  path: '/profile',
  tags: ['スタッフポータル'],
  summary: 'プロフィール下書き取得',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: ProfileDraftResponseSchema } },
      description: 'プロフィール下書き',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const saveProfileRoute = createRoute({
  method: 'put',
  path: '/profile',
  tags: ['スタッフポータル'],
  summary: 'プロフィール下書き保存',
  security: [{ Bearer: [] }],
  request: {
    body: { content: { 'application/json': { schema: saveProfileDraftSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ProfileDraftResponseSchema } },
      description: '保存成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

const submitProfileRoute = createRoute({
  method: 'post',
  path: '/profile/{id}/submit',
  tags: ['スタッフポータル'],
  summary: 'プロフィール下書き承認申請',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
    body: { content: { 'application/json': { schema: submitDraftSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ProfileDraftResponseSchema } },
      description: '申請成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

// --- Schedule Draft ---

const getScheduleRoute = createRoute({
  method: 'get',
  path: '/schedule',
  tags: ['スタッフポータル'],
  summary: 'スケジュール下書き取得',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: ScheduleDraftResponseSchema } },
      description: 'スケジュール下書き',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const saveScheduleRoute = createRoute({
  method: 'put',
  path: '/schedule',
  tags: ['スタッフポータル'],
  summary: 'スケジュール下書き保存',
  security: [{ Bearer: [] }],
  request: {
    body: {
      content: { 'application/json': { schema: saveScheduleDraftSchema } },
      required: true,
    },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ScheduleDraftResponseSchema } },
      description: '保存成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

const submitScheduleRoute = createRoute({
  method: 'post',
  path: '/schedule/{id}/submit',
  tags: ['スタッフポータル'],
  summary: 'スケジュール下書き承認申請',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
    body: { content: { 'application/json': { schema: submitDraftSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ScheduleDraftResponseSchema } },
      description: '申請成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

// --- Image Management ---

const listImagesRoute = createRoute({
  method: 'get',
  path: '/images',
  tags: ['スタッフポータル'],
  summary: '画像一覧取得',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: z.array(StaffImageSchema) } },
      description: '画像一覧',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const uploadImageRoute = createRoute({
  method: 'post',
  path: '/images',
  tags: ['スタッフポータル'],
  summary: '画像アップロード',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: StaffImageSchema } },
      description: 'アップロード成功',
    },
    400: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'リクエストエラー',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
  },
})

const deleteImageRoute = createRoute({
  method: 'delete',
  path: '/images/{imageId}',
  tags: ['スタッフポータル'],
  summary: '画像削除',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ imageId: z.string().openapi({ description: '画像 ID' }) }),
  },
  responses: {
    204: { description: '削除成功' },
    400: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'リクエストエラー',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
  },
})

const setMainImageRoute = createRoute({
  method: 'put',
  path: '/images/main',
  tags: ['スタッフポータル'],
  summary: 'メイン画像設定',
  security: [{ Bearer: [] }],
  request: {
    body: { content: { 'application/json': { schema: setMainImageSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: MessageResponseSchema } },
      description: '設定成功',
    },
    400: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'リクエストエラー',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
  },
})

const updateCropRoute = createRoute({
  method: 'put',
  path: '/images/{imageId}/crop',
  tags: ['スタッフポータル'],
  summary: '画像クロップ位置更新',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ imageId: z.string().openapi({ description: '画像 ID' }) }),
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
    400: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'リクエストエラー',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
  },
})

// ============================================================
// Handlers
// ============================================================

// --- Profile Draft ---

portalRoutes.openapi(getProfileRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/profile`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data, 200)
})

portalRoutes.openapi(saveProfileRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/profile`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status === 409 ? 409 : 401)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data, 200)
})

portalRoutes.openapi(submitProfileRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const { id } = c.req.valid('param')
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/profile/${id}/submit`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status === 409 ? 409 : 401)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data, 200)
})

// --- Schedule Draft ---

portalRoutes.openapi(getScheduleRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/schedule`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data, 200)
})

portalRoutes.openapi(saveScheduleRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/schedule`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status === 409 ? 409 : 401)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data, 200)
})

portalRoutes.openapi(submitScheduleRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const { id } = c.req.valid('param')
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/schedule/${id}/submit`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status === 409 ? 409 : 401)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data, 200)
})

// --- Image Management ---

portalRoutes.openapi(listImagesRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/images`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data: StaffImageResponse[] = await res.json()
  return c.json(data.map(toImage), 200)
})

portalRoutes.openapi(uploadImageRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = await c.req.raw.arrayBuffer()
  const contentType = c.req.header('content-type') || ''

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/images`, {
    method: 'POST',
    headers: {
      'Content-Type': contentType,
      Authorization: authHeader,
    },
    body,
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to upload image' }))
    return c.json(error as { error: string }, 400)
  }

  const data: StaffImageResponse = await res.json()
  return c.json(toImage(data), 200)
})

portalRoutes.openapi(deleteImageRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const { imageId } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/images/${imageId}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete image' }))
    return c.json(error as { error: string }, 400)
  }

  return c.body(null, 204)
})

portalRoutes.openapi(setMainImageRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/images/main`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to set main image' }))
    return c.json(error as { error: string }, 400)
  }

  return c.json({ message: 'ok' }, 200)
})

portalRoutes.openapi(updateCropRoute, async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const { imageId } = c.req.valid('param')
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/images/${imageId}/crop`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to update crop position' }))
    return c.json(error as { error: string }, 400)
  }

  const data: StaffImageResponse = await res.json()
  return c.json(toImage(data), 200)
})

export { portalRoutes }
