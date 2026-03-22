<script setup lang="ts">
  import { onMounted } from 'vue'
  import { useStaffApi } from '@/composables/useStaffApi'

  const { staffList, isLoading, error, fetchStaffs } = useStaffApi()

  onMounted(() => {
    fetchStaffs()
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
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Our Team</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Staff
        </h1>
        <span class="block w-12 h-px bg-primary mt-6" />
      </div>

      <!-- ローディング -->
      <div v-if="isLoading" class="flex justify-center py-20">
        <div class="w-6 h-6 border border-primary border-t-transparent rounded-full animate-spin" />
      </div>

      <!-- エラー -->
      <div v-else-if="error" class="text-center py-20">
        <p class="text-destructive text-sm tracking-wide">{{ error }}</p>
        <button
          class="mt-6 text-xs tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground border border-border px-8 py-3 transition-colors duration-300"
          @click="fetchStaffs"
        >
          Retry
        </button>
      </div>

      <!-- データなし -->
      <div v-else-if="staffList.length === 0" class="text-center py-20">
        <p class="text-muted-foreground tracking-widest text-sm uppercase">
          No staff information available
        </p>
      </div>

      <!-- スタッフカード一覧（縦一列・横長カード） -->
      <div v-else class="max-w-4xl mx-auto space-y-6">
        <router-link
          v-for="(staff, index) in staffList"
          :key="staff.id"
          :to="`/staff/${staff.id}`"
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :visible="{ opacity: 1, y: 0, transition: { duration: 600, delay: index * 150 } }"
          class="group block"
        >
          <div
            class="border border-border bg-card overflow-hidden transition-all duration-500 hover:border-primary/50 hover:shadow-lg hover:shadow-primary/5 flex flex-col sm:flex-row sm:h-52"
          >
            <!-- スタッフ画像 -->
            <div
              class="sm:w-48 md:w-56 flex-shrink-0 aspect-[4/3] sm:aspect-auto overflow-hidden bg-secondary"
            >
              <img
                v-if="staff.imageUrl"
                :src="staff.imageUrl"
                :alt="staff.name"
                class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
                :style="{
                  objectPosition: staff.imageCropPosition
                    ? staff.imageCropPosition
                        .split(' ')
                        .map((v: string) => v + '%')
                        .join(' ')
                    : '50% 50%',
                }"
              />
              <div v-else class="w-full h-full flex items-center justify-center min-h-[160px]">
                <span class="text-muted-foreground text-sm tracking-widest uppercase"
                  >No Photo</span
                >
              </div>
            </div>

            <!-- スタッフ情報 -->
            <div class="p-6 flex flex-col justify-center flex-1 min-w-0">
              <p class="text-[11px] tracking-[0.2em] uppercase text-primary mb-2">
                {{ staff.role }}
              </p>
              <h2
                class="text-lg font-light tracking-[0.1em] text-foreground group-hover:text-primary transition-colors duration-300"
              >
                {{ staff.name }}
              </h2>
              <p class="mt-3 text-xs text-muted-foreground leading-relaxed line-clamp-3">
                {{ staff.bio }}
              </p>

              <!-- View Profile リンク表示 -->
              <div
                class="mt-4 flex items-center gap-2 text-[11px] tracking-[0.2em] uppercase text-muted-foreground group-hover:text-foreground transition-colors duration-300"
              >
                <span>View Profile</span>
                <span
                  class="inline-block w-4 h-px bg-current transition-all duration-300 group-hover:w-8"
                />
              </div>
            </div>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>
