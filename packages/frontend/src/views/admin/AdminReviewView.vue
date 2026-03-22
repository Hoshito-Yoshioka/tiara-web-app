<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useAdminReviewApi } from '@/composables/useAdminReviewApi'
  import type { ProfileDraft, ScheduleDraft } from '@/types/staffPortal'

  const {
    pendingProfiles,
    pendingSchedules,
    isLoading,
    error,
    fetchPendingProfiles,
    fetchPendingSchedules,
    reviewProfileDraft,
    reviewScheduleDraft,
  } = useAdminReviewApi()

  const activeTab = ref<'profiles' | 'schedules'>('profiles')
  const reviewComment = ref('')
  const reviewingId = ref<string | null>(null)

  const dayLabels = ['日', '月', '火', '水', '木', '金', '土']

  onMounted(async () => {
    await Promise.all([fetchPendingProfiles(), fetchPendingSchedules()])
  })

  /** プロフィール下書きを承認 */
  async function approveProfile(draft: ProfileDraft) {
    if (!draft.id) return
    if (!confirm(`${draft.name} のプロフィール変更を承認しますか？ライブデータに反映されます。`))
      return
    await reviewProfileDraft(draft.id, 'approved', reviewComment.value)
    reviewComment.value = ''
    reviewingId.value = null
  }

  /** プロフィール下書きを却下 */
  async function rejectProfile(draft: ProfileDraft) {
    if (!draft.id) return
    if (!reviewComment.value.trim()) {
      alert('却下理由を入力してください')
      return
    }
    await reviewProfileDraft(draft.id, 'rejected', reviewComment.value)
    reviewComment.value = ''
    reviewingId.value = null
  }

  /** スケジュール下書きを承認 */
  async function approveSchedule(draft: ScheduleDraft) {
    if (!draft.id) return
    if (!confirm('スケジュール変更を承認しますか？ライブデータに反映されます。')) return
    await reviewScheduleDraft(draft.id, 'approved', reviewComment.value)
    reviewComment.value = ''
    reviewingId.value = null
  }

  /** スケジュール下書きを却下 */
  async function rejectSchedule(draft: ScheduleDraft) {
    if (!draft.id) return
    if (!reviewComment.value.trim()) {
      alert('却下理由を入力してください')
      return
    }
    await reviewScheduleDraft(draft.id, 'rejected', reviewComment.value)
    reviewComment.value = ''
    reviewingId.value = null
  }

  /** レビューパネルの表示切り替え */
  function toggleReview(id: string) {
    if (reviewingId.value === id) {
      reviewingId.value = null
      reviewComment.value = ''
    } else {
      reviewingId.value = id
      reviewComment.value = ''
    }
  }
</script>

