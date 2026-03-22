<script setup lang="ts">
  import { RouterView, useRoute } from 'vue-router'
  import TheHeader from '@/components/layout/TheHeader.vue'
  import TheFooter from '@/components/layout/TheFooter.vue'

  const route = useRoute()
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background">
    <!-- Admin ルートでは専用レイアウトを使用（各ビュー内で AdminLayout を使用） -->
    <template v-if="route.meta.layout === 'admin' || route.meta.layout === 'portal'">
      <RouterView />
    </template>
    <!-- 公開ページでは通常のヘッダー/フッターレイアウト -->
    <template v-else>
      <TheHeader />
      <main class="flex-1">
        <RouterView />
      </main>
      <TheFooter />
    </template>
  </div>
</template>

<style scoped></style>
