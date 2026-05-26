import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type { Shop, ShopResponse } from '../../types/shop'
import { updateShopSchema } from '../../schemas'
import { ShopSchema, ErrorResponseSchema } from '../../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminShopRoutes = new OpenAPIHono<AuthEnv>()

// 全ルートに認証ミドルウェアを適用
adminShopRoutes.use('/*', authMiddleware)

/** Backend のレスポンスを Frontend 向けに変換（フィールド選択） */
function toShop(raw: ShopResponse): Shop {
  return {
    id: raw.id,
    name: raw.name,
    address: raw.address,
    openingTime: raw.openingTime,
    closingTime: raw.closingTime,
  }
}

// --- Route definition ---

const updateShopRoute = createRoute({
  method: 'put',
  path: '/{id}',
  tags: ['管理者 - 店舗'],
  summary: '店舗情報更新',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: '店舗 ID' }) }),
    body: { content: { 'application/json': { schema: updateShopSchema } }, required: true },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ShopSchema } },
      description: '更新成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

// --- Handler ---

adminShopRoutes.openapi(updateShopRoute, async (c) => {
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader') as string

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/shops/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: authHeader,
    },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const error = await res.json()
    return c.json(error as { error: string }, 500)
  }

  const data: ShopResponse = await res.json()
  return c.json(toShop(data), 200)
})

export { adminShopRoutes }
