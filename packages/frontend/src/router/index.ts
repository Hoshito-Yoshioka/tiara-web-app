import type { Router, RouteRecordRaw, RouterScrollBehavior } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useStaffAuthStore } from '@/stores/staffAuth'

// 管理画面のベースパス。本番では VITE_ADMIN_BASE_PATH 環境変数で変更し、URLからの推測を困難にする。
const adminBase = import.meta.env.VITE_ADMIN_BASE_PATH || '/admin'

// ページ遷移のたびにスクロール位置をトップへリセット
export const scrollBehavior: RouterScrollBehavior = () => ({ top: 0, behavior: 'smooth' })

// ルーターの生成は vite-ssg（main.ts の ViteSSG）が行うため、ここではルート定義のみ持つ
export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue'),
  },
  {
    path: '/shop',
    name: 'shop',
    component: () => import('@/views/ShopView.vue'),
  },
  {
    path: '/staff',
    name: 'staff',
    component: () => import('@/views/StaffView.vue'),
  },
  {
    path: '/staff/:id',
    name: 'staff-detail',
    component: () => import('@/views/StaffDetailView.vue'),
  },
  {
    path: '/schedule',
    name: 'schedule',
    component: () => import('@/views/ScheduleView.vue'),
  },
  {
    path: '/price',
    name: 'price',
    component: () => import('@/views/PriceView.vue'),
  },
  {
    path: '/access',
    name: 'access',
    component: () => import('@/views/AccessView.vue'),
  },
  // ==============================
  // Admin Routes
  // layout: 'admin' で公開ページと異なるレイアウトを使用
  // requiresAuth: true で認証ガードを適用
  // ==============================
  {
    path: adminBase,
    redirect: `${adminBase}/login`,
  },
  {
    path: `${adminBase}/login`,
    name: 'admin-login',
    component: () => import('@/views/admin/AdminLoginView.vue'),
    meta: { layout: 'admin' },
  },
  {
    path: `${adminBase}/shop`,
    name: 'admin-shop-edit',
    component: () => import('@/views/admin/AdminShopEditView.vue'),
    meta: { layout: 'admin', requiresAuth: true },
  },
  {
    path: `${adminBase}/staffs`,
    name: 'admin-staff-list',
    component: () => import('@/views/admin/AdminStaffListView.vue'),
    meta: { layout: 'admin', requiresAuth: true },
  },
  {
    path: `${adminBase}/staffs/new`,
    name: 'admin-staff-new',
    component: () => import('@/views/admin/AdminStaffEditView.vue'),
    meta: { layout: 'admin', requiresAuth: true },
  },
  {
    path: `${adminBase}/staffs/:id/edit`,
    name: 'admin-staff-edit',
    component: () => import('@/views/admin/AdminStaffEditView.vue'),
    meta: { layout: 'admin', requiresAuth: true },
  },
  {
    path: `${adminBase}/menu`,
    name: 'admin-menu-edit',
    component: () => import('@/views/admin/AdminMenuView.vue'),
    meta: { layout: 'admin', requiresAuth: true },
  },
  {
    path: `${adminBase}/profile-reviews`,
    name: 'admin-profile-reviews',
    component: () => import('@/views/admin/AdminProfileReviewView.vue'),
    meta: { layout: 'admin', requiresAuth: true },
  },
  {
    path: `${adminBase}/schedule-reviews`,
    name: 'admin-schedule-reviews',
    component: () => import('@/views/admin/AdminScheduleReviewView.vue'),
    meta: { layout: 'admin', requiresAuth: true },
  },
  // ==============================
  // Staff Portal Routes
  // layout: 'portal' でポータル専用レイアウトを使用
  // requiresStaffAuth: true でスタッフ認証ガードを適用
  // ==============================
  {
    path: '/mypage',
    redirect: '/mypage/login',
  },
  {
    path: '/mypage/login',
    name: 'portal-login',
    component: () => import('@/views/portal/StaffLoginView.vue'),
    meta: { layout: 'portal' },
  },
  {
    path: '/mypage/dashboard',
    name: 'portal-dashboard',
    component: () => import('@/views/portal/StaffDashboardView.vue'),
    meta: { layout: 'portal', requiresStaffAuth: true },
  },
  {
    path: '/mypage/profile',
    name: 'portal-profile',
    component: () => import('@/views/portal/StaffProfileEditView.vue'),
    meta: { layout: 'portal', requiresStaffAuth: true },
  },
  {
    path: '/mypage/schedule',
    name: 'portal-schedule',
    component: () => import('@/views/portal/StaffScheduleEditView.vue'),
    meta: { layout: 'portal', requiresStaffAuth: true },
  },
]

// 認証ナビゲーションガード
// requiresAuth が設定されたルートでは、localStorage のトークンを確認する。
// 期限切れ時は verify（staff は refresh を含む）を実行し、失効していればログインへ戻す。
//
// 複数タブ対応: localStorage は複数タブで自動同期される。タブ複製時など
// 複数タブが同時に verify/refresh を実行しても、localStorage の更新を
// 各ストアが検出して再同期するため、不整合は解決される。
export function setupRouterGuards(router: Router): void {
  router.beforeEach(async (to) => {
    // 管理者認証ガード
    if (to.meta.requiresAuth) {
      const authStore = useAuthStore()
      if (!authStore.token) {
        return { name: 'admin-login' }
      }

      const isValid = await authStore.verify()
      if (!isValid) {
        return { name: 'admin-login' }
      }
    }

    // スタッフ認証ガード
    if (to.meta.requiresStaffAuth) {
      const staffAuthStore = useStaffAuthStore()
      if (!staffAuthStore.token) {
        return { name: 'portal-login' }
      }

      const isValid = await staffAuthStore.verify()
      if (!isValid) {
        return { name: 'portal-login' }
      }
    }
  })
}
