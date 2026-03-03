import { serve } from '@hono/node-server'
import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { shopRoutes } from './routes/shop'
import { staffRoutes } from './routes/staff'
import { scheduleRoutes } from './routes/schedule'

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

// Shop routes — /api/shops 配下にマウント
app.route('/api/shops', shopRoutes)

// Staff routes — /api/staffs 配下にマウント
app.route('/api/staffs', staffRoutes)

// Schedule routes — /api/schedules 配下にマウント
app.route('/api/schedules', scheduleRoutes)

// Node.js HTTP サーバーとして起動
const port = Number(process.env.PORT) || 3001

serve({ fetch: app.fetch, port }, (info) => {
  console.log(`🚀 Tiara BFF is running on http://localhost:${info.port}`)
})
