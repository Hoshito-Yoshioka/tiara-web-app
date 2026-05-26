import { Hono } from 'hono'
import type {
  MenuCategory,
  MenuCategoryWithItems,
  MenuCategoryWithItemsResponse,
  MenuItem,
  MenuCategoryResponse,
  MenuItemResponse,
} from '../types/menu'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'

const menuRoutes = new Hono()

function toMenuCategory(raw: MenuCategoryResponse): MenuCategory {
  return {
    id: raw.ID,
    name: raw.Name,
    description: raw.Description,
    sortOrder: raw.SortOrder,
  }
}

function toMenuItem(raw: MenuItemResponse): MenuItem {
  return {
    id: raw.ID,
    categoryId: raw.CategoryID,
    name: raw.Name,
    price: raw.Price,
    description: raw.Description,
    sortOrder: raw.SortOrder,
  }
}

function toMenuCategoryWithItems(raw: MenuCategoryWithItemsResponse): MenuCategoryWithItems {
  return {
    category: toMenuCategory(raw.Category),
    items: raw.Items ? raw.Items.map(toMenuItem) : [],
  }
}

/** GET /api/menus — カテゴリ＋アイテム一覧（公開） */
menuRoutes.get('/', async (c) => {
  const res = await fetch(`${BACKEND_URL}/api/v1/menus`)
  if (!res.ok) {
    return c.json({ error: 'Failed to fetch menus from backend' }, 502)
  }
  const data: MenuCategoryWithItemsResponse[] = await res.json()
  const menus: MenuCategoryWithItems[] = (data ?? []).map(toMenuCategoryWithItems)
  return c.json(menus)
})

export { menuRoutes, toMenuCategory, toMenuItem, toMenuCategoryWithItems }
