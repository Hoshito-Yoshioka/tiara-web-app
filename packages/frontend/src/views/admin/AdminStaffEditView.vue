<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import AdminLayout from '@/components/layout/AdminLayout.vue'
  import { useStaffApi } from '@/composables/useStaffApi'
  import { useShopApi } from '@/composables/useShopApi'
  import { useAdminApi } from '@/composables/useAdminApi'
  import { Save, ArrowLeft, Plus, X } from 'lucide-vue-next'

  const route = useRoute()
  const router = useRouter()

  const { staffDetail, fetchStaffById } = useStaffApi()
  const { shops, fetchShops } = useShopApi()
  const { createStaff, updateStaff, isLoading: isSaving, error: saveError } = useAdminApi()

  const isEditMode = computed(() => !!route.params.id)
  const pageTitle = computed(() => (isEditMode.value ? 'スタッフ編集' : 'スタッフ新規作成'))

  const DAYS = ['日', '月', '火', '水', '木', '金', '土']

  const form = ref({
    name: '',
    role: '',
    bio: '',
    imageUrl: '',
    sortOrder: 0,
    shopId: '',
  })

  interface ScheduleForm {
    dayOfWeek: number
    startTime: string
    endTime: string
  }

  const schedules = ref<ScheduleForm[]>([])
  const isLoading = ref(false)
  const successMessage = ref<string | null>(null)

  /**
   * Backend から返る時刻文字列を "HH:MM" 形式に変換する。
   * "0000-01-01T18:00:00Z" → "18:00"
   */
  function formatTimeForInput(timeStr: string): string {
    if (timeStr.includes('T')) {
      const match = timeStr.match(/T(\d{2}:\d{2})/)
      return match ? match[1] : timeStr
    }
    return timeStr.slice(0, 5)
  }

  onMounted(async () => {
    isLoading.value = true
    await fetchShops()

    if (isEditMode.value) {
      // 編集モード: 既存データでフォームを初期化
      await fetchStaffById(route.params.id as string)
      if (staffDetail.value) {
        const { staff, schedules: staffSchedules } = staffDetail.value
        form.value = {
          name: staff.name,
          role: staff.role,
          bio: staff.bio,
          imageUrl: staff.imageUrl,
          sortOrder: staff.sortOrder,
          shopId: staff.shopId,
        }
        schedules.value = staffSchedules.map((s) => ({
          dayOfWeek: s.dayOfWeek,
          startTime: formatTimeForInput(s.startTime),
          endTime: formatTimeForInput(s.endTime),
        }))
      }
    } else {
      // 新規作成モード: デフォルト店舗を設定
      if (shops.value.length > 0) {
        form.value.shopId = shops.value[0].id
      }
    }

    isLoading.value = false
  })

  /** 出勤スケジュールを追加 */
  function addSchedule() {
    schedules.value.push({
      dayOfWeek: 1,
      startTime: '20:00',
      endTime: '02:00',
    })
  }

  /** 出勤スケジュールを削除 */
  function removeSchedule(index: number) {
    schedules.value.splice(index, 1)
  }

  async function handleSubmit() {
    successMessage.value = null

    const scheduleData = schedules.value.map((s) => ({
      dayOfWeek: s.dayOfWeek,
      startTime: s.startTime,
      endTime: s.endTime,
    }))

    if (isEditMode.value) {
      const result = await updateStaff(route.params.id as string, {
        ...form.value,
        schedules: scheduleData,
      })
      if (result) {
        successMessage.value = 'スタッフ情報を更新しました'
        window.scrollTo({ top: 0, behavior: 'smooth' })
        setTimeout(() => {
          successMessage.value = null
        }, 3000)
      }
    } else {
      const result = await createStaff({
        ...form.value,
        schedules: scheduleData,
      })
      if (result) {
        router.push({ name: 'admin-staff-list' })
      }
    }
  }
</script>

<template>
  <AdminLayout>
    <div
      class="max-w-2xl"
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 400 } }"
    >
      <!-- ヘッダー -->
      <div class="flex items-center gap-4 mb-8">
        <button
          @click="router.push({ name: 'admin-staff-list' })"
          class="p-2 text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft :size="20" />
        </button>
        <h1 class="text-xl font-light tracking-wider text-foreground">{{ pageTitle }}</h1>
      </div>

      <div v-if="isLoading" class="text-muted-foreground text-sm">読み込み中...</div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-6">
        <!-- エラー表示 -->
        <div v-if="saveError" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3">
          <p class="text-sm text-red-400">{{ saveError }}</p>
        </div>

        <!-- 成功メッセージ -->
        <div
          v-if="successMessage"
          class="bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3"
        >
          <p class="text-sm text-green-400">{{ successMessage }}</p>
        </div>

        <!-- 所属店舗（新規作成時のみ表示） -->
        <div v-if="!isEditMode" class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">所属店舗</label>
          <select
            v-model="form.shopId"
            required
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors [color-scheme:dark]"
          >
            <option v-for="shop in shops" :key="shop.id" :value="shop.id">
              {{ shop.name }}
            </option>
          </select>
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">名前</label>
          <input
            v-model="form.name"
            type="text"
            required
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">役職</label>
          <input
            v-model="form.role"
            type="text"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">自己紹介</label>
          <textarea
            v-model="form.bio"
            rows="4"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors resize-none"
          ></textarea>
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">画像URL</label>
          <input
            v-model="form.imageUrl"
            type="url"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors"
            placeholder="https://..."
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">表示順</label>
          <input
            v-model.number="form.sortOrder"
            type="number"
            min="0"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
          />
        </div>

        <!-- 出勤スケジュール -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <label class="text-xs text-muted-foreground tracking-wider uppercase">
              出勤スケジュール
            </label>
            <button
              type="button"
              @click="addSchedule"
              class="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              <Plus :size="14" />
              追加
            </button>
          </div>

          <div
            v-for="(schedule, index) in schedules"
            :key="index"
            class="flex items-center gap-3 bg-white/5 border border-white/10 rounded-lg px-4 py-3"
          >
            <select
              v-model.number="schedule.dayOfWeek"
              class="bg-zinc-900 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30"
            >
              <option v-for="(day, i) in DAYS" :key="i" :value="i">{{ day }}</option>
            </select>

            <input
              v-model="schedule.startTime"
              type="time"
              class="bg-zinc-900 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
            />

            <span class="text-muted-foreground text-sm">〜</span>

            <input
              v-model="schedule.endTime"
              type="time"
              class="bg-zinc-900 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
            />

            <button
              type="button"
              @click="removeSchedule(index)"
              class="p-1 text-muted-foreground hover:text-red-400 transition-colors"
            >
              <X :size="16" />
            </button>
          </div>

          <p v-if="schedules.length === 0" class="text-xs text-muted-foreground">
            出勤スケジュールが設定されていません
          </p>
        </div>

        <!-- 保存ボタン -->
        <button
          type="submit"
          :disabled="isSaving"
          class="flex items-center gap-2 bg-white text-black rounded-lg px-6 py-3 text-sm font-medium tracking-wider uppercase hover:bg-white/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Save :size="16" />
          {{ isSaving ? '保存中...' : '保存' }}
        </button>
      </form>
    </div>
  </AdminLayout>
</template>
