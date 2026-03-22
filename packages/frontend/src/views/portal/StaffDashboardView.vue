<script setup lang="ts">
  import { onMounted } from 'vue'
  import StaffPortalLayout from '@/components/layout/StaffPortalLayout.vue'
  import { useStaffPortalApi } from '@/composables/useStaffPortalApi'

  const { profileDraft, scheduleDraft, isLoading, fetchMyProfileDraft, fetchMyScheduleDraft } =
    useStaffPortalApi()

  onMounted(async () => {
    await Promise.all([fetchMyProfileDraft(), fetchMyScheduleDraft()])
  })

  /** ステータス表示のラベルマッピング */
  function statusLabel(status: string): string {
    const map: Record<string, string> = {
      draft: '下書き',
      pending: '承認待ち',
      approved: '承認済み',
      rejected: '却下',
      '': '未作成',
    }
    return map[status] ?? status
  }

  /** ステータスに応じた色クラス */
  function statusColor(status: string): string {
    const map: Record<string, string> = {
      draft: 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20',
      pending: 'text-blue-400 bg-blue-400/10 border-blue-400/20',
      approved: 'text-green-400 bg-green-400/10 border-green-400/20',
      rejected: 'text-red-400 bg-red-400/10 border-red-400/20',
      '': 'text-muted-foreground bg-white/5 border-white/10',
    }
    return map[status] ?? 'text-muted-foreground bg-white/5 border-white/10'
  }
</script>

<template>
  <StaffPortalLayout>
    <div
      v-motion
      :initial="{ opacity: 0, y: 20 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
    >
      <h2 class="text-lg font-light tracking-wider text-foreground mb-6">ダッシュボード</h2>

      <div v-if="isLoading" class="text-sm text-muted-foreground">読み込み中...</div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- プロフィール下書きカード -->
        <RouterLink
          :to="{ name: 'portal-profile' }"
          class="block border border-white/10 rounded-lg p-6 hover:border-white/20 transition-colors"
        >
          <h3 class="text-sm font-medium tracking-wider text-foreground mb-3">プロフィール</h3>
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <span class="text-xs text-muted-foreground">ステータス:</span>
              <span
                :class="[
                  'text-xs px-2 py-0.5 rounded border',
                  statusColor(profileDraft?.status ?? ''),
                ]"
              >
                {{ statusLabel(profileDraft?.status ?? '') }}
              </span>
            </div>
            <p v-if="profileDraft?.name" class="text-sm text-muted-foreground">
              {{ profileDraft.name }} / {{ profileDraft.role }}
            </p>
            <p
              v-if="profileDraft?.adminComment"
              class="text-xs text-red-400 mt-2 border-l-2 border-red-400/30 pl-2"
            >
              管理者コメント: {{ profileDraft.adminComment }}
            </p>
          </div>
        </RouterLink>

        <!-- スケジュール下書きカード -->
        <RouterLink
          :to="{ name: 'portal-schedule' }"
          class="block border border-white/10 rounded-lg p-6 hover:border-white/20 transition-colors"
        >
          <h3 class="text-sm font-medium tracking-wider text-foreground mb-3">スケジュール</h3>
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <span class="text-xs text-muted-foreground">ステータス:</span>
              <span
                :class="[
                  'text-xs px-2 py-0.5 rounded border',
                  statusColor(scheduleDraft?.status ?? ''),
                ]"
              >
                {{ statusLabel(scheduleDraft?.status ?? '') }}
              </span>
            </div>
            <p class="text-sm text-muted-foreground">
              登録枠: {{ scheduleDraft?.items?.length ?? 0 }}件
            </p>
            <p
              v-if="scheduleDraft?.adminComment"
              class="text-xs text-red-400 mt-2 border-l-2 border-red-400/30 pl-2"
            >
              管理者コメント: {{ scheduleDraft.adminComment }}
            </p>
          </div>
        </RouterLink>
      </div>
    </div>
  </StaffPortalLayout>
</template>
