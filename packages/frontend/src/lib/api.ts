/**
 * API クライアントの共通設定。
 * すべての API 呼び出しはこの基底関数を通すことで、
 * ベース URL やエラーハンドリングを一元管理する。
 *
 * 開発環境: BFF は localhost:3001 で起動する。
 *           BFF 側で CORS が設定済み（origin: localhost:5173）。
 * 本番環境: VITE_API_BASE_URL で BFF の URL を指定する。
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001'

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
  } catch (_e) {
    throw new Error(
      `BFF (${API_BASE_URL}) に接続できません。BFF サーバーが起動しているか確認してください。`
    )
  }

  // 204 No Content のハンドリング（DELETE 等のレスポンスボディがないケース）
  if (res.status === 204) {
    return undefined as T
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

/**
 * マルチパートファイルアップロード用の API クライアント。
 * Content-Type を自動設定させるため、明示的に設定しない。
 */
export async function apiUpload<T>(
  path: string,
  body: FormData,
  options?: RequestInit
): Promise<T> {
  const url = `${API_BASE_URL}${path}`

  let res: Response
  try {
    res = await fetch(url, {
      method: 'POST',
      ...options,
      headers: {
        // Content-Type を設定しない（ブラウザが boundary 付きで自動設定）
        ...options?.headers,
      },
      body,
    })
  } catch (_e) {
    throw new Error(
      `BFF (${API_BASE_URL}) に接続できません。BFF サーバーが起動しているか確認してください。`
    )
  }

  const contentType = res.headers.get('content-type') ?? ''
  if (!contentType.includes('application/json')) {
    throw new Error(`BFF から JSON ではなく ${contentType || 'unknown'} が返されました。`)
  }

  if (!res.ok) {
    throw new Error(`API Error: ${res.status} ${res.statusText}`)
  }

  return res.json() as Promise<T>
}
