<script setup lang="ts">
  import { ref, onMounted, watch } from 'vue'
  import StaffPortalLayout from '@/components/layout/StaffPortalLayout.vue'
  import { useStaffPortalApi } from '@/composables/useStaffPortalApi'

  const {
    profileDraft,
    isLoading,
    error,
    saveError,
    fetchMyProfileDraft,
    saveProfileDraft,
    submitProfileDraft,
  } = useStaffPortalApi()

  // フォーム
  const name = ref('')
  const role = ref('')
  const bio = ref('')
  const imageUrl = ref('')
  const imageCropPosition = ref('center')
  const isSaving = ref(false)
  const successMessage = ref<string | null>(null)

  onMounted(async () => {
    await fetchMyProfileDraft()
    if (profileDraft.value) {
      name.value = profileDraft.value.name
      role.value = profileDraft.value.role
      bio.value = profileDraft.value.bio
      imageUrl.value = profileDraft.value.imageUrl
      imageCropPosition.value = profileDraft.value.imageCropPosition || 'center'
    }
  })

  // エラー時にスクロールトップ
  watch(saveError, (v) => {
    if (v) window.scrollTo({ top: 0, behavior: 'smooth' })
  })

  /** 下書き保存 */
  async function handleSave() {
    isSaving.value = true
    successMessage.value = null

    const result = await saveProfileDraft({
      name: name.value,
      role: role.value,
      bio: bio.value,
      imageUrl: imageUrl.value,
      imageCropPosition: imageCropPosition.value,
    })

    if (result) {
      successMessage.value = '下書きを保存しました'
      setTimeout(() => (successMessage.value = null), 3000)
    }
    isSaving.value = false
  }

  /** 承認申請 */
  async function handleSubmit() {
    if (!profileDraft.value?.id) return
    if (!confirm('この内容で承認申請しますか？申請後は管理者が承認するまで編集できません。')) return

    isSaving.value = true
    successMessage.value = null

    // 先に保存してから申請
    const saved = await saveProfileDraft({
      name: name.value,
      role: role.value,
      bio: bio.value,
      imageUrl: imageUrl.value,
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

  /** 編集可能かどうか */
  function isEditable(): boolean {
    return (
      !profileDraft.value?.status ||
      profileDraft.value.status === 'draft' ||
      profileDraft.value.status === 'rejected'
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
        <p class="text-sm text-blue-400">承認待ちです。管理者が確認するまでお待ちください。</p>
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
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
            placeholder="名前を入力"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">役職</label>
          <input
            v-model="role"
            type="text"
            :disabled="!isEditable()"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
            placeholder="役職を入力"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">自己紹介</label>
          <textarea
            v-model="bio"
            :disabled="!isEditable()"
            rows="4"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors resize-none disabled:opacity-50"
            placeholder="自己紹介を入力"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">画像URL</label>
          <input
            v-model="imageUrl"
            type="text"
            :disabled="!isEditable()"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
            placeholder="画像URLを入力"
          />
        </div>

        <div class="space-y-2">
          <label class="text-xs text-muted-foreground tracking-wider uppercase">
            画像切り抜き位置
          </label>
          <select
            v-model="imageCropPosition"
            :disabled="!isEditable()"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-3 text-sm text-foreground focus:outline-none focus:border-white/30 transition-colors disabled:opacity-50"
          >
            <option value="center">中央</option>
            <option value="top">上</option>
            <option value="bottom">下</option>
          </select>
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

          <button
            type="button"
            :disabled="isSaving || !profileDraft?.id"
            @click="handleSubmit"
            class="bg-white text-black rounded-lg px-6 py-3 text-sm font-medium tracking-wider hover:bg-white/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            承認申請
          </button>
        </div>
      </form>
    </div>
  </StaffPortalLayout>
</template>
