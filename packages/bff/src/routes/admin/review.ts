import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type { ProfileDraftResponse, ScheduleDraftResponse } from '../../types/staffPortal'
import { reviewDraftSchema } from '../../schemas'
import {
  ProfileDraftResponseSchema,
  ScheduleDraftResponseSchema,
  ErrorResponseSchema,
  MessageResponseSchema,
} from '../../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminReviewRoutes = new OpenAPIHono<AuthEnv>()

// 全ルートに認証ミドルウェアを適用
adminReviewRoutes.use('/*', authMiddleware)

// ============================================================
// Route definitions — Profile Draft Review
// ============================================================

const listProfileDraftsRoute = createRoute({
  method: 'get',
  path: '/profiles',
  tags: ['管理者 - レビュー'],
  summary: '承認待ちプロフィール下書き一覧',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: z.array(ProfileDraftResponseSchema) } },
      description: '下書き一覧',
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

const getProfileDraftRoute = createRoute({
  method: 'get',
  path: '/profiles/{id}',
  tags: ['管理者 - レビュー'],
  summary: 'プロフィール下書き単体取得',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ProfileDraftResponseSchema } },
      description: '下書き詳細',
    },
    404: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '見つからない',
    },
  },
})

const reviewProfileRoute = createRoute({
  method: 'put',
  path: '/profiles/{id}',
  tags: ['管理者 - レビュー'],
  summary: 'プロフィール下書きレビュー',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
    body: { content: { 'application/json': { schema: reviewDraftSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ProfileDraftResponseSchema } },
      description: 'レビュー成功',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

const editProfileContentRoute = createRoute({
  method: 'put',
  path: '/profiles/{id}/content',
  tags: ['管理者 - レビュー'],
  summary: 'プロフィール下書き内容修正',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ProfileDraftResponseSchema } },
      description: '修正成功',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

// ============================================================
// Route definitions — Schedule Draft Review
// ============================================================

const listScheduleDraftsRoute = createRoute({
  method: 'get',
  path: '/schedules',
  tags: ['管理者 - レビュー'],
  summary: '承認待ちスケジュール下書き一覧',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: z.array(ScheduleDraftResponseSchema) } },
      description: '下書き一覧',
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

const listApprovedSchedulesRoute = createRoute({
  method: 'get',
  path: '/schedules/approved',
  tags: ['管理者 - レビュー'],
  summary: '承認済み（未反映）スケジュール下書き一覧',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: z.array(ScheduleDraftResponseSchema) } },
      description: '承認済み一覧',
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

const getScheduleDraftRoute = createRoute({
  method: 'get',
  path: '/schedules/{id}',
  tags: ['管理者 - レビュー'],
  summary: 'スケジュール下書き単体取得',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ScheduleDraftResponseSchema } },
      description: '下書き詳細',
    },
    404: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '見つからない',
    },
  },
})

const reviewScheduleRoute = createRoute({
  method: 'put',
  path: '/schedules/{id}',
  tags: ['管理者 - レビュー'],
  summary: 'スケジュール下書きレビュー',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
    body: { content: { 'application/json': { schema: reviewDraftSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ScheduleDraftResponseSchema } },
      description: 'レビュー成功',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

const editScheduleContentRoute = createRoute({
  method: 'put',
  path: '/schedules/{id}/content',
  tags: ['管理者 - レビュー'],
  summary: 'スケジュール下書き内容修正',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ScheduleDraftResponseSchema } },
      description: '修正成功',
    },
    409: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '競合エラー',
    },
  },
})

const publishScheduleRoute = createRoute({
  method: 'post',
  path: '/schedules/{id}/publish',
  tags: ['管理者 - レビュー'],
  summary: '承認済みスケジュールを反映',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '下書き ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: MessageResponseSchema } },
      description: '反映成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証エラー',
    },
  },
})

// ============================================================
// Handlers — Profile Draft Review
// ============================================================

adminReviewRoutes.openapi(listProfileDraftsRoute, async (c) => {
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/profiles`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data: ProfileDraftResponse[] = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(getProfileDraftRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/profiles/${id}`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 404)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(reviewProfileRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/profiles/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 409)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(editProfileContentRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')
  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/profiles/${id}/content`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 409)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data, 200)
})

// ============================================================
// Handlers — Schedule Draft Review
// ============================================================

adminReviewRoutes.openapi(listScheduleDraftsRoute, async (c) => {
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/schedules`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data: ScheduleDraftResponse[] = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(listApprovedSchedulesRoute, async (c) => {
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/schedules/approved`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data: ScheduleDraftResponse[] = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(getScheduleDraftRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/schedules/${id}`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 404)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(reviewScheduleRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/schedules/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 409)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(editScheduleContentRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')
  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/schedules/${id}/content`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 409)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data, 200)
})

adminReviewRoutes.openapi(publishScheduleRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/reviews/schedules/${id}/publish`, {
    method: 'POST',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 401)
  }

  const data = (await res.json()) as { message: string }
  return c.json(data, 200)
})

export { adminReviewRoutes }
