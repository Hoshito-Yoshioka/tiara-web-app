<script setup lang="ts">
  import { onMounted, computed } from 'vue'
  import { ExternalLink, CalendarDays, UserRound } from 'lucide-vue-next'
  import { useScheduleApi } from '@/composables/useScheduleApi'

  const { scheduleData, isLoading, error, fetchSchedules } = useScheduleApi()

  const shopScheduleUrl = 'https://www.pokepara.jp/_hokkaido/m824/a1990/shop24197/'

  const linkedStaffs = computed(() => {
    return scheduleData.value
      .filter(
        (item) => item.staff.externalScheduleUrl && item.staff.externalScheduleUrl.trim() !== ''
      )
      .sort((a, b) => a.staff.sortOrder - b.staff.sortOrder)
  })

  onMounted(() => {
    fetchSchedules()
  })
</script>

<template>
  <div class="min-h-screen bg-background pt-24 pb-28">
    <div class="container mx-auto px-6">
      <div
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 700 } }"
        class="flex flex-col items-center mb-16 text-center"
      >
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Attendance</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Schedule
        </h1>
        <p class="mt-5 text-sm md:text-base text-muted-foreground max-w-2xl leading-relaxed">
          出勤情報は外部サイトにて公開しています。以下のリンクから、店舗全体またはスタッフ個別の
          最新スケジュールをご確認ください。
        </p>
      </div>

      <div v-if="isLoading" class="flex justify-center py-20">
        <div class="w-6 h-6 border border-primary border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else-if="error" class="text-center py-20">
        <p class="text-destructive text-sm tracking-wide">{{ error }}</p>
        <button
          class="mt-6 text-xs tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground border border-border px-8 py-3 transition-colors duration-300"
          @click="fetchSchedules"
        >
          Retry
        </button>
      </div>

      <div v-else class="max-w-6xl mx-auto space-y-14">
        <section
          v-motion
          :initial="{ opacity: 0, y: 25 }"
          :enter="{ opacity: 1, y: 0, transition: { duration: 700, delay: 150 } }"
          class="rounded-2xl border border-primary/30 bg-gradient-to-br from-card via-card to-primary/5 p-7 md:p-9"
        >
          <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-6">
            <div class="space-y-3">
              <div
                class="inline-flex items-center gap-2 text-primary text-xs tracking-[0.22em] uppercase"
              >
                <CalendarDays class="w-4 h-4" />
                店舗全体のスケジュール
              </div>
              <h2 class="text-xl md:text-2xl font-light tracking-[0.08em] text-foreground">
                ポケパラ掲載ページで最新の出勤情報を見る
              </h2>
              <p class="text-sm text-muted-foreground leading-relaxed max-w-2xl">
                各スタッフの出勤日・詳細情報は、ポケパラの公式掲載ページで最新状態を確認できます。
              </p>
            </div>
            <a
              :href="shopScheduleUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-full border border-primary/50 text-primary hover:bg-primary hover:text-primary-foreground transition-colors duration-300 text-sm tracking-wide"
            >
              店舗ページへ移動
              <ExternalLink class="w-4 h-4" />
            </a>
          </div>
        </section>

        <section
          v-motion
          :initial="{ opacity: 0, y: 25 }"
          :enter="{ opacity: 1, y: 0, transition: { duration: 700, delay: 250 } }"
          class="space-y-6"
        >
          <div class="flex items-center justify-between gap-4">
            <h3
              class="text-base md:text-lg tracking-[0.14em] uppercase text-foreground flex items-center gap-2"
            >
              <UserRound class="w-4 h-4 text-primary" />
              スタッフ個別スケジュール
            </h3>
            <span class="text-xs text-muted-foreground">
              {{ linkedStaffs.length }} 名が外部ページに紐づいています
            </span>
          </div>

          <div v-if="linkedStaffs.length > 0" class="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
            <article
              v-for="item in linkedStaffs"
              :key="item.staff.id"
              class="group rounded-xl border border-white/10 bg-card/70 p-5 hover:border-primary/40 hover:bg-card transition-colors duration-300"
            >
              <div class="flex items-center gap-3 mb-4">
                <img
                  v-if="item.staff.imageUrl"
                  :src="item.staff.imageUrl"
                  :alt="item.staff.name"
                  class="w-12 h-12 rounded-full object-cover border border-white/15"
                  :style="{
                    objectPosition: item.staff.imageCropPosition
                      .split(' ')
                      .map((v: string) => v + '%')
                      .join(' '),
                  }"
                />
                <div
                  v-else
                  class="w-12 h-12 rounded-full bg-secondary border border-white/15 flex items-center justify-center text-sm text-muted-foreground"
                >
                  {{ item.staff.name.charAt(0) }}
                </div>
                <div class="min-w-0">
                  <p class="text-foreground text-sm font-light tracking-wide truncate">
                    {{ item.staff.name }}
                  </p>
                  <p class="text-[11px] text-muted-foreground tracking-[0.14em] uppercase truncate">
                    {{ item.staff.role }}
                  </p>
                </div>
              </div>

              <a
                :href="item.staff.externalScheduleUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 text-sm text-primary group-hover:text-primary/80 transition-colors"
              >
                {{ item.staff.name }} の出勤情報を見る
                <ExternalLink class="w-4 h-4" />
              </a>
            </article>
          </div>

          <div
            v-else
            class="rounded-xl border border-white/10 bg-card/60 p-8 text-center text-sm text-muted-foreground"
          >
            現在、個別ページが登録されているスタッフはありません。
            <br />
            店舗全体のスケジュールページをご利用ください。
          </div>
        </section>
      </div>
    </div>
  </div>
</template>
