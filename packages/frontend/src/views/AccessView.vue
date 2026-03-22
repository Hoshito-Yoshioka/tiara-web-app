<script setup lang="ts">
  import { MapPin, Clock, Train } from 'lucide-vue-next'

  /** 店舗情報 */
  const shopInfo = {
    name: 'BAR Tiara',
    postalCode: '〒040-0011',
    address: '北海道函館市本町１−２８ 第５大栄ビル',
    mapEmbedUrl:
      'https://maps.google.com/maps?q=%E5%8C%97%E6%B5%B7%E9%81%93%E5%87%BD%E9%A4%A8%E5%B8%82%E6%9C%AC%E7%94%BA1-28+%E7%AC%AC5%E5%A4%A7%E6%A0%84%E3%83%93%E3%83%AB&t=m&z=17&output=embed&iwloc=',
    mapLinkUrl:
      'https://www.google.com/maps/search/?api=1&query=北海道函館市本町１−２８+第５大栄ビル',
    hours: '20:00 – 02:00',
    access: [
      '函館市電「中央病院前」電停より徒歩約3分',
      '函館市電「五稜郭公園前」電停より徒歩約4分',
    ],
  } as const
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
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Location</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Access
        </h1>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <div class="max-w-4xl mx-auto">
        <!-- Google Map 埋め込み -->
        <div
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :enter="{ opacity: 1, y: 0, transition: { duration: 700, delay: 200 } }"
          class="mb-16"
        >
          <div class="border border-border overflow-hidden">
            <iframe
              :src="shopInfo.mapEmbedUrl"
              width="100%"
              height="450"
              style="border: 0"
              allowfullscreen
              loading="lazy"
              referrerpolicy="no-referrer-when-downgrade"
              title="BAR Tiara の所在地"
              class="w-full h-[300px] md:h-[450px]"
            />
          </div>
        </div>

        <!-- 店舗情報カード -->
        <div
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :visible="{ opacity: 1, y: 0, transition: { duration: 600, delay: 300 } }"
          class="border border-border bg-card p-8 md:p-12"
        >
          <!-- 店舗名 -->
          <div class="flex flex-col items-center mb-10">
            <h2 class="text-2xl md:text-3xl font-light tracking-[0.15em] uppercase text-foreground">
              {{ shopInfo.name }}
            </h2>
            <span class="block w-8 h-px bg-primary mt-4" />
          </div>

          <div class="space-y-8">
            <!-- 住所 -->
            <div class="flex items-start gap-4">
              <MapPin class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
              <div>
                <p class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground mb-1">
                  Address
                </p>
                <p class="text-sm text-foreground leading-relaxed">
                  {{ shopInfo.postalCode }}<br />
                  {{ shopInfo.address }}
                </p>
                <a
                  :href="shopInfo.mapLinkUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="mt-2 inline-flex items-center gap-1 text-[11px] tracking-[0.15em] uppercase text-primary hover:text-primary/80 transition-colors duration-300"
                >
                  Google Maps で開く
                  <span class="inline-block w-3 h-px bg-current" />
                </a>
              </div>
            </div>

            <!-- 営業時間 -->
            <div class="flex items-start gap-4">
              <Clock class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
              <div>
                <p class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground mb-1">
                  Hours
                </p>
                <p class="text-sm text-foreground">
                  {{ shopInfo.hours }}
                </p>
              </div>
            </div>

            <!-- アクセス方法 -->
            <div class="flex items-start gap-4">
              <Train class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
              <div>
                <p class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground mb-1">
                  Transportation
                </p>
                <ul class="space-y-1">
                  <li
                    v-for="(line, idx) in shopInfo.access"
                    :key="idx"
                    class="text-sm text-foreground leading-relaxed"
                  >
                    {{ line }}
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
