import { Hono } from 'hono'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminAccountRoutes = new Hono()

/**
 * 管理者用スタッフアカウント管理ルート。
 * Backend の /admin/staff-accounts/* エンドポイントへプロキシする。
 * Admin JWT トークンが必要。
 */

/** スタッフアカウントレスポンス型 */
interface StaffAccountResponse {
  id: string
  staffId: string
  username: string
  createdAt: string
  updatedAt: string
}

/** GET /api/admin/staff-accounts — スタッフアカウント一覧 */
adminAccountRoutes.get('/', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/admin/staff-accounts`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 401 | 500)
  }

  const data: StaffAccountResponse[] = await res.json()
  return c.json(data)
})

/** GET /api/admin/staff-accounts/staff/:staffId — スタッフIDでアカウント取得 */
adminAccountRoutes.get('/staff/:staffId', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const staffId = c.req.param('staffId')
  const res = await fetch(`${BACKEND_URL}/admin/staff-accounts/staff/${staffId}`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 401 | 500)
  }

  const data = await res.json()
  return c.json(data)
})

/** POST /api/admin/staff-accounts — スタッフアカウント作成 */
adminAccountRoutes.post('/', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/admin/staff-accounts`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401)
  }

  const data: StaffAccountResponse = await res.json()
  return c.json(data, 201)
})

/** PUT /api/admin/staff-accounts/:id — スタッフアカウント更新 */
adminAccountRoutes.put('/:id', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const id = c.req.param('id')
  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/admin/staff-accounts/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401)
  }

  const data: StaffAccountResponse = await res.json()
  return c.json(data)
})

/** DELETE /api/admin/staff-accounts/:id — スタッフアカウント削除 */
adminAccountRoutes.delete('/:id', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/admin/staff-accounts/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401)
  }

  return c.body(null, 204)
})

export { adminAccountRoutes }
