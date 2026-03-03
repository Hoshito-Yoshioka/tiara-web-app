/** Shop のレスポンス型（Backend API のレスポンスに対応） */
export interface ShopResponse {
  ID: string
  Name: string
  Address: string
  OpeningTime: string
  ClosingTime: string
  CreatedAt: string
  UpdatedAt: string
}

/** BFF から Frontend へ返す Shop 型 */
export interface Shop {
  id: string
  name: string
  address: string
  openingTime: string
  closingTime: string
}
