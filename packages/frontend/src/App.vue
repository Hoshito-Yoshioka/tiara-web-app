<script setup lang="ts">
  import { RouterView } from 'vue-router'
  import TheHeader from '@/components/layout/TheHeader.vue'
  import TheFooter from '@/components/layout/TheFooter.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background">
    <TheHeader />
    <!-- pt-16 でヘッダー分の余白を確保（ヒーローセクションは各ビューで個別に対応） -->
    <main class="flex-1">
      <RouterView />
    </main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
