import { createBrowserRouter, Navigate } from 'react-router-dom'
import AuthGuard from './AuthGuard'
import MainLayout from '@/layouts/MainLayout'
import LoginPage from '@/pages/login'
import DashboardPage from '@/pages/dashboard'
import DevicesPage from '@/pages/devices'
import RoutesPage from '@/pages/routes'
import UpstreamsPage from '@/pages/upstreams'
import RateLimitsPage from '@/pages/rate-limits'
import CircuitBreakersPage from '@/pages/circuit-breakers'
import BansPage from '@/pages/bans'
import UsersPage from '@/pages/users'
import RequestLogsPage from '@/pages/request-logs'
import AuditLogsPage from '@/pages/audit-logs'

export const router = createBrowserRouter([
  { path: '/login', element: <LoginPage /> },
  {
    path: '/',
    element: (
      <AuthGuard>
        <MainLayout />
      </AuthGuard>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'devices', element: <DevicesPage /> },
      { path: 'routes', element: <RoutesPage /> },
      { path: 'upstreams', element: <UpstreamsPage /> },
      { path: 'rate-limits', element: <RateLimitsPage /> },
      { path: 'circuit-breakers', element: <CircuitBreakersPage /> },
      { path: 'bans', element: <BansPage /> },
      { path: 'users', element: <UsersPage /> },
      { path: 'request-logs', element: <RequestLogsPage /> },
      { path: 'audit-logs', element: <AuditLogsPage /> },
    ],
  },
  { path: '*', element: <Navigate to="/dashboard" replace /> },
])
