<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  import AdminLayout from '@/components/layout/AdminLayout.vue'
  import { useAdminReviewApi } from '@/composables/useAdminReviewApi'
  import type { ProfileDraft, StaffImageForDraft } from '@/types/staffPortal'
  import { Check, X, ChevronDown, ChevronUp, UserCheck, ImageOff, Star } from 'lucide-vue-next'

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

  /** 画像読み込みエラー（ドラフトのメイン画像用） */
  const imageLoadError = ref(false)

  /** 各画像の読み込みエラー追跡（imageId → boolean） */
  const imageLoadErrors = ref<Record<string, boolean>>({})

  /** BFF ベース URL（画像プロキシ先） */
  const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001'

  /** 画像URLをフルパスに解決 */
  function resolveImageUrl(url: string): string {
    if (!url) return ''
    if (url.startsWith('http://') || url.startsWith('https://')) return url
    return `${API_BASE}${url.startsWith('/') ? '' : '/'}${url}`
  }

  /** 現在編集中の画像プレビューURL */
  const previewImageUrl = computed(() => resolveImageUrl(editForm.value.imageUrl))

  /** 現在展開中のドラフトに紐づくスタッフ画像一覧 */
  const currentDraftImages = computed<StaffImageForDraft[]>(() => {
    if (!expandedId.value) return []
    const draft = pendingProfiles.value.find((d) => d.id === expandedId.value)
    return draft?.images ?? []
  })

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
      imageLoadError.value = false
      imageLoadErrors.value = {}
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

      <!-- 空状態 -->
      <div
        v-else-if="pendingProfiles.length === 0"
        class="flex flex-col items-center justify-center py-24 border border-dashed border-white/10 rounded-lg bg-white/[0.02]"
      >
        <div class="w-16 h-16 rounded-full bg-white/5 flex items-center justify-center mb-5">
          <UserCheck :size="28" class="text-muted-foreground/40" />
        </div>
        <p class="text-sm font-medium text-muted-foreground mb-1">承認待ちの申請はありません</p>
        <p class="text-xs text-muted-foreground/50">
          スタッフがポータルからプロフィール変更を申請するとここに表示されます
        </p>
      </div>

      <!-- ドラフト一覧 -->
      <div v-else class="space-y-4">
        <div
          v-for="draft in pendingProfiles"
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
                  {{ draft.staffName || draft.name }}
                </h3>
                <span class="text-xs text-muted-foreground">→</span>
                <span class="text-sm text-foreground/80">{{ draft.name }}</span>
                <span class="text-xs text-muted-foreground bg-white/5 px-2 py-0.5 rounded">{{
                  draft.role
                }}</span>
              </div>
              <p class="text-xs text-muted-foreground/70 mt-1.5">
                申請:
                {{ draft.submittedAt ? new Date(draft.submittedAt).toLocaleString('ja-JP') : '-' }}
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
            <!-- 編集可能フォーム -->
            <div class="grid sm:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-xs text-muted-foreground/70 tracking-wider uppercase"
                  >名前</label
                >
                <input
                  v-model="editForm.name"
                  type="text"
                  class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2.5 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
                />
              </div>

              <div class="space-y-1.5">
                <label class="text-xs text-muted-foreground/70 tracking-wider uppercase"
                  >役職</label
                >
                <input
                  v-model="editForm.role"
                  type="text"
                  class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2.5 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
                />
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-xs text-muted-foreground/70 tracking-wider uppercase"
                >自己紹介</label
              >
              <textarea
                v-model="editForm.bio"
                rows="3"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2.5 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors resize-none"
              />
            </div>

            <div class="space-y-3">
              <label class="text-xs text-muted-foreground/70 tracking-wider uppercase"
                >プロフィール画像</label
              >

              <!-- ドラフトの申請画像プレビュー + URL 編集 -->
              <div class="flex items-start gap-4">
                <div
                  class="relative w-28 h-28 shrink-0 rounded-xl overflow-hidden border border-white/10 bg-white/5"
                >
                  <img
                    v-if="editForm.imageUrl && !imageLoadError"
                    :src="previewImageUrl"
                    alt="申請中のプロフィール画像"
                    class="w-full h-full object-cover"
                    :style="{
                      objectPosition: editForm.imageCropPosition || 'center',
                    }"
                    @error="imageLoadError = true"
                  />
                  <div
                    v-else
                    class="w-full h-full flex flex-col items-center justify-center gap-1.5 text-muted-foreground/40"
                  >
                    <ImageOff :size="24" />
                    <span class="text-[10px]">No Image</span>
                  </div>
                  <!-- 申請画像バッジ -->
                  <span
                    class="absolute bottom-1 left-1 text-[9px] bg-blue-500/80 text-white px-1.5 py-0.5 rounded"
                    >申請</span
                  >
                </div>

                <div class="flex-1 space-y-1.5">
                  <label class="text-xs text-muted-foreground/70 tracking-wider uppercase"
                    >画像URL</label
                  >
                  <input
                    v-model="editForm.imageUrl"
                    type="text"
                    class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2.5 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors"
                    @input="imageLoadError = false"
                  />
                  <p v-if="editForm.imageCropPosition" class="text-[11px] text-muted-foreground/50">
                    トリミング位置: {{ editForm.imageCropPosition }}
                  </p>
                </div>
              </div>

              <!-- スタッフの全画像ギャラリー -->
              <div v-if="currentDraftImages.length > 0" class="space-y-2 pt-1">
                <label class="text-xs text-muted-foreground/70 tracking-wider uppercase"
                  >登録済み画像一覧
                  <span class="text-muted-foreground/40"
                    >({{ currentDraftImages.length }}枚)</span
                  ></label
                >
                <div class="flex flex-wrap gap-2">
                  <div
                    v-for="img in currentDraftImages"
                    :key="img.id"
                    class="relative w-20 h-20 rounded-lg overflow-hidden border bg-white/5 transition-all"
                    :class="
                      img.isMain
                        ? 'border-yellow-500/50 ring-1 ring-yellow-500/30'
                        : 'border-white/10'
                    "
                  >
                    <img
                      v-if="!imageLoadErrors[img.id]"
                      :src="resolveImageUrl(img.imageUrl)"
                      :alt="`スタッフ画像 ${img.sortOrder}`"
                      class="w-full h-full object-cover"
                      :style="{ objectPosition: img.cropPosition || 'center' }"
                      @error="imageLoadErrors[img.id] = true"
                    />
                    <div
                      v-else
                      class="w-full h-full flex items-center justify-center text-muted-foreground/30"
                    >
                      <ImageOff :size="16" />
                    </div>
                    <!-- メイン画像バッジ -->
                    <span
                      v-if="img.isMain"
                      class="absolute top-0.5 right-0.5 text-yellow-400"
                      title="メイン画像"
                    >
                      <Star :size="12" class="fill-yellow-400" />
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 修正バッジ -->
            <div
              v-if="hasContentChanged(draft)"
              class="text-xs text-yellow-400 bg-yellow-400/10 border border-yellow-400/20 rounded-lg px-3 py-2"
            >
              内容が修正されています。承認時に修正内容が反映されます。
            </div>

            <!-- コメント -->
            <div class="space-y-1.5">
              <label class="text-xs text-muted-foreground/70 tracking-wider uppercase">
                コメント（却下時は必須）
              </label>
              <textarea
                v-model="reviewComment"
                rows="2"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground/40 focus:outline-none focus:border-white/30 transition-colors resize-none"
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
                {{ hasContentChanged(draft) ? '修正して承認' : '承認' }}
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
    </div>
  </AdminLayout>
</template>
