import { writeFile } from 'node:fs/promises'
import { fileURLToPath } from 'node:url'

const SITE_URL = (process.env.SITE_URL || 'https://tiara-hakodate.com').replace(/\/+$/, '')
const API_BASE_URL = (
  process.env.SITEMAP_API_BASE_URL ||
  process.env.VITE_API_BASE_URL ||
  'http://localhost:3001'
).replace(/\/+$/, '')

const outputPath = fileURLToPath(new URL('../public/sitemap.xml', import.meta.url))
const now = new Date().toISOString().slice(0, 10)

const fixedRoutes = [
  { path: '/', changefreq: 'weekly', priority: '1.0' },
  { path: '/shop', changefreq: 'weekly', priority: '0.8' },
  { path: '/staff', changefreq: 'daily', priority: '0.9' },
  { path: '/schedule', changefreq: 'daily', priority: '0.8' },
  { path: '/price', changefreq: 'monthly', priority: '0.7' },
  { path: '/access', changefreq: 'monthly', priority: '0.7' },
]

function escapeXml(value) {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&apos;')
}

async function fetchStaffPaths() {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/staffs`)
    if (!response.ok) {
      console.warn(`[sitemap] staff fetch failed: ${response.status} ${response.statusText}`)
      return []
    }

    const body = await response.json()
    const staffs = Array.isArray(body) ? body : Array.isArray(body?.data) ? body.data : null
    if (!staffs) {
      console.warn('[sitemap] unexpected staff response shape')
      return []
    }

    return staffs
      .map((staff) => (typeof staff?.id === 'string' ? staff.id : null))
      .filter((id) => id)
      .map((id) => ({
        path: `/staff/${encodeURIComponent(id)}`,
        changefreq: 'weekly',
        priority: '0.7',
      }))
  } catch (error) {
    console.warn(
      `[sitemap] staff fetch error: ${error instanceof Error ? error.message : String(error)}`
    )
    return []
  }
}

function toUrlNode(route) {
  const loc = escapeXml(`${SITE_URL}${route.path}`)
  return [
    '  <url>',
    `    <loc>${loc}</loc>`,
    `    <lastmod>${now}</lastmod>`,
    `    <changefreq>${route.changefreq}</changefreq>`,
    `    <priority>${route.priority}</priority>`,
    '  </url>',
  ].join('\n')
}

async function main() {
  const dynamicStaffRoutes = await fetchStaffPaths()
  const routes = [...fixedRoutes, ...dynamicStaffRoutes]

  const xml = [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
    ...routes.map(toUrlNode),
    '</urlset>',
    '',
  ].join('\n')

  await writeFile(outputPath, xml, 'utf8')
  console.log(`[sitemap] generated ${routes.length} urls -> ${outputPath}`)
}

main().catch((error) => {
  console.error('[sitemap] generation failed', error)
  process.exitCode = 1
})
