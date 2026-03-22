import { serve } from '@hono/node-server'
import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { shopRoutes } from './routes/shop'
import { staffRoutes } from './routes/staff'
import { scheduleRoutes } from './routes/schedule'
import { menuRoutes } from './routes/menu'
import { authRoutes } from './routes/auth'
import { adminShopRoutes } from './routes/admin/shop'
import { adminStaffRoutes } from './routes/admin/staff'
import { adminMenuRoutes } from './routes/admin/menu'

const app = new Hono()

// CORS: Frontend (Vite dev server) からのアクセスを許可
app.use(
  '/*',
  cors({
    origin: 'http://localhost:5173',
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE'],
    allowHeaders: ['Content-Type', 'Authorization'],
  })
)

// Health check
app.get('/', (c) => {
  return c.json({ message: 'Tiara BFF is running' })
})

// Public routes
app.route('/api/shops', shopRoutes)
app.route('/api/staffs', staffRoutes)
app.route('/api/schedules', scheduleRoutes)
app.route('/api/menus', menuRoutes)

// Auth routes
app.route('/api/auth', authRoutes)

// Admin routes (認証ミドルウェアは各ルートファイル内で適用)
app.route('/api/admin/shops', adminShopRoutes)
app.route('/api/admin/staffs', adminStaffRoutes)
app.route('/api/admin/menu', adminMenuRoutes)

// Static file proxy: アップロード画像を Backend から配信
const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:1323'
app.get('/uploads/*', async (c) => {
  const path = c.req.path
  const res = await fetch(`${BACKEND_URL}${path}`)
  if (!res.ok) {
    return c.notFound()
  }
  const contentType = res.headers.get('content-type') || 'application/octet-stream'
  const body = await res.arrayBuffer()
  return new Response(body, {
    status: 200,
    headers: {
      'Content-Type': contentType,
      'Cache-Control': 'public, max-age=86400',
    },
  })
})

// Node.js HTTP サーバーとして起動
const port = Number(process.env.PORT) || 3001

serve({ fetch: app.fetch, port }, (info) => {
  console.log(`🚀 Tiara BFF is running on http://localhost:${info.port}`)
})
