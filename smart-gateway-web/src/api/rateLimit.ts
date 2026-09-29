import { del, get, post, put } from './http'
import type { Page, RateLimitPolicy, RateLimitQuery } from '@/types'

export const listRateLimits = (params: RateLimitQuery) =>
  get<Page<RateLimitPolicy>>('/api/rate-limits', params)

export const createRateLimit = (data: Partial<RateLimitPolicy>) =>
  post<void>('/api/rate-limits', data)

export const updateRateLimit = (id: number, data: Partial<RateLimitPolicy>) =>
  put<void>(`/api/rate-limits/${id}`, data)

export const deleteRateLimit = (id: number) => del<void>(`/api/rate-limits/${id}`)
