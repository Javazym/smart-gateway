import { del, get, post } from './http'
import type { Ban, BanCreate, BanQuery, Page } from '@/types'

export const listBans = (params: BanQuery) => get<Page<Ban>>('/api/bans', params)

export const createBan = (data: BanCreate) => post<void>('/api/bans', data)

export const unban = (id: number) => post<void>(`/api/bans/${id}/unban`)

export const deleteBan = (id: number) => del<void>(`/api/bans/${id}`)
