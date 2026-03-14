import { Hono } from 'hono'
import type { Shop, ShopResponse } from '../types/shop'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const shopRoutes = new Hono()

/**
 * Backend のレスポンス（PascalCase）を Frontend 向け（camelCase）に変換する。
 * BFF の責務の一つである「データ整形」をここで担う。
 */
function toShop(raw: ShopResponse): Shop {
  return {
    id: raw.ID,
    name: raw.Name,
    address: raw.Address,
    openingTime: raw.OpeningTime,
    closingTime: raw.ClosingTime,
  }
}

/** GET /api/shops — 店舗一覧を取得 */
shopRoutes.get('/', async (c) => {
  const res = await fetch(`${BACKEND_URL}/shops`)

  if (!res.ok) {
    return c.json({ error: 'Failed to fetch shops from backend' }, 502)
  }

  const data: ShopResponse[] = await res.json()
  const shops: Shop[] = data.map(toShop)

  return c.json(shops)
})

/** GET /api/shops/:id — 店舗詳細を取得 */
shopRoutes.get('/:id', async (c) => {
  const id = c.req.param('id')

  const res = await fetch(`${BACKEND_URL}/shops/${id}`)

  if (!res.ok) {
    return c.json({ error: 'Shop not found' }, 404)
  }

  const data: ShopResponse = await res.json()
  const shop: Shop = toShop(data)

  return c.json(shop)
})

export { shopRoutes }
