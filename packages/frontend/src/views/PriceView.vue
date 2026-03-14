<script setup lang="ts">
  import { onMounted } from 'vue'
  import { useMenuApi } from '@/composables/useMenuApi'

  const { menuList, isLoading, error, fetchMenus } = useMenuApi()

  onMounted(() => {
    fetchMenus()
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
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">System</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Price
        </h1>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <!-- ローディング -->
      <div v-if="isLoading" class="flex justify-center py-20">
        <p class="text-sm text-muted-foreground tracking-wider">Loading...</p>
      </div>

      <!-- エラー -->
      <div v-else-if="error" class="max-w-3xl mx-auto text-center py-12">
        <p class="text-sm text-red-400">{{ error }}</p>
      </div>

      <!-- 料金カテゴリ -->
      <div v-else class="max-w-3xl mx-auto space-y-16">
        <section
          v-for="(entry, catIdx) in menuList"
          :key="entry.category.id"
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :visible="{ opacity: 1, y: 0, transition: { duration: 600, delay: catIdx * 150 } }"
        >
          <!-- カテゴリヘッダー -->
          <div class="flex items-center gap-3 mb-8">
            <span class="w-1 h-4 bg-primary block" />
            <h2 class="text-lg font-light tracking-[0.15em] uppercase text-foreground">
              {{ entry.category.name }}
            </h2>
            <span class="flex-1 h-px bg-border" />
          </div>

          <!-- カテゴリ説明 -->
          <p
            v-if="entry.category.description"
            class="text-xs text-muted-foreground mb-6 tracking-wide"
          >
            {{ entry.category.description }}
          </p>

          <!-- 料金リスト -->
          <div class="border border-border bg-card overflow-hidden">
            <div
              v-for="(item, idx) in entry.items"
              :key="item.id"
              class="flex items-baseline justify-between px-6 md:px-8 py-4 transition-colors duration-300 hover:bg-secondary/50"
              :class="idx < entry.items.length - 1 ? 'border-b border-border' : ''"
            >
              <div class="flex-1 min-w-0 pr-4">
                <p class="text-sm text-foreground tracking-wide">{{ item.name }}</p>
                <p v-if="item.description" class="text-[11px] text-muted-foreground mt-0.5">
                  {{ item.description }}
                </p>
              </div>
              <p class="text-sm text-primary font-light tracking-wider whitespace-nowrap">
                {{ item.price }}
              </p>
            </div>
          </div>
        </section>
      </div>

      <!-- 注意書き -->
      <div
        v-motion
        :initial="{ opacity: 0 }"
        :visible="{ opacity: 1, transition: { duration: 500, delay: 600 } }"
        class="max-w-3xl mx-auto mt-16 text-center"
      >
        <p class="text-[11px] text-muted-foreground tracking-wide leading-relaxed">
          ※ 表示価格はすべて税込です。<br />
          ※ 料金は予告なく変更する場合がございます。
        </p>
      </div>
    </div>
  </div>
</template>
