import { Hono } from 'hono'
import { zValidator } from '@hono/zod-validator'
import { loginSchema } from '../schemas'
import type { StaffLoginResponse } from '../types/staffPortal'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const staffAuthRoutes = new Hono()

/** POST /api/staff-auth/login — スタッフログイン（Backend へプロキシ） */
staffAuthRoutes.post('/login', zValidator('json', loginSchema), async (c) => {
  const body = c.req.valid('json')

  const res = await fetch(`${BACKEND_URL}/staff-auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 401)
  }

  const data: StaffLoginResponse = await res.json()
  return c.json(data)
})

/** GET /api/staff-auth/verify — スタッフトークン検証（Backend へプロキシ） */
staffAuthRoutes.get('/verify', async (c) => {
  const authHeader = c.req.header('Authorization')

  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/portal/auth/verify`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    return c.json({ error: 'Invalid token' }, 401)
  }

  const data = await res.json()
  return c.json(data)
})

export { staffAuthRoutes }
