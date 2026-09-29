import axios, { type AxiosError } from 'axios'
import { message } from 'antd'
import type { ApiResponse } from '@/types'

const TOKEN_KEY = 'sg_token'

export const getToken = () => localStorage.getItem(TOKEN_KEY)
export const setToken = (t: string) => localStorage.setItem(TOKEN_KEY, t)
export const clearToken = () => localStorage.removeItem(TOKEN_KEY)

const http = axios.create({
  baseURL: '/',
  timeout: 15000,
})

http.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse
    if (body && typeof body.code === 'number' && body.code !== 0) {
      message.error(body.message || '请求失败')
      return Promise.reject(new Error(body.message || '请求失败'))
    }
    return response
  },
  (error: AxiosError<ApiResponse>) => {
    const status = error.response?.status
    const msg = error.response?.data?.message
    if (status === 401) {
      clearToken()
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
      message.error('登录已过期，请重新登录')
    } else if (status === 403) {
      message.error(msg || '无权限操作')
    } else if (status === 404) {
      message.error(msg || '资源不存在')
    } else {
      message.error(msg || error.message || '网络错误')
    }
    return Promise.reject(error)
  },
)

export async function get<T>(url: string, params?: object): Promise<T> {
  const res = await http.get<ApiResponse<T>>(url, { params })
  return res.data.data as T
}

export async function post<T>(url: string, data?: object): Promise<T> {
  const res = await http.post<ApiResponse<T>>(url, data)
  return res.data.data as T
}

export async function put<T>(url: string, data?: object): Promise<T> {
  const res = await http.put<ApiResponse<T>>(url, data)
  return res.data.data as T
}

export async function del<T>(url: string): Promise<T> {
  const res = await http.delete<ApiResponse<T>>(url)
  return res.data.data as T
}

export default http
