import { useAuthStore } from '@/store/auth'

// 前端 RBAC：admin 全部，operator 除用户管理写外全部，viewer 只读。
export function usePermission() {
  const role = useAuthStore((s) => s.role)
  const canWrite = role === 'admin' || role === 'operator'
  const canManageUsers = role === 'admin'
  return { role, canWrite, canManageUsers }
}