<template>
  <div
    v-motion
    :initial="{ opacity: 0, y: 20 }"
    :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
  >
    <h2 class="text-lg font-light tracking-wider text-foreground mb-6">承認管理</h2>

    <!-- タブ -->
    <div class="flex gap-4 mb-6 border-b border-white/10">
      <button
        @click="activeTab = 'profiles'"
        :class="[
          'pb-3 text-sm tracking-wider transition-colors border-b-2',
          activeTab === 'profiles'
            ? 'text-foreground border-white'
            : 'text-muted-foreground border-transparent hover:text-foreground',
        ]"
      >
        プロフィール
        <span
          v-if="pendingProfiles.length > 0"
          class="ml-1 text-xs bg-blue-500/20 text-blue-400 px-1.5 py-0.5 rounded"
        >
          {{ pendingProfiles.length }}
        </span>
      </button>
      <button
        @click="activeTab = 'schedules'"
        :class="[
          'pb-3 text-sm tracking-wider transition-colors border-b-2',
          activeTab === 'schedules'
            ? 'text-foreground border-white'
            : 'text-muted-foreground border-transparent hover:text-foreground',
        ]"
      >
        スケジュール
        <span
          v-if="pendingSchedules.length > 0"
          class="ml-1 text-xs bg-blue-500/20 text-blue-400 px-1.5 py-0.5 rounded"
        >
          {{ pendingSchedules.length }}
        </span>
      </button>
    </div>

    <div v-if="isLoading" class="text-sm text-muted-foreground">読み込み中...</div>

    <div v-if="error" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4">
      <p class="text-sm text-red-400">{{ error }}</p>
    </div>

    <!-- プロフィール下書き一覧 -->
    <div v-if="activeTab === 'profiles'" class="space-y-4">
      <div v-if="pendingProfiles.length === 0" class="text-sm text-muted-foreground">
        承認待ちのプロフィール変更はありません
      </div>

      <div
        v-for="draft in pendingProfiles"
        :key="draft.id"
        class="border border-white/10 rounded-lg p-6 space-y-3"
      >
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-sm font-medium text-foreground">{{ draft.name }}</h3>
            <p class="text-xs text-muted-foreground">{{ draft.role }}</p>
          </div>
          <span
            class="text-xs text-blue-400 bg-blue-400/10 border border-blue-400/20 px-2 py-0.5 rounded"
          >
            承認待ち
          </span>
        </div>

        <p v-if="draft.bio" class="text-sm text-muted-foreground whitespace-pre-wrap">
          {{ draft.bio }}
        </p>

        <div v-if="draft.imageUrl" class="text-xs text-muted-foreground">
          画像: {{ draft.imageUrl }}
        </div>

        <p class="text-xs text-muted-foreground">
          申請日:
          {{ draft.submittedAt ? new Date(draft.submittedAt).toLocaleString('ja-JP') : '-' }}
        </p>

        <!-- レビューパネル -->
        <div class="pt-2">
          <button
            @click="toggleReview(draft.id!)"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
          >
            {{ reviewingId === draft.id ? '閉じる' : 'レビューする' }}
          </button>

          <div v-if="reviewingId === draft.id" class="mt-3 space-y-3">
            <textarea
              v-model="reviewComment"
              rows="2"
              class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors resize-none"
              placeholder="コメント（却下時は必須）"
            />
            <div class="flex gap-3">
              <button
                @click="approveProfile(draft)"
                class="bg-green-600 text-white rounded-lg px-4 py-2 text-sm font-medium hover:bg-green-500 transition-colors"
              >
                承認
              </button>
              <button
                @click="rejectProfile(draft)"
                class="bg-red-600 text-white rounded-lg px-4 py-2 text-sm font-medium hover:bg-red-500 transition-colors"
              >
                却下
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- スケジュール下書き一覧 -->
    <div v-if="activeTab === 'schedules'" class="space-y-4">
      <div v-if="pendingSchedules.length === 0" class="text-sm text-muted-foreground">
        承認待ちのスケジュール変更はありません
      </div>

      <div
        v-for="draft in pendingSchedules"
        :key="draft.id"
        class="border border-white/10 rounded-lg p-6 space-y-3"
      >
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-sm font-medium text-foreground">スタッフID: {{ draft.staffId }}</h3>
            <p class="text-xs text-muted-foreground">出勤枠: {{ draft.items.length }}件</p>
          </div>
          <span
            class="text-xs text-blue-400 bg-blue-400/10 border border-blue-400/20 px-2 py-0.5 rounded"
          >
            承認待ち
          </span>
        </div>

        <!-- スケジュール内容 -->
        <div class="space-y-1">
          <div v-for="item in draft.items" :key="item.id" class="text-sm text-muted-foreground">
            {{ dayLabels[item.dayOfWeek] }}曜日: {{ item.startTime }} 〜 {{ item.endTime }}
          </div>
        </div>

        <p class="text-xs text-muted-foreground">
          申請日:
          {{ draft.submittedAt ? new Date(draft.submittedAt).toLocaleString('ja-JP') : '-' }}
        </p>

        <!-- レビューパネル -->
        <div class="pt-2">
          <button
            @click="toggleReview(draft.id!)"
            class="text-xs text-muted-foreground hover:text-foreground transition-colors"
          >
            {{ reviewingId === draft.id ? '閉じる' : 'レビューする' }}
          </button>

          <div v-if="reviewingId === draft.id" class="mt-3 space-y-3">
            <textarea
              v-model="reviewComment"
              rows="2"
              class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors resize-none"
              placeholder="コメント（却下時は必須）"
            />
            <div class="flex gap-3">
              <button
                @click="approveSchedule(draft)"
                class="bg-green-600 text-white rounded-lg px-4 py-2 text-sm font-medium hover:bg-green-500 transition-colors"
              >
                承認
              </button>
              <button
                @click="rejectSchedule(draft)"
                class="bg-red-600 text-white rounded-lg px-4 py-2 text-sm font-medium hover:bg-red-500 transition-colors"
              >
                却下
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
