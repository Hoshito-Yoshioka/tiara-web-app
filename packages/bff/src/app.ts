import { OpenAPIHono } from '@hono/zod-openapi'
import { swaggerUI } from '@hono/swagger-ui'
import { cors } from 'hono/cors'
import { secureHeaders } from 'hono/secure-headers'
import { shopRoutes } from './routes/shop'
import { staffRoutes } from './routes/staff'
import { scheduleRoutes } from './routes/schedule'
import { menuRoutes } from './routes/menu'
import { authRoutes } from './routes/auth'
import { staffAuthRoutes } from './routes/staff-auth'
import { portalRoutes } from './routes/portal'
import { adminShopRoutes } from './routes/admin/shop'
import { adminStaffRoutes } from './routes/admin/staff'
import { adminMenuRoutes } from './routes/admin/menu'
import { adminReviewRoutes } from './routes/admin/review'
import { adminAccountRoutes } from './routes/admin/account'

const redocCsp = [
  "default-src 'self'",
  "script-src 'self' https://cdn.redocly.com 'unsafe-inline'",
  "style-src 'self' 'unsafe-inline' https://cdn.redocly.com",
  "img-src 'self' data: https://cdn.redocly.com",
  "font-src 'self' data: https://cdn.redocly.com",
  "connect-src 'self' https://cdn.redocly.com",
  "frame-ancestors 'self'",
].join('; ')

function renderRedocPage(specUrl: string) {
  return `<!doctype html>
<html lang="ja">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="color-scheme" content="dark" />
    <title>Tiara BFF API Reference</title>
    <style>
      :root {
        color-scheme: dark;
        --bg: #060709;
        --bg-soft: rgba(12, 14, 20, 0.82);
        --border: rgba(255, 255, 255, 0.08);
        --text: #f4efe3;
        --muted: rgba(244, 239, 227, 0.72);
        --accent: #c8a96a;
      }

      * {
        box-sizing: border-box;
      }

      html,
      body {
        margin: 0;
        min-height: 100%;
        background:
          radial-gradient(circle at top left, rgba(200, 169, 106, 0.18), transparent 32%),
          radial-gradient(circle at top right, rgba(92, 122, 255, 0.12), transparent 28%),
          linear-gradient(180deg, #0b0c10 0%, #050608 100%);
        color: var(--text);
        font-family: "Segoe UI", system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
      }

      body::before {
        content: '';
        position: fixed;
        inset: 0;
        pointer-events: none;
        background-image: linear-gradient(rgba(255, 255, 255, 0.02) 1px, transparent 1px),
          linear-gradient(90deg, rgba(255, 255, 255, 0.02) 1px, transparent 1px);
        background-size: 48px 48px;
        mask-image: linear-gradient(180deg, rgba(0, 0, 0, 0.7), transparent 92%);
      }

      .page {
        position: relative;
        min-height: 100vh;
        padding: 40px 24px 32px;
      }

      .shell {
        width: min(1280px, 100%);
        margin: 0 auto;
      }

      .hero {
        display: flex;
        flex-wrap: wrap;
        align-items: end;
        justify-content: space-between;
        gap: 16px;
        margin-bottom: 20px;
      }

      .eyebrow {
        margin: 0 0 8px;
        color: var(--accent);
        font-size: 0.75rem;
        letter-spacing: 0.28em;
        text-transform: uppercase;
      }

      h1 {
        margin: 0;
        font-family: Georgia, 'Times New Roman', serif;
        font-size: clamp(2rem, 4vw, 3.5rem);
        line-height: 1.02;
        font-weight: 700;
      }

      p {
        margin: 10px 0 0;
        max-width: 60ch;
        color: var(--muted);
        line-height: 1.7;
      }

      .links {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
      }

      .links a {
        color: var(--text);
        text-decoration: none;
        border: 1px solid var(--border);
        background: rgba(255, 255, 255, 0.04);
        border-radius: 999px;
        padding: 10px 14px;
        backdrop-filter: blur(12px);
        transition: transform 0.18s ease, background 0.18s ease, border-color 0.18s ease;
      }

      .links a:hover {
        transform: translateY(-1px);
        background: rgba(255, 255, 255, 0.08);
        border-color: rgba(200, 169, 106, 0.45);
      }

      .panel {
        border: 1px solid var(--border);
        border-radius: 24px;
        overflow: hidden;
        background: var(--bg-soft);
        box-shadow: 0 24px 72px rgba(0, 0, 0, 0.45);
      }

      #redoc {
        min-height: calc(100vh - 180px);
      }
    </style>
    <script src="https://cdn.redocly.com/redoc/v2.5.1/bundles/redoc.standalone.js"></script>
  </head>
  <body>
    <main class="page">
      <div class="shell">
        <section class="hero">
          <div>
            <p class="eyebrow">Tiara BFF API</p>
            <h1>Readable API reference</h1>
            <p>
              Swagger UI より読みやすいレイアウトで、エンドポイントとスキーマを確認できます。
            </p>
          </div>
          <div class="links">
            <a href="/api/v1/docs">Swagger UI</a>
            <a href="/api/v1/doc">OpenAPI JSON</a>
          </div>
        </section>

        <section class="panel">
          <div id="redoc"></div>
        </section>
      </div>
    </main>

    <script>
      Redoc.init(${JSON.stringify(specUrl)}, {
        scrollYOffset: 24,
        hideDownloadButton: false,
        theme: {
          colors: {
            primary: {
              main: '#c8a96a',
            },
          },
        },
      }, document.getElementById('redoc'))
    </script>
  </body>
</html>`
}

