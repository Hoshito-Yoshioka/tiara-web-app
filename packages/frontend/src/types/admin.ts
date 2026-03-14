/** 管理者ログインフォームの入力型 */
export interface LoginInput {
  username: string
  password: string
}

/** 管理者ログインレスポンス型 */
export interface LoginResponse {
  token: string
}

/** 店舗更新フォームの入力型 */
export interface UpdateShopInput {
  name: string
  address: string
  openingTime: string
  closingTime: string
}

/** スタッフ作成フォームの入力型 */
export interface CreateStaffInput {
  shopId: string
  name: string
  role: string
  bio: string
  imageUrl: string
  sortOrder: number
  schedules: ScheduleInput[]
}

/** スタッフ更新フォームの入力型 */
export interface UpdateStaffInput {
  name: string
  role: string
  bio: string
  imageUrl: string
  sortOrder: number
  schedules: ScheduleInput[]
}

/** スケジュール入力型 */
export interface ScheduleInput {
  dayOfWeek: number
  startTime: string
  endTime: string
}

/** メニューカテゴリ作成フォームの入力型 */
export interface CreateMenuCategoryInput {
  name: string
  description: string
  sortOrder: number
}

/** メニューカテゴリ更新フォームの入力型 */
export interface UpdateMenuCategoryInput {
  name: string
  description: string
  sortOrder: number
}

/** メニューアイテム作成フォームの入力型 */
export interface CreateMenuItemInput {
  categoryId: string
  name: string
  price: string
  description: string
  sortOrder: number
}

/** メニューアイテム更新フォームの入力型 */
export interface UpdateMenuItemInput {
  name: string
  price: string
  description: string
  sortOrder: number
}
