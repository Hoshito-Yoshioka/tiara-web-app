<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import AdminLayout from '@/components/layout/AdminLayout.vue'
  import { useShopApi } from '@/composables/useShopApi'
  import { useAdminApi } from '@/composables/useAdminApi'
  import { Save } from 'lucide-vue-next'

  const { shops, fetchShops, isLoading: isFetching } = useShopApi()
  const { updateShop, isLoading: isSaving, error: saveError } = useAdminApi()

  const form = ref({
    name: '',
    address: '',
    openingTime: '',
    closingTime: '',
  })
  const shopId = ref('')
  const successMessage = ref<string | null>(null)

  /**
   * Backend から返る時刻文字列 "0000-01-01T18:00:00Z" を
   * HTML time input 用の "HH:MM" 形式に変換するヘルパー。
   */
  function formatTimeForInput(timeStr: string): string {
    if (timeStr.includes('T')) {
      const match = timeStr.match(/T(\d{2}:\d{2})/)
      return match ? match[1] : timeStr
    }
    return timeStr.slice(0, 5)
  }

  onMounted(async () => {
    await fetchShops()
    // 最初の店舗データでフォームを初期化
    if (shops.value.length > 0) {
      const shop = shops.value[0]
      shopId.value = shop.id
      form.value = {
        name: shop.name,
        address: shop.address,
        openingTime: formatTimeForInput(shop.openingTime),
        closingTime: formatTimeForInput(shop.closingTime),
      }
    }
  })

  async function handleSubmit() {
    successMessage.value = null
    const result = await updateShop(shopId.value, form.value)
    if (result) {
      successMessage.value = '店舗情報を更新しました'
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
    }
  }
</script>

<template>
  <AdminLayout>
    <div
      class="max-w-2xl"
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 400 } }"
    >
      <h1 class="text-xl font-light tracking-wider text-foreground mb-8">店舗情報編集</h1>

      <div v-if="isFetching" class="text-muted-foreground text-sm">読み込み中...</div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-6">
        <!-- エラー表示 -->
        <div v-if="saveError" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3">
          <p class="text-sm text-red-400">{{ saveError }}</p>
        </div>

        <!-- 成功メッセージ -->
        <div
          v-if="successMessage"
          class="bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3"
        >
          <p class="text-sm text-green-400">{{ successMessage }}</p>
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">店舗名</label>
          <input
            v-model="form.name"
            type="text"
            required
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">住所</label>
          <input
            v-model="form.address"
            type="text"
            required
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <label class="text-xs text-muted-foreground tracking-wider uppercase">開店時間</label>
            <input
              v-model="form.openingTime"
              type="time"
              required
              class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors [color-scheme:dark]"
            />
          </div>
          <div class="space-y-2">
            <label class="text-xs text-muted-foreground tracking-wider uppercase">閉店時間</label>
            <input
              v-model="form.closingTime"
              type="time"
              required
              class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors [color-scheme:dark]"
            />
          </div>
        </div>

        <button
          type="submit"
          :disabled="isSaving"
          class="flex items-center gap-2 bg-white text-black rounded-lg px-6 py-3 text-sm font-medium tracking-wider uppercase hover:bg-white/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Save :size="16" />
          {{ isSaving ? '保存中...' : '保存' }}
        </button>
      </form>
    </div>
  </AdminLayout>
</template>
