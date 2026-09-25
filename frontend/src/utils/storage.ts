export const TOKEN_KEY = 'wjecoffeetaste_token'

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

const PENDING_FAVORITE_KEY = 'wjecoffeetaste_pending_favorite_bean'

// Remember a bean the user wanted to favorite before being asked to log in.
export function setPendingFavoriteBean(id: number) {
  localStorage.setItem(PENDING_FAVORITE_KEY, String(id))
}

export function takePendingFavoriteBean(): number | null {
  const raw = localStorage.getItem(PENDING_FAVORITE_KEY)
  localStorage.removeItem(PENDING_FAVORITE_KEY)
  const id = Number(raw)
  return raw !== null && Number.isInteger(id) && id > 0 ? id : null
}
