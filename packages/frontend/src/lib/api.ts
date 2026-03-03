/**
 * API クライアントの共通設定。
 * すべての API 呼び出しはこの基底関数を通すことで、
 * ベース URL やエラーハンドリングを一元管理する。
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'

export async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE_URL}${path}`

  let res: Response
  try {
    res = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    })
  } catch (e) {
    throw new Error(
      `BFF (${API_BASE_URL}) に接続できません。BFF サーバーが起動しているか確認してください。`
    )
  }

  // BFF ではなく別サーバー（Vite 等）が応答した場合を検知
  const contentType = res.headers.get('content-type') ?? ''
  if (!contentType.includes('application/json')) {
    throw new Error(
      `BFF から JSON ではなく ${contentType || 'unknown'} が返されました。BFF が正しいポートで起動しているか確認してください。`
    )
  }

  if (!res.ok) {
    throw new Error(`API Error: ${res.status} ${res.statusText}`)
  }

  return res.json() as Promise<T>
}
