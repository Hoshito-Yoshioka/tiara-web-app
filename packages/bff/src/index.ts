import { serve } from '@hono/node-server'
import { app } from './app'

// Node.js HTTP サーバーとして起動
const port = Number(process.env.PORT) || 3001

serve({ fetch: app.fetch, port }, (info) => {
  console.log(`🚀 Tiara BFF is running on http://localhost:${info.port}`)
})
