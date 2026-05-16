import { describe, it, expect, vi, beforeEach } from 'vitest'
import { apiFetch, apiUpload } from '@/lib/api'

// fetch のモック
const fetchMock = vi.fn()
vi.stubGlobal('fetch', fetchMock)

/** JSON レスポンスを返すヘルパー */
function jsonResponse(body: unknown, status = 200, statusText = 'OK'): Response {
  return new Response(JSON.stringify(body), {
    status,
    statusText,
    headers: { 'content-type': 'application/json' },
  })
}

/** 非 JSON レスポンスを返すヘルパー */
function htmlResponse(): Response {
  return new Response('<html></html>', {
    status: 200,
    statusText: 'OK',
    headers: { 'content-type': 'text/html' },
  })
}

beforeEach(() => {
  fetchMock.mockReset()
})

// ============================================================
// apiFetch
// ============================================================
describe('apiFetch', () => {
  it('API_BASE_URL + path で fetch を呼び出す', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ ok: true }))

    await apiFetch('/api/health')

    expect(fetchMock).toHaveBeenCalledOnce()
    const [url] = fetchMock.mock.calls[0]
    expect(url).toBe('http://localhost:3001/api/health')
  })

  it('デフォルトで Content-Type: application/json を付与する', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ ok: true }))

    await apiFetch('/api/test')

    const [, init] = fetchMock.mock.calls[0]
    expect(init.headers['Content-Type']).toBe('application/json')
  })

  it('options.headers でカスタムヘッダーを追加できる', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ ok: true }))

    await apiFetch('/api/test', {
      headers: { Authorization: 'Bearer token123' },
    })

    const [, init] = fetchMock.mock.calls[0]
    expect(init.headers['Content-Type']).toBe('application/json')
    expect(init.headers['Authorization']).toBe('Bearer token123')
  })

  it('options.headers で Content-Type を上書きできる', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ ok: true }))

    await apiFetch('/api/test', {
      headers: { 'Content-Type': 'text/plain' },
    })

    const [, init] = fetchMock.mock.calls[0]
    expect(init.headers['Content-Type']).toBe('text/plain')
  })

  it('204 No Content のとき undefined を返す', async () => {
    fetchMock.mockResolvedValue(
      new Response(null, {
        status: 204,
        statusText: 'No Content',
        headers: { 'content-type': 'application/json' },
      })
    )

    const result = await apiFetch('/api/items/1', { method: 'DELETE' })

    expect(result).toBeUndefined()
  })

  it('非 JSON レスポンスでエラーをスローする', async () => {
    fetchMock.mockResolvedValue(htmlResponse())

    await expect(apiFetch('/api/test')).rejects.toThrow(
      'BFF から JSON ではなく text/html が返されました'
    )
  })

  it('非 ok レスポンス（4xx/5xx）でステータスコード付きエラーをスローする', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ error: 'not found' }, 404, 'Not Found'))

    await expect(apiFetch('/api/missing')).rejects.toThrow('API Error: 404 Not Found')
  })

  it('ネットワークエラーで接続エラーメッセージをスローする', async () => {
    fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))

    await expect(apiFetch('/api/test')).rejects.toThrow(
      'BFF (http://localhost:3001) に接続できません'
    )
  })
})

// ============================================================
// apiUpload
// ============================================================
describe('apiUpload', () => {
  it('method: POST と FormData で fetch を呼び出す', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ id: 1 }))
    const formData = new FormData()
    formData.append('file', new Blob(['test']), 'test.png')

    await apiUpload('/api/upload', formData)

    expect(fetchMock).toHaveBeenCalledOnce()
    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe('http://localhost:3001/api/upload')
    expect(init.method).toBe('POST')
    expect(init.body).toBe(formData)
  })

  it('Content-Type を明示設定しない（ブラウザが boundary 付きで自動設定）', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ id: 1 }))
    const formData = new FormData()

    await apiUpload('/api/upload', formData)

    const [, init] = fetchMock.mock.calls[0]
    expect(init.headers).not.toHaveProperty('Content-Type')
  })

  it('options.headers でカスタムヘッダーを渡せる', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ id: 1 }))
    const formData = new FormData()

    await apiUpload('/api/upload', formData, {
      headers: { Authorization: 'Bearer token123' },
    })

    const [, init] = fetchMock.mock.calls[0]
    expect(init.headers['Authorization']).toBe('Bearer token123')
    expect(init.headers).not.toHaveProperty('Content-Type')
  })

  it('非 JSON レスポンスでエラーをスローする', async () => {
    fetchMock.mockResolvedValue(htmlResponse())
    const formData = new FormData()

    await expect(apiUpload('/api/upload', formData)).rejects.toThrow(
      'BFF から JSON ではなく text/html が返されました'
    )
  })

  it('非 ok レスポンスでステータスコード付きエラーをスローする', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ error: 'too large' }, 413, 'Payload Too Large'))
    const formData = new FormData()

    await expect(apiUpload('/api/upload', formData)).rejects.toThrow(
      'API Error: 413 Payload Too Large'
    )
  })

  it('ネットワークエラーで接続エラーメッセージをスローする', async () => {
    fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))
    const formData = new FormData()

    await expect(apiUpload('/api/upload', formData)).rejects.toThrow(
      'BFF (http://localhost:3001) に接続できません'
    )
  })
})
