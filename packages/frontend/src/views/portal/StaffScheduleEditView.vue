<script setup lang="ts">
  import { ref, onMounted, watch } from 'vue'
  import StaffPortalLayout from '@/components/layout/StaffPortalLayout.vue'
  import { useStaffPortalApi } from '@/composables/useStaffPortalApi'
  import type { ScheduleDraftItem } from '@/types/staffPortal'

  const {
    scheduleDraft,
    isLoading,
    error,
    saveError,
    fetchMyScheduleDraft,
    saveScheduleDraft,
    submitScheduleDraft,
  } = useStaffPortalApi()

  const items = ref<ScheduleDraftItem[]>([])
  const isSaving = ref(false)
  const successMessage = ref<string | null>(null)

  const dayLabels = ['日', '月', '火', '水', '木', '金', '土']

  onMounted(async () => {
    await fetchMyScheduleDraft()
    if (scheduleDraft.value?.items) {
      items.value = scheduleDraft.value.items.map((item) => ({
        dayOfWeek: item.dayOfWeek,
        startTime: item.startTime,
        endTime: item.endTime,
      }))
    }
  })

  watch(saveError, (v) => {
    if (v) window.scrollTo({ top: 0, behavior: 'smooth' })
  })

  /** 出勤枠を追加 */
  function addItem() {
    items.value.push({ dayOfWeek: 1, startTime: '20:00', endTime: '01:00' })
  }

  /** 出勤枠を削除 */
  function removeItem(index: number) {
    items.value.splice(index, 1)
  }

  /** 下書き保存 */
  async function handleSave() {
    isSaving.value = true
    successMessage.value = null

    const result = await saveScheduleDraft(items.value)
    if (result) {
      successMessage.value = '下書きを保存しました'
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isSaving.value = false
  }

  /** 承認申請 */
  async function handleSubmit() {
    if (!confirm('この内容で承認申請しますか？')) return

    isSaving.value = true
    successMessage.value = null

    // 先に保存してから申請
    const saved = await saveScheduleDraft(items.value)
    if (saved?.id) {
      const result = await submitScheduleDraft(saved.id)
      if (result) {
        successMessage.value = '承認申請しました。管理者の確認をお待ちください。'
      }
    }
    isSaving.value = false
  }

  /** 編集可能かどうか（approved 以外は全て編集可能） */
  function isEditable(): boolean {
    return (
      !scheduleDraft.value?.status ||
      scheduleDraft.value.status === 'draft' ||
      scheduleDraft.value.status === 'rejected' ||
      scheduleDraft.value.status === 'pending'
    )
  }

  /** 提出可能かどうか（draft / rejected のみ） */
  function isSubmittable(): boolean {
    return (
      !!scheduleDraft.value?.id &&
      (scheduleDraft.value.status === 'draft' || scheduleDraft.value.status === 'rejected')
    )
  }
</script>

<template>
  <StaffPortalLayout>
    <div
      v-motion
      :initial="{ opacity: 0, y: 20 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
    >
      <h2 class="text-lg font-light tracking-wider text-foreground mb-6">スケジュール編集</h2>

      <div v-if="isLoading" class="text-sm text-muted-foreground">読み込み中...</div>

      <div v-if="error" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4">
        <p class="text-sm text-red-400">{{ error }}</p>
      </div>

      <div
        v-if="saveError"
        class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-red-400">{{ saveError }}</p>
      </div>

      <div
        v-if="successMessage"
        class="bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-green-400">{{ successMessage }}</p>
      </div>

      <div
        v-if="scheduleDraft?.status === 'pending'"
        class="bg-blue-500/10 border border-blue-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-blue-400">
          承認待ちです。編集して保存すると申請が取り下げられ、再提出が必要になります。
        </p>
      </div>

      <div
        v-if="scheduleDraft?.status === 'rejected' && scheduleDraft.adminComment"
        class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-red-400">
          <strong>却下理由:</strong> {{ scheduleDraft.adminComment }}
        </p>
        <p class="text-xs text-red-400/70 mt-1">内容を修正して再申請してください。</p>
      </div>

      <!-- スケジュール一覧 -->
      <div v-if="!isLoading" class="space-y-4">
        <div
          v-for="(item, index) in items"
          :key="index"
          class="flex items-center gap-3 border border-white/10 rounded-lg p-4"
        >
          <!-- 曜日 -->
          <select
            v-model.number="item.dayOfWeek"
            :disabled="!isEditable()"
            class="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
          >
            <option v-for="(label, i) in dayLabels" :key="i" :value="i">{{ label }}</option>
          </select>

          <!-- 開始時間 -->
          <input
            v-model="item.startTime"
            type="time"
            :disabled="!isEditable()"
            class="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
          />

          <span class="text-muted-foreground text-sm">〜</span>

          <!-- 終了時間 -->
          <input
            v-model="item.endTime"
            type="time"
            :disabled="!isEditable()"
            class="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
          />

          <!-- 削除ボタン -->
          <button
            v-if="isEditable()"
            type="button"
            @click="removeItem(index)"
            class="text-red-400 hover:text-red-300 transition-colors text-sm ml-auto"
          >
            削除
          </button>
        </div>

        <!-- 追加ボタン -->
        <button
          v-if="isEditable()"
          type="button"
          @click="addItem"
          class="w-full border border-dashed border-white/20 rounded-lg py-3 text-sm text-muted-foreground hover:text-foreground hover:border-white/40 transition-colors"
        >
          + 出勤枠を追加
        </button>

        <!-- ボタン群 -->
        <div v-if="isEditable()" class="flex gap-4 pt-4">
          <button
            type="button"
            :disabled="isSaving"
            @click="handleSave"
            class="bg-white/10 text-foreground border border-white/20 rounded-lg px-6 py-3 text-sm font-medium tracking-wider hover:bg-white/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isSaving ? '保存中...' : '下書き保存' }}
          </button>

          <button
            type="button"
            :disabled="isSaving || !isSubmittable()"
            @click="handleSubmit"
            class="bg-white text-black rounded-lg px-6 py-3 text-sm font-medium tracking-wider hover:bg-white/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ scheduleDraft?.status === 'pending' ? '提出済み' : '承認申請' }}
          </button>
        </div>
      </div>
    </div>
  </StaffPortalLayout>
</template>
