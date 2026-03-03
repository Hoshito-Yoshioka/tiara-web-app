<script setup lang="ts">
  import { Wine, Coffee, Beer, Sparkles } from 'lucide-vue-next'

  /** 料金カテゴリの型 */
  interface PriceItem {
    name: string
    price: string
    description?: string
  }

  interface PriceCategory {
    title: string
    icon: typeof Wine
    items: PriceItem[]
  }

  /** 料金データ */
  const categories: PriceCategory[] = [
    {
      title: 'System',
      icon: Sparkles,
      items: [
        {
          name: 'チャージ',
          price: '¥1,000',
          description: 'お一人様あたり',
        },
        {
          name: 'セット料金（60分）',
          price: '¥3,000',
          description: 'チャージ＋ドリンク2杯',
        },
        {
          name: '延長（30分）',
          price: '¥1,500',
          description: 'ドリンク1杯付き',
        },
      ],
    },
    {
      title: 'Cocktails',
      icon: Wine,
      items: [
        { name: 'スタンダードカクテル', price: '¥800〜' },
        { name: 'プレミアムカクテル', price: '¥1,200〜' },
        { name: 'オリジナルカクテル', price: '¥1,000〜' },
        { name: 'ノンアルコールカクテル', price: '¥700〜' },
      ],
    },
    {
      title: 'Whisky & Spirits',
      icon: Coffee,
      items: [
        { name: 'ハウスウイスキー', price: '¥800' },
        { name: 'プレミアムウイスキー', price: '¥1,200〜' },
        { name: 'ブランデー', price: '¥1,000〜' },
        { name: 'ジン / ウォッカ / ラム', price: '¥800〜' },
      ],
    },
    {
      title: 'Beer & Wine',
      icon: Beer,
      items: [
        { name: '生ビール', price: '¥700' },
        { name: 'クラフトビール', price: '¥900〜' },
        { name: 'グラスワイン（赤・白）', price: '¥800〜' },
        { name: 'シャンパン（グラス）', price: '¥1,500〜' },
        { name: 'シャンパン（ボトル）', price: '¥8,000〜' },
      ],
    },
  ]
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

      <!-- 料金カテゴリ -->
      <div class="max-w-3xl mx-auto space-y-16">
        <section
          v-for="(category, catIdx) in categories"
          :key="category.title"
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :visible="{ opacity: 1, y: 0, transition: { duration: 600, delay: catIdx * 150 } }"
        >
          <!-- カテゴリヘッダー -->
          <div class="flex items-center gap-3 mb-8">
            <component :is="category.icon" class="w-4 h-4 text-primary" />
            <h2 class="text-lg font-light tracking-[0.15em] uppercase text-foreground">
              {{ category.title }}
            </h2>
            <span class="flex-1 h-px bg-border" />
          </div>

          <!-- 料金リスト -->
          <div class="border border-border bg-card overflow-hidden">
            <div
              v-for="(item, idx) in category.items"
              :key="item.name"
              class="flex items-baseline justify-between px-6 md:px-8 py-4 transition-colors duration-300 hover:bg-secondary/50"
              :class="idx < category.items.length - 1 ? 'border-b border-border' : ''"
            >
              <div class="flex-1 min-w-0 pr-4">
                <p class="text-sm text-foreground tracking-wide">
                  {{ item.name }}
                </p>
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
