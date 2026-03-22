<script setup lang="ts">
  import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import AdminLayout from '@/components/layout/AdminLayout.vue'
  import { useStaffApi } from '@/composables/useStaffApi'
  import { useShopApi } from '@/composables/useShopApi'
  import { useAdminApi } from '@/composables/useAdminApi'
  import { useAdminAccountApi } from '@/composables/useAdminAccountApi'
  import { Save, ArrowLeft, Plus, X, Move, Upload, Star, Trash2, KeyRound } from 'lucide-vue-next'
  import type { StaffImage } from '@/types/staff'

  const route = useRoute()
  const router = useRouter()

  const { staffDetail, fetchStaffById } = useStaffApi()
  const { shops, fetchShops } = useShopApi()
  const {
    createStaff,
    updateStaff,
    uploadStaffImage,
    deleteStaffImage,
    setMainImage,
    isLoading: isSaving,
    error: saveError,
  } = useAdminApi()

  const isEditMode = computed(() => !!route.params.id)
  const pageTitle = computed(() => (isEditMode.value ? 'スタッフ編集' : 'スタッフ新規作成'))

  const DAYS = ['日', '月', '火', '水', '木', '金', '土']

  const form = ref({
    name: '',
    role: '',
    bio: '',
    imageUrl: '',
    imageCropPosition: '50 50',
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

  // --- Staff Account Management ---
  const {
    account: staffAccount,
    error: accountError,
    fetchAccountByStaffId,
    createAccount,
    updateAccount,
    deleteAccount,
  } = useAdminAccountApi()

  const accountForm = ref({
    username: '',
    password: '',
  })
  const accountSuccess = ref<string | null>(null)
  const isAccountSaving = ref(false)

  /** スタッフ画像一覧（リアクティブ） */
  const images = ref<StaffImage[]>([])
  /** 画像アップロード中フラグ */
  const isUploading = ref(false)
  /** ファイル入力の ref */
  const fileInputRef = ref<HTMLInputElement | null>(null)

  /** メイン画像の URL（プレビュー用） */
  const mainImageUrl = computed(() => {
    const main = images.value.find((img) => img.isMain)
    return main?.imageUrl ?? images.value[0]?.imageUrl ?? ''
  })

  /** エラー表示時にページトップへ自動スクロール */
  watch(saveError, (newVal) => {
    if (newVal) {
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  })

  // --- 画像クロップ位置のドラッグ操作 ---
  /** object-position を "X% Y%" 形式の CSS 文字列に変換 */
  const cropPositionStyle = computed(() => {
    const [x, y] = form.value.imageCropPosition.split(' ').map(Number)
    return `${x}% ${y}%`
  })

  /** ドラッグ中フラグ */
  const isDragging = ref(false)
  /** ドラッグ開始位置 */
  let dragStartX = 0
  let dragStartY = 0
  /** ドラッグ開始時の cropPosition */
  let dragStartPosX = 50
  let dragStartPosY = 50

  function parseCropPosition(): { x: number; y: number } {
    const parts = form.value.imageCropPosition.split(' ').map(Number)
    return { x: parts[0] ?? 50, y: parts[1] ?? 50 }
  }

  function clamp(val: number, min: number, max: number): number {
    return Math.min(Math.max(val, min), max)
  }

  function onDragStart(e: MouseEvent | TouchEvent) {
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
    // 感度: プレビュー枠幅(192px)分のドラッグで 0→100% 移動
    const sensitivity = 100 / 192
    const newX = clamp(Math.round(dragStartPosX - deltaX * sensitivity), 0, 100)
    const newY = clamp(Math.round(dragStartPosY - deltaY * sensitivity), 0, 100)
    form.value.imageCropPosition = `${newX} ${newY}`
  }

  function onDragEnd() {
    isDragging.value = false
    document.removeEventListener('mousemove', onDragMove)
    document.removeEventListener('mouseup', onDragEnd)
    document.removeEventListener('touchmove', onDragMove)
    document.removeEventListener('touchend', onDragEnd)
  }

  /** クロップ位置をセンターにリセット */
  function resetCropPosition() {
    form.value.imageCropPosition = '50 50'
  }

  // --- 画像管理 ---
  /** ファイル選択ダイアログを開く */
  function triggerFileInput() {
    fileInputRef.value?.click()
  }

  /** 画像ファイルをアップロード */
  async function handleImageUpload(event: Event) {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    if (!file || !isEditMode.value) return

    isUploading.value = true
    const staffId = route.params.id as string
    const result = await uploadStaffImage(staffId, file)
    if (result) {
      images.value.push(result)
      // アップロード後に imageUrl を同期（最初の画像は自動でメインに）
      if (images.value.length === 1) {
        form.value.imageUrl = result.imageUrl
      }
    }
    isUploading.value = false
    // input をリセット
    target.value = ''
  }

  /** 画像を削除 */
  async function handleDeleteImage(imageId: string) {
    if (!confirm('この画像を削除しますか？')) return
    const staffId = route.params.id as string
    const success = await deleteStaffImage(staffId, imageId)
    if (success) {
      images.value = images.value.filter((img) => img.id !== imageId)
    }
  }

  /** メイン画像を設定 */
  async function handleSetMain(imageId: string) {
    const staffId = route.params.id as string
    const success = await setMainImage(staffId, imageId)
    if (success) {
      images.value = images.value.map((img) => ({
        ...img,
        isMain: img.id === imageId,
      }))
    }
  }

  onBeforeUnmount(() => {
    // クリーンアップ
    document.removeEventListener('mousemove', onDragMove)
    document.removeEventListener('mouseup', onDragEnd)
    document.removeEventListener('touchmove', onDragMove)
    document.removeEventListener('touchend', onDragEnd)
  })

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
        const { staff, schedules: staffSchedules, images: staffImages } = staffDetail.value
        form.value = {
          name: staff.name,
          role: staff.role,
          bio: staff.bio,
          imageUrl: staff.imageUrl,
          imageCropPosition: staff.imageCropPosition || '50 50',
          sortOrder: staff.sortOrder,
          shopId: staff.shopId,
        }
        schedules.value = staffSchedules.map((s) => ({
          dayOfWeek: s.dayOfWeek,
          startTime: formatTimeForInput(s.startTime),
          endTime: formatTimeForInput(s.endTime),
        }))
        images.value = staffImages ?? []
      }
      // スタッフアカウント情報を取得
      await fetchAccountByStaffId(route.params.id as string)
      if (staffAccount.value) {
        accountForm.value.username = staffAccount.value.username
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

  // --- Staff Account Handlers ---

  /** スタッフアカウントを保存（作成 or 更新） */
  async function handleSaveAccount() {
    isAccountSaving.value = true
    accountSuccess.value = null

    if (staffAccount.value) {
      // 更新
      const result = await updateAccount(
        staffAccount.value.id,
        accountForm.value.username,
        accountForm.value.password
      )
      if (result) {
        accountSuccess.value = 'アカウント情報を更新しました'
        accountForm.value.password = '' // パスワード欄をクリア
        setTimeout(() => (accountSuccess.value = null), 3000)
      }
    } else {
      // 作成
      const staffId = route.params.id as string
      const result = await createAccount(
        staffId,
        accountForm.value.username,
        accountForm.value.password
      )
      if (result) {
        accountSuccess.value = 'ポータルアカウントを作成しました'
        accountForm.value.password = '' // パスワード欄をクリア
        setTimeout(() => (accountSuccess.value = null), 3000)
      }
    }
    isAccountSaving.value = false
  }

  /** スタッフアカウントを削除 */
  async function handleDeleteAccount() {
    if (!staffAccount.value) return
    if (
      !confirm(
        'このスタッフのポータルアカウントを削除しますか？\nスタッフはマイページにログインできなくなります。'
      )
    )
      return

    const success = await deleteAccount(staffAccount.value.id)
    if (success) {
      accountForm.value = { username: '', password: '' }
      accountSuccess.value = 'アカウントを削除しました'
      setTimeout(() => (accountSuccess.value = null), 3000)
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

        <!-- 画像管理セクション（編集モードのみ） -->
        <div v-if="isEditMode" class="space-y-4">
          <div class="flex items-center justify-between">
            <label class="text-xs text-muted-foreground tracking-wider uppercase">画像管理</label>
            <button
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

          <!-- アップロード済み画像一覧 -->
          <div v-if="images.length > 0" class="grid grid-cols-3 sm:grid-cols-4 gap-3">
            <div
              v-for="img in images"
              :key="img.id"
              class="relative group rounded-lg overflow-hidden border-2 transition-colors"
              :class="img.isMain ? 'border-primary' : 'border-white/10 hover:border-white/20'"
            >
              <div class="aspect-square overflow-hidden bg-secondary">
                <img :src="img.imageUrl" :alt="'staff image'" class="w-full h-full object-cover" />
              </div>
              <!-- メインバッジ -->
              <div
                v-if="img.isMain"
                class="absolute top-1.5 left-1.5 bg-primary text-primary-foreground text-[9px] tracking-wider uppercase px-2 py-0.5 rounded-full"
              >
                メイン
              </div>
              <!-- 操作ボタン（ホバーで表示） -->
              <div
                class="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2"
              >
                <button
                  v-if="!img.isMain"
                  type="button"
                  @click="handleSetMain(img.id)"
                  class="p-2 bg-white/20 rounded-full hover:bg-white/30 transition-colors"
                  title="メイン画像に設定"
                >
                  <Star :size="14" class="text-white" />
                </button>
                <button
                  type="button"
                  @click="handleDeleteImage(img.id)"
                  class="p-2 bg-red-500/40 rounded-full hover:bg-red-500/60 transition-colors"
                  title="削除"
                >
                  <Trash2 :size="14" class="text-white" />
                </button>
              </div>
            </div>
          </div>
          <p v-else class="text-xs text-muted-foreground">
            画像がアップロードされていません。「画像を追加」から画像をアップロードしてください。
          </p>
        </div>

        <!-- 画像表示位置プレビュー（ドラッグで切り抜き位置を調整） -->
        <div v-if="mainImageUrl || form.imageUrl" class="space-y-4">
          <div class="flex items-center gap-3">
            <p class="text-[11px] tracking-wider uppercase text-muted-foreground">
              表示プレビュー — ドラッグで表示位置を調整できます
            </p>
            <button
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
                :class="
                  isDragging ? 'border-primary cursor-grabbing' : 'border-primary/40 cursor-grab'
                "
                @mousedown="onDragStart"
                @touchstart="onDragStart"
              >
                <img
                  :src="mainImageUrl || form.imageUrl"
                  :alt="form.name || 'プレビュー'"
                  class="w-full h-full object-cover pointer-events-none"
                  :style="{ objectPosition: cropPositionStyle }"
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <div
                  class="absolute inset-0 border-2 border-dashed border-primary/30 pointer-events-none"
                />
                <div
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
                  :src="mainImageUrl || form.imageUrl"
                  :alt="form.name || 'プレビュー'"
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
            位置: {{ form.imageCropPosition.split(' ')[0] }}% /
            {{ form.imageCropPosition.split(' ')[1] }}%
          </p>
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

        <!-- ポータルアカウント管理（編集モードのみ） -->
        <div v-if="isEditMode" class="space-y-4 border-t border-white/10 pt-6">
          <div class="flex items-center gap-2">
            <KeyRound :size="16" class="text-muted-foreground" />
            <label class="text-xs text-muted-foreground tracking-wider uppercase">
              ポータルアカウント
            </label>
            <span
              v-if="staffAccount"
              class="text-[10px] bg-green-500/20 text-green-400 rounded-full px-2 py-0.5"
            >
              登録済み
            </span>
            <span v-else class="text-[10px] bg-zinc-500/20 text-zinc-400 rounded-full px-2 py-0.5">
              未登録
            </span>
          </div>

          <!-- アカウント成功メッセージ -->
          <div
            v-if="accountSuccess"
            class="bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3 text-sm text-green-400"
          >
            {{ accountSuccess }}
          </div>

          <!-- アカウントエラーメッセージ -->
          <div
            v-if="accountError"
            class="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 text-sm text-red-400"
          >
            {{ accountError }}
          </div>

          <div class="space-y-3">
            <div class="space-y-2">
              <label class="text-xs text-muted-foreground/80">ユーザー名</label>
              <input
                v-model="accountForm.username"
                type="text"
                placeholder="ログインID"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder-white/30 focus:outline-none focus:border-white/30 transition-colors"
              />
            </div>
            <div class="space-y-2">
              <label class="text-xs text-muted-foreground/80">パスワード</label>
              <input
                v-model="accountForm.password"
                type="password"
                :placeholder="staffAccount ? '変更する場合のみ入力' : 'パスワードを設定'"
                class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder-white/30 focus:outline-none focus:border-white/30 transition-colors"
              />
            </div>
          </div>

          <div class="flex items-center gap-3">
            <button
              type="button"
              :disabled="
                isAccountSaving || !accountForm.username || (!staffAccount && !accountForm.password)
              "
              @click="handleSaveAccount"
              class="flex items-center gap-2 bg-white/10 text-foreground rounded-lg px-4 py-2.5 text-xs font-medium tracking-wider uppercase hover:bg-white/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Save :size="14" />
              {{
                isAccountSaving ? '保存中...' : staffAccount ? 'アカウント更新' : 'アカウント作成'
              }}
            </button>
            <button
              v-if="staffAccount"
              type="button"
              @click="handleDeleteAccount"
              class="flex items-center gap-2 bg-red-500/10 text-red-400 rounded-lg px-4 py-2.5 text-xs font-medium tracking-wider uppercase hover:bg-red-500/20 transition-colors"
            >
              <Trash2 :size="14" />
              アカウント削除
            </button>
          </div>
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
