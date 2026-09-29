import { get } from './http'
import type { AuditLog, AuditLogQuery, Page, RequestLog, RequestLogQuery } from '@/types'

export const listRequestLogs = (params: RequestLogQuery) =>
  get<Page<RequestLog>>('/api/request-logs', params)

export const listAuditLogs = (params: AuditLogQuery) =>
  get<Page<AuditLog>>('/api/audit-logs', params)
