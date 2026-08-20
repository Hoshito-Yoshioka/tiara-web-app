import path from 'path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import type { ViteSSGOptions } from 'vite-ssg'

// SSG プリレンダリング時にスタッフ詳細ページの URL を API から取得する。
// scripts/generate-sitemap.mjs と同じ環境変数（SITEMAP_API_BASE_URL）を使用する。
async function fetchStaffPaths(): Promise<string[]> {
  const apiBaseUrl = (
    process.env.SITEMAP_API_BASE_URL ||
    process.env.VITE_API_BASE_URL ||
    'http://localhost:3001'
  ).replace(/\/+$/, '')

  try {
    const response = await fetch(`${apiBaseUrl}/api/v1/staffs`)
    if (!response.ok) {
      console.warn(`[ssg] staff fetch failed: ${response.status} ${response.statusText}`)
      return []
    }
    const body = await response.json()
    const staffs = Array.isArray(body) ? body : Array.isArray(body?.data) ? body.data : []
    return staffs
      .map((staff: { id?: unknown }) => (typeof staff?.id === 'string' ? staff.id : null))
      .filter((id: string | null): id is string => !!id)
      .map((id: string) => `/staff/${encodeURIComponent(id)}`)
  } catch (error) {
    console.warn(`[ssg] staff fetch error: ${error instanceof Error ? error.message : error}`)
    return []
  }
}

const ssgOptions: ViteSSGOptions = {
  formatting: 'minify',
  // Nginx の try_files ($uri/ + index.html) で配信できるよう /price/index.html 形式で出力する
  dirStyle: 'nested',
  // 公開ページのみプリレンダリングする（管理画面・スタッフポータルは除外）
  async includedRoutes(paths) {
    const adminBase = process.env.VITE_ADMIN_BASE_PATH || '/admin'
    const publicPaths = paths.filter(
      (p) => !p.includes(':') && !p.startsWith(adminBase) && !p.startsWith('/mypage')
    )
    const staffPaths = await fetchStaffPaths()
    return [...publicPaths, ...staffPaths]
  },
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  ssgOptions,
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  test: {
    globals: true,
    environment: 'happy-dom',
  },
  server: {
    proxy: {
      // アップロード画像を BFF 経由で Backend から配信
      '/uploads': {
        target: 'http://localhost:3001',
        changeOrigin: true,
      },
    },
  },
})
