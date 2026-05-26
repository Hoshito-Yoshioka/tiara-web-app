import { describe, it, expect, vi, beforeEach } from 'vitest'
import { OpenAPIHono } from '@hono/zod-openapi'
import { staffRoutes } from '../routes/staff'
import { menuRoutes } from '../routes/menu'
import { shopRoutes } from '../routes/shop'
import { scheduleRoutes } from '../routes/schedule'

// グローバル fetch をモック化
const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

beforeEach(() => {
  mockFetch.mockReset()
})

// --- Staff Routes ---

describe('staffRoutes', () => {
  const app = new OpenAPIHono()
  app.route('/api/v1/staffs', staffRoutes)

  describe('GET /api/v1/staffs', () => {
    it('PascalCase → camelCase 変換を行う', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            ID: 'staff-1',
            ShopID: 'shop-1',
            Name: 'Yuki',
            Role: 'Cast',
            Bio: 'bio text',
            ImageURL: '/img/yuki.jpg',
            ImageCropPosition: '50 50',
            SortOrder: 1,
          },
        ],
      })

      const res = await app.request('/api/v1/staffs')
      expect(res.status).toBe(200)

      const body = await res.json()
      expect(body).toHaveLength(1)
      expect(body[0]).toEqual({
        id: 'staff-1',
        shopId: 'shop-1',
        name: 'Yuki',
        role: 'Cast',
        bio: 'bio text',
        imageUrl: '/img/yuki.jpg',
        imageCropPosition: '50 50',
        sortOrder: 1,
      })
    })

    it('Backend エラー時 → 502', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 500 })

      const res = await app.request('/api/v1/staffs')
      expect(res.status).toBe(502)
    })
  })

  describe('GET /api/v1/staffs/:id', () => {
    it('StaffWithSchedules の PascalCase → camelCase 変換を行う', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          Staff: {
            ID: 'staff-1',
            ShopID: 'shop-1',
            Name: 'Yuki',
            Role: 'Cast',
            Bio: '',
            ImageURL: '/img.jpg',
            ImageCropPosition: '50 50',
            SortOrder: 1,
          },
          Schedules: [
            { ID: 'sch-1', StaffID: 'staff-1', DayOfWeek: 1, StartTime: '20:00', EndTime: '02:00' },
          ],
          Images: [
            {
              ID: 'img-1',
              StaffID: 'staff-1',
              ImageURL: '/img.jpg',
              IsMain: true,
              SortOrder: 0,
              CropPosition: '30 70',
            },
          ],
        }),
      })

      const res = await app.request('/api/v1/staffs/staff-1')
      expect(res.status).toBe(200)

      const body = await res.json()
      expect(body.staff.id).toBe('staff-1')
      expect(body.staff.shopId).toBe('shop-1')
      expect(body.schedules[0].staffId).toBe('staff-1')
      expect(body.schedules[0].dayOfWeek).toBe(1)
      expect(body.images[0].imageUrl).toBe('/img.jpg')
      expect(body.images[0].isMain).toBe(true)
      expect(body.images[0].cropPosition).toBe('30 70')
    })

    it('Images が null の場合 → 空配列', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          Staff: {
            ID: 's1',
            ShopID: 'sh1',
            Name: 'A',
            Role: '',
            Bio: '',
            ImageURL: '',
            ImageCropPosition: '',
            SortOrder: 0,
          },
          Schedules: null,
          Images: null,
        }),
      })

      const res = await app.request('/api/v1/staffs/s1')
      expect(res.status).toBe(200)

      const body = await res.json()
      expect(body.schedules).toEqual([])
      expect(body.images).toEqual([])
    })

    it('CropPosition 未設定 → デフォルト "50 50"', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          Staff: {
            ID: 's1',
            ShopID: 'sh1',
            Name: 'A',
            Role: '',
            Bio: '',
            ImageURL: '',
            ImageCropPosition: '',
            SortOrder: 0,
          },
          Schedules: [],
          Images: [
            {
              ID: 'img-1',
              StaffID: 's1',
              ImageURL: '/img.jpg',
              IsMain: false,
              SortOrder: 0,
              CropPosition: undefined,
            },
          ],
        }),
      })

      const res = await app.request('/api/v1/staffs/s1')
      const body = await res.json()
      expect(body.images[0].cropPosition).toBe('50 50')
    })

    it('Backend 404 → 404', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 404 })

      const res = await app.request('/api/v1/staffs/nonexistent')
      expect(res.status).toBe(404)
    })
  })
})

// --- Menu Routes ---

