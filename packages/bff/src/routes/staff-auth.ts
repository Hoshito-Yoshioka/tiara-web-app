import { OpenAPIHono, createRoute } from '@hono/zod-openapi'
import { loginSchema } from '../schemas'
import type { StaffLoginResponse } from '../types/staffPortal'
import {
  StaffTokenResponseSchema,
  ErrorResponseSchema,
  StaffVerifyResponseSchema,
} from '../schemas/responses'

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

export { staffAuthRoutes }
