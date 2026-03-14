/** BFF から受け取るメニューカテゴリ型 */
export interface MenuCategory {
  id: string
  name: string
  description: string
  sortOrder: number
}

/** BFF から受け取るメニューアイテム型 */
export interface MenuItem {
  id: string
  categoryId: string
  name: string
  price: string
  description: string
  sortOrder: number
}

/** BFF から受け取るカテゴリ＋アイテムの集約型 */
export interface MenuCategoryWithItems {
  category: MenuCategory
  items: MenuItem[]
}
