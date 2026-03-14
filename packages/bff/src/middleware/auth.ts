import { createMiddleware } from 'hono/factory'

/** authMiddleware が Context に設定するカスタム変数の型定義 */
export type AuthEnv = {
  Variables: {
    authHeader: string
  }
}

/**
 * BFF の認証ミドルウェア。
 * Frontend から送られた Authorization ヘッダーの存在を確認し、
 * Backend への転送用にコンテキストに保持する。
 * 実際のトークン検証は Backend の JWT ミドルウェアで行うため、
 * BFF ではヘッダーの存在チェックのみを担う。
 */
export const authMiddleware = createMiddleware<AuthEnv>(async (c, next) => {
  const authHeader = c.req.header('Authorization')

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return c.json({ error: 'Authorization header is required' }, 401)
  }

  // トークンをコンテキストに保存（admin ルートで Backend へ転送する際に使用）
  c.set('authHeader', authHeader)

  await next()
})
