<script setup lang="ts">
  import { onMounted } from 'vue'
  import Button from '@/components/ui/button.vue'
  import Card from '@/components/ui/card.vue'
  import CardHeader from '@/components/ui/card-header.vue'
  import CardTitle from '@/components/ui/card-title.vue'
  import CardContent from '@/components/ui/card-content.vue'
  import CardFooter from '@/components/ui/card-footer.vue'
  import { useShopApi } from '@/composables/useShopApi'

  const { shops, fetchShops } = useShopApi()

  onMounted(() => {
    fetchShops()
  })

  /** TIME型のISO文字列から "HH:MM" を抽出 */
  function formatTime(raw: string): string {
    const date = new Date(raw)
    const hours = date.getUTCHours().toString().padStart(2, '0')
    const minutes = date.getUTCMinutes().toString().padStart(2, '0')
    return `${hours}:${minutes}`
  }
</script>

<template>
  <div class="bg-background text-foreground">
    <!-- Hero Section: h-screen でヘッダー下から全画面 -->
    <section
      class="relative h-screen flex items-center justify-center text-center bg-cover bg-center"
      style="
        background-image: url('https://images.unsplash.com/photo-1543782245-ae5265494d30?q=80&w=2940&auto=format&fit=crop');
      "
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

        <h1 class="tracking-[0.6em] text-4xl md:text-6xl font-light uppercase mb-6">TIARA</h1>
        <p class="text-sm md:text-base tracking-[0.2em] text-white/70 max-w-md uppercase">
          Experience the finest selection of beverages<br />in a sophisticated ambiance.
        </p>

        <span class="block w-12 h-px bg-white/30 my-8" />

        <Button
          variant="outline"
          class="tracking-[0.2em] text-xs uppercase border-white/50 text-white bg-transparent hover:bg-white hover:text-black transition-all duration-500 px-10 py-5"
          as-child
        >
          <router-link to="/shop">Explore Our Bars</router-link>
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

    <!-- Shops Section -->
    <section class="container mx-auto py-28 px-6">
      <!-- セクションヘッダー -->
      <div
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :visible="{ opacity: 1, y: 0, transition: { duration: 700 } }"
        class="flex flex-col items-center mb-16"
      >
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Our Locations</span>
        <h2 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Premier Bars
        </h2>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <!-- カードグリッド -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="(shop, index) in shops"
          :key="shop.id"
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :visible="{ opacity: 1, y: 0, transition: { duration: 600, delay: index * 150 } }"
        >
          <Card
            class="flex flex-col h-full rounded-none border-border bg-card hover:border-primary/50 transition-colors duration-500 group"
          >
            <CardHeader class="pb-3">
              <CardTitle class="text-base font-light tracking-[0.1em] text-foreground">
                {{ shop.name }}
              </CardTitle>
            </CardHeader>

            <CardContent class="flex-grow pt-0">
              <div class="space-y-2 text-xs text-muted-foreground">
                <p class="flex items-center gap-2">
                  <span class="w-1 h-1 rounded-full bg-primary inline-block flex-shrink-0" />
                  {{ shop.address }}
                </p>
                <p class="flex items-center gap-2">
                  <span class="w-1 h-1 rounded-full bg-primary inline-block flex-shrink-0" />
                  {{ formatTime(shop.openingTime) }} – {{ formatTime(shop.closingTime) }}
                </p>
              </div>
            </CardContent>

            <CardFooter class="pt-4">
              <Button
                variant="ghost"
                class="w-full text-[11px] tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground border border-transparent hover:border-white/20 transition-all duration-300 rounded-none"
                as-child
              >
                <router-link to="/shop">View Details</router-link>
              </Button>
            </CardFooter>
          </Card>
        </div>
      </div>
    </section>
  </div>
</template>
