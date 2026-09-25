import request from '@/utils/request'
import type { CoffeeBean } from '@/constants/bean'
import type { PageData } from '@/types/api'

export interface BeanQuery {
  page?: number
  page_size?: number
  origin?: string
  process?: string
  keyword?: string
}

export function listBeans(params: BeanQuery) {
  return request.get<never, PageData<CoffeeBean>>('/beans', { params })
}
export function listFavoriteBeans(params: BeanQuery) {
  return request.get<never, PageData<CoffeeBean>>('/beans/favorites', { params })
}
export function favoriteBean(id: number) {
  return request.post<never, { favorite: { id: number; user_id: number; bean_id: number }; favorite_count: number }>(`/beans/${id}/favorite`)
}
export function unfavoriteBean(id: number) {
  return request.delete<never, { unfavorited: boolean; favorite_count: number }>(`/beans/${id}/favorite`)
}
export function createBean(payload: Partial<CoffeeBean>) { return request.post<never, CoffeeBean>('/beans', payload) }
export function updateBean(id: number, payload: Partial<CoffeeBean>) { return request.put<never, CoffeeBean>(`/beans/${id}`, payload) }
export function deleteBean(id: number) { return request.delete<never, { deleted: boolean }>(`/beans/${id}`) }
