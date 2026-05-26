import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import { createStaffAccountSchema, updateStaffAccountSchema } from '../../schemas'
import { StaffAccountResponseSchema, ErrorResponseSchema } from '../../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminAccountRoutes = new OpenAPIHono<AuthEnv>()

// 全ルートに認証ミドルウェアを適用
adminAccountRoutes.use('/*', authMiddleware)

// ============================================================
// Route definitions
// ============================================================

const listAccountsRoute = createRoute({
  method: 'get',
  path: '/',
  tags: ['管理者 - アカウント'],
  summary: 'スタッフアカウント一覧',
  security: [{ Bearer: [] }],
  responses: {
    200: {
      content: { 'application/json': { schema: z.array(StaffAccountResponseSchema) } },
      description: 'アカウント一覧',
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

const getAccountByStaffRoute = createRoute({
  method: 'get',
  path: '/staff/{staffId}',
  tags: ['管理者 - アカウント'],
  summary: 'スタッフ ID でアカウント取得',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ staffId: z.string().openapi({ description: 'スタッフ ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffAccountResponseSchema } },
      description: 'アカウント詳細',
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

const createAccountRoute = createRoute({
  method: 'post',
  path: '/',
  tags: ['管理者 - アカウント'],
  summary: 'スタッフアカウント作成',
  security: [{ Bearer: [] }],
  request: {
    body: {
      content: { 'application/json': { schema: createStaffAccountSchema } },
      required: true,
    },
  },
  responses: {
    201: {
      content: { 'application/json': { schema: StaffAccountResponseSchema } },
      description: '作成成功',
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

const updateAccountRoute = createRoute({
  method: 'put',
  path: '/{id}',
  tags: ['管理者 - アカウント'],
  summary: 'スタッフアカウント更新',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'アカウント ID' }) }),
    body: {
      content: { 'application/json': { schema: updateStaffAccountSchema } },
      required: true,
    },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: StaffAccountResponseSchema } },
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

const deleteAccountRoute = createRoute({
  method: 'delete',
  path: '/{id}',
  tags: ['管理者 - アカウント'],
  summary: 'スタッフアカウント削除',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'アカウント ID' }) }),
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

// ============================================================
// Handlers
// ============================================================

/** スタッフアカウントレスポンス型 */
interface StaffAccountResponse {
  id: string
  staffId: string
  username: string
  createdAt: string
  updatedAt: string
}

adminAccountRoutes.openapi(listAccountsRoute, async (c) => {
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staff-accounts`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data: StaffAccountResponse[] = await res.json()
  return c.json(data, 200)
})

adminAccountRoutes.openapi(getAccountByStaffRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { staffId } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staff-accounts/staff/${staffId}`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status >= 500 ? 500 : 401)
  }

  const data = await res.json()
  return c.json(data, 200)
})

adminAccountRoutes.openapi(createAccountRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staff-accounts`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status === 401 ? 401 : 400)
  }

  const data: StaffAccountResponse = await res.json()
  return c.json(data, 201)
})

adminAccountRoutes.openapi(updateAccountRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staff-accounts/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status === 401 ? 401 : 400)
  }

  const data: StaffAccountResponse = await res.json()
  return c.json(data, 200)
})

adminAccountRoutes.openapi(deleteAccountRoute, async (c) => {
  const authHeader = c.get('authHeader')
  const { id } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/staff-accounts/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, res.status === 401 ? 401 : 400)
  }

  return c.body(null, 204)
})

export { adminAccountRoutes }
