export const TOKEN_KEY = 'wjecoffeetaste_token'
export const PENDING_BEAN_FAVORITE_KEY = 'wjecoffeetaste_pending_bean_favorite'

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

// Anonymous users that click "favorite" are sent to login; the target bean id
// is persisted here so it can be favorited automatically after they return.
export function setPendingBeanFavorite(beanId: number) {
  localStorage.setItem(PENDING_BEAN_FAVORITE_KEY, String(beanId))
}

export function takePendingBeanFavorite(): number | null {
  const raw = localStorage.getItem(PENDING_BEAN_FAVORITE_KEY)
  if (raw === null) return null
  localStorage.removeItem(PENDING_BEAN_FAVORITE_KEY)
  const id = Number(raw)
  return Number.isFinite(id) && id > 0 ? id : null
}
