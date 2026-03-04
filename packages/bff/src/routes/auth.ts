import { Hono } from 'hono'
import type { LoginRequest, LoginResponse } from '../types/admin'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const authRoutes = new Hono()

/** POST /api/auth/login — ログイン（Backend へプロキシ） */
authRoutes.post('/login', async (c) => {
  const body = await c.req.json<LoginRequest>()

  const res = await fetch(`${BACKEND_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error, res.status as 401)
  }

  const data: LoginResponse = await res.json()
  return c.json(data)
})

/** GET /api/auth/verify — トークン検証（Backend へプロキシ） */
authRoutes.get('/verify', async (c) => {
  const authHeader = c.req.header('Authorization')

  if (!authHeader) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  const res = await fetch(`${BACKEND_URL}/admin/auth/verify`, {
    headers: { Authorization: authHeader },
  })

  if (!res.ok) {
    return c.json({ error: 'Invalid token' }, 401)
  }

  const data = await res.json()
  return c.json(data)
})

export { authRoutes }
