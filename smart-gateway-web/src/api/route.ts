import { del, get, post, put } from './http'
import type { Page, Route, RouteQuery } from '@/types'

export const listRoutes = (params: RouteQuery) => get<Page<Route>>('/api/routes', params)

export const createRoute = (data: Partial<Route>) => post<void>('/api/routes', data)

export const updateRoute = (id: number, data: Partial<Route>) => put<void>(`/api/routes/${id}`, data)

export const deleteRoute = (id: number) => del<void>(`/api/routes/${id}`)
