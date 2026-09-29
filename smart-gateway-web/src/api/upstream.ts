import { del, get, post, put } from './http'
import type { GroupQuery, Page, Upstream, UpstreamGroup } from '@/types'

// 分组
export const listGroups = (params: GroupQuery) =>
  get<Page<UpstreamGroup>>('/api/upstream-groups', params)

export const createGroup = (data: Partial<UpstreamGroup>) =>
  post<void>('/api/upstream-groups', data)

export const updateGroup = (id: number, data: Partial<UpstreamGroup>) =>
  put<void>(`/api/upstream-groups/${id}`, data)

export const deleteGroup = (id: number) => del<void>(`/api/upstream-groups/${id}`)

// 节点（列表接口返回数组，非分页）
export const listUpstreams = (groupId?: number) =>
  get<Upstream[]>('/api/upstreams', groupId ? { group_id: groupId } : undefined)

export const createUpstream = (data: Partial<Upstream>) =>
  post<void>('/api/upstreams', data)

export const updateUpstream = (id: number, data: Partial<Upstream>) =>
  put<void>(`/api/upstreams/${id}`, data)

export const deleteUpstream = (id: number) => del<void>(`/api/upstreams/${id}`)
