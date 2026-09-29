import { del, get, put } from './http'
import type { Device, DeviceQuery, Page } from '@/types'

export const listDevices = (params: DeviceQuery) => get<Page<Device>>('/api/devices', params)

export const getDevice = (id: number) => get<Device>(`/api/devices/${id}`)

export const updateDeviceStatus = (id: number, status: string) =>
  put<void>(`/api/devices/${id}/status`, { status })

export const deleteDevice = (id: number) => del<void>(`/api/devices/${id}`)
