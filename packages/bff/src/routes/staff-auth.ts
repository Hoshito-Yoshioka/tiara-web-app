import { OpenAPIHono, createRoute } from '@hono/zod-openapi'
import { loginSchema } from '../schemas'
import type { StaffLoginResponse } from '../types/staffPortal'
import {
  StaffTokenResponseSchema,
  ErrorResponseSchema,
  StaffVerifyResponseSchema,
  RefreshTokenResponseSchema,
} from '../schemas/responses'
import { z } from '@hono/zod-openapi'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const staffAuthRoutes = new OpenAPIHono()

// --- Route definitions ---

const staffLoginRoute = createRoute({
  method: 'post',
  path: '/login',
  tags: ['スタッフ認証'],
  summary: 'スタッフログイン',
  request: {
    body: { content: { 'application/json': { schema: loginSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffTokenResponseSchema } },
      description: 'ログイン成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証失敗',
    },
  },
})

const staffVerifyRoute = createRoute({
  method: 'get',
  path: '/verify',
  tags: ['スタッフ認証'],
  summary: 'スタッフトークン検証',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: StaffVerifyResponseSchema } },
      description: 'トークン有効',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'トークン無効',
    },
  },
})

// --- Handlers ---

staffAuthRoutes.openapi(staffLoginRoute, async (c) => {
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/staff-auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 401)
  }

  const data: StaffLoginResponse = await res.json()
  return c.json(data, 200)
})

staffAuthRoutes.openapi(staffVerifyRoute, async (c) => {
  const authHeader = c.req.header('Authorization')

  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/api/v1/portal/auth/verify`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    return c.json({ error: 'Invalid token' }, 401)
  }

  const data = (await res.json()) as { status: string; staffId: string }
  return c.json(data, 200)
})

const staffRefreshRoute = createRoute({
  method: 'post',
  path: '/refresh',
  tags: ['スタッフ認証'],
  summary: 'アクセストークン再発行',
  request: {
    body: {
      content: {
        'application/json': {
          schema: z.object({ refreshToken: z.string() }),
        },
      },
      required: true,
    },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: RefreshTokenResponseSchema } },
      description: '再発行成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'リフレッシュトークン無効',
    },
  },
})

const staffLogoutRoute = createRoute({
  method: 'post',
  path: '/logout',
  tags: ['スタッフ認証'],
  summary: 'ログアウト（リフレッシュトークン失効）',
  request: {
    body: {
      content: {
        'application/json': {
          schema: z.object({ refreshToken: z.string() }),
        },
      },
      required: true,
    },
  },
  responses: {
    204: { description: 'ログアウト成功' },
  },
})

staffAuthRoutes.openapi(staffRefreshRoute, async (c) => {
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/staff-auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    return c.json({ error: 'Invalid or expired refresh token' }, 401)
  }

  const data = (await res.json()) as { token: string; refreshToken: string }
  return c.json(data, 200)
})

staffAuthRoutes.openapi(staffLogoutRoute, async (c) => {
  const body = c.req.valid('json')

  await fetch(`${BACKEND_URL}/api/v1/staff-auth/logout`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  return new Response(null, { status: 204 })
})

export { staffAuthRoutes }
