import { describe, it, expect } from 'vitest'
import { app } from '../app'

describe('API docs routes', () => {
  it('/api/v1/redoc は Redoc 用 HTML を返す', async () => {
    const res = await app.request('/api/v1/redoc')

    expect(res.status).toBe(200)
    expect(res.headers.get('content-type')).toContain('text/html')
    expect(res.headers.get('content-security-policy')).toContain('cdn.redocly.com')

    const body = await res.text()
    expect(body).toContain('Redoc.init')
    expect(body).toContain('/api/v1/doc')
  })
})
