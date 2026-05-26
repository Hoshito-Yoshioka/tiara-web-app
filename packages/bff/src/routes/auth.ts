import { OpenAPIHono, createRoute } from '@hono/zod-openapi'
import { loginSchema } from '../schemas'
import type { LoginResponse } from '../types/admin'
import {
  TokenResponseSchema,
  ErrorResponseSchema,
  AdminVerifyResponseSchema,
} from '../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const authRoutes = new OpenAPIHono()

// --- Route definitions ---

const loginRoute = createRoute({
  method: 'post',
  path: '/login',
  tags: ['認証'],
  summary: '管理者ログイン',
  request: {
    body: { content: { 'application/json': { schema: loginSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: TokenResponseSchema } },
      description: 'ログイン成功',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '認証失敗',
    },
  },
})

const verifyRoute = createRoute({
  method: 'get',
  path: '/verify',
  tags: ['認証'],
  summary: '管理者トークン検証',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: AdminVerifyResponseSchema } },
      description: 'トークン有効',
    },
    401: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'トークン無効',
    },
  },
})

// --- Handlers ---

authRoutes.openapi(loginRoute, async (c) => {
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 401)
  }

  const data: LoginResponse = await res.json()
  return c.json(data, 200)
})

authRoutes.openapi(verifyRoute, async (c) => {
  const authHeader = c.req.header('Authorization')

  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/auth/verify`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    return c.json({ error: 'Invalid token' }, 401)
  }

  const data = (await res.json()) as { status: string }
  return c.json(data, 200)
})

export { authRoutes }
