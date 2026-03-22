import { Hono } from 'hono'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type {
  Staff,
  StaffResponse,
  StaffImage,
  StaffImageResponse,
  StaffSchedule,
  StaffScheduleResponse,
  StaffWithSchedules,
  StaffWithSchedulesResponse,
} from '../../types/staff'
import type { CreateStaffRequest, UpdateStaffRequest } from '../../types/admin'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminStaffRoutes = new Hono<AuthEnv>()

// 全ルートに認証ミドルウェアを適用
adminStaffRoutes.use('/*', authMiddleware)

/** Backend の StaffResponse を Frontend 向けに変換 */
function toStaff(raw: StaffResponse): Staff {
  return {
    id: raw.ID,
    shopId: raw.ShopID,
    name: raw.Name,
    role: raw.Role,
    bio: raw.Bio,
    imageUrl: raw.ImageURL,
    imageCropPosition: raw.ImageCropPosition,
    sortOrder: raw.SortOrder,
  }
}

/** Backend の StaffScheduleResponse を Frontend 向けに変換 */
function toSchedule(raw: StaffScheduleResponse): StaffSchedule {
  return {
    id: raw.ID,
    staffId: raw.StaffID,
    dayOfWeek: raw.DayOfWeek,
    startTime: raw.StartTime,
    endTime: raw.EndTime,
  }
}

/** Backend の StaffImageResponse を Frontend 向けに変換 */
function toImage(raw: StaffImageResponse): StaffImage {
  return {
    id: raw.ID,
    staffId: raw.StaffID,
    imageUrl: raw.ImageURL,
    isMain: raw.IsMain,
    sortOrder: raw.SortOrder,
  }
}

/** Backend の StaffWithSchedulesResponse を Frontend 向けに変換 */
function toStaffWithSchedules(data: StaffWithSchedulesResponse): StaffWithSchedules {
  return {
    staff: toStaff(data.Staff),
    schedules: data.Schedules ? data.Schedules.map(toSchedule) : [],
    images: data.Images ? data.Images.map(toImage) : [],
  }
}

/** POST /api/admin/staffs — スタッフ新規作成 */
adminStaffRoutes.post('/', async (c) => {
  const body = await c.req.json<CreateStaffRequest>()
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/admin/staffs`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 500)
  }

  const data: StaffWithSchedulesResponse = await res.json()
  return c.json(toStaffWithSchedules(data))
})

/** PUT /api/admin/staffs/:id — スタッフ情報更新 */
adminStaffRoutes.put('/:id', async (c) => {
  const id = c.req.param('id')
  const body = await c.req.json<UpdateStaffRequest>()
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/admin/staffs/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 500)
  }

  const data: StaffWithSchedulesResponse = await res.json()
  return c.json(toStaffWithSchedules(data))
})

/** DELETE /api/admin/staffs/:id — スタッフ削除 */
adminStaffRoutes.delete('/:id', async (c) => {
  const id = c.req.param('id')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/admin/staffs/${id}`, {
    method: 'DELETE',
    headers: {
      Authorization: authHeader,
    },
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete staff' }))
    return c.json(error, res.status as 500)
  }

  return new Response(null, { status: 204 })
})

// ============================================================
// Staff Image Management
// ============================================================

/** POST /api/admin/staffs/:id/images — 画像アップロード（multipart/form-data をそのまま Backend へ転送） */
adminStaffRoutes.post('/:id/images', async (c) => {
  const id = c.req.param('id')
  const authHeader = c.get('authHeader') as string

  // multipart body をそのまま Backend へ転送
  const body = await c.req.raw.arrayBuffer()
  const contentType = c.req.header('content-type') || ''

  const res = await fetch(`${BACKEND_URL}/admin/staffs/${id}/images`, {
    method: 'POST',
    headers: {
      'Content-Type': contentType,
      Authorization: authHeader,
    },
    body,
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to upload image' }))
    return c.json(error, res.status as 500)
  }

  const data: StaffImageResponse = await res.json()
  return c.json(toImage(data))
})

/** DELETE /api/admin/staffs/:id/images/:imageId — 画像を削除 */
adminStaffRoutes.delete('/:id/images/:imageId', async (c) => {
  const id = c.req.param('id')
  const imageId = c.req.param('imageId')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/admin/staffs/${id}/images/${imageId}`, {
    method: 'DELETE',
    headers: {
      Authorization: authHeader,
    },
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete image' }))
    return c.json(error, res.status as 500)
  }

  return new Response(null, { status: 204 })
})

/** PUT /api/admin/staffs/:id/images/main — メイン画像を設定 */
adminStaffRoutes.put('/:id/images/main', async (c) => {
  const id = c.req.param('id')
  const body = await c.req.json()
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/admin/staffs/${id}/images/main`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to set main image' }))
    return c.json(error, res.status as 500)
  }

  return c.json({ message: 'ok' })
})

export { adminStaffRoutes }
