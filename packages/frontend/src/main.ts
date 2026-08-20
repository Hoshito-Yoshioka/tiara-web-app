import { ViteSSG } from 'vite-ssg'
import { createPinia } from 'pinia'
import { MotionPlugin } from '@vueuse/motion'
import './style.css'
import App from './App.vue'
import { routes, scrollBehavior, setupRouterGuards } from './router'

// vite-ssg のエントリポイント。
// ビルド時（vite-ssg build）は公開ルートをプリレンダリングして静的 HTML を生成し、
// ブラウザではその HTML にハイドレートする。開発時は従来どおり SPA として動作する。
export const createApp = ViteSSG(
  App,
  {
    routes,
    scrollBehavior,
    base: import.meta.env.BASE_URL,
  },
  ({ app, router }) => {
    app.use(createPinia())
    app.use(MotionPlugin)
    setupRouterGuards(router)
  }
)
