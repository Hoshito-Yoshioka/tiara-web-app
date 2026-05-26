import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import type { Shop, ShopResponse } from '../types/shop'
import { ShopSchema, ErrorResponseSchema } from '../schemas/responses'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const shopRoutes = new OpenAPIHono()

/**
 * Backend のレスポンスを Frontend 向けに変換する。
 * Backend が handler DTO で camelCase + "HH:MM" 形式を返すため、
 * BFF はフィールド選択のみ行う。
 */
function toShop(raw: ShopResponse): Shop {
  return {
    id: raw.id,
    name: raw.name,
    address: raw.address,
    openingTime: raw.openingTime,
    closingTime: raw.closingTime,
  }
}

// --- Route definitions ---

const listShopsRoute = createRoute({
  method: 'get',
  path: '/',
  tags: ['店舗'],
  summary: '店舗一覧取得',
  responses: {
    200: {
      content: { 'application/json': { schema: z.array(ShopSchema) } },
      description: '店舗一覧',
    },
    502: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'Backend エラー',
    },
  },
})

const getShopRoute = createRoute({
  method: 'get',
  path: '/{id}',
  tags: ['店舗'],
  summary: '店舗詳細取得',
  request: {
    params: z.object({ id: z.string().openapi({ description: '店舗 ID' }) }),
  },
  responses: {
    200: {
      content: { 'application/json': { schema: ShopSchema } },
      description: '店舗詳細',
    },
    404: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: '店舗が見つからない',
    },
  },
})

// --- Handlers ---

shopRoutes.openapi(listShopsRoute, async (c) => {
  const res = await fetch(`${BACKEND_URL}/api/v1/shops`)

  if (!res.ok) {
    return c.json({ error: 'Failed to fetch shops from backend' }, 502)
  }

  const data: ShopResponse[] = await res.json()
  const shops: Shop[] = data.map(toShop)

  return c.json(shops, 200)
})

shopRoutes.openapi(getShopRoute, async (c) => {
  const { id } = c.req.valid('param')

  const res = await fetch(`${BACKEND_URL}/api/v1/shops/${id}`)

  if (!res.ok) {
    return c.json({ error: 'Shop not found' }, 404)
  }

  const data: ShopResponse = await res.json()
  const shop: Shop = toShop(data)

  return c.json(shop, 200)
})

export { shopRoutes }
