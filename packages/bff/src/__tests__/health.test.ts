import { describe, it, expect } from 'vitest'
import { Hono } from 'hono'

describe('BFF Health Check', () => {
  const app = new Hono()
  app.get('/', (c) => c.json({ message: 'Tiara BFF is running' }))

  it('GET / はステータス 200 と正しいメッセージを返す', async () => {
    const res = await app.request('/')
    expect(res.status).toBe(200)
    const body = await res.json()
    expect(body.message).toBe('Tiara BFF is running')
  })
})
