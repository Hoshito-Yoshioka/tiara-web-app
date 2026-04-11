<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import AdminLayout from '@/components/layout/AdminLayout.vue'
  import { useMenuApi } from '@/composables/useMenuApi'
  import { useAdminApi } from '@/composables/useAdminApi'
  import { Plus, Pencil, Trash2, X, Check } from 'lucide-vue-next'
  import type { MenuCategory, MenuItem } from '@/types/menu'

  const { menuList, fetchMenus, isLoading } = useMenuApi()
  const {
    createMenuCategory,
    updateMenuCategory,
    deleteMenuCategory,
    createMenuItem,
    updateMenuItem,
    deleteMenuItem,
    isLoading: isSaving,
    error: saveError,
  } = useAdminApi()

  // ============================================================
  // カテゴリ編集状態
  // ============================================================
  interface CategoryForm {
    name: string
    description: string
    sortOrder: number
  }

  const successMessage = ref<string | null>(null)

  /** 成功メッセージを表示してページ上部にスクロールする */
  function showSuccess(message: string) {
    successMessage.value = message
    window.scrollTo({ top: 0, behavior: 'smooth' })
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  }

  // ============================================================
  // 料金フォーマットヘルパー
  // ============================================================

  /** DB保存値（"¥1,200〜"）から ¥ を除去して編集用テキスト（"1,200〜"）にする */
  function stripYenPrefix(price: string): string {
    return price.replace(/^¥/, '')
  }

  /**
   * 入力中のテキストを整形する。
   * - 数字部分にカンマを自動挿入（既存カンマは除去してから再整形）
   * - 末尾の "〜" 等の非数字サフィックスは保持
   */
  function formatPriceText(raw: string): string {
    const match = raw.match(/^([0-9,]+)(.*)$/)
    if (!match) return raw
    const digits = match[1].replace(/,/g, '')
    const suffix = match[2]
    const formatted = digits.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
    return formatted + suffix
  }

  /** フォーム内の price フィールドをリアルタイム整形する */
  function onPriceInput(form: { price: string }) {
    form.price = formatPriceText(form.price)
  }

  /** 編集用テキスト（"1,200〜"）を DB 保存用（"¥1,200〜"）に変換する */
  function toStoragePrice(editPrice: string): string {
    return '¥' + editPrice
  }

  const editingCategoryId = ref<string | null>(null)
  const categoryForm = ref<CategoryForm>({ name: '', description: '', sortOrder: 0 })
  const showNewCategoryForm = ref(false)
  const newCategoryForm = ref<CategoryForm>({ name: '', description: '', sortOrder: 0 })

  function startEditCategory(cat: MenuCategory) {
    editingCategoryId.value = cat.id
    categoryForm.value = { name: cat.name, description: cat.description, sortOrder: cat.sortOrder }
  }

  function cancelEditCategory() {
    editingCategoryId.value = null
  }

  async function saveCategory(id: string) {
    const result = await updateMenuCategory(id, categoryForm.value)
    if (result) {
      editingCategoryId.value = null
      await fetchMenus()
      showSuccess('カテゴリを更新しました')
    }
  }

  async function handleDeleteCategory(id: string, name: string) {
    if (!confirm(`カテゴリ「${name}」と配下のメニューをすべて削除しますか？`)) return
    const ok = await deleteMenuCategory(id)
    if (ok) {
      await fetchMenus()
      showSuccess(`カテゴリ「${name}」を削除しました`)
    }
  }

  async function handleCreateCategory() {
    const result = await createMenuCategory(newCategoryForm.value)
    if (result) {
      showNewCategoryForm.value = false
      newCategoryForm.value = { name: '', description: '', sortOrder: 0 }
      await fetchMenus()
      showSuccess('カテゴリを作成しました')
    }
  }

  // ============================================================
  // アイテム編集状態
  // ============================================================
  interface ItemForm {
    name: string
    price: string
    description: string
    sortOrder: number
  }

  const editingItemId = ref<string | null>(null)
  const itemForm = ref<ItemForm>({ name: '', price: '', description: '', sortOrder: 0 })
  const showNewItemCategoryId = ref<string | null>(null)
  const newItemForm = ref<ItemForm>({ name: '', price: '', description: '', sortOrder: 0 })

  function startEditItem(item: MenuItem) {
    editingItemId.value = item.id
    itemForm.value = {
      name: item.name,
      price: stripYenPrefix(item.price),
      description: item.description,
      sortOrder: item.sortOrder,
    }
  }

  function cancelEditItem() {
    editingItemId.value = null
  }

  async function saveItem(id: string) {
    const result = await updateMenuItem(id, {
      ...itemForm.value,
      price: toStoragePrice(itemForm.value.price),
    })
    if (result) {
      editingItemId.value = null
      await fetchMenus()
      showSuccess('メニューアイテムを更新しました')
    }
  }

  async function handleDeleteItem(id: string, name: string) {
    if (!confirm(`「${name}」を削除しますか？`)) return
    const ok = await deleteMenuItem(id)
    if (ok) {
      await fetchMenus()
      showSuccess(`「${name}」を削除しました`)
    }
  }

  function startAddItem(categoryId: string, currentCount: number) {
    showNewItemCategoryId.value = categoryId
    newItemForm.value = { name: '', price: '', description: '', sortOrder: currentCount + 1 }
  }

  async function handleCreateItem(categoryId: string) {
    const result = await createMenuItem({
      categoryId,
      ...newItemForm.value,
      price: toStoragePrice(newItemForm.value.price),
    })
    if (result) {
      showNewItemCategoryId.value = null
      newItemForm.value = { name: '', price: '', description: '', sortOrder: 0 }
      await fetchMenus()
      showSuccess('メニューアイテムを追加しました')
    }
  }

  /** カテゴリの表示順ドロップダウン選択肢を生成 */
  function categorySortOrderOptions(
    currentSortOrder: number,
    excludeId?: string
  ): { value: number; label: string }[] {
    const categories = menuList.value.map((e) => e.category).filter((c) => c.id !== excludeId)
    const usedOrders = new Map<number, string>()
    for (const c of categories) {
      usedOrders.set(c.sortOrder, c.name)
    }
    const allOrders = [...categories.map((c) => c.sortOrder), currentSortOrder]
    const maxOrder = Math.max(...allOrders, 0) + 1
    const options: { value: number; label: string }[] = []
    for (let i = 0; i <= maxOrder; i++) {
      if (i === currentSortOrder && excludeId) {
        options.push({ value: i, label: `${i}（現在）` })
      } else if (usedOrders.has(i)) {
        options.push({ value: i, label: `${i}（${usedOrders.get(i)}と入れ替える）` })
      } else {
        options.push({ value: i, label: `${i}` })
      }
    }
    return options
  }

  /** アイテムの表示順ドロップダウン選択肢を生成 */
  function itemSortOrderOptions(
    categoryId: string,
    currentSortOrder: number,
    excludeId?: string
  ): { value: number; label: string }[] {
    const entry = menuList.value.find((e) => e.category.id === categoryId)
    const items = (entry?.items ?? []).filter((it) => it.id !== excludeId)
    const usedOrders = new Map<number, string>()
    for (const it of items) {
      usedOrders.set(it.sortOrder, it.name)
    }
    const allOrders = [...items.map((it) => it.sortOrder), currentSortOrder]
    const maxOrder = Math.max(...allOrders, 0) + 1
    const options: { value: number; label: string }[] = []
    for (let i = 0; i <= maxOrder; i++) {
      if (i === currentSortOrder && excludeId) {
        options.push({ value: i, label: `${i}（現在）` })
      } else if (usedOrders.has(i)) {
        options.push({ value: i, label: `${i}（${usedOrders.get(i)}と入れ替える）` })
      } else {
        options.push({ value: i, label: `${i}` })
      }
    }
    return options
  }

  onMounted(() => fetchMenus())
