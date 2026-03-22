<script setup lang="ts">
  import { onMounted, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useStaffApi } from '@/composables/useStaffApi'
  import { ArrowLeft, Calendar, Clock } from 'lucide-vue-next'
  import type { StaffSchedule } from '@/types/staff'

  const route = useRoute()
  const router = useRouter()
  const { staffDetail, isLoading, error, fetchStaffById } = useStaffApi()

  /** 曜日ラベル（0=日 … 6=土） */
  const DAY_LABELS = ['日', '月', '火', '水', '木', '金', '土'] as const

  onMounted(() => {
    const id = route.params.id as string
    fetchStaffById(id)
  })

  /** TIME型のISO文字列から "HH:MM" を抽出 */
  function formatTime(raw: string): string {
    const date = new Date(raw)
    const hours = date.getUTCHours().toString().padStart(2, '0')
    const minutes = date.getUTCMinutes().toString().padStart(2, '0')
    return `${hours}:${minutes}`
  }

  /** スケジュールを曜日順にソートして返す */
  const sortedSchedules = computed<StaffSchedule[]>(() => {
    if (!staffDetail.value?.schedules) return []
    return [...staffDetail.value.schedules].sort((a, b) => a.dayOfWeek - b.dayOfWeek)
  })

  /**
   * 今日(JST)を基準に、今週（日〜土）の各曜日の日付を返す。
   * index: 0=日, 1=月, … 6=土
   */
  const weekDates = computed(() => {
    const now = new Date()
    const jstNow = new Date(now.getTime() + (9 * 60 - now.getTimezoneOffset()) * 60 * 1000)
    const todayDay = jstNow.getUTCDay()
    const dates: Date[] = []
    for (let i = 0; i < 7; i++) {
      const d = new Date(jstNow)
      d.setUTCDate(d.getUTCDate() - todayDay + i)
      dates.push(d)
    }
    return dates
  })

  /** Date を "M/D" 形式にフォーマット */
  function formatDateShort(d: Date): string {
    return `${d.getUTCMonth() + 1}/${d.getUTCDate()}`
  }
</script>

<template>
  <div class="min-h-screen bg-background pt-24 pb-28">
    <div class="container mx-auto px-6">
      <!-- 戻るリンク -->
      <div
        v-motion
        :initial="{ opacity: 0, x: -10 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500 } }"
        class="mb-12"
      >
        <button
          class="flex items-center gap-2 text-xs tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
          @click="router.push('/staff')"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>Back to Staff</span>
        </button>
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
          @click="fetchStaffById(route.params.id as string)"
        >
          Retry
        </button>
      </div>

      <!-- データなし -->
      <div v-else-if="!staffDetail" class="text-center py-20">
        <p class="text-muted-foreground tracking-widest text-sm uppercase">Staff not found</p>
      </div>

      <!-- スタッフ詳細 -->
      <div v-else class="max-w-4xl mx-auto">
        <!-- プロフィールセクション -->
        <div
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :enter="{ opacity: 1, y: 0, transition: { duration: 700 } }"
          class="grid grid-cols-1 md:grid-cols-2 gap-10 md:gap-14 mb-20"
        >
          <!-- 写真 -->
          <div class="aspect-[3/4] overflow-hidden border border-border bg-secondary">
            <img
              v-if="staffDetail.staff.imageUrl"
              :src="staffDetail.staff.imageUrl"
              :alt="staffDetail.staff.name"
              class="w-full h-full object-cover"
            />
            <div v-else class="w-full h-full flex items-center justify-center">
              <span class="text-muted-foreground text-sm tracking-widest uppercase">No Photo</span>
            </div>
          </div>

          <!-- プロフィール情報 -->
          <div class="flex flex-col justify-center">
            <p class="text-primary text-[11px] tracking-[0.3em] uppercase mb-3">
              {{ staffDetail.staff.role }}
            </p>
            <h1 class="text-3xl md:text-4xl font-light tracking-[0.15em] text-foreground mb-6">
              {{ staffDetail.staff.name }}
            </h1>
            <span class="block w-10 h-px bg-primary mb-8" />
            <p class="text-sm text-muted-foreground leading-[1.8] whitespace-pre-wrap">
              {{ staffDetail.staff.bio }}
            </p>
          </div>
        </div>

        <!-- スケジュールセクション -->
        <div
          v-if="sortedSchedules.length > 0"
          v-motion
          :initial="{ opacity: 0, y: 30 }"
          :visible="{ opacity: 1, y: 0, transition: { duration: 600, delay: 200 } }"
        >
          <!-- セクションヘッダー -->
          <div class="flex flex-col items-center mb-12 text-center">
            <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Schedule</span>
            <h2 class="text-2xl font-light tracking-[0.2em] uppercase text-foreground">
              Weekly Schedule
            </h2>
            <span class="block w-8 h-px bg-primary mt-5" />
          </div>

          <!-- スケジュールテーブル -->
          <div class="border border-border bg-card overflow-hidden">
            <table class="w-full">
              <thead>
                <tr class="border-b border-border">
                  <th class="px-6 py-4 text-left">
                    <div class="flex items-center gap-2">
                      <Calendar class="w-3.5 h-3.5 text-primary" />
                      <span
                        class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground font-normal"
                        >Day</span
                      >
                    </div>
                  </th>
                  <th class="px-6 py-4 text-left">
                    <div class="flex items-center gap-2">
                      <Clock class="w-3.5 h-3.5 text-primary" />
                      <span
                        class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground font-normal"
                        >Hours</span
                      >
                    </div>
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="schedule in sortedSchedules"
                  :key="schedule.id"
                  class="border-b border-border last:border-b-0 transition-colors duration-300 hover:bg-secondary/50"
                >
                  <td class="px-6 py-4">
                    <div class="flex items-baseline gap-2">
                      <span
                        class="text-sm tracking-[0.1em] font-light"
                        :class="
                          schedule.dayOfWeek === 0
                            ? 'text-red-400'
                            : schedule.dayOfWeek === 6
                              ? 'text-blue-400'
                              : 'text-foreground'
                        "
                      >
                        {{ DAY_LABELS[schedule.dayOfWeek] }}曜日
                      </span>
                      <span class="text-xs text-muted-foreground/60">
                        {{ formatDateShort(weekDates[schedule.dayOfWeek]) }}
                      </span>
                    </div>
                  </td>
                  <td class="px-6 py-4">
                    <span class="text-sm text-muted-foreground tracking-wider">
                      {{ formatTime(schedule.startTime) }} – {{ formatTime(schedule.endTime) }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
