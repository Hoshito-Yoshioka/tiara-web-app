import { z } from 'zod'

// --- Auth ---

export const loginSchema = z.object({
  username: z.string().min(1, 'ユーザー名は必須です'),
  password: z.string().min(1, 'パスワードは必須です'),
})

// --- Shop ---

export const updateShopSchema = z.object({
  name: z.string().min(1, '店舗名は必須です'),
  address: z.string().min(1, '住所は必須です'),
  openingTime: z.string().min(1, '開店時間は必須です'),
  closingTime: z.string().min(1, '閉店時間は必須です'),
})

// --- Staff ---

const scheduleInputSchema = z.object({
  dayOfWeek: z.number().int().min(0).max(6),
  startTime: z.string().regex(/^\d{2}:\d{2}$/, '時刻は HH:MM 形式で指定してください'),
  endTime: z.string().regex(/^\d{2}:\d{2}$/, '時刻は HH:MM 形式で指定してください'),
})

export const createStaffSchema = z.object({
  shopId: z.string().uuid('有効な UUID を指定してください'),
  name: z.string().min(1, '名前は必須です'),
  role: z.string().min(1, '役職は必須です'),
  bio: z.string(),
  imageUrl: z.string(),
  imageCropPosition: z.string(),
  sortOrder: z.number().int(),
  schedules: z.array(scheduleInputSchema),
})

export const updateStaffSchema = z.object({
  name: z.string().min(1, '名前は必須です'),
  role: z.string().min(1, '役職は必須です'),
  bio: z.string(),
  imageUrl: z.string(),
  imageCropPosition: z.string(),
  sortOrder: z.number().int(),
  schedules: z.array(scheduleInputSchema),
})

// --- Menu ---

export const createMenuCategorySchema = z.object({
  name: z.string().min(1, 'カテゴリ名は必須です'),
  description: z.string(),
  sortOrder: z.number().int(),
})

export const updateMenuCategorySchema = z.object({
  name: z.string().min(1, 'カテゴリ名は必須です'),
  description: z.string(),
  sortOrder: z.number().int(),
})

export const createMenuItemSchema = z.object({
  categoryId: z.string().uuid('有効な UUID を指定してください'),
  name: z.string().min(1, 'アイテム名は必須です'),
  price: z.string().min(1, '価格は必須です'),
  description: z.string(),
  sortOrder: z.number().int(),
})

export const updateMenuItemSchema = z.object({
  name: z.string().min(1, 'アイテム名は必須です'),
  price: z.string().min(1, '価格は必須です'),
  description: z.string(),
  sortOrder: z.number().int(),
})

// --- Staff Portal ---

export const saveProfileDraftSchema = z.object({
  name: z.string().min(1, '名前は必須です'),
  role: z.string().min(1, '役職は必須です'),
  bio: z.string(),
  imageUrl: z.string(),
  imageCropPosition: z.string(),
})

const scheduleDraftItemSchema = z.object({
  dayOfWeek: z.number().int().min(0).max(6),
  startTime: z.string().regex(/^\d{2}:\d{2}$/, '時刻は HH:MM 形式で指定してください'),
  endTime: z.string().regex(/^\d{2}:\d{2}$/, '時刻は HH:MM 形式で指定してください'),
})

export const saveScheduleDraftSchema = z.object({
  items: z.array(scheduleDraftItemSchema),
})

// --- Admin Review ---

export const reviewDraftSchema = z.object({
  status: z.enum(['approved', 'rejected'], {
    errorMap: () => ({ message: 'ステータスは approved または rejected を指定してください' }),
  }),
  adminComment: z.string(),
})

// --- Staff Account ---

export const createStaffAccountSchema = z.object({
  staffId: z.string().uuid('有効な UUID を指定してください'),
  username: z.string().min(1, 'ユーザー名は必須です'),
  password: z.string().min(1, 'パスワードは必須です'),
})

export const updateStaffAccountSchema = z.object({
  username: z.string().min(1, 'ユーザー名は必須です'),
  password: z.string(),
})

// --- Image ---

export const setMainImageSchema = z.object({
  imageId: z.string().uuid('有効な UUID を指定してください'),
})

export const updateCropPositionSchema = z.object({
  cropPosition: z.string().min(1, 'クロップ位置は必須です'),
})