describe('menuRoutes', () => {
  const app = new OpenAPIHono()
  app.route('/api/v1/menus', menuRoutes)

  describe('GET /api/v1/menus', () => {
    it('カテゴリ＋アイテムの PascalCase → camelCase 変換を行う', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            Category: { ID: 'cat-1', Name: 'Whisky', Description: 'Single malt', SortOrder: 1 },
            Items: [
              {
                ID: 'item-1',
                CategoryID: 'cat-1',
                Name: 'Macallan 12',
                Price: '¥1,500',
                Description: '',
                SortOrder: 1,
              },
            ],
          },
        ],
      })

      const res = await app.request('/api/v1/menus')
      expect(res.status).toBe(200)

      const body = await res.json()
      expect(body).toHaveLength(1)
      expect(body[0].category).toEqual({
        id: 'cat-1',
        name: 'Whisky',
        description: 'Single malt',
        sortOrder: 1,
      })
      expect(body[0].items[0]).toEqual({
        id: 'item-1',
        categoryId: 'cat-1',
        name: 'Macallan 12',
        price: '¥1,500',
        description: '',
        sortOrder: 1,
      })
    })

    it('Items が null の場合 → 空配列', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            Category: { ID: 'cat-1', Name: 'Beer', Description: '', SortOrder: 1 },
            Items: null,
          },
        ],
      })

      const res = await app.request('/api/v1/menus')
      const body = await res.json()
      expect(body[0].items).toEqual([])
    })

    it('Backend レスポンスが null → 空配列', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => null,
      })

      const res = await app.request('/api/v1/menus')
      const body = await res.json()
      expect(body).toEqual([])
    })

    it('Backend エラー時 → 502', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 500 })

      const res = await app.request('/api/v1/menus')
      expect(res.status).toBe(502)
    })
  })
})

// --- Shop Routes ---

describe('shopRoutes', () => {
  const app = new OpenAPIHono()
  app.route('/api/v1/shops', shopRoutes)

  describe('GET /api/v1/shops', () => {
    it('Backend レスポンスのフィールド選択を行う', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            id: 'shop-1',
            name: 'New Club TIARA',
            address: '北海道函館市本町１−２８',
            openingTime: '20:00',
            closingTime: '05:00',
          },
        ],
      })

      const res = await app.request('/api/v1/shops')
      expect(res.status).toBe(200)

      const body = await res.json()
      expect(body[0]).toEqual({
        id: 'shop-1',
        name: 'New Club TIARA',
        address: '北海道函館市本町１−２８',
        openingTime: '20:00',
        closingTime: '05:00',
      })
    })

    it('Backend エラー時 → 502', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 500 })

      const res = await app.request('/api/v1/shops')
      expect(res.status).toBe(502)
    })
  })

  describe('GET /api/v1/shops/:id', () => {
    it('単一店舗のフィールド選択を行う', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'shop-1',
          name: 'New Club TIARA',
          address: '函館',
          openingTime: '20:00',
          closingTime: '05:00',
        }),
      })

      const res = await app.request('/api/v1/shops/shop-1')
      expect(res.status).toBe(200)

      const body = await res.json()
      expect(body.id).toBe('shop-1')
      expect(body.openingTime).toBe('20:00')
    })

    it('Backend 404 → 404', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 404 })

      const res = await app.request('/api/v1/shops/nonexistent')
      expect(res.status).toBe(404)
    })
  })
})

// --- Schedule Routes ---

describe('scheduleRoutes', () => {
  const app = new OpenAPIHono()
  app.route('/api/v1/schedules', scheduleRoutes)

  describe('GET /api/v1/schedules', () => {
    it('全スタッフスケジュールの PascalCase → camelCase 変換を行う', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            Staff: {
              ID: 's1',
              ShopID: 'sh1',
              Name: 'Yuki',
              Role: 'Cast',
              Bio: '',
              ImageURL: '',
              ImageCropPosition: '',
              SortOrder: 1,
            },
            Schedules: [
              { ID: 'sch-1', StaffID: 's1', DayOfWeek: 5, StartTime: '21:00', EndTime: '03:00' },
            ],
            Images: [
              {
                ID: 'img-1',
                StaffID: 's1',
                ImageURL: '/img.jpg',
                IsMain: true,
                SortOrder: 0,
                CropPosition: '50 50',
              },
            ],
          },
        ],
      })

      const res = await app.request('/api/v1/schedules')
      expect(res.status).toBe(200)

      const body = await res.json()
      expect(body).toHaveLength(1)
      expect(body[0].staff.name).toBe('Yuki')
      expect(body[0].schedules[0].dayOfWeek).toBe(5)
      expect(body[0].images[0].isMain).toBe(true)
    })

    it('Backend エラー時 → 502', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 500 })

      const res = await app.request('/api/v1/schedules')
      expect(res.status).toBe(502)
    })
  })
})
