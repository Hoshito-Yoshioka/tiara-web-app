<script setup lang="ts">
  import { computed } from 'vue'
  import { useHead } from '@unhead/vue'
  import Button from '@/components/ui/button.vue'
  import { MapPin, Phone } from 'lucide-vue-next'
  import { useShopApi } from '@/composables/useShopApi'
  import { usePageMeta } from '@/composables/usePageMeta'

  /** ヒーロー背景画像。LCP 改善のため preload する */
  const HERO_IMAGE_URL =
    'https://images.unsplash.com/photo-1566417713940-fe7c737a9ef2?q=80&w=2940&auto=format&fit=crop'

  const { shops, fetchShops } = useShopApi()

  usePageMeta()
  useHead({
    link: [{ rel: 'preload', as: 'image', href: HERO_IMAGE_URL }],
  })

  // async setup で取得することで、SSG ビルド時に店舗情報が HTML に含まれる
  await fetchShops()

  /** 一店舗のみ想定 */
  const shop = computed(() => shops.value[0] ?? null)
</script>

<template>
  <div class="bg-background text-foreground">
    <!-- Hero Section: h-screen でヘッダー下から全画面 -->
    <section
      class="relative h-screen flex items-center justify-center text-center bg-cover bg-center"
      :style="{ backgroundImage: `url('${HERO_IMAGE_URL}')` }"
    >
      <!-- オーバーレイ -->
      <div class="absolute inset-0 bg-black/70" />

      <!-- コンテンツ -->
      <div
        v-motion
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 1000, delay: 200 } }"
        class="relative z-10 text-white px-6 flex flex-col items-center"
      >
        <!-- ゴールドラインアクセント -->
        <span class="block w-12 h-px bg-primary mb-8" />

        <p class="text-[11px] tracking-[0.4em] uppercase text-white/50 mb-4">New Club</p>
        <h1 class="tracking-[0.6em] text-4xl md:text-6xl font-light uppercase mb-6">TIARA</h1>
        <p class="text-sm md:text-base tracking-[0.2em] text-white/70 max-w-md uppercase">
          An exclusive space where elegance meets<br />the art of hospitality.
        </p>

        <span class="block w-12 h-px bg-white/30 my-8" />

        <Button
          variant="outline"
          class="tracking-[0.2em] text-xs uppercase border-white/50 text-white bg-transparent hover:bg-white hover:text-black transition-all duration-500 px-10 py-5"
          as-child
        >
          <router-link to="/shop">Enter TIARA</router-link>
        </Button>
      </div>

      <!-- スクロールインジケーター -->
      <div
        class="absolute bottom-10 left-1/2 -translate-x-1/2 flex flex-col items-center gap-2 text-white/40"
      >
        <span class="text-[10px] tracking-[0.3em] uppercase">Scroll</span>
        <span class="block w-px h-10 bg-white/20 animate-pulse" />
      </div>
    </section>

    <!-- Club Section: 一店舗フィーチャーレイアウト -->
    <section class="container mx-auto py-28 px-6">
      <!-- セクションヘッダー -->
      <div
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :visibleOnce="{ opacity: 1, y: 0, transition: { duration: 700 } }"
        class="flex flex-col items-center mb-16"
      >
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Our Club</span>
        <h2 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          New Club TIARA
        </h2>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <!-- 単一店舗フィーチャーカード -->
      <div
        v-if="shop"
        v-motion
        :initial="{ opacity: 0, y: 30 }"
        :visibleOnce="{ opacity: 1, y: 0, transition: { duration: 600 } }"
        class="max-w-2xl mx-auto"
      >
        <div
          class="border border-border bg-card p-8 md:p-12 hover:border-primary/50 transition-colors duration-500"
        >
          <!-- 店舗名 -->
          <div class="flex flex-col items-center mb-10">
            <h3 class="text-2xl md:text-3xl font-light tracking-[0.15em] uppercase text-foreground">
              {{ shop.name }}
            </h3>
            <span class="block w-8 h-px bg-primary mt-4" />
          </div>

          <!-- 店舗情報 -->
          <div class="space-y-6">
            <div class="flex items-start gap-4">
              <MapPin class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
              <p class="text-sm text-muted-foreground leading-relaxed">
                {{ shop.address }}
              </p>
            </div>
            <div class="flex items-start gap-4">
              <Phone class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
              <p class="text-sm text-muted-foreground">
                <a
                  href="tel:080-2818-2554"
                  class="hover:text-primary transition-colors duration-300"
                  >080-2818-2554</a
                >
              </p>
            </div>
          </div>

          <!-- CTAボタン -->
          <div class="mt-10 flex justify-center">
            <Button
              variant="ghost"
              class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground border border-border hover:border-white/20 transition-all duration-300 rounded-none px-10 py-4"
              as-child
            >
              <router-link to="/shop">View Details</router-link>
            </Button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
