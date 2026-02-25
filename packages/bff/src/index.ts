import { Hono } from 'hono'

const app = new Hono()

app.get('/', (c) => {
  return c.json({ message: 'Hello from Tiara BFF!' })
})

export default app
