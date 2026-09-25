import { defineStore } from 'pinia'
import { ref } from 'vue'
import { favoriteBean, listBeans, listFavoriteBeans, unfavoriteBean, type BeanListParams } from '@/api/bean'
import type { BeanItem } from '@/constants/bean'

export const useBeanStore = defineStore('bean', () => {
  const beans = ref<BeanItem[]>([])
  const total = ref(0)
  const favorites = ref<BeanItem[]>([])
  const favoriteTotal = ref(0)

  async function load(params: BeanListParams = {}) {
    const res = await listBeans(params)
    beans.value = res.list
    total.value = res.total
  }

  async function loadFavorites(params: BeanListParams = {}) {
    const res = await listFavoriteBeans(params)
    favorites.value = res.list
    favoriteTotal.value = res.total
  }

  // Clear per-account view data (call on logout).
  function resetFavorites() {
    favorites.value = []
    favoriteTotal.value = 0
    beans.value = beans.value.map((b) => ({ ...b, is_favorited: false }))
  }

  // Apply the new favorite state locally so every visible card updates at once.
  function applyFavorite(id: number, favorited: boolean) {
    const patch = (b: BeanItem): BeanItem => ({
      ...b,
      is_favorited: favorited,
      favorite_count: Math.max(0, b.favorite_count + (favorited ? 1 : -1)),
    })
    beans.value = beans.value.map((b) => (b.id === id ? patch(b) : b))
    if (favorited) {
      // Keep the state consistent if the bean is already shown in "我的收藏".
      favorites.value = favorites.value.map((b) => (b.id === id ? patch(b) : b))
    } else {
      favorites.value = favorites.value.filter((b) => b.id !== id)
      favoriteTotal.value = Math.max(0, favoriteTotal.value - 1)
    }
  }

  async function addFavorite(id: number) {
    applyFavorite(id, true)
    try {
      await favoriteBean(id)
    } catch (err) {
      applyFavorite(id, false)
      throw err
    }
  }

  async function removeFavorite(id: number) {
    const prevBeans = beans.value
    const prevFavorites = favorites.value
    const prevTotal = favoriteTotal.value
    applyFavorite(id, false)
    try {
      await unfavoriteBean(id)
    } catch (err) {
      beans.value = prevBeans
      favorites.value = prevFavorites
      favoriteTotal.value = prevTotal
      throw err
    }
  }

  return {
    beans,
    total,
    favorites,
    favoriteTotal,
    load,
    loadFavorites,
    resetFavorites,
    addFavorite,
    removeFavorite,
  }
})
