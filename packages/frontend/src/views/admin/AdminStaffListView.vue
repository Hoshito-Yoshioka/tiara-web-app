<script setup lang="ts">
  import { onMounted } from 'vue'
  import { RouterLink } from 'vue-router'
  import AdminLayout from '@/components/layout/AdminLayout.vue'
  import { useStaffApi } from '@/composables/useStaffApi'
  import { useAdminApi } from '@/composables/useAdminApi'
  import { Plus, Pencil, Trash2 } from 'lucide-vue-next'

  const { staffList, fetchStaffs, isLoading } = useStaffApi()
  const { deleteStaff, isLoading: isDeleting } = useAdminApi()

  onMounted(() => {
    fetchStaffs()
  })

  async function handleDelete(id: string, name: string) {
    if (!confirm(`「${name}」を削除してもよろしいですか？`)) return

    const success = await deleteStaff(id)
    if (success) {
      await fetchStaffs()
    }
  }
</script>

<template>
  <AdminLayout>
    <div
      class="max-w-4xl"
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 400 } }"
    >
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-xl font-light tracking-wider text-foreground">スタッフ管理</h1>
        <RouterLink
          to="/admin/staffs/new"
          class="flex items-center gap-2 bg-white text-black rounded-lg px-4 py-2 text-sm font-medium tracking-wider uppercase hover:bg-white/90 transition-colors"
        >
          <Plus :size="16" />
          新規作成
        </RouterLink>
      </div>

      <div v-if="isLoading" class="text-muted-foreground text-sm">読み込み中...</div>

      <div v-else class="space-y-3">
        <div
          v-for="staff in staffList"
          :key="staff.id"
          class="flex items-center justify-between bg-white/5 border border-white/10 rounded-lg px-5 py-4 hover:border-white/20 transition-colors"
        >
          <div class="flex items-center gap-4">
            <img
              v-if="staff.imageUrl"
              :src="staff.imageUrl"
              :alt="staff.name"
              class="w-10 h-10 rounded-full object-cover"
            />
            <div
              v-else
              class="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center text-xs text-muted-foreground"
            >
              {{ staff.name.charAt(0) }}
            </div>
            <div>
              <p class="text-sm text-foreground font-medium">{{ staff.name }}</p>
              <p class="text-xs text-muted-foreground">{{ staff.role }}</p>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <RouterLink
              :to="`/admin/staffs/${staff.id}/edit`"
              class="p-2 text-muted-foreground hover:text-foreground transition-colors"
              title="編集"
            >
              <Pencil :size="16" />
            </RouterLink>
            <button
              @click="handleDelete(staff.id, staff.name)"
              :disabled="isDeleting"
              class="p-2 text-muted-foreground hover:text-red-400 transition-colors disabled:opacity-50"
              title="削除"
            >
              <Trash2 :size="16" />
            </button>
          </div>
        </div>

        <div v-if="staffList.length === 0" class="text-center py-12 text-muted-foreground text-sm">
          スタッフが登録されていません
        </div>
      </div>
    </div>
  </AdminLayout>
</template>
