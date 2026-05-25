/** Backend API のレスポンス型（handler DTOにより camelCase + "HH:MM" 形式） */
export interface ShopResponse {
  id: string
  name: string
  address: string
  openingTime: string
  closingTime: string
}

/** BFF から Frontend へ返す Shop 型 */
export interface Shop {
  id: string
  name: string
  address: string
  openingTime: string
  closingTime: string
}
