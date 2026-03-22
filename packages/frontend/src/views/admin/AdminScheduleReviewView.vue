<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useAdminReviewApi } from '@/composables/useAdminReviewApi'
  import type { ScheduleDraft, ScheduleDraftItem } from '@/types/staffPortal'
  import { Check, X, ChevronDown, ChevronUp, Plus, Trash2 } from 'lucide-vue-next'

  const {
    pendingSchedules,
    isLoading,
    error,
    fetchPendingSchedules,
    updateScheduleDraftContent,
    reviewScheduleDraft,
  } = useAdminReviewApi()

  const expandedId = ref<string | null>(null)
  const reviewComment = ref('')
  const editItems = ref<ScheduleDraftItem[]>([])
  const successMessage = ref<string | null>(null)
  const isProcessing = ref(false)

  const dayLabels = ['日', '月', '火', '水', '木', '金', '土']

  /** 元のアイテムのJSON文字列（変更検知用） */
  const originalItemsJson = ref('')

  onMounted(async () => {
    await fetchPendingSchedules()
  })

  /** ドラフトを展開/閉じる */
  function toggleExpand(draft: ScheduleDraft) {
    if (expandedId.value === draft.id) {
      expandedId.value = null
      reviewComment.value = ''
    } else {
      expandedId.value = draft.id!
      reviewComment.value = ''
      // アイテムをコピー
      editItems.value = draft.items.map((item) => ({
        dayOfWeek: item.dayOfWeek,
        startTime: item.startTime,
        endTime: item.endTime,
      }))
      originalItemsJson.value = JSON.stringify(editItems.value)
    }
  }

  /** 出勤枠を追加 */
  function addItem() {
    editItems.value.push({ dayOfWeek: 1, startTime: '20:00', endTime: '01:00' })
  }

  /** 出勤枠を削除 */
  function removeItem(index: number) {
    editItems.value.splice(index, 1)
  }

  /** 内容が変更されたか */
  function hasContentChanged(): boolean {
    return JSON.stringify(editItems.value) !== originalItemsJson.value
  }

  /** 承認 */
  async function handleApprove(draft: ScheduleDraft) {
    if (!draft.id) return
    if (
      !confirm(
        `${draft.staffName || 'スタッフ'} のスケジュール変更を承認しますか？ライブデータに反映されます。`
      )
    )
      return

    isProcessing.value = true

    // 内容が変更されていれば先に保存
    if (hasContentChanged()) {
      const updated = await updateScheduleDraftContent(draft.id, editItems.value)
      if (!updated) {
        isProcessing.value = false
        return
      }
    }

    const result = await reviewScheduleDraft(draft.id, 'approved', reviewComment.value)
    if (result) {
      successMessage.value = `${draft.staffName || 'スタッフ'} のスケジュール変更を承認しました`
      expandedId.value = null
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isProcessing.value = false
  }

  /** 却下 */
  async function handleReject(draft: ScheduleDraft) {
    if (!draft.id) return
    if (!reviewComment.value.trim()) {
      alert('却下理由を入力してください')
      return
    }

    isProcessing.value = true
    const result = await reviewScheduleDraft(draft.id, 'rejected', reviewComment.value)
    if (result) {
      successMessage.value = `${draft.staffName || 'スタッフ'} のスケジュール変更を却下しました`
      expandedId.value = null
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isProcessing.value = false
  }
</script>

<template>
  <div
    v-motion
    :initial="{ opacity: 0, y: 20 }"
    :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
  >
    <h2 class="text-lg font-light tracking-wider text-foreground mb-6">シフト承認</h2>

    <div v-if="isLoading" class="text-sm text-muted-foreground">読み込み中...</div>

    <div v-if="error" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4">
      <p class="text-sm text-red-400">{{ error }}</p>
    </div>

    <div
      v-if="successMessage"
      class="bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3 mb-4"
    >
      <p class="text-sm text-green-400">{{ successMessage }}</p>
    </div>

    <div v-if="!isLoading && pendingSchedules.length === 0" class="text-sm text-muted-foreground">
      承認待ちのシフト変更はありません
    </div>

    <!-- ドラフト一覧 -->
    <div class="space-y-4">
      <div
        v-for="draft in pendingSchedules"
        :key="draft.id"
        class="border border-white/10 rounded-lg overflow-hidden"
      >
        <!-- ヘッダー -->
        <button
          @click="toggleExpand(draft)"
          class="w-full flex items-center justify-between p-6 text-left hover:bg-white/5 transition-colors"
        >
          <div class="flex-1">
            <div class="flex items-center gap-3">
              <h3 class="text-sm font-medium text-foreground">
                {{ draft.staffName || 'スタッフ' }}
              </h3>
              <span class="text-xs text-muted-foreground">出勤枠: {{ draft.items.length }}件</span>
            </div>
            <p class="text-xs text-muted-foreground mt-1">
              申請日:
              {{ draft.submittedAt ? new Date(draft.submittedAt).toLocaleString('ja-JP') : '-' }}
            </p>
          </div>
          <div class="flex items-center gap-3">
            <span
              class="text-xs text-blue-400 bg-blue-400/10 border border-blue-400/20 px-2 py-0.5 rounded"
            >
              承認待ち
            </span>
            <component :is="expandedId === draft.id ? ChevronUp : ChevronDown" :size="16" />
          </div>
        </button>

        <!-- 展開エリア -->
        <div v-if="expandedId === draft.id" class="border-t border-white/10 p-6 space-y-4">
          <!-- スケジュール編集 -->
          <div class="space-y-3">
            <div
              v-for="(item, index) in editItems"
              :key="index"
              class="flex items-center gap-3 bg-white/5 border border-white/10 rounded-lg px-4 py-3"
            >
              <select
                v-model.number="item.dayOfWeek"
                class="bg-zinc-900 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30"
              >
                <option v-for="(day, i) in dayLabels" :key="i" :value="i">{{ day }}</option>
              </select>

              <input
                v-model="item.startTime"
                type="time"
                class="bg-zinc-900 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
              />

              <span class="text-muted-foreground text-sm">〜</span>

              <input
                v-model="item.endTime"
                type="time"
                class="bg-zinc-900 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
              />

              <button
                type="button"
                @click="removeItem(index)"
                class="p-1 text-muted-foreground hover:text-red-400 transition-colors"
              >
                <Trash2 :size="16" />
              </button>
            </div>

            <button
              type="button"
              @click="addItem"
              class="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              <Plus :size="14" />
              出勤枠を追加
            </button>
          </div>

          <!-- 修正バッジ -->
          <div
            v-if="hasContentChanged()"
            class="text-xs text-yellow-400 bg-yellow-400/10 border border-yellow-400/20 rounded-lg px-3 py-2"
          >
            内容が修正されています。承認時に修正内容が反映されます。
          </div>

          <!-- コメント -->
          <div class="space-y-2">
            <label class="text-xs text-muted-foreground tracking-wider uppercase">
              コメント（却下時は必須）
            </label>
            <textarea
              v-model="reviewComment"
              rows="2"
              class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors resize-none"
              placeholder="コメントを入力"
            />
          </div>

          <!-- アクションボタン -->
          <div class="flex gap-3 pt-2">
            <button
              :disabled="isProcessing"
              @click="handleApprove(draft)"
              class="flex items-center gap-2 bg-green-600 text-white rounded-lg px-4 py-2.5 text-sm font-medium hover:bg-green-500 transition-colors disabled:opacity-50"
            >
              <Check :size="14" />
              {{ hasContentChanged() ? '修正して承認' : '承認' }}
            </button>
            <button
              :disabled="isProcessing"
              @click="handleReject(draft)"
              class="flex items-center gap-2 bg-red-600 text-white rounded-lg px-4 py-2.5 text-sm font-medium hover:bg-red-500 transition-colors disabled:opacity-50"
            >
              <X :size="14" />
              却下
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
