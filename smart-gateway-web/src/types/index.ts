// 与后端 internal/model 对齐的接口定义（json 字段使用 snake_case）

export type Role = 'admin' | 'operator' | 'viewer'

// ---------------- 通用响应 ----------------

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
  details?: string[]
  traceId?: string
}

export interface Page<T> {
  list: T[]
  total: number
  page_num: number
  page_size: number
}

export interface PageQuery {
  page?: number
  size?: number
}

// ---------------- 认证 ----------------

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  username: string
  role: Role
}

// ---------------- 用户 ----------------

export interface User {
  id: number
  username: string
  role: Role
  status: number
  last_login_at: string | null
  created_at: string
  updated_at: string
}

export interface UserQuery extends PageQuery {
  username?: string
  status?: number
}

export interface UserCreate {
  username: string
  password: string
  role?: Role
  status?: number
}

export interface UserUpdate {
  role?: Role
  status?: number
  password?: string
}

// ---------------- 上游 ----------------

export interface UpstreamGroup {
  id: number
  name: string
  lb_algorithm: string
  health_check_path: string | null
  health_check_interval: number
  status: number
  created_at: string
  updated_at: string
  upstreams?: Upstream[]
}

export interface Upstream {
  id: number
  group_id: number
  name: string
  scheme: string
  host: string
  port: number
  weight: number
  max_conns: number
  health_status: string
  status: number
  created_at: string
  updated_at: string
}

export interface GroupQuery extends PageQuery {
  status?: number
}

// ---------------- 路由 ----------------

export interface Route {
  id: number
  name: string
  match_type: string
  match_pattern: string
  upstream_group_id: number
  rate_limit_id: number
  circuit_breaker_id: number
  strip_prefix: string
  priority: number
  status: number
  created_at: string
  updated_at: string
}

export interface RouteQuery extends PageQuery {
  status?: number
  name?: string
}

// ---------------- 设备 ----------------

export interface Device {
  id: number
  fingerprint_id: string
  ja4_hash: string | null
  ja3_hash: string | null
  os_fingerprint: string | null
  http_header_hash: string | null
  user_agent: string | null
  last_seen_ip: string | null
  device_status: string
  ban_reason: string | null
  first_seen_at: string
  last_seen_at: string
  created_at: string
  updated_at: string
}

export interface DeviceQuery extends PageQuery {
  status?: string
}

// ---------------- 封禁 ----------------

export interface Ban {
  id: number
  device_id: number | null
  ban_type: string
  target_value: string
  reason: string | null
  operator_id: number | null
  status: number
  banned_at: string
  unbanned_at: string | null
  created_at: string
}

export interface BanQuery extends PageQuery {
  status?: number
}

export interface BanCreate {
  device_id?: number
  ban_type: string
  target_value: string
  reason?: string
}

// ---------------- 限流 ----------------

export interface RateLimitPolicy {
  id: number
  name: string
  algorithm: string
  rate: number
  burst: number
  window_seconds: number
  key_type: string
  status: number
  created_at: string
  updated_at: string
}

export interface RateLimitQuery extends PageQuery {
  status?: number
}

// ---------------- 熔断 ----------------

export interface CircuitBreaker {
  id: number
  name: string
  failure_threshold: number
  success_threshold: number
  timeout_ms: number
  open_duration_seconds: number
  half_open_max_requests: number
  status: number
  created_at: string
  updated_at: string
}

export interface CircuitBreakerQuery extends PageQuery {
  status?: number
}

// ---------------- 请求日志 ----------------

export interface RequestLog {
  id: number
  request_id: string
  route_id: number
  device_fingerprint_id: string
  source_ip: string
  method: string
  path: string
  upstream_id: number
  status_code: number
  latency_ms: number
  response_size: number
  received_at: string
  created_at: string
}

export interface RequestLogQuery extends PageQuery {
  device_fingerprint_id?: string
  upstream_id?: number
  route_id?: number
  source_ip?: string
}

// ---------------- 审计日志 ----------------

export interface AuditLog {
  id: number
  operator_id: number | null
  username: string | null
  action: string
  target_type: string | null
  target_id: string | null
  detail: string | null
  source_ip: string | null
  created_at: string
}

export interface AuditLogQuery extends PageQuery {
  username?: string
  action?: string
}
