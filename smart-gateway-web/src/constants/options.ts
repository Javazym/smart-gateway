// 下拉选项与状态标签映射，集中管理，供表单与表格复用

export interface Option {
  label: string
  value: string | number
}

export const STATUS_OPTIONS: Option[] = [
  { label: '启用', value: 1 },
  { label: '停用', value: 0 },
]

export const DEVICE_STATUS_OPTIONS: Option[] = [
  { label: '正常', value: 'normal' },
  { label: '可疑', value: 'suspected' },
  { label: '已封禁', value: 'banned' },
]

export const LB_ALGORITHM_OPTIONS: Option[] = [
  { label: '轮询', value: 'round_robin' },
  { label: '加权轮询', value: 'weighted_rr' },
  { label: 'IP哈希', value: 'ip_hash' },
  { label: '最小连接数', value: 'least_conn' },
  { label: '一致性哈希', value: 'consistent_hash' },
]

export const RATE_ALGORITHM_OPTIONS: Option[] = [
  { label: '令牌桶', value: 'token_bucket' },
  { label: '漏桶', value: 'leaky_bucket' },
  { label: '滑动窗口', value: 'sliding_window' },
  { label: '固定窗口', value: 'fixed_window' },
]

export const KEY_TYPE_OPTIONS: Option[] = [
  { label: 'IP', value: 'ip' },
  { label: '设备指纹', value: 'device' },
  { label: '路由', value: 'route' },
]

export const MATCH_TYPE_OPTIONS: Option[] = [
  { label: '路径', value: 'path' },
  { label: '前缀', value: 'prefix' },
  { label: '域名', value: 'domain' },
  { label: '请求头', value: 'header' },
]

export const BAN_TYPE_OPTIONS: Option[] = [
  { label: '设备', value: 'device' },
  { label: 'IP', value: 'ip' },
]

export const ROLE_OPTIONS: Option[] = [
  { label: '超级管理员', value: 'admin' },
  { label: '运维', value: 'operator' },
  { label: '只读', value: 'viewer' },
]

// ---------------- 状态标签映射 ----------------

export interface TagMeta {
  text: string
  color: string
}

export const STATUS_MAP: Record<string, TagMeta> = {
  '1': { text: '启用', color: 'green' },
  '0': { text: '停用', color: 'red' },
}

export const DEVICE_STATUS_MAP: Record<string, TagMeta> = {
  normal: { text: '正常', color: 'green' },
  suspected: { text: '可疑', color: 'orange' },
  banned: { text: '已封禁', color: 'red' },
}

export const ROLE_MAP: Record<string, TagMeta> = {
  admin: { text: '超级管理员', color: 'red' },
  operator: { text: '运维', color: 'blue' },
  viewer: { text: '只读', color: 'default' },
}

export const HEALTH_MAP: Record<string, TagMeta> = {
  healthy: { text: '健康', color: 'green' },
  unhealthy: { text: '异常', color: 'red' },
  unknown: { text: '未知', color: 'default' },
}

export const BAN_STATUS_MAP: Record<string, TagMeta> = {
  '1': { text: '生效', color: 'red' },
  '0': { text: '已解除', color: 'default' },
}