</script>

<template>
  <AdminLayout>
    <div
      class="max-w-4xl"
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 400 } }"
    >
      <!-- ページヘッダー -->
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-xl font-light tracking-wider text-foreground">メニュー・料金管理</h1>
        <button
          @click="showNewCategoryForm = true"
          class="flex items-center gap-2 bg-white text-black rounded-lg px-4 py-2 text-sm font-medium tracking-wider uppercase hover:bg-white/90 transition-colors"
        >
          <Plus :size="16" />
          カテゴリ追加
        </button>
      </div>

      <!-- グローバルエラー -->
      <div
        v-if="saveError"
        class="mb-6 bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3"
      >
        <p class="text-sm text-red-400">{{ saveError }}</p>
      </div>

      <!-- 成功メッセージ -->
      <div
        v-if="successMessage"
        class="mb-6 bg-green-500/10 border border-green-500/20 rounded-lg px-4 py-3"
      >
        <p class="text-sm text-green-400">{{ successMessage }}</p>
      </div>

      <div v-if="isLoading" class="text-muted-foreground text-sm">読み込み中...</div>

      <div v-else class="space-y-8">
        <!-- 新規カテゴリフォーム -->
        <div
          v-if="showNewCategoryForm"
          class="border border-white/20 rounded-lg p-5 bg-white/5 space-y-3"
        >
          <p class="text-xs text-muted-foreground tracking-wider uppercase mb-3">新規カテゴリ</p>
          <input
            v-model="newCategoryForm.name"
            placeholder="カテゴリ名（例: Cocktails）"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 text-sm text-foreground focus:outline-none focus:border-white/30"
          />
          <input
            v-model="newCategoryForm.description"
            placeholder="説明（任意）"
            class="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 text-sm text-foreground focus:outline-none focus:border-white/30"
          />
          <select
            v-model.number="newCategoryForm.sortOrder"
            class="w-auto bg-white/5 border border-white/10 rounded-lg px-4 py-2 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
          >
            <option
              v-for="opt in categorySortOrderOptions(newCategoryForm.sortOrder)"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
          <div class="flex gap-2 pt-1">
            <button
              @click="handleCreateCategory"
              :disabled="isSaving || !newCategoryForm.name"
              class="flex items-center gap-1.5 bg-white text-black rounded-lg px-4 py-2 text-xs font-medium hover:bg-white/90 transition-colors disabled:opacity-50"
            >
              <Check :size="14" /> 作成
            </button>
            <button
              @click="showNewCategoryForm = false"
              class="flex items-center gap-1.5 text-muted-foreground hover:text-foreground text-xs transition-colors px-2"
            >
              <X :size="14" /> キャンセル
            </button>
          </div>
        </div>

        <!-- カテゴリ一覧 -->
        <div
          v-for="entry in menuList"
          :key="entry.category.id"
          class="border border-white/10 rounded-lg overflow-hidden"
        >
          <!-- カテゴリヘッダー -->
          <div class="flex items-center gap-3 px-5 py-4 bg-white/5">
            <!-- 編集モード -->
            <template v-if="editingCategoryId === entry.category.id">
              <input
                v-model="categoryForm.name"
                class="flex-1 bg-white/10 border border-white/20 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/40"
              />
              <input
                v-model="categoryForm.description"
                placeholder="説明"
                class="flex-1 bg-white/10 border border-white/20 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/40"
              />
              <select
                v-model.number="categoryForm.sortOrder"
                class="w-auto bg-white/10 border border-white/20 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/40 [color-scheme:dark]"
              >
                <option
                  v-for="opt in categorySortOrderOptions(categoryForm.sortOrder, entry.category.id)"
                  :key="opt.value"
                  :value="opt.value"
                >
                  {{ opt.label }}
                </option>
              </select>
              <button
                @click="saveCategory(entry.category.id)"
                :disabled="isSaving"
                class="p-1.5 text-green-400 hover:text-green-300 transition-colors"
              >
                <Check :size="16" />
              </button>
              <button
                @click="cancelEditCategory"
                class="p-1.5 text-muted-foreground hover:text-foreground transition-colors"
              >
                <X :size="16" />
              </button>
            </template>
            <!-- 表示モード -->
            <template v-else>
              <h2 class="flex-1 text-sm font-medium tracking-wider text-foreground uppercase">
                {{ entry.category.name }}
              </h2>
              <span v-if="entry.category.description" class="text-xs text-muted-foreground flex-1">
                {{ entry.category.description }}
              </span>
              <span class="text-xs text-muted-foreground">順: {{ entry.category.sortOrder }}</span>
              <button
                @click="startEditCategory(entry.category)"
                class="p-1.5 text-muted-foreground hover:text-foreground transition-colors"
              >
                <Pencil :size="14" />
              </button>
              <button
                @click="handleDeleteCategory(entry.category.id, entry.category.name)"
                class="p-1.5 text-muted-foreground hover:text-red-400 transition-colors"
              >
                <Trash2 :size="14" />
              </button>
            </template>
          </div>

          <!-- アイテム一覧 -->
          <div class="divide-y divide-white/5">
            <div v-for="item in entry.items" :key="item.id" class="px-5 py-3">
              <!-- アイテム編集モード -->
              <template v-if="editingItemId === item.id">
                <div class="grid grid-cols-2 gap-2 mb-2">
                  <input
                    v-model="itemForm.name"
                    placeholder="アイテム名"
                    class="bg-white/5 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30"
                  />
                  <div
                    class="flex items-center bg-white/5 border border-white/10 rounded focus-within:border-white/30"
                  >
                    <span class="pl-3 text-sm text-muted-foreground select-none">¥</span>
                    <input
                      v-model="itemForm.price"
                      @input="onPriceInput(itemForm)"
                      placeholder="800〜"
                      class="flex-1 bg-transparent py-1.5 px-1.5 text-sm text-foreground focus:outline-none"
                    />
                  </div>
                  <input
                    v-model="itemForm.description"
                    placeholder="説明（任意）"
                    class="bg-white/5 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30"
                  />
                  <select
                    v-model.number="itemForm.sortOrder"
                    class="bg-white/5 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
                  >
                    <option
                      v-for="opt in itemSortOrderOptions(entry.category.id, itemForm.sortOrder, item.id)"
                      :key="opt.value"
                      :value="opt.value"
                    >
                      {{ opt.label }}
                    </option>
                  </select>
                </div>
                <div class="flex gap-2">
                  <button
                    @click="saveItem(item.id)"
                    :disabled="isSaving"
                    class="flex items-center gap-1 text-xs text-green-400 hover:text-green-300 transition-colors"
                  >
                    <Check :size="13" /> 保存
                  </button>
                  <button
                    @click="cancelEditItem"
                    class="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
                  >
                    <X :size="13" /> キャンセル
                  </button>
                </div>
              </template>
              <!-- アイテム表示モード -->
              <template v-else>
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-4 flex-1 min-w-0">
                    <p class="text-sm text-foreground">{{ item.name }}</p>
                    <p v-if="item.description" class="text-xs text-muted-foreground">
                      {{ item.description }}
                    </p>
                  </div>
                  <div class="flex items-center gap-3">
                    <p class="text-sm text-primary font-light tracking-wider">{{ item.price }}</p>
                    <span class="text-xs text-muted-foreground">順: {{ item.sortOrder }}</span>
                    <button
                      @click="startEditItem(item)"
                      class="p-1 text-muted-foreground hover:text-foreground transition-colors"
                    >
                      <Pencil :size="13" />
                    </button>
                    <button
                      @click="handleDeleteItem(item.id, item.name)"
                      class="p-1 text-muted-foreground hover:text-red-400 transition-colors"
                    >
                      <Trash2 :size="13" />
                    </button>
                  </div>
                </div>
              </template>
            </div>

            <!-- 新規アイテムフォーム -->
            <div v-if="showNewItemCategoryId === entry.category.id" class="px-5 py-3 bg-white/5">
              <div class="grid grid-cols-2 gap-2 mb-2">
                <input
                  v-model="newItemForm.name"
                  placeholder="アイテム名 *"
                  class="bg-white/5 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30"
                />
                <div
                  class="flex items-center bg-white/5 border border-white/10 rounded focus-within:border-white/30"
                >
                  <span class="pl-3 text-sm text-muted-foreground select-none">¥</span>
                  <input
                    v-model="newItemForm.price"
                    @input="onPriceInput(newItemForm)"
                    placeholder="800〜 *"
                    class="flex-1 bg-transparent py-1.5 px-1.5 text-sm text-foreground focus:outline-none"
                  />
                </div>
                <input
                  v-model="newItemForm.description"
                  placeholder="説明（任意）"
                  class="bg-white/5 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30"
                />
                <select
                  v-model.number="newItemForm.sortOrder"
                  class="bg-white/5 border border-white/10 rounded px-3 py-1.5 text-sm text-foreground focus:outline-none focus:border-white/30 [color-scheme:dark]"
                >
                  <option
                    v-for="opt in itemSortOrderOptions(entry.category.id, newItemForm.sortOrder)"
                    :key="opt.value"
                    :value="opt.value"
                  >
                    {{ opt.label }}
                  </option>
                </select>
              </div>
              <div class="flex gap-2">
                <button
                  @click="handleCreateItem(entry.category.id)"
                  :disabled="isSaving || !newItemForm.name || !newItemForm.price"
                  class="flex items-center gap-1 text-xs bg-white text-black rounded px-3 py-1.5 hover:bg-white/90 transition-colors disabled:opacity-50"
                >
                  <Check :size="13" /> 追加
                </button>
                <button
                  @click="showNewItemCategoryId = null"
                  class="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
                >
                  <X :size="13" /> キャンセル
                </button>
              </div>
            </div>

            <!-- アイテム追加ボタン -->
            <div class="px-5 py-2" v-if="showNewItemCategoryId !== entry.category.id">
              <button
                @click="startAddItem(entry.category.id, entry.items.length)"
                class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
              >
                <Plus :size="13" /> アイテムを追加
              </button>
            </div>
          </div>
        </div>

        <div v-if="menuList.length === 0" class="text-center py-12 text-muted-foreground text-sm">
          カテゴリが登録されていません
        </div>
      </div>
    </div>
  </AdminLayout>
</template>
