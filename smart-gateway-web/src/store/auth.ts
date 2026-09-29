import { create } from 'zustand'
import { clearToken, getToken, setToken } from '@/api/http'
import type { Role } from '@/types'

interface AuthState {
  token: string | null
  username: string
  role: Role | ''
  login: (token: string, username: string, role: Role) => void
  logout: () => void
}

// 单一数据源：token 存 localStorage('sg_token')，由 http.ts 拦截器读取
export const useAuthStore = create<AuthState>((set) => ({
  token: getToken(),
  username: localStorage.getItem('sg_username') || '',
  role: (localStorage.getItem('sg_role') as Role | null) || '',
  login: (token, username, role) => {
    setToken(token)
    localStorage.setItem('sg_username', username)
    localStorage.setItem('sg_role', role)
    set({ token, username, role })
  },
  logout: () => {
    clearToken()
    localStorage.removeItem('sg_username')
    localStorage.removeItem('sg_role')
    set({ token: null, username: '', role: '' })
  },
}))
