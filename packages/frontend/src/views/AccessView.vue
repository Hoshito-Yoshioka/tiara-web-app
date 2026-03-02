<script setup lang="ts">
  // TODO: 店舗の住所・地図情報を定義する
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
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Location</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Access
        </h1>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <!-- TODO: 地図・アクセス情報実装予定 -->
      <p class="text-center text-muted-foreground tracking-widest text-sm uppercase">Coming Soon</p>
    </div>
  </div>
</template>
