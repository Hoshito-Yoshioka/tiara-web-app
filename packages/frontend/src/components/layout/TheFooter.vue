<script setup lang="ts">
  import { RouterLink } from 'vue-router'

  interface NavItem {
    label: string
    to: string
  }

  const navItems: NavItem[] = [
    { label: 'HOME', to: '/' },
    { label: 'SHOP', to: '/shop' },
    { label: 'STAFF', to: '/staff' },
    { label: 'SCHEDULE', to: '/schedule' },
    { label: 'PRICE', to: '/price' },
    { label: 'ACCESS', to: '/access' },
  ]

  const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="bg-black border-t border-white/10 py-16">
    <div class="container mx-auto px-6">
      <div class="flex flex-col items-center gap-10">
        <!-- ロゴ -->
        <span class="tracking-[0.6em] text-xl font-light text-foreground uppercase"> TIARA </span>

        <!-- ナビリンク -->
        <nav class="flex flex-wrap justify-center gap-x-8 gap-y-4">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="text-[11px] tracking-[0.25em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
          >
            {{ item.label }}
          </RouterLink>
        </nav>

        <!-- ディバイダー -->
        <div class="w-12 h-px bg-white/20" />

        <!-- コピーライト -->
        <p class="text-[11px] text-muted-foreground tracking-widest uppercase">
          &copy; {{ currentYear }}&nbsp;TIARA. ALL RIGHTS RESERVED.
        </p>
      </div>
    </div>
  </footer>
</template>
