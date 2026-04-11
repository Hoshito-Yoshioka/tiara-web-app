import { Hono } from 'hono'
import type {
  SaveProfileDraftRequest,
  SaveScheduleDraftRequest,
  ProfileDraftResponse,
  ScheduleDraftResponse,
} from '../types/staffPortal'
import type { StaffImageResponse, StaffImage } from '../types/staff'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const portalRoutes = new Hono()

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

// --- Image Management ---

/** GET /api/portal/images — 自分の画像一覧を取得 */
portalRoutes.get('/images', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/portal/images`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 400 | 401 | 500)
  }

  const data: StaffImageResponse[] = await res.json()
  return c.json(data.map(toImage))
})

/** POST /api/portal/images — 画像アップロード（multipart/form-data をそのまま Backend へ転送） */
portalRoutes.post('/images', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = await c.req.raw.arrayBuffer()
  const contentType = c.req.header('content-type') || ''

  const res = await fetch(`${BACKEND_URL}/portal/images`, {
    method: 'POST',
    headers: {
      'Content-Type': contentType,
      Authorization: authHeader,
    },
    body,
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to upload image' }))
    return c.json(error, res.status as 400 | 500)
  }

  const data: StaffImageResponse = await res.json()
  return c.json(toImage(data))
})

/** DELETE /api/portal/images/:imageId — 画像を削除 */
portalRoutes.delete('/images/:imageId', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const imageId = c.req.param('imageId')

  const res = await fetch(`${BACKEND_URL}/portal/images/${imageId}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete image' }))
    return c.json(error, res.status as 400 | 500)
  }

  return new Response(null, { status: 204 })
})

/** PUT /api/portal/images/main — メイン画像を設定 */
portalRoutes.put('/images/main', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/portal/images/main`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to set main image' }))
    return c.json(error, res.status as 400 | 500)
  }

  return c.json({ message: 'ok' })
})

/** PUT /api/portal/images/:imageId/crop — 画像のクロップ位置を更新 */
portalRoutes.put('/images/:imageId/crop', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const imageId = c.req.param('imageId')
  const body = await c.req.json()

  const res = await fetch(`${BACKEND_URL}/portal/images/${imageId}/crop`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to update crop position' }))
    return c.json(error, res.status as 400 | 500)
  }

  const data: StaffImageResponse = await res.json()
  return c.json(toImage(data))
})

export { portalRoutes }
