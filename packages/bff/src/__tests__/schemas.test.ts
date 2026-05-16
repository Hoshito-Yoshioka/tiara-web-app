import { describe, it, expect } from 'vitest'
import {
  loginSchema,
  updateShopSchema,
  createStaffSchema,
  createMenuCategorySchema,
  createMenuItemSchema,
  reviewDraftSchema,
  createStaffAccountSchema,
  saveProfileDraftSchema,
  saveScheduleDraftSchema,
  setMainImageSchema,
  updateCropPositionSchema,
} from '../schemas'

describe('loginSchema', () => {
  it('valid input', () => {
    const result = loginSchema.safeParse({ username: 'admin', password: 'pass' })
    expect(result.success).toBe(true)
  })
  it('rejects empty username', () => {
    const result = loginSchema.safeParse({ username: '', password: 'pass' })
    expect(result.success).toBe(false)
  })
  it('rejects missing password', () => {
    const result = loginSchema.safeParse({ username: 'admin' })
    expect(result.success).toBe(false)
  })
})

describe('updateShopSchema', () => {
  it('valid input', () => {
    const result = updateShopSchema.safeParse({
      name: 'New Club TIARA',
      address: '北海道函館市',
      openingTime: '20:00',
      closingTime: '05:00',
    })
    expect(result.success).toBe(true)
  })
  it('rejects empty name', () => {
    const result = updateShopSchema.safeParse({
      name: '',
      address: '北海道函館市',
      openingTime: '20:00',
      closingTime: '05:00',
    })
    expect(result.success).toBe(false)
  })
})

describe('createStaffSchema', () => {
  it('valid input', () => {
    const result = createStaffSchema.safeParse({
      shopId: '550e8400-e29b-41d4-a716-446655440000',
      name: 'テスト',
      role: 'キャスト',
      bio: '',
      imageUrl: '',
      imageCropPosition: '50 50',
      sortOrder: 1,
      schedules: [{ dayOfWeek: 1, startTime: '20:00', endTime: '05:00' }],
    })
    expect(result.success).toBe(true)
  })
  it('rejects invalid UUID', () => {
    const result = createStaffSchema.safeParse({
      shopId: 'not-uuid',
      name: 'テスト',
      role: 'キャスト',
      bio: '',
      imageUrl: '',
      imageCropPosition: '',
      sortOrder: 1,
      schedules: [],
    })
    expect(result.success).toBe(false)
  })
  it('rejects invalid schedule time format', () => {
    const result = createStaffSchema.safeParse({
      shopId: '550e8400-e29b-41d4-a716-446655440000',
      name: 'テスト',
      role: 'キャスト',
      bio: '',
      imageUrl: '',
      imageCropPosition: '',
      sortOrder: 1,
      schedules: [{ dayOfWeek: 1, startTime: '8pm', endTime: '05:00' }],
    })
    expect(result.success).toBe(false)
  })
})

describe('createMenuCategorySchema', () => {
  it('rejects empty name', () => {
    const result = createMenuCategorySchema.safeParse({ name: '', description: '', sortOrder: 0 })
    expect(result.success).toBe(false)
  })
})

describe('createMenuItemSchema', () => {
  it('valid input', () => {
    const result = createMenuItemSchema.safeParse({
      categoryId: '550e8400-e29b-41d4-a716-446655440000',
      name: 'ビール',
      price: '500',
      description: '',
      sortOrder: 1,
    })
    expect(result.success).toBe(true)
  })
})

describe('reviewDraftSchema', () => {
  it('accepts approved', () => {
    const result = reviewDraftSchema.safeParse({ status: 'approved', adminComment: '' })
    expect(result.success).toBe(true)
  })
  it('rejects invalid status', () => {
    const result = reviewDraftSchema.safeParse({ status: 'pending', adminComment: '' })
    expect(result.success).toBe(false)
  })
})

describe('createStaffAccountSchema', () => {
  it('rejects non-UUID staffId', () => {
    const result = createStaffAccountSchema.safeParse({
      staffId: 'abc',
      username: 'user',
      password: 'pass',
    })
    expect(result.success).toBe(false)
  })
})

describe('saveProfileDraftSchema', () => {
  it('valid input', () => {
    const result = saveProfileDraftSchema.safeParse({
      name: 'テスト',
      role: 'キャスト',
      bio: '',
      imageUrl: '',
      imageCropPosition: '50 50',
    })
    expect(result.success).toBe(true)
  })
})

describe('saveScheduleDraftSchema', () => {
  it('valid input', () => {
    const result = saveScheduleDraftSchema.safeParse({
      items: [{ dayOfWeek: 0, startTime: '20:00', endTime: '05:00' }],
    })
    expect(result.success).toBe(true)
  })
  it('rejects dayOfWeek > 6', () => {
    const result = saveScheduleDraftSchema.safeParse({
      items: [{ dayOfWeek: 7, startTime: '20:00', endTime: '05:00' }],
    })
    expect(result.success).toBe(false)
  })
})

describe('setMainImageSchema', () => {
  it('rejects non-UUID', () => {
    const result = setMainImageSchema.safeParse({ imageId: 'abc' })
    expect(result.success).toBe(false)
  })
})

describe('updateCropPositionSchema', () => {
  it('rejects empty', () => {
    const result = updateCropPositionSchema.safeParse({ cropPosition: '' })
    expect(result.success).toBe(false)
  })
})
