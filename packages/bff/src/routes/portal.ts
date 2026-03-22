import { Hono } from 'hono'
import type {
  SaveProfileDraftRequest,
  SaveScheduleDraftRequest,
  ProfileDraftResponse,
  ScheduleDraftResponse,
} from '../types/staffPortal'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const portalRoutes = new Hono()

/**
 * スタッフポータルルート
 * Backend の /portal/* エンドポイントへプロキシする。
 * Backend がすでに camelCase JSON を返すため、変換は不要。
 */

// --- Profile Draft ---

/** GET /api/portal/profile — 自分のプロフィール下書きを取得 */
portalRoutes.get('/profile', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/portal/profile`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401 | 500)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data)
})

/** PUT /api/portal/profile — プロフィール下書きを保存 */
portalRoutes.put('/profile', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = await c.req.json<SaveProfileDraftRequest>()

  const res = await fetch(`${BACKEND_URL}/portal/profile`, {
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

  const data: ProfileDraftResponse = await res.json()
  return c.json(data)
})

/** POST /api/portal/profile/:id/submit — プロフィール下書きを承認申請 */
portalRoutes.post('/profile/:id/submit', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/portal/profile/${id}/submit`, {
    method: 'POST',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data)
})

// --- Schedule Draft ---

/** GET /api/portal/schedule — 自分のスケジュール下書きを取得 */
portalRoutes.get('/schedule', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/portal/schedule`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401 | 500)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data)
})

/** PUT /api/portal/schedule — スケジュール下書きを保存 */
portalRoutes.put('/schedule', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = await c.req.json<SaveScheduleDraftRequest>()

  const res = await fetch(`${BACKEND_URL}/portal/schedule`, {
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

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data)
})

/** POST /api/portal/schedule/:id/submit — スケジュール下書きを承認申請 */
portalRoutes.post('/schedule/:id/submit', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/portal/schedule/${id}/submit`, {
    method: 'POST',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data)
})

export { portalRoutes }
