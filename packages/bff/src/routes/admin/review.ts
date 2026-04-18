import { Hono } from 'hono'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type {
  ReviewDraftRequest,
  ProfileDraftResponse,
  ScheduleDraftResponse,
} from '../../types/staffPortal'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminReviewRoutes = new Hono<AuthEnv>()

/**
 * 管理者用レビュールート
 * Backend の /admin/reviews/* エンドポイントへプロキシする。
 * Admin JWT トークンが必要。
 */

// 全ルートに認証ミドルウェアを適用
adminReviewRoutes.use('/*', authMiddleware)

// --- Profile Draft Review ---

/** GET /api/admin/reviews/profiles — 承認待ちプロフィール下書き一覧 */
adminReviewRoutes.get('/profiles', async (c) => {
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/reviews/profiles`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 401 | 500)
  }

  const data: ProfileDraftResponse[] = await res.json()
  return c.json(data)
})

/** GET /api/admin/reviews/profiles/:id — プロフィール下書き単体取得 */
adminReviewRoutes.get('/profiles/:id', async (c) => {
  const authHeader = c.get('authHeader')

  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/admin/reviews/profiles/${id}`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401 | 404)
  }

  const data: ProfileDraftResponse = await res.json()
  return c.json(data)
})

/** PUT /api/admin/reviews/profiles/:id — プロフィール下書きをレビュー */
adminReviewRoutes.put('/profiles/:id', async (c) => {
  const authHeader = c.get('authHeader')

  const id = c.req.param('id')
  const body = await c.req.json<ReviewDraftRequest>()

  const res = await fetch(`${BACKEND_URL}/admin/reviews/profiles/${id}`, {
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

/** PUT /api/admin/reviews/profiles/:id/content — プロフィール下書きの内容を修正 */
adminReviewRoutes.put('/profiles/:id/content', async (c) => {
  const authHeader = c.get('authHeader')

  const id = c.req.param('id')
  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/admin/reviews/profiles/${id}/content`, {
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

// --- Schedule Draft Review ---

/** GET /api/admin/reviews/schedules — 承認待ちスケジュール下書き一覧 */
adminReviewRoutes.get('/schedules', async (c) => {
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/reviews/schedules`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 401 | 500)
  }

  const data: ScheduleDraftResponse[] = await res.json()
  return c.json(data)
})

/** GET /api/admin/reviews/schedules/approved — 承認済み（未反映）スケジュール下書き一覧 */
adminReviewRoutes.get('/schedules/approved', async (c) => {
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/reviews/schedules/approved`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 401 | 500)
  }

  const data: ScheduleDraftResponse[] = await res.json()
  return c.json(data)
})

/** GET /api/admin/reviews/schedules/:id — スケジュール下書き単体取得 */
adminReviewRoutes.get('/schedules/:id', async (c) => {
  const authHeader = c.get('authHeader')

  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/admin/reviews/schedules/${id}`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401 | 404)
  }

  const data: ScheduleDraftResponse = await res.json()
  return c.json(data)
})

/** PUT /api/admin/reviews/schedules/:id — スケジュール下書きをレビュー */
adminReviewRoutes.put('/schedules/:id', async (c) => {
  const authHeader = c.get('authHeader')

  const id = c.req.param('id')
  const body = await c.req.json<ReviewDraftRequest>()

  const res = await fetch(`${BACKEND_URL}/admin/reviews/schedules/${id}`, {
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

/** PUT /api/admin/reviews/schedules/:id/content — スケジュール下書きの内容を修正 */
adminReviewRoutes.put('/schedules/:id/content', async (c) => {
  const authHeader = c.get('authHeader')

  const id = c.req.param('id')
  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/admin/reviews/schedules/${id}/content`, {
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

/** POST /api/admin/reviews/schedules/:id/publish — 承認済みスケジュールを店舗ページに反映 */
adminReviewRoutes.post('/schedules/:id/publish', async (c) => {
  const authHeader = c.get('authHeader')

  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/admin/reviews/schedules/${id}/publish`, {
    method: 'POST',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401)
  }

  const data = await res.json()
  return c.json(data)
})

export { adminReviewRoutes }
