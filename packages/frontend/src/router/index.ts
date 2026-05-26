import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  // ページ遷移のたびにスクロール位置をトップへリセット
  scrollBehavior: () => ({ top: 0, behavior: 'smooth' }),
  routes: [
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
    // NOTE: /schedule は一時的に非公開（admin/portal のスケジュール操作は有効）
    // {
    //   path: '/schedule',
    //   name: 'schedule',
    //   component: () => import('@/views/ScheduleView.vue'),
    // },
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
      path: '/admin',
      redirect: '/admin/login',
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('@/views/admin/AdminLoginView.vue'),
      meta: { layout: 'admin' },
    },
    {
      path: '/admin/shop',
      name: 'admin-shop-edit',
      component: () => import('@/views/admin/AdminShopEditView.vue'),
      meta: { layout: 'admin', requiresAuth: true },
    },
    {
      path: '/admin/staffs',
      name: 'admin-staff-list',
      component: () => import('@/views/admin/AdminStaffListView.vue'),
      meta: { layout: 'admin', requiresAuth: true },
    },
    {
      path: '/admin/staffs/new',
      name: 'admin-staff-new',
      component: () => import('@/views/admin/AdminStaffEditView.vue'),
      meta: { layout: 'admin', requiresAuth: true },
    },
    {
      path: '/admin/staffs/:id/edit',
      name: 'admin-staff-edit',
      component: () => import('@/views/admin/AdminStaffEditView.vue'),
      meta: { layout: 'admin', requiresAuth: true },
    },
    {
      path: '/admin/menu',
      name: 'admin-menu-edit',
      component: () => import('@/views/admin/AdminMenuView.vue'),
      meta: { layout: 'admin', requiresAuth: true },
    },
    {
      path: '/admin/profile-reviews',
      name: 'admin-profile-reviews',
      component: () => import('@/views/admin/AdminProfileReviewView.vue'),
      meta: { layout: 'admin', requiresAuth: true },
    },
    {
      path: '/admin/schedule-reviews',
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
  ],
})

// 認証ナビゲーションガード
// requiresAuth が設定されたルートでは、localStorage のトークンを確認する。
// Pinia の初期化前でも動作するよう localStorage を直接参照。
router.beforeEach((to) => {
  // 管理者認証ガード
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('tiara_admin_token')
    if (!token) {
      return { name: 'admin-login' }
    }
  }

  // スタッフ認証ガード
  if (to.meta.requiresStaffAuth) {
    const token = localStorage.getItem('tiara_staff_token')
    if (!token) {
      return { name: 'portal-login' }
    }
  }
})

export default router
