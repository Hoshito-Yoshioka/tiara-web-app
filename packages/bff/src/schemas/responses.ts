import { z } from '@hono/zod-openapi'

// --- Common ---

export const ErrorResponseSchema = z
  .object({
    error: z.string(),
  })
  .openapi('ErrorResponse')

export const AdminVerifyResponseSchema = z
  .object({
    status: z.string(),
  })
  .openapi('AdminVerifyResponse')

export const StaffVerifyResponseSchema = z
  .object({
    status: z.string(),
    staffId: z.string(),
  })
  .openapi('StaffVerifyResponse')

export const MessageResponseSchema = z
  .object({
    message: z.string(),
  })
  .openapi('MessageResponse')

// --- Shop ---

export const ShopSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    address: z.string(),
    openingTime: z.string(),
    closingTime: z.string(),
  })
  .openapi('Shop')

// --- Staff ---

export const StaffSchema = z
  .object({
    id: z.string(),
    shopId: z.string(),
    name: z.string(),
    role: z.string(),
    bio: z.string(),
    imageUrl: z.string(),
    externalScheduleUrl: z.string(),
    imageCropPosition: z.string(),
    sortOrder: z.number(),
  })
  .openapi('Staff')

export const StaffScheduleSchema = z
  .object({
    id: z.string(),
    staffId: z.string(),
    dayOfWeek: z.number(),
    startTime: z.string(),
    endTime: z.string(),
  })
  .openapi('StaffSchedule')

export const StaffImageSchema = z
  .object({
    id: z.string(),
    staffId: z.string(),
    imageUrl: z.string(),
    isMain: z.boolean(),
    sortOrder: z.number(),
    cropPosition: z.string(),
  })
  .openapi('StaffImage')

export const StaffWithSchedulesSchema = z
  .object({
    staff: StaffSchema,
    schedules: z.array(StaffScheduleSchema),
    images: z.array(StaffImageSchema),
  })
  .openapi('StaffWithSchedules')

export const PaginationSchema = z
  .object({
    page: z.number(),
    perPage: z.number(),
    totalCount: z.number(),
    totalPages: z.number(),
  })
  .openapi('Pagination')

export const PaginatedStaffsSchema = z
  .object({
    data: z.array(StaffSchema),
    pagination: PaginationSchema,
  })
  .openapi('PaginatedStaffs')

// --- Menu ---

export const MenuCategorySchema = z
  .object({
    id: z.string(),
    name: z.string(),
    description: z.string(),
    sortOrder: z.number(),
  })
  .openapi('MenuCategory')

export const MenuItemSchema = z
  .object({
    id: z.string(),
    categoryId: z.string(),
    name: z.string(),
    price: z.string(),
    description: z.string(),
    sortOrder: z.number(),
  })
  .openapi('MenuItem')

export const MenuCategoryWithItemsSchema = z
  .object({
    category: MenuCategorySchema,
    items: z.array(MenuItemSchema),
  })
  .openapi('MenuCategoryWithItems')

// --- Auth ---

export const TokenResponseSchema = z
  .object({
    token: z.string(),
  })
  .openapi('TokenResponse')

export const StaffTokenResponseSchema = z
  .object({
    token: z.string(),
    staffId: z.string(),
  })
  .openapi('StaffTokenResponse')

// --- Portal: Profile Draft ---

export const StaffImageForDraftSchema = z
  .object({
    id: z.string(),
    staffId: z.string(),
    imageUrl: z.string(),
    isMain: z.boolean(),
    sortOrder: z.number(),
    cropPosition: z.string(),
  })
  .openapi('StaffImageForDraft')

export const ProfileDraftResponseSchema = z
  .object({
    id: z.string().optional(),
    staffId: z.string(),
    name: z.string(),
    role: z.string(),
    bio: z.string(),
    imageUrl: z.string(),
    externalScheduleUrl: z.string(),
    imageCropPosition: z.string(),
    status: z.string(),
    adminComment: z.string(),
    submittedAt: z.string().optional(),
    reviewedAt: z.string().optional(),
    createdAt: z.string().optional(),
    updatedAt: z.string().optional(),
    images: z.array(StaffImageForDraftSchema).optional(),
  })
  .openapi('ProfileDraftResponse')

// --- Portal: Schedule Draft ---

export const ScheduleDraftItemSchema = z
  .object({
    id: z.string().optional(),
    dayOfWeek: z.number(),
    startTime: z.string(),
    endTime: z.string(),
  })
  .openapi('ScheduleDraftItem')

export const ScheduleDraftResponseSchema = z
  .object({
    id: z.string().optional(),
    staffId: z.string(),
    status: z.string(),
    adminComment: z.string(),
    submittedAt: z.string().optional(),
    reviewedAt: z.string().optional(),
    createdAt: z.string().optional(),
    updatedAt: z.string().optional(),
    items: z.array(ScheduleDraftItemSchema),
  })
  .openapi('ScheduleDraftResponse')

// --- Staff Account ---

export const StaffAccountResponseSchema = z
  .object({
    id: z.string(),
    staffId: z.string(),
    username: z.string(),
    createdAt: z.string(),
    updatedAt: z.string(),
  })
  .openapi('StaffAccountResponse')
