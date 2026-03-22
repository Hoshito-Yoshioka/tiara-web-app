<script setup lang="ts">
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useStaffAuthStore } from '@/stores/staffAuth'

  const staffAuthStore = useStaffAuthStore()
  const router = useRouter()

  const username = ref('')
  const password = ref('')
  const error = ref<string | null>(null)
  const isLoading = ref(false)

  async function handleLogin() {
    error.value = null
    isLoading.value = true

    try {
      await staffAuthStore.login(username.value, password.value)
      router.push({ name: 'portal-dashboard' })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'ログインに失敗しました'
    } finally {
      isLoading.value = false
    }
  }
</script>

<template>
  <div
    class="min-h-screen bg-background flex items-center justify-center px-4"
    v-motion
    :initial="{ opacity: 0, y: 20 }"
    :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
  >
    <div class="w-full max-w-sm">
      <!-- ロゴ -->
      <div class="text-center mb-8">
        <h1 class="tracking-[0.5em] text-lg font-light text-foreground uppercase mb-2">TIARA</h1>
        <p class="text-xs text-muted-foreground tracking-wider">STAFF PORTAL</p>
      </div>

      <!-- ログインフォーム -->
      <form @submit.prevent="handleLogin" class="space-y-6">
        <!-- エラー表示 -->
        <div v-if="error" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3">
          <p class="text-sm text-red-400">{{ error }}</p>
        </div>

        <div class="space-y-2">
          <label for="username" class="text-xs text-muted-foreground tracking-wider uppercase">
            Username
          </label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            autocomplete="username"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors"
            placeholder="スタッフユーザー名を入力"
          />
        </div>

        <div class="space-y-2">
          <label for="password" class="text-xs text-muted-foreground tracking-wider uppercase">
            Password
          </label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="current-password"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors"
            placeholder="パスワードを入力"
          />
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full bg-white text-black rounded-lg py-3 text-sm font-medium tracking-wider uppercase hover:bg-white/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isLoading ? 'ログイン中...' : 'ログイン' }}
        </button>
      </form>

      <!-- サイトへ戻るリンク -->
      <div class="mt-8 text-center">
        <RouterLink
          to="/"
          class="text-xs text-muted-foreground hover:text-foreground transition-colors"
        >
          ← サイトへ戻る
        </RouterLink>
      </div>
    </div>
  </div>
</template>
