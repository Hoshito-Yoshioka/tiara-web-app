import { describe, it, expect } from 'vitest'
import { Hono } from 'hono'
import { authMiddleware, type AuthEnv } from '../middleware/auth'

describe('authMiddleware', () => {
  const app = new Hono<AuthEnv>()
  app.use('/*', authMiddleware)
  app.get('/test', (c) => c.json({ authHeader: c.get('authHeader') }))

  it('有効な Bearer トークンがある場合、リクエストを通過させる', async () => {
    const res = await app.request('/test', {
      headers: { Authorization: 'Bearer test-token-123' },
    })
    expect(res.status).toBe(200)
    const body = await res.json()
    expect(body.authHeader).toBe('Bearer test-token-123')
  })

  it('Authorization ヘッダーがない場合、401 を返す', async () => {
    const res = await app.request('/test')
    expect(res.status).toBe(401)
    const body = await res.json()
    expect(body.error).toBe('Authorization header is required')
  })

  it('Bearer プレフィックスがない場合、401 を返す', async () => {
    const res = await app.request('/test', {
      headers: { Authorization: 'Basic invalid-token' },
    })
    expect(res.status).toBe(401)
  })

  it('空の Authorization ヘッダーの場合、401 を返す', async () => {
    const res = await app.request('/test', {
      headers: { Authorization: '' },
    })
    expect(res.status).toBe(401)
  })
})
