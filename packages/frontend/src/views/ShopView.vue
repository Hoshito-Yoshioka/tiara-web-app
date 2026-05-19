<script setup lang="ts">
  import { onMounted } from 'vue'
  import { useShopApi } from '@/composables/useShopApi'
  import { Phone, MapPin } from 'lucide-vue-next'

  const { shops, isLoading, error, fetchShops } = useShopApi()

  onMounted(() => {
    fetchShops()
  })
</script>

<template>
  <div class="min-h-screen bg-background pt-24 pb-28">
    <div class="container mx-auto px-6">
      <!-- セクションヘッダー -->
      <div
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 700 } }"
        class="flex flex-col items-center mb-20 text-center"
      >
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Our Club</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Shop
        </h1>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <!-- ローディング -->
      <div v-if="isLoading" class="flex justify-center py-20">
        <div class="w-6 h-6 border border-primary border-t-transparent rounded-full animate-spin" />
      </div>

      <!-- エラー -->
      <div v-else-if="error" class="text-center py-20">
        <p class="text-destructive text-sm tracking-wide">{{ error }}</p>
        <button
          class="mt-6 text-xs tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground border border-border px-8 py-3 transition-colors duration-300"
          @click="fetchShops"
        >
          Retry
        </button>
      </div>

      <!-- 店舗データなし -->
      <div v-else-if="shops.length === 0" class="text-center py-20">
        <p class="text-muted-foreground tracking-widest text-sm uppercase">
          No shop information available
        </p>
      </div>

      <!-- 店舗リスト -->
      <div v-else class="space-y-24">
        <article
          v-for="(shop, index) in shops"
          :key="shop.id"
          v-motion
          :initial="{ opacity: 0, y: 40 }"
          :visibleOnce="{ opacity: 1, y: 0, transition: { duration: 600, delay: index * 200 } }"
          class="max-w-3xl mx-auto"
        >
          <!-- 店舗名 -->
          <div class="flex flex-col items-center mb-10">
            <h2 class="text-2xl md:text-3xl font-light tracking-[0.15em] uppercase text-foreground">
              {{ shop.name }}
            </h2>
            <span class="block w-8 h-px bg-primary mt-4" />
          </div>

          <!-- 店舗情報カード -->
          <div class="border border-border bg-card p-8 md:p-12 space-y-6">
            <!-- 住所 -->
            <div class="flex items-start gap-4">
              <MapPin class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
              <div>
                <p class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground mb-1">
                  Address
                </p>
                <p class="text-sm text-foreground leading-relaxed">
                  {{ shop.address }}
                </p>
              </div>
            </div>

            <!-- 電話番号 -->
            <div class="flex items-start gap-4">
              <Phone class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
              <div>
                <p class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground mb-1">
                  Phone
                </p>
                <p class="text-sm text-foreground">
                  <a
                    href="tel:0138-83-6040"
                    class="hover:text-primary transition-colors duration-300"
                    >0138-83-6040</a
                  >
                </p>
              </div>
            </div>
          </div>
        </article>
      </div>
    </div>
  </div>
</template>
