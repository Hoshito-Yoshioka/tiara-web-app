<script setup lang="ts">
  import { computed } from 'vue'
  import { RouterView, useRoute } from 'vue-router'
  import { useHead } from '@unhead/vue'
  import TheHeader from '@/components/layout/TheHeader.vue'
  import TheFooter from '@/components/layout/TheFooter.vue'
  import { DEFAULT_TITLE, DEFAULT_DESCRIPTION } from '@/lib/seo'

  const route = useRoute()

  // サイト共通のデフォルトメタ。各公開ページは usePageMeta() で上書きする。
  // 管理画面・スタッフポータルは検索エンジンにインデックスさせない。
  useHead({
    title: DEFAULT_TITLE,
    meta: computed(() => {
      const meta: { name: string; content: string }[] = [
        { name: 'description', content: DEFAULT_DESCRIPTION },
      ]
      if (route.meta.layout === 'admin' || route.meta.layout === 'portal') {
        meta.push({ name: 'robots', content: 'noindex, nofollow' })
      }
      return meta
    }),
  })
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background">
    <!-- Admin ルートでは専用レイアウトを使用（各ビュー内で AdminLayout を使用） -->
    <template v-if="route.meta.layout === 'admin' || route.meta.layout === 'portal'">
      <RouterView />
    </template>
    <!-- 公開ページでは通常のヘッダー/フッターレイアウト -->
    <!-- Suspense: 公開ビューは async setup で API データを取得するため（SSG でも HTML に内容が入る） -->
    <template v-else>
      <TheHeader />
      <main class="flex-1">
        <RouterView v-slot="{ Component }">
          <Suspense>
            <component :is="Component" />
            <template #fallback>
              <div class="min-h-screen flex items-center justify-center">
                <div
                  class="w-6 h-6 border border-primary border-t-transparent rounded-full animate-spin"
                />
              </div>
            </template>
          </Suspense>
        </RouterView>
      </main>
      <TheFooter />
    </template>
  </div>
</template>

<style scoped></style>
