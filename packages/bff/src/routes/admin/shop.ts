import { Hono } from 'hono'
import { zValidator } from '@hono/zod-validator'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type { Shop, ShopResponse } from '../../types/shop'
import { updateShopSchema } from '../../schemas'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminShopRoutes = new Hono<AuthEnv>()

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

/** PUT /api/admin/shops/:id — 店舗情報を更新 */
adminShopRoutes.put('/:id', zValidator('json', updateShopSchema), async (c) => {
  const id = c.req.param('id')
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
    return c.json(error, res.status as 500)
  }

  const data: ShopResponse = await res.json()
  return c.json(toShop(data))
})

export { adminShopRoutes }