export const app = new OpenAPIHono()

// Security headers
app.use('/*', secureHeaders())

// CORS: Frontend からのアクセスを許可（環境変数で設定可能）
const CORS_ORIGIN = process.env.CORS_ORIGIN || 'http://localhost:5173'
app.use(
  '/*',
  cors({
    origin: CORS_ORIGIN,
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE'],
    allowHeaders: ['Content-Type', 'Authorization'],
  })
)

// Health check
app.get('/', (c) => {
  return c.json({ message: 'Tiara BFF is running' })
})

// Public routes
app.route('/api/v1/shops', shopRoutes)
app.route('/api/v1/staffs', staffRoutes)
app.route('/api/v1/schedules', scheduleRoutes)
app.route('/api/v1/menus', menuRoutes)

// Auth routes
app.route('/api/v1/auth', authRoutes)
app.route('/api/v1/staff-auth', staffAuthRoutes)

// Staff Portal routes (スタッフ専用)
app.route('/api/v1/portal', portalRoutes)

// Admin routes (認証ミドルウェアは各ルートファイル内で適用)
app.route('/api/v1/admin/shops', adminShopRoutes)
app.route('/api/v1/admin/staffs', adminStaffRoutes)
app.route('/api/v1/admin/menu', adminMenuRoutes)
app.route('/api/v1/admin/reviews', adminReviewRoutes)
app.route('/api/v1/admin/staff-accounts', adminAccountRoutes)

// OpenAPI: Bearer 認証スキームを登録
app.openAPIRegistry.registerComponent('securitySchemes', 'Bearer', {
  type: 'http',
  scheme: 'bearer',
  bearerFormat: 'JWT',
})

// OpenAPI JSON ドキュメント
app.doc('/api/v1/doc', {
  openapi: '3.0.0',
  info: {
    title: 'Tiara BFF API',
    version: '1.0.0',
    description: 'Bar Tiara の BFF API 仕様書',
  },
})

// より読みやすい API ドキュメント
app.get('/api/v1/redoc', (c) => {
  c.header('Content-Security-Policy', redocCsp)
  return c.html(renderRedocPage('/api/v1/doc'))
})

// Swagger UI
app.get('/api/v1/docs', swaggerUI({ url: '/api/v1/doc' }))

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
