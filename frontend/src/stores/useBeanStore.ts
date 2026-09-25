import { defineStore } from 'pinia'
import { ref } from 'vue'
import { listBeans, listFavoriteBeans, favoriteBean, unfavoriteBean, type BeanQuery } from '@/api/bean'
import type { CoffeeBean } from '@/constants/bean'

export type BeanViewMode = 'all' | 'favorites'

export const useBeanStore = defineStore('bean', () => {
  const beans = ref<CoffeeBean[]>([])
  const total = ref(0)
  const mode = ref<BeanViewMode>('all')
  // bean ids with an in-flight favorite request, to disable repeated clicks
  const pendingIds = ref<Set<number>>(new Set())

  async function load(params: BeanQuery & { mode?: BeanViewMode } = {}) {
    const { mode: viewMode = 'all', ...query } = params
    mode.value = viewMode
    const res = viewMode === 'favorites' ? await listFavoriteBeans(query) : await listBeans(query)
    beans.value = res.list
    total.value = res.total
  }

  function isPending(id: number) {
    return pendingIds.value.has(id)
  }

  function patchFavorite(id: number, isFavorite: boolean, favoriteCount: number) {
    const target = beans.value.find((b) => b.id === id)
    if (target) {
      target.is_favorite = isFavorite
      target.favorite_count = favoriteCount
    }
    if (!isFavorite && mode.value === 'favorites') {
      beans.value = beans.value.filter((b) => b.id !== id)
      total.value = Math.max(0, total.value - 1)
    }
  }

  async function toggleFavorite(id: number): Promise<boolean> {
    const target = beans.value.find((b) => b.id === id)
    const next = !target?.is_favorite
    pendingIds.value.add(id)
    try {
      if (next) {
        const res = await favoriteBean(id)
        patchFavorite(id, true, res.favorite_count)
      } else {
        const res = await unfavoriteBean(id)
        patchFavorite(id, false, res.favorite_count)
      }
      return next
    } finally {
      pendingIds.value.delete(id)
    }
  }

  return { beans, total, mode, pendingIds, load, isPending, toggleFavorite, patchFavorite }
})
