<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useAdminReviewApi } from '@/composables/useAdminReviewApi'
  import type { ProfileDraft } from '@/types/staffPortal'
  import { Check, X, ChevronDown, ChevronUp } from 'lucide-vue-next'

  const {
    pendingProfiles,
    isLoading,
    error,
    fetchPendingProfiles,
    updateProfileDraftContent,
    reviewProfileDraft,
  } = useAdminReviewApi()

  /** 展開中のドラフトID */
  const expandedId = ref<string | null>(null)

  /** レビューコメント */
  const reviewComment = ref('')

  /** 編集用フォーム（展開時にコピー） */
  const editForm = ref({
    name: '',
    role: '',
    bio: '',
    imageUrl: '',
    imageCropPosition: '',
  })

  /** 成功メッセージ */
  const successMessage = ref<string | null>(null)

  /** 処理中フラグ */
  const isProcessing = ref(false)

  onMounted(async () => {
    await fetchPendingProfiles()
  })

  /** ドラフトを展開/閉じる */
  function toggleExpand(draft: ProfileDraft) {
    if (expandedId.value === draft.id) {
      expandedId.value = null
      reviewComment.value = ''
    } else {
      expandedId.value = draft.id!
      reviewComment.value = ''
      // 編集フォームにコピー
      editForm.value = {
        name: draft.name,
        role: draft.role,
        bio: draft.bio,
        imageUrl: draft.imageUrl,
        imageCropPosition: draft.imageCropPosition,
      }
    }
  }

  /** 承認（修正があれば先に保存） */
  async function handleApprove(draft: ProfileDraft) {
    if (!draft.id) return
    if (
      !confirm(
        `${draft.staffName || draft.name} のプロフィール変更を承認しますか？ライブデータに反映されます。`
      )
    )
      return

    isProcessing.value = true

    // 内容が変更されていれば先に保存
    if (hasContentChanged(draft)) {
      const updated = await updateProfileDraftContent(draft.id, editForm.value)
      if (!updated) {
        isProcessing.value = false
        return
      }
    }

    const result = await reviewProfileDraft(draft.id, 'approved', reviewComment.value)
    if (result) {
      successMessage.value = `${draft.staffName || draft.name} のプロフィール変更を承認しました`
      expandedId.value = null
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isProcessing.value = false
  }

  /** 却下 */
  async function handleReject(draft: ProfileDraft) {
    if (!draft.id) return
    if (!reviewComment.value.trim()) {
      alert('却下理由を入力してください')
      return
    }

    isProcessing.value = true
    const result = await reviewProfileDraft(draft.id, 'rejected', reviewComment.value)
    if (result) {
      successMessage.value = `${draft.staffName || draft.name} のプロフィール変更を却下しました`
      expandedId.value = null
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isProcessing.value = false
  }

  /** 内容が変更されたか判定 */
  function hasContentChanged(draft: ProfileDraft): boolean {
    return (
      editForm.value.name !== draft.name ||
      editForm.value.role !== draft.role ||
      editForm.value.bio !== draft.bio ||
      editForm.value.imageUrl !== draft.imageUrl ||
      editForm.value.imageCropPosition !== draft.imageCropPosition
    )
  }
</script>

<template>
  <div
    v-motion
    :initial="{ opacity: 0, y: 20 }"
    :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
  >
    <h2 class="text-lg font-light tracking-wider text-foreground mb-6">プロフィール承認</h2>

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

    <div v-if="!isLoading && pendingProfiles.length === 0" class="text-sm text-muted-foreground">
      承認待ちのプロフィール変更はありません
    </div>

    <!-- ドラフト一覧 -->
    <div class="space-y-4">
      <div
        v-for="draft in pendingProfiles"
        :key="draft.id"
        class="border border-white/10 rounded-lg overflow-hidden"
      >
        <!-- ヘッダー（クリックで展開） -->
        <button
          @click="toggleExpand(draft)"
          class="w-full flex items-center justify-between p-6 text-left hover:bg-white/5 transition-colors"
        >
          <div class="flex-1">
            <div class="flex items-center gap-3">
              <span v-if="draft.staffName" class="text-xs text-muted-foreground">
                {{ draft.staffName }}
              </span>
              <span class="text-xs text-muted-foreground">→</span>
              <h3 class="text-sm font-medium text-foreground">{{ draft.name }}</h3>
              <span class="text-xs text-muted-foreground">/</span>
              <span class="text-xs text-muted-foreground">{{ draft.role }}</span>
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
          <!-- 編集可能フォーム -->
          <div class="grid gap-4">
            <div class="space-y-2">
              <label class="text-xs text-muted-foreground tracking-wider uppercase">名前</label>
              <input
                v-model="editForm.name"
                type="text"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
              />
            </div>

            <div class="space-y-2">
              <label class="text-xs text-muted-foreground tracking-wider uppercase">役職</label>
              <input
                v-model="editForm.role"
                type="text"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
              />
            </div>

            <div class="space-y-2">
              <label class="text-xs text-muted-foreground tracking-wider uppercase">自己紹介</label>
              <textarea
                v-model="editForm.bio"
                rows="3"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors resize-none"
              />
            </div>

            <div class="space-y-2">
              <label class="text-xs text-muted-foreground tracking-wider uppercase">画像URL</label>
              <input
                v-model="editForm.imageUrl"
                type="text"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
              />
            </div>
          </div>

          <!-- 修正があればバッジ表示 -->
          <div
            v-if="hasContentChanged(draft)"
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
              {{ hasContentChanged(draft) ? '修正して承認' : '承認' }}
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
