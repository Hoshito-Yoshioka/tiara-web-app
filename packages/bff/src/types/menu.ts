/** Backend から返るメニューカテゴリの PascalCase レスポンス型 */
export interface MenuCategoryResponse {
  ID: string
  Name: string
  Description: string
  SortOrder: number
  CreatedAt: string
  UpdatedAt: string
}

/** Backend から返るメニューアイテムの PascalCase レスポンス型 */
export interface MenuItemResponse {
  ID: string
  CategoryID: string
  Name: string
  Price: string
  Description: string
  SortOrder: number
  CreatedAt: string
  UpdatedAt: string
}

/** Backend から返るカテゴリ＋アイテムの集約レスポンス型 */
export interface MenuCategoryWithItemsResponse {
  Category: MenuCategoryResponse
  Items: MenuItemResponse[] | null
}

/** BFF から Frontend へ返すメニューカテゴリ型 */
export interface MenuCategory {
  id: string
  name: string
  description: string
  sortOrder: number
}

/** BFF から Frontend へ返すメニューアイテム型 */
export interface MenuItem {
  id: string
  categoryId: string
  name: string
  price: string
  description: string
  sortOrder: number
}

/** BFF から Frontend へ返すカテゴリ＋アイテムの集約型 */
export interface MenuCategoryWithItems {
  category: MenuCategory
  items: MenuItem[]
}
