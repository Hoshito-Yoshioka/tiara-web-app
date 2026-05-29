<script setup lang="ts">
  import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
  import StaffPortalLayout from '@/components/layout/StaffPortalLayout.vue'
  import { useStaffPortalApi } from '@/composables/useStaffPortalApi'
  import { Upload, Star, Trash2, Move } from 'lucide-vue-next'
  import type { StaffImage } from '@/types/staff'

  const {
    profileDraft,
    isLoading,
    error,
    saveError,
    fetchMyProfileDraft,
    saveProfileDraft,
    submitProfileDraft,
    fetchMyImages,
    uploadMyImage,
    deleteMyImage,
    setMyMainImage,
    updateMyImageCropPosition,
  } = useStaffPortalApi()

  // フォーム
  const name = ref('')
  const role = ref('')
  const bio = ref('')
  const imageUrl = ref('')
  const externalScheduleUrl = ref('')
  const imageCropPosition = ref('50 50')
  const isSaving = ref(false)
  const successMessage = ref<string | null>(null)

  /** 最後に保存した値のスナップショット（未保存変更の検知用） */
  const savedSnapshot = ref({ name: '', role: '', bio: '' })

  /** フォームに未保存の変更があるか */
  const hasUnsavedChanges = computed(() => {
    return (
      name.value !== savedSnapshot.value.name ||
      role.value !== savedSnapshot.value.role ||
      bio.value !== savedSnapshot.value.bio
    )
  })

  /** スタッフ画像一覧（リアクティブ） */
  const images = ref<StaffImage[]>([])
  /** 画像一覧の初回読み込み完了フラグ */
  const imagesLoaded = ref(false)
  /** 画像アップロード中フラグ */
  const isUploading = ref(false)
  /** ファイル入力の ref */
  const fileInputRef = ref<HTMLInputElement | null>(null)
  /** 選択中の画像ID（クロップ編集対象） */
  const selectedImageId = ref<string | null>(null)

  /** メイン画像の URL（プレビュー用） */
  const mainImageUrl = computed(() => {
    const main = images.value.find((img) => img.isMain)
    return main?.imageUrl ?? images.value[0]?.imageUrl ?? ''
  })

  /** 選択中の画像オブジェクト */
  const selectedImage = computed(() => {
    if (!selectedImageId.value) return null
    return images.value.find((img) => img.id === selectedImageId.value) ?? null
  })

  /** 選択中画像のURL */
  const selectedImageUrl = computed(() => {
    return selectedImage.value?.imageUrl ?? mainImageUrl.value
  })

  /** 選択中画像のクロップ位置 */
  const selectedCropPosition = computed(() => {
    return selectedImage.value?.cropPosition ?? '50 50'
  })

  onMounted(async () => {
    await fetchMyProfileDraft()
    if (profileDraft.value) {
      name.value = profileDraft.value.name
      role.value = profileDraft.value.role
      bio.value = profileDraft.value.bio
      imageUrl.value = profileDraft.value.imageUrl
      externalScheduleUrl.value = profileDraft.value.externalScheduleUrl || ''
      imageCropPosition.value = profileDraft.value.imageCropPosition || '50 50'
      savedSnapshot.value = { name: name.value, role: role.value, bio: bio.value }
    }
    // 画像一覧を取得
    images.value = await fetchMyImages()
    imagesLoaded.value = true
    // メイン画像を初期選択
    const mainImg = images.value.find((img) => img.isMain) ?? images.value[0]
    if (mainImg) {
      selectedImageId.value = mainImg.id
    }
  })

  // エラー時にスクロールトップ
  watch(saveError, (v) => {
    if (v) window.scrollTo({ top: 0, behavior: 'smooth' })
  })

  // --- 画像クロップ位置のドラッグ操作（選択中画像に対して） ---
  /** object-position を "X% Y%" 形式の CSS 文字列に変換 */
  const cropPositionStyle = computed(() => {
    const [x, y] = selectedCropPosition.value.split(' ').map(Number)
    return `${x}% ${y}%`
  })

  /** ドラッグ中フラグ */
  const isDragging = ref(false)
  let dragStartX = 0
  let dragStartY = 0
  let dragStartPosX = 50
  let dragStartPosY = 50

  function parseCropPosition(): { x: number; y: number } {
    const parts = selectedCropPosition.value.split(' ').map(Number)
    return { x: parts[0] ?? 50, y: parts[1] ?? 50 }
  }

  function clamp(val: number, min: number, max: number): number {
    return Math.min(Math.max(val, min), max)
  }

  /** 選択中の画像の cropPosition を更新するヘルパー */
  function updateSelectedCropPosition(newPos: string) {
    if (!selectedImageId.value) return
    const idx = images.value.findIndex((img) => img.id === selectedImageId.value)
    if (idx >= 0) {
      images.value[idx] = { ...images.value[idx], cropPosition: newPos }
    }
    // メイン画像の場合は imageCropPosition も同期
    const img = images.value[idx]
    if (img?.isMain) {
      imageCropPosition.value = newPos
    }
  }

  function onDragStart(e: MouseEvent | TouchEvent) {
    if (!isEditable() || !selectedImageId.value) return
    e.preventDefault()
    isDragging.value = true
    const point = 'touches' in e ? e.touches[0] : e
    dragStartX = point.clientX
    dragStartY = point.clientY
    const pos = parseCropPosition()
    dragStartPosX = pos.x
    dragStartPosY = pos.y
    document.addEventListener('mousemove', onDragMove)
    document.addEventListener('mouseup', onDragEnd)
    document.addEventListener('touchmove', onDragMove, { passive: false })
    document.addEventListener('touchend', onDragEnd)
  }

  function onDragMove(e: MouseEvent | TouchEvent) {
    if (!isDragging.value) return
    e.preventDefault()
    const point = 'touches' in e ? e.touches[0] : e
    const deltaX = point.clientX - dragStartX
    const deltaY = point.clientY - dragStartY
    const sensitivity = 100 / 192
    const newX = clamp(Math.round(dragStartPosX - deltaX * sensitivity), 0, 100)
    const newY = clamp(Math.round(dragStartPosY - deltaY * sensitivity), 0, 100)
    updateSelectedCropPosition(`${newX} ${newY}`)
  }

  async function onDragEnd() {
    isDragging.value = false
    document.removeEventListener('mousemove', onDragMove)
    document.removeEventListener('mouseup', onDragEnd)
    document.removeEventListener('touchmove', onDragMove)
    document.removeEventListener('touchend', onDragEnd)
    // ドラッグ終了時に自動保存
    if (selectedImageId.value) {
      await updateMyImageCropPosition(selectedImageId.value, selectedCropPosition.value)
    }
  }

  async function resetCropPosition() {
    updateSelectedCropPosition('50 50')
    if (selectedImageId.value) {
      await updateMyImageCropPosition(selectedImageId.value, '50 50')
    }
  }

  onBeforeUnmount(() => {
    document.removeEventListener('mousemove', onDragMove)
    document.removeEventListener('mouseup', onDragEnd)
    document.removeEventListener('touchmove', onDragMove)
    document.removeEventListener('touchend', onDragEnd)
  })

  // --- 画像管理 ---
  function triggerFileInput() {
    fileInputRef.value?.click()
  }

  async function handleImageUpload(event: Event) {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    if (!file) return

    isUploading.value = true
    const result = await uploadMyImage(file)
    if (result) {
      images.value.push(result)
      // 最初の画像はメインとして imageUrl にセット
      if (images.value.length === 1) {
        imageUrl.value = result.imageUrl
      }
    }
    isUploading.value = false
    target.value = ''
  }

  async function handleDeleteImage(imageId: string) {
    if (!confirm('この画像を削除しますか？')) return
    const success = await deleteMyImage(imageId)
    if (success) {
      images.value = images.value.filter((img) => img.id !== imageId)
      if (selectedImageId.value === imageId) {
        const mainImg = images.value.find((img) => img.isMain) ?? images.value[0]
        selectedImageId.value = mainImg?.id ?? null
      }
    }
  }

  async function handleSetMain(imageId: string) {
    const success = await setMyMainImage(imageId)
    if (success) {
      images.value = images.value.map((img) => ({
        ...img,
        isMain: img.id === imageId,
      }))
    }
  }

  /** 下書き保存（メイン画像URLを自動同期） */
  async function handleSave() {
    isSaving.value = true
    successMessage.value = null

    // メイン画像の URL を imageUrl に自動同期
    const currentImageUrl = mainImageUrl.value || imageUrl.value

    const result = await saveProfileDraft({
      name: name.value,
      role: role.value,
      bio: bio.value,
      imageUrl: currentImageUrl,
      externalScheduleUrl: externalScheduleUrl.value,
      imageCropPosition: imageCropPosition.value,
    })

    if (result) {
      imageUrl.value = currentImageUrl
      savedSnapshot.value = { name: name.value, role: role.value, bio: bio.value }
      successMessage.value = '下書きを保存しました'
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isSaving.value = false
  }

  /** 承認申請 */
  async function handleSubmit() {
    if (!profileDraft.value?.id) return
    if (hasUnsavedChanges.value) {
      alert('未保存の変更があります。先に「下書き保存」を行ってから承認申請してください。')
      return
    }
    if (!confirm('この内容で承認申請しますか？')) return

    isSaving.value = true
    successMessage.value = null

    // メイン画像の URL を imageUrl に自動同期
    const currentImageUrl = mainImageUrl.value || imageUrl.value

    const saved = await saveProfileDraft({
      name: name.value,
      role: role.value,
      bio: bio.value,
      imageUrl: currentImageUrl,
      externalScheduleUrl: externalScheduleUrl.value,
      imageCropPosition: imageCropPosition.value,
    })

    if (saved?.id) {
      const result = await submitProfileDraft(saved.id)
      if (result) {
        successMessage.value = '承認申請しました。管理者の確認をお待ちください。'
      }
    }
    isSaving.value = false
  }

  /** 編集可能かどうか（approved 以外は全て編集可能） */
  function isEditable(): boolean {
    return (
      !profileDraft.value?.status ||
      profileDraft.value.status === 'draft' ||
      profileDraft.value.status === 'rejected' ||
      profileDraft.value.status === 'pending'
    )
  }

  /** 提出可能かどうか（draft / rejected のみ） */
  function isSubmittable(): boolean {
    return (
      !!profileDraft.value?.id &&
      (profileDraft.value.status === 'draft' || profileDraft.value.status === 'rejected')
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
      <h2 class="text-lg font-light tracking-wider text-foreground mb-6">プロフィール編集</h2>

      <!-- ローディング -->
      <div v-if="isLoading" class="text-sm text-muted-foreground">読み込み中...</div>

      <!-- エラー -->
      <div v-if="error" class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4">
        <p class="text-sm text-red-400">{{ error }}</p>
      </div>

      <!-- 保存エラー -->
      <div
        v-if="saveError"
        class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-red-400">{{ saveError }}</p>
      </div>

      <!-- 成功メッセージ -->
      <div
        v-if="successMessage"
        class="bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-green-400">{{ successMessage }}</p>
      </div>

      <!-- ステータス表示 -->
      <div
        v-if="profileDraft?.status === 'pending'"
        class="bg-blue-500/10 border border-blue-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-blue-400">
          承認待ちです。編集して保存すると申請が取り下げられ、再提出が必要になります。
        </p>
      </div>

      <div
        v-if="profileDraft?.status === 'rejected' && profileDraft.adminComment"
        class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 mb-4"
      >
        <p class="text-sm text-red-400">
          <strong>却下理由:</strong> {{ profileDraft.adminComment }}
        </p>
        <p class="text-xs text-red-400/70 mt-1">内容を修正して再申請してください。</p>
      </div>

      <!-- フォーム -->
      <form v-if="!isLoading" @submit.prevent="handleSave" class="space-y-6">
        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">名前</label>
          <input
            v-model="name"
            type="text"
            :disabled="!isEditable()"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/70 focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
            placeholder="名前を入力"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">役職</label>
          <input
            v-model="role"
            type="text"
            :disabled="!isEditable()"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/70 focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
            placeholder="役職を入力"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">自己紹介</label>
          <textarea
            v-model="bio"
            :disabled="!isEditable()"
            rows="4"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/70 focus:outline-none focus:border-white/30 transition-colors resize-none disabled:opacity-50"
            placeholder="自己紹介を入力"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">
            外部スケジュールURL（ポケパラ）
          </label>
          <input
            v-model="externalScheduleUrl"
            type="url"
            :disabled="!isEditable()"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/70 focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
            placeholder="https://www.pokepara.jp/..."
          />
          <p class="text-[11px] text-muted-foreground/70">
            スタッフ個別の出勤情報ページURLを入力してください。承認後に公開ページへ反映されます。
          </p>
        </div>

        <!-- 画像管理セクション -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <label class="text-xs text-muted-foreground tracking-wider uppercase">画像管理</label>
            <button
              v-if="isEditable()"
              type="button"
              @click="triggerFileInput"
              :disabled="isUploading"
              class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
            >
              <Upload :size="14" />
              {{ isUploading ? 'アップロード中...' : '画像を追加' }}
            </button>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*"
              class="hidden"
              @change="handleImageUpload"
            />
          </div>

          <!-- アップロード済み画像一覧（クリックで選択→クロップ編集） -->
          <div v-if="images.length > 0" class="grid grid-cols-3 sm:grid-cols-4 gap-3">
            <div
              v-for="img in images"
              :key="img.id"
              class="relative group rounded-lg overflow-hidden border-2 transition-colors cursor-pointer"
              :class="
                selectedImageId === img.id
                  ? 'border-primary ring-2 ring-primary/30'
                  : img.isMain
                    ? 'border-primary/50'
                    : 'border-white/10 hover:border-white/20'
              "
              @click="selectedImageId = img.id"
            >
              <div class="aspect-square overflow-hidden bg-secondary">
                <img
                  :src="img.imageUrl"
                  alt="staff image"
                  class="w-full h-full object-cover"
                  :style="{
                    objectPosition: (img.cropPosition ?? '50 50')
                      .split(' ')
                      .map((v: string) => v + '%')
                      .join(' '),
                  }"
                />
              </div>
              <!-- メインバッジ -->
              <div
                v-if="img.isMain"
                class="absolute top-1.5 left-1.5 bg-primary text-primary-foreground text-[9px] tracking-wider uppercase px-2 py-0.5 rounded-full"
              >
                メイン
              </div>
              <!-- 選択中インジケーター -->
              <div
                v-if="selectedImageId === img.id"
                class="absolute top-1.5 right-1.5 bg-primary text-primary-foreground text-[9px] tracking-wider uppercase px-2 py-0.5 rounded-full"
              >
                編集中
              </div>
              <!-- 操作ボタン（ホバーで表示） -->
              <div
                v-if="isEditable()"
                class="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2"
              >
                <button
                  v-if="!img.isMain"
                  type="button"
                  @click.stop="handleSetMain(img.id)"
                  class="p-2 bg-white/20 rounded-full hover:bg-white/30 transition-colors"
                  title="メイン画像に設定"
                >
                  <Star :size="14" class="text-white" />
                </button>
                <button
                  type="button"
                  @click.stop="handleDeleteImage(img.id)"
                  class="p-2 bg-red-500/40 rounded-full hover:bg-red-500/60 transition-colors"
                  title="削除"
                >
                  <Trash2 :size="14" class="text-white" />
                </button>
              </div>
            </div>
          </div>
          <p v-else-if="imagesLoaded" class="text-xs text-muted-foreground">
            画像がアップロードされていません。「画像を追加」から画像をアップロードしてください。
          </p>
        </div>

        <!-- 画像表示位置プレビュー（ドラッグで切り抜き位置を調整 — 選択中の画像に適用） -->
        <div v-if="selectedImageUrl || imageUrl" class="space-y-4">
          <div class="flex items-center gap-3">
            <p class="text-[11px] tracking-wider uppercase text-muted-foreground">
              表示プレビュー — 画像を選択してドラッグで表示位置を調整
            </p>
            <button
              v-if="isEditable()"
              type="button"
              @click="resetCropPosition"
              class="text-[10px] tracking-wider text-muted-foreground/60 hover:text-foreground underline transition-colors"
            >
              中央にリセット
            </button>
          </div>
          <div class="flex flex-wrap gap-6">
            <!-- カード表示プレビュー（ドラッグ対応） -->
            <div class="space-y-2">
              <p class="text-[10px] tracking-wider uppercase text-muted-foreground/70">
                スタッフ一覧（カード）
              </p>
              <div
                class="relative w-48 h-52 overflow-hidden rounded-lg border-2 bg-secondary transition-colors select-none"
                :class="[
                  isDragging
                    ? 'border-primary cursor-grabbing'
                    : selectedImageId
                      ? 'border-primary/40 cursor-grab'
                      : 'border-white/20',
                  !isEditable() && 'pointer-events-none opacity-70',
                ]"
                @mousedown="onDragStart"
                @touchstart="onDragStart"
              >
                <img
                  :src="selectedImageUrl || imageUrl"
                  :alt="name || 'プレビュー'"
                  class="w-full h-full object-cover pointer-events-none"
                  :style="{ objectPosition: cropPositionStyle }"
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <div
                  class="absolute inset-0 border-2 border-dashed border-primary/30 pointer-events-none"
                />
                <div
                  v-if="selectedImageId"
                  class="absolute top-2 right-2 bg-black/60 rounded-full p-1.5 pointer-events-none"
                >
                  <Move class="w-3 h-3 text-white/80" />
                </div>
              </div>
            </div>
            <!-- 詳細ページプレビュー（連動表示） -->
            <div class="space-y-2">
              <p class="text-[10px] tracking-wider uppercase text-muted-foreground/70">
                スタッフ詳細ページ
              </p>
              <div
                class="relative w-36 overflow-hidden rounded-lg border-2 border-primary/40 bg-secondary"
                style="aspect-ratio: 3/4"
              >
                <img
                  :src="selectedImageUrl || imageUrl"
                  :alt="name || 'プレビュー'"
                  class="w-full h-full object-cover"
                  :style="{ objectPosition: cropPositionStyle }"
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <div
                  class="absolute inset-0 border-2 border-dashed border-primary/30 pointer-events-none"
                />
              </div>
            </div>
          </div>
          <!-- 座標表示 -->
          <p class="text-[10px] text-muted-foreground/50 tracking-wider">
            位置: {{ selectedCropPosition.split(' ')[0] }}% /
            {{ selectedCropPosition.split(' ')[1] }}%
          </p>
        </div>

        <!-- ボタン群 -->
        <div v-if="isEditable()" class="flex gap-4">
          <button
            type="submit"
            :disabled="isSaving"
            class="bg-white/10 text-foreground border border-white/20 rounded-lg px-6 py-3 text-sm font-medium tracking-wider hover:bg-white/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isSaving ? '保存中...' : '下書き保存' }}
          </button>

          <div class="flex flex-col gap-1">
            <button
              type="button"
              :disabled="isSaving || !isSubmittable()"
              @click="handleSubmit"
              class="bg-white text-black rounded-lg px-6 py-3 text-sm font-medium tracking-wider hover:bg-white/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ profileDraft?.status === 'pending' ? '提出済み' : '承認申請' }}
            </button>
            <p
              v-if="
                !isSubmittable() &&
                profileDraft?.status !== 'pending' &&
                profileDraft?.status !== 'approved'
              "
              class="text-[11px] text-muted-foreground/70"
            >
              ⮶ 先に「下書き保存」を行ってください
            </p>
            <p
              v-else-if="hasUnsavedChanges && isSubmittable()"
              class="text-[11px] text-amber-400/80"
            >
              ⚠ 未保存の変更があります
            </p>
          </div>
        </div>
      </form>
    </div>
  </StaffPortalLayout>
</template>
