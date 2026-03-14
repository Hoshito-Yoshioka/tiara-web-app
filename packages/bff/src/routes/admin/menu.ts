import { Hono } from 'hono'
import { authMiddleware, type AuthEnv } from '../../middleware/auth'
import type { MenuCategoryResponse, MenuItemResponse } from '../../types/menu'
import type {
  CreateMenuCategoryRequest,
  UpdateMenuCategoryRequest,
  CreateMenuItemRequest,
  UpdateMenuItemRequest,
} from '../../types/admin'
import { toMenuCategory, toMenuItem, toMenuCategoryWithItems } from '../menu'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const adminMenuRoutes = new Hono<AuthEnv>()

adminMenuRoutes.use('/*', authMiddleware)

// ============================================================
// MenuCategory CRUD
// ============================================================

/** POST /api/admin/menu/categories */
adminMenuRoutes.post('/categories', async (c) => {
  const body = await c.req.json<CreateMenuCategoryRequest>()
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/menu/categories`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json(await res.json(), res.status as 500)
  }
  const data: MenuCategoryResponse = await res.json()
  return c.json(toMenuCategory(data), 201)
})

/** PUT /api/admin/menu/categories/:id */
adminMenuRoutes.put('/categories/:id', async (c) => {
  const id = c.req.param('id')
  const body = await c.req.json<UpdateMenuCategoryRequest>()
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/menu/categories/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json(await res.json(), res.status as 500)
  }
  const data: MenuCategoryResponse = await res.json()
  return c.json(toMenuCategory(data))
})

/** DELETE /api/admin/menu/categories/:id */
adminMenuRoutes.delete('/categories/:id', async (c) => {
  const id = c.req.param('id')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/menu/categories/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete category' }))
    return c.json(error, res.status as 500)
  }
  return new Response(null, { status: 204 })
})

// ============================================================
// MenuItem CRUD
// ============================================================

/** POST /api/admin/menu/items */
adminMenuRoutes.post('/items', async (c) => {
  const body = await c.req.json<CreateMenuItemRequest>()
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/menu/items`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json(await res.json(), res.status as 500)
  }
  const data: MenuItemResponse = await res.json()
  return c.json(toMenuItem(data), 201)
})

/** PUT /api/admin/menu/items/:id */
adminMenuRoutes.put('/items/:id', async (c) => {
  const id = c.req.param('id')
  const body = await c.req.json<UpdateMenuItemRequest>()
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/menu/items/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: authHeader },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    return c.json(await res.json(), res.status as 500)
  }
  const data: MenuItemResponse = await res.json()
  return c.json(toMenuItem(data))
})

/** DELETE /api/admin/menu/items/:id */
adminMenuRoutes.delete('/items/:id', async (c) => {
  const id = c.req.param('id')
  const authHeader = c.get('authHeader')

  const res = await fetch(`${BACKEND_URL}/admin/menu/items/${id}`, {
    method: 'DELETE',
    headers: { Authorization: authHeader },
  })
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Failed to delete item' }))
    return c.json(error, res.status as 500)
  }
  return new Response(null, { status: 204 })
})

export { adminMenuRoutes }
