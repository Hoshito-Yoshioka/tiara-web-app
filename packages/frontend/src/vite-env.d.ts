/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  /** 管理画面のベースパス。本番では推測困難な文字列に変更する。例: /mgmt */
  readonly VITE_ADMIN_BASE_PATH: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
