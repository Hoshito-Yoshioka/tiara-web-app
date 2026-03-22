<script setup lang="ts">
  import { useRouter } from 'vue-router'
  import { useStaffAuthStore } from '@/stores/staffAuth'

  const router = useRouter()
  const staffAuth = useStaffAuthStore()

  function handleLogout() {
    staffAuth.logout()
    router.push({ name: 'portal-login' })
  }
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- ヘッダー -->
    <header class="border-b border-white/10 px-6 py-4">
      <div class="max-w-4xl mx-auto flex items-center justify-between">
        <div>
          <h1 class="tracking-[0.3em] text-sm font-light text-foreground uppercase">TIARA</h1>
          <p class="text-[10px] text-muted-foreground tracking-wider mt-0.5">STAFF PORTAL</p>
        </div>
        <div class="flex items-center gap-4">
          <RouterLink
            :to="{ name: 'portal-dashboard' }"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            active-class="text-foreground"
          >
            ダッシュボード
          </RouterLink>
          <RouterLink
            :to="{ name: 'portal-profile' }"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            active-class="text-foreground"
          >
            プロフィール
          </RouterLink>
          <RouterLink
            :to="{ name: 'portal-schedule' }"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            active-class="text-foreground"
          >
            スケジュール
          </RouterLink>
          <button
            @click="handleLogout"
            class="text-xs text-muted-foreground hover:text-red-400 transition-colors ml-4"
          >
            ログアウト
          </button>
        </div>
      </div>
    </header>

    <!-- コンテンツ -->
    <main class="max-w-4xl mx-auto px-6 py-8">
      <slot />
    </main>
  </div>
</template>
