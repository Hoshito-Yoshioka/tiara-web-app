<script setup lang="ts">
  import { RouterLink, useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/auth'
  import {
    LogOut,
    Store,
    Users,
    Home,
    UtensilsCrossed,
    UserCheck,
    CalendarCheck,
  } from 'lucide-vue-next'

  const authStore = useAuthStore()
  const router = useRouter()

  function handleLogout() {
    authStore.logout()
    router.push({ name: 'admin-login' })
  }
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Admin Header -->
    <header class="bg-black/95 border-b border-white/10 sticky top-0 z-50">
      <div class="container mx-auto px-6 h-14 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <RouterLink
            to="/admin/shop"
            class="tracking-[0.3em] text-sm font-light text-foreground uppercase"
          >
            TIARA ADMIN
          </RouterLink>

          <nav class="hidden sm:flex items-center gap-4">
            <RouterLink
              to="/admin/shop"
              class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
              active-class="text-foreground"
            >
              <Store :size="14" />
              店舗編集
            </RouterLink>
            <RouterLink
              to="/admin/staffs"
              class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
              active-class="text-foreground"
            >
              <Users :size="14" />
              スタッフ管理
            </RouterLink>
            <RouterLink
              to="/admin/menu"
              class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
              active-class="text-foreground"
            >
              <UtensilsCrossed :size="14" />
              メニュー管理
            </RouterLink>
            <RouterLink
              to="/admin/profile-reviews"
              class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
              active-class="text-foreground"
            >
              <UserCheck :size="14" />
              プロフィール承認
            </RouterLink>
            <RouterLink
              to="/admin/schedule-reviews"
              class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
              active-class="text-foreground"
            >
              <CalendarCheck :size="14" />
              シフト承認
            </RouterLink>
          </nav>
        </div>

        <div class="flex items-center gap-4">
          <RouterLink
            to="/"
            class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
          >
            <Home :size="14" />
            <span class="hidden sm:inline">サイトへ</span>
          </RouterLink>
          <button
            @click="handleLogout"
            class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-red-400 transition-colors"
          >
            <LogOut :size="14" />
            <span class="hidden sm:inline">ログアウト</span>
          </button>
        </div>
      </div>

      <!-- モバイルナビ -->
      <nav class="sm:hidden flex items-center gap-4 px-6 pb-3 overflow-x-auto">
        <RouterLink
          to="/admin/shop"
          class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
          active-class="text-foreground"
        >
          <Store :size="14" />
          店舗編集
        </RouterLink>
        <RouterLink
          to="/admin/staffs"
          class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
          active-class="text-foreground"
        >
          <Users :size="14" />
          スタッフ管理
        </RouterLink>
        <RouterLink
          to="/admin/menu"
          class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
          active-class="text-foreground"
        >
          <UtensilsCrossed :size="14" />
          メニュー管理
        </RouterLink>
        <RouterLink
          to="/admin/profile-reviews"
          class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
          active-class="text-foreground"
        >
          <UserCheck :size="14" />
          プロフィール承認
        </RouterLink>
        <RouterLink
          to="/admin/schedule-reviews"
          class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
          active-class="text-foreground"
        >
          <CalendarCheck :size="14" />
          シフト承認
        </RouterLink>
      </nav>
    </header>

    <!-- Content -->
    <main class="container mx-auto px-6 py-8">
      <slot />
    </main>
  </div>
</template>
