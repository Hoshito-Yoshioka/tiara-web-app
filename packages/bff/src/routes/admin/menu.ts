import { OpenAPIHono, createRoute, z } from '@hono/zod-openapi'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type { MenuCategoryResponse, MenuItemResponse } from '../../types/menu'
import {
  createMenuCategorySchema,
  updateMenuCategorySchema,
  createMenuItemSchema,
  updateMenuItemSchema,
} from '../../schemas'
import { MenuCategorySchema, MenuItemSchema, ErrorResponseSchema } from '../../schemas/responses'
import { toMenuCategory, toMenuItem } from '../menu'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminMenuRoutes = new OpenAPIHono<AuthEnv>()

adminMenuRoutes.use('/*', authMiddleware)

// ============================================================
// Route definitions — MenuCategory
// ============================================================

const createCategoryRoute = createRoute({
  method: 'post',
  path: '/categories',
  tags: ['管理者 - メニュー'],
  summary: 'カテゴリ作成',
  security: [{ Bearer: [] }],
  request: {
    body: {
      content: { 'application/json': { schema: createMenuCategorySchema } },
      required: true,
    },
  },
  responses: {
    201: {
      content: { 'application/json': { schema: MenuCategorySchema } },
      description: '作成成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const updateCategoryRoute = createRoute({
  method: 'put',
  path: '/categories/{id}',
  tags: ['管理者 - メニュー'],
  summary: 'カテゴリ更新',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'カテゴリ ID' }) }),
    body: {
      content: { 'application/json': { schema: updateMenuCategorySchema } },
      required: true,
    },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: MenuCategorySchema } },
      description: '更新成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const deleteCategoryRoute = createRoute({
  method: 'delete',
  path: '/categories/{id}',
  tags: ['管理者 - メニュー'],
  summary: 'カテゴリ削除',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'カテゴリ ID' }) }),
  },
  responses: {
    204: { description: '削除成功' },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

// ============================================================
// Route definitions — MenuItem
// ============================================================

const createItemRoute = createRoute({
  method: 'post',
  path: '/items',
  tags: ['管理者 - メニュー'],
  summary: 'アイテム作成',
  security: [{ Bearer: [] }],
  request: {
    body: {
      content: { 'application/json': { schema: createMenuItemSchema } },
      required: true,
    },
  },
  responses: {
    201: {
      content: { 'application/json': { schema: MenuItemSchema } },
      description: '作成成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const updateItemRoute = createRoute({
  method: 'put',
  path: '/items/{id}',
  tags: ['管理者 - メニュー'],
  summary: 'アイテム更新',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'アイテム ID' }) }),
    body: {
      content: { 'application/json': { schema: updateMenuItemSchema } },
      required: true,
    },
  },
  responses: {
    200: {
      content: { 'application/json': { schema: MenuItemSchema } },
      description: '更新成功',
    },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

const deleteItemRoute = createRoute({
  method: 'delete',
  path: '/items/{id}',
  tags: ['管理者 - メニュー'],
  summary: 'アイテム削除',
  security: [{ Bearer: [] }],
  request: {
    params: z.object({ id: z.string().openapi({ description: 'アイテム ID' }) }),
  },
  responses: {
    204: { description: '削除成功' },
    500: {
      content: { 'application/json': { schema: ErrorResponseSchema } },
      description: 'サーバーエラー',
    },
  },
})

// ============================================================
// Handlers — MenuCategory
// ============================================================

adminMenuRoutes.openapi(createCategoryRoute, async (c) => {
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/menu/categories`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json((await res.json()) as { error: string }, 500)
  }
  const data: MenuCategoryResponse = await res.json()
  return c.json(toMenuCategory(data), 201)
})

adminMenuRoutes.openapi(updateCategoryRoute, async (c) => {
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/menu/categories/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json((await res.json()) as { error: string }, 500)
  }
  const data: MenuCategoryResponse = await res.json()
  return c.json(toMenuCategory(data), 200)
})

adminMenuRoutes.openapi(deleteCategoryRoute, async (c) => {
  const { id } = c.req.valid('param')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/menu/categories/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete category' }))
    return c.json(error as { error: string }, 500)
  }
  return c.body(null, 204)
})

// ============================================================
// Handlers — MenuItem
// ============================================================

adminMenuRoutes.openapi(createItemRoute, async (c) => {
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/menu/items`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json((await res.json()) as { error: string }, 500)
  }
  const data: MenuItemResponse = await res.json()
  return c.json(toMenuItem(data), 201)
})

adminMenuRoutes.openapi(updateItemRoute, async (c) => {
  const { id } = c.req.valid('param')
  const body = c.req.valid('json')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/menu/items/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json((await res.json()) as { error: string }, 500)
  }
  const data: MenuItemResponse = await res.json()
  return c.json(toMenuItem(data), 200)
})

adminMenuRoutes.openapi(deleteItemRoute, async (c) => {
  const { id } = c.req.valid('param')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/api/v1/admin/menu/items/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete item' }))
    return c.json(error as { error: string }, 500)
  }
  return c.body(null, 204)
})

export { adminMenuRoutes }
