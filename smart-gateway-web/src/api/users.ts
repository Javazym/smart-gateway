import { del, get, post, put } from './http'
import type { Page, User, UserCreate, UserQuery, UserUpdate } from '@/types'

export const listUsers = (params: UserQuery) => get<Page<User>>('/api/users', params)

export const createUser = (data: UserCreate) => post<void>('/api/users', data)

export const updateUser = (id: number, data: UserUpdate) => put<void>(`/api/users/${id}`, data)

export const deleteUser = (id: number) => del<void>(`/api/users/${id}`)
