<script setup lang="ts">
  import { ref, onMounted, onUnmounted } from 'vue'
  import { RouterLink, useRoute } from 'vue-router'
  import { Menu, X } from 'lucide-vue-next'

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

  const isMenuOpen = ref(false)
  const isScrolled = ref(false)
  const route = useRoute()

  const handleScroll = () => {
    isScrolled.value = window.scrollY > 40
  }

  const toggleMenu = () => {
    isMenuOpen.value = !isMenuOpen.value
  }

  const closeMenu = () => {
    isMenuOpen.value = false
  }

  onMounted(() => {
    window.addEventListener('scroll', handleScroll, { passive: true })
  })

  onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
  })
</script>

<template>
  <header
    class="fixed top-0 left-0 right-0 z-50 transition-all duration-500"
    :class="isScrolled ? 'bg-black/95 backdrop-blur-sm border-b border-white/10' : 'bg-transparent'"
  >
    <div class="container mx-auto px-6 h-16 flex items-center justify-between">
      <!-- ロゴ -->
      <RouterLink
        to="/"
        class="tracking-[0.5em] text-base font-light text-foreground hover:text-primary transition-colors duration-500 uppercase"
        @click="closeMenu"
      >
        TIARA
      </RouterLink>

      <!-- デスクトップナビ -->
      <nav class="hidden md:flex items-center gap-10">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="relative text-[11px] tracking-[0.25em] uppercase transition-colors duration-300"
          :class="
            route.path === item.to ? 'text-primary' : 'text-muted-foreground hover:text-foreground'
          "
        >
          {{ item.label }}
          <!-- アクティブなリンクの下線 -->
          <span
            v-if="route.path === item.to"
            class="absolute -bottom-1 left-0 right-0 h-px bg-primary"
          />
        </RouterLink>
      </nav>

      <!-- モバイル ハンバーガーボタン -->
      <button
        class="md:hidden text-foreground p-2 -mr-2 transition-opacity duration-300 hover:opacity-70"
        :aria-label="isMenuOpen ? 'メニューを閉じる' : 'メニューを開く'"
        @click="toggleMenu"
      >
        <Transition
          enter-active-class="transition-all duration-200"
          enter-from-class="opacity-0 rotate-90"
          enter-to-class="opacity-100 rotate-0"
          leave-active-class="transition-all duration-200"
          leave-from-class="opacity-100 rotate-0"
          leave-to-class="opacity-0 rotate-90"
          mode="out-in"
        >
          <X v-if="isMenuOpen" :size="20" />
          <Menu v-else :size="20" />
        </Transition>
      </button>
    </div>

    <!-- モバイルメニュー -->
    <Transition
      enter-active-class="transition-all duration-500 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-300 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="isMenuOpen"
        class="md:hidden absolute top-16 left-0 right-0 bg-black/98 border-b border-white/10 py-10"
      >
        <nav class="flex flex-col items-center gap-8">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="text-sm tracking-[0.3em] uppercase transition-colors duration-300"
            :class="
              route.path === item.to
                ? 'text-primary'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="closeMenu"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
      </div>
    </Transition>
  </header>
</template>
