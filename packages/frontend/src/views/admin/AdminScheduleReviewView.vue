<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import AdminLayout from '@/components/layout/AdminLayout.vue'
  import { useAdminReviewApi } from '@/composables/useAdminReviewApi'
  import type { ScheduleDraft, ScheduleDraftItem } from '@/types/staffPortal'
  import {
    Check,
    X,
    ChevronDown,
    ChevronUp,
    Plus,
    Trash2,
    CalendarCheck,
    Send,
  } from 'lucide-vue-next'

  const {
    pendingSchedules,
    approvedSchedules,
    isLoading,
    error,
    fetchPendingSchedules,
    fetchApprovedSchedules,
    updateScheduleDraftContent,
    reviewScheduleDraft,
    publishScheduleDraft,
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
    await Promise.all([fetchPendingSchedules(), fetchApprovedSchedules()])
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
        `${draft.staffName || 'スタッフ'} のスケジュール変更を承認しますか？（店舗ページへの反映は別途行います）`
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
      // 承認済みリストを再取得
      await fetchApprovedSchedules()
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

  /** 店舗ページに反映 */
  async function handlePublish(draft: ScheduleDraft) {
    if (!draft.id) return
    if (!confirm(`${draft.staffName || 'スタッフ'} のシフトを店舗ページに反映しますか？`)) return

    isProcessing.value = true
    const result = await publishScheduleDraft(draft.id)
    if (result) {
      successMessage.value = `${draft.staffName || 'スタッフ'} のシフトを店舗ページに反映しました`
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isProcessing.value = false
  }
</script>

<template>
  <AdminLayout>
    <div
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 400 } }"
    >
      <!-- ステータスメッセージ -->
      <div v-if="error" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-6">
        <p class="text-sm text-red-400">{{ error }}</p>
      </div>
      <div
        v-if="successMessage"
        class="bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3 mb-6"
      >
        <p class="text-sm text-green-400">{{ successMessage }}</p>
      </div>

      <!-- ローディング -->
      <div v-if="isLoading" class="flex items-center justify-center py-20">
        <p class="text-sm text-muted-foreground">読み込み中...</p>
      </div>

      <template v-else>
        <!-- ============ 承認待ちセクション ============ -->
        <section>
          <h2 class="text-xs font-medium tracking-widest uppercase text-muted-foreground mb-4">
            承認待ち
          </h2>

          <div
            v-if="pendingSchedules.length === 0"
            class="flex flex-col items-center justify-center py-16 border border-dashed border-white/10 rounded-lg bg-white/[0.02]"
          >
            <div class="w-14 h-14 rounded-full bg-white/5 flex items-center justify-center mb-4">
              <CalendarCheck :size="24" class="text-muted-foreground/40" />
            </div>
            <p class="text-sm font-medium text-muted-foreground mb-1">承認待ちの申請はありません</p>
            <p class="text-xs text-muted-foreground/50">
              スタッフがポータルからシフト変更を申請するとここに表示されます
            </p>
          </div>

          <!-- ドラフト一覧 -->
          <div class="space-y-4">
            <div
              v-for="draft in pendingSchedules"
              :key="draft.id"
              class="border border-white/10 rounded-lg bg-card overflow-hidden"
            >
              <!-- カードヘッダー -->
              <button
                @click="toggleExpand(draft)"
                class="w-full flex items-center justify-between px-6 py-5 text-left hover:bg-white/[0.03] transition-colors"
              >
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 flex-wrap">
                    <h3 class="text-sm font-medium text-foreground">
                      {{ draft.staffName || 'スタッフ' }}
                    </h3>
                    <span class="text-xs text-muted-foreground bg-white/5 px-2 py-0.5 rounded"
                      >出勤枠: {{ draft.items.length }}件</span
                    >
                  </div>
                  <p class="text-xs text-muted-foreground/70 mt-1.5">
                    申請:
                    {{
                      draft.submittedAt ? new Date(draft.submittedAt).toLocaleString('ja-JP') : '-'
                    }}
                  </p>
                </div>
                <div class="flex items-center gap-3 ml-4 shrink-0">
                  <span
                    class="text-xs text-blue-400 bg-blue-500/10 border border-blue-500/20 px-2.5 py-1 rounded-full font-medium"
                  >
                    承認待ち
                  </span>
                  <component
                    :is="expandedId === draft.id ? ChevronUp : ChevronDown"
                    :size="16"
                    class="text-muted-foreground"
                  />
                </div>
              </button>

              <!-- 展開エリア -->
              <div
                v-if="expandedId === draft.id"
                class="border-t border-white/5 bg-white/[0.02] px-6 py-5 space-y-5"
              >
                <!-- スケジュール編集 -->
                <div class="space-y-3">
                  <div
                    v-for="(item, index) in editItems"
                    :key="index"
                    class="flex items-center gap-3 bg-white/[0.04] border border-white/10 rounded-lg px-4 py-3"
                  >
                    <select
                      v-model.number="item.dayOfWeek"
                      class="bg-zinc-900 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
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
                    class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 text-sm text-foreground placeholder:text-muted-foreground/70 focus:outline-none focus:border-white/30 transition-colors resize-none"
                    placeholder="コメントを入力"
                  />
                </div>

                <!-- アクションボタン -->
                <div class="flex gap-3 pt-1">
                  <button
                    :disabled="isProcessing"
                    @click="handleApprove(draft)"
                    class="flex items-center gap-2 bg-green-600 text-white rounded-lg px-5 py-2.5 text-sm font-medium hover:bg-green-500 transition-colors disabled:opacity-50"
                  >
                    <Check :size="14" />
                    {{ hasContentChanged() ? '修正して承認' : '承認' }}
                  </button>
                  <button
                    :disabled="isProcessing"
                    @click="handleReject(draft)"
                    class="flex items-center gap-2 bg-white/5 border border-white/10 text-foreground rounded-lg px-5 py-2.5 text-sm font-medium hover:bg-red-600 hover:border-red-600 hover:text-white transition-colors disabled:opacity-50"
                  >
                    <X :size="14" />
                    却下
                  </button>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- ============ 承認済み（未反映）セクション ============ -->
        <section class="mt-12">
          <h2
            class="text-xs font-medium tracking-widest uppercase text-muted-foreground mb-4 flex items-center gap-2"
          >
            <Send :size="14" class="text-primary" />
            店舗反映待ち
          </h2>

          <div
            v-if="approvedSchedules.length === 0"
            class="flex flex-col items-center justify-center py-12 border border-dashed border-white/10 rounded-lg bg-white/[0.02]"
          >
            <div class="w-12 h-12 rounded-full bg-white/5 flex items-center justify-center mb-3">
              <Send :size="18" class="text-muted-foreground/40" />
            </div>
            <p class="text-sm font-medium text-muted-foreground mb-1">反映待ちなし</p>
            <p class="text-xs text-muted-foreground/50">承認済みで未反映のシフトはありません</p>
          </div>

          <div class="space-y-4">
            <div
              v-for="draft in approvedSchedules"
              :key="draft.id"
              class="border border-primary/20 bg-primary/[0.04] rounded-lg overflow-hidden"
            >
              <div class="px-6 py-5">
                <div class="flex items-center justify-between mb-4">
                  <div class="min-w-0">
                    <div class="flex items-center gap-2 flex-wrap">
                      <h4 class="text-sm font-medium text-foreground">
                        {{ draft.staffName || 'スタッフ' }}
                      </h4>
                      <span class="text-xs text-muted-foreground bg-white/5 px-2 py-0.5 rounded"
                        >出勤枠: {{ draft.items.length }}件</span
                      >
                      <span
                        class="text-xs text-green-400 bg-green-500/10 border border-green-500/20 px-2.5 py-1 rounded-full font-medium"
                      >
                        承認済み
                      </span>
                    </div>
                    <p class="text-xs text-muted-foreground/70 mt-1.5">
                      申請:
                      {{
                        draft.submittedAt
                          ? new Date(draft.submittedAt).toLocaleString('ja-JP')
                          : '-'
                      }}
                    </p>
                  </div>
                  <button
                    :disabled="isProcessing"
                    @click="handlePublish(draft)"
                    class="flex items-center gap-2 bg-primary text-primary-foreground rounded-lg px-5 py-2.5 text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
                  >
                    <Send :size="14" />
                    店舗に反映
                  </button>
                </div>

                <!-- スケジュールプレビュー -->
                <div class="flex flex-wrap gap-2">
                  <div
                    v-for="(item, index) in draft.items"
                    :key="index"
                    class="bg-white/5 border border-white/10 rounded px-3 py-1.5 text-xs text-muted-foreground"
                  >
                    {{ dayLabels[item.dayOfWeek] }} {{ item.startTime }}〜{{ item.endTime }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </template>
    </div>
  </AdminLayout>
</template>
