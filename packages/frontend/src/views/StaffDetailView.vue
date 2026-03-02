<script setup lang="ts">
  import { useRoute } from 'vue-router'

  const route = useRoute()
  // TODO: route.params.id を使って API からスタッフ詳細データを取得する
</script>

<template>
  <div class="min-h-screen bg-background pt-24 pb-28">
    <div class="container mx-auto px-6">
      <div
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 700 } }"
        class="flex flex-col items-center mb-20 text-center"
      >
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Profile</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Staff Detail
        </h1>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <!-- TODO: スタッフ詳細・出勤スケジュール実装予定 -->
      <p class="text-center text-muted-foreground tracking-widest text-sm uppercase">Coming Soon</p>
    </div>
  </div>
</template>
