import { del, get, post, put } from './http'
import type { CircuitBreaker, CircuitBreakerQuery, Page } from '@/types'

export const listCircuitBreakers = (params: CircuitBreakerQuery) =>
  get<Page<CircuitBreaker>>('/api/circuit-breakers', params)

export const createCircuitBreaker = (data: Partial<CircuitBreaker>) =>
  post<void>('/api/circuit-breakers', data)

export const updateCircuitBreaker = (id: number, data: Partial<CircuitBreaker>) =>
  put<void>(`/api/circuit-breakers/${id}`, data)

export const deleteCircuitBreaker = (id: number) =>
  del<void>(`/api/circuit-breakers/${id}`)
