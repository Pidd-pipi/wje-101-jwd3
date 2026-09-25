import request from '@/utils/request'
import type { BeanItem, CoffeeBean } from '@/constants/bean'
import type { PageData } from '@/types/api'

export interface BeanListParams {
  page?: number
  page_size?: number
  origin?: string
  process?: string
  keyword?: string
}

export function listBeans(params: BeanListParams) {
  return request.get<never, PageData<BeanItem>>('/beans', { params })
}
export function listFavoriteBeans(params: BeanListParams) {
  return request.get<never, PageData<BeanItem>>('/beans/favorites', { params })
}
export function favoriteBean(id: number) {
  return request.post<never, { id: number; user_id: number; bean_id: number; created_at: string }>(`/beans/${id}/favorite`)
}
export function unfavoriteBean(id: number) {
  return request.delete<never, { unfavorited: boolean }>(`/beans/${id}/favorite`)
}
export function createBean(payload: Partial<CoffeeBean>) { return request.post<never, CoffeeBean>('/beans', payload) }
export function updateBean(id: number, payload: Partial<CoffeeBean>) { return request.put<never, CoffeeBean>(`/beans/${id}`, payload) }
export function deleteBean(id: number) { return request.delete<never, { deleted: boolean }>(`/beans/${id}`) }
