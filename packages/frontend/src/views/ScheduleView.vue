<script setup lang="ts">
  import { onMounted, computed } from 'vue'
  import { useScheduleApi } from '@/composables/useScheduleApi'
  import { Circle, Minus } from 'lucide-vue-next'

  const { scheduleData, isLoading, error, fetchSchedules } = useScheduleApi()

  /** 曜日ラベル（0=日 … 6=土） */
  const DAY_LABELS = ['日', '月', '火', '水', '木', '金', '土'] as const

  /** 曜日の色クラス */
  function dayColorClass(dayIndex: number): string {
    if (dayIndex === 0) return 'text-red-400'
    if (dayIndex === 6) return 'text-blue-400'
    return 'text-foreground'
  }

  /**
   * 今日(JST)を基準に、今週（日〜土）の各曜日の日付を返す。
   * index: 0=日, 1=月, … 6=土
   */
  const weekDates = computed(() => {
    const now = new Date()
    // JSTオフセット (+9h)
    const jstNow = new Date(now.getTime() + (9 * 60 - now.getTimezoneOffset()) * 60 * 1000)
    const todayDay = jstNow.getUTCDay() // 0=日
    const dates: Date[] = []
    for (let i = 0; i < 7; i++) {
      const d = new Date(jstNow)
      d.setUTCDate(d.getUTCDate() - todayDay + i)
      dates.push(d)
    }
    return dates
  })

  /** Date を "M/D" 形式にフォーマット */
  function formatDate(d: Date): string {
    return `${d.getUTCMonth() + 1}/${d.getUTCDate()}`
  }

  /**
   * スタッフごとに出勤曜日を Set で持つマップを生成。
   * key: staffId, value: Set<dayOfWeek>
   */
  const scheduleMap = computed(() => {
    const map = new Map<string, Set<number>>()
    for (const item of scheduleData.value) {
      const days = new Set<number>()
      for (const s of item.schedules) {
        days.add(s.dayOfWeek)
      }
      map.set(item.staff.id, days)
    }
    return map
  })

  /** TIME型のISO文字列から "HH:MM" を抽出 */
  function formatTime(raw: string): string {
    const date = new Date(raw)
    const hours = date.getUTCHours().toString().padStart(2, '0')
    const minutes = date.getUTCMinutes().toString().padStart(2, '0')
    return `${hours}:${minutes}`
  }

  /** スタッフの特定曜日のスケジュールを取得 */
  function getScheduleForDay(staffId: string, dayOfWeek: number) {
    const item = scheduleData.value.find((d) => d.staff.id === staffId)
    return item?.schedules.find((s) => s.dayOfWeek === dayOfWeek)
  }

  /** 今日(JST)の曜日 */
  const todayDayOfWeek = computed(() => {
    const now = new Date()
    const jstNow = new Date(now.getTime() + (9 * 60 - now.getTimezoneOffset()) * 60 * 1000)
    return jstNow.getUTCDay()
  })

  onMounted(() => {
    fetchSchedules()
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
        <span class="text-primary text-[11px] tracking-[0.4em] uppercase mb-4">Attendance</span>
        <h1 class="text-3xl md:text-4xl font-light tracking-[0.2em] uppercase text-foreground">
          Schedule
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
          @click="fetchSchedules"
        >
          Retry
        </button>
      </div>

      <!-- データなし -->
      <div v-else-if="scheduleData.length === 0" class="text-center py-20">
        <p class="text-muted-foreground tracking-widest text-sm uppercase">
          No schedule information available
        </p>
      </div>

      <!-- スケジュールマトリクス -->
      <div
        v-else
        v-motion
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 700, delay: 200 } }"
        class="max-w-4xl mx-auto"
      >
        <!-- テーブル（横スクロール対応） -->
        <div class="border border-border bg-card overflow-x-auto">
          <table class="w-full min-w-[640px]">
            <thead>
              <tr class="border-b border-border">
                <!-- スタッフ名ヘッダー -->
                <th class="px-4 md:px-6 py-4 text-left sticky left-0 bg-card z-10">
                  <span
                    class="text-[11px] tracking-[0.2em] uppercase text-muted-foreground font-normal"
                  >
                    Staff
                  </span>
                </th>
                <!-- 曜日＋日付ヘッダー -->
                <th
                  v-for="(label, idx) in DAY_LABELS"
                  :key="idx"
                  class="px-3 md:px-4 py-4 text-center"
                  :class="idx === todayDayOfWeek ? 'bg-primary/5' : ''"
                >
                  <div class="flex flex-col items-center gap-0.5">
                    <span
                      class="text-[11px] tracking-[0.15em] font-normal"
                      :class="dayColorClass(idx)"
                    >
                      {{ label }}
                    </span>
                    <span class="text-[10px] text-muted-foreground">
                      {{ formatDate(weekDates[idx]) }}
                    </span>
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="item in scheduleData"
                :key="item.staff.id"
                class="border-b border-border last:border-b-0 transition-colors duration-300 hover:bg-secondary/50"
              >
                <!-- スタッフ名＋画像 -->
                <td class="px-4 md:px-6 py-4 sticky left-0 bg-card z-10">
                  <router-link
                    :to="`/staff/${item.staff.id}`"
                    class="group flex items-center gap-3"
                  >
                    <img
                      v-if="item.staff.imageUrl"
                      :src="item.staff.imageUrl"
                      :alt="item.staff.name"
                      class="w-9 h-9 rounded-full object-cover flex-shrink-0 border border-border"
                    />
                    <div
                      v-else
                      class="w-9 h-9 rounded-full bg-secondary flex items-center justify-center flex-shrink-0 border border-border"
                    >
                      <span class="text-[10px] text-muted-foreground">{{
                        item.staff.name.charAt(0)
                      }}</span>
                    </div>
                    <div class="flex flex-col min-w-0">
                      <span
                        class="text-sm font-light tracking-[0.05em] text-foreground group-hover:text-primary transition-colors duration-300"
                      >
                        {{ item.staff.name }}
                      </span>
                      <span
                        class="text-[10px] tracking-[0.15em] uppercase text-muted-foreground mt-0.5"
                      >
                        {{ item.staff.role }}
                      </span>
                    </div>
                  </router-link>
                </td>
                <!-- 各曜日のセル -->
                <td
                  v-for="(_, dayIdx) in DAY_LABELS"
                  :key="dayIdx"
                  class="px-3 md:px-4 py-4 text-center"
                  :class="dayIdx === todayDayOfWeek ? 'bg-primary/5' : ''"
                >
                  <template v-if="scheduleMap.get(item.staff.id)?.has(dayIdx)">
                    <div class="flex flex-col items-center gap-0.5">
                      <Circle class="w-3.5 h-3.5 text-primary fill-primary" />
                      <span class="text-[9px] text-muted-foreground tracking-wide">
                        {{ formatTime(getScheduleForDay(item.staff.id, dayIdx)!.startTime) }}
                      </span>
                    </div>
                  </template>
                  <template v-else>
                    <Minus class="w-3.5 h-3.5 text-muted-foreground/30 mx-auto" />
                  </template>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 凡例 -->
        <div
          v-motion
          :initial="{ opacity: 0 }"
          :visible="{ opacity: 1, transition: { duration: 500, delay: 400 } }"
          class="mt-6 flex items-center justify-center gap-6 text-[11px] text-muted-foreground tracking-wide"
        >
          <div class="flex items-center gap-2">
            <Circle class="w-3 h-3 text-primary fill-primary" />
            <span>出勤</span>
          </div>
          <div class="flex items-center gap-2">
            <Minus class="w-3 h-3 text-muted-foreground/30" />
            <span>休み</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
