import { useMemo } from 'react'
import { Avatar, Dropdown, Layout, Menu, Space, Tag } from 'antd'
import {
  ApartmentOutlined,
  AuditOutlined,
  ClusterOutlined,
  DashboardOutlined,
  DesktopOutlined,
  FileSearchOutlined,
  LogoutOutlined,
  SafetyOutlined,
  StopOutlined,
  TeamOutlined,
  ThunderboltOutlined,
  UserOutlined,
} from '@ant-design/icons'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'

const { Header, Sider, Content } = Layout

const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: '实时监控' },
  { key: '/devices', icon: <DesktopOutlined />, label: '设备管理' },
  { key: '/routes', icon: <ApartmentOutlined />, label: '路由管理' },
  { key: '/upstreams', icon: <ClusterOutlined />, label: '上游管理' },
  { key: '/rate-limits', icon: <ThunderboltOutlined />, label: '限流策略' },
  { key: '/circuit-breakers', icon: <SafetyOutlined />, label: '熔断配置' },
  { key: '/bans', icon: <StopOutlined />, label: '封禁管理' },
  { key: '/users', icon: <TeamOutlined />, label: '用户管理' },
  { key: '/request-logs', icon: <FileSearchOutlined />, label: '请求日志' },
  { key: '/audit-logs', icon: <AuditOutlined />, label: '审计日志' },
]

const roleText: Record<string, string> = {
  admin: '超级管理员',
  operator: '运维',
  viewer: '只读',
}

export default function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const username = useAuthStore((s) => s.username)
  const role = useAuthStore((s) => s.role)
  const logout = useAuthStore((s) => s.logout)

  const selectedKey = useMemo(() => {
    const match = menuItems.find((i) => location.pathname.startsWith(i.key))
    return match?.key ?? '/dashboard'
  }, [location.pathname])

  const onLogout = () => {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" width={208}>
        <div
          style={{
            height: 56,
            color: '#fff',
            fontWeight: 600,
            fontSize: 16,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          网关检测系统
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            justifyContent: 'flex-end',
            alignItems: 'center',
            boxShadow: '0 1px 4px rgba(0,21,41,0.08)',
          }}
        >
          <Dropdown
            placement="bottomRight"
            menu={{
              items: [{ key: 'logout', icon: <LogoutOutlined />, label: '退出登录' }],
              onClick: onLogout,
            }}
          >
            <Space style={{ cursor: 'pointer' }}>
              <Avatar size="small" icon={<UserOutlined />} />
              <span>{username || '未登录'}</span>
              <Tag color="blue">{role ? roleText[role] ?? role : '-'}</Tag>
            </Space>
          </Dropdown>
        </Header>
        <Content style={{ margin: 16 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
