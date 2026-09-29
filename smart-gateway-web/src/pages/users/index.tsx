import { useState } from 'react'
import { Button, Input, Popconfirm, Select, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useMutation, useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import { createUser, deleteUser, listUsers, updateUser } from '@/api/users'
import type { User, UserCreate, UserQuery, UserUpdate } from '@/types'
import FormModal, { type FieldConfig } from '@/components/FormModal'
import StatusTag from '@/components/StatusTag'
import { ROLE_MAP, ROLE_OPTIONS, STATUS_MAP, STATUS_OPTIONS } from '@/constants/options'
import { usePage } from '@/hooks/usePage'
import { usePermission } from '@/hooks/usePermission'

export default function UsersPage() {
  const { canManageUsers } = usePermission()
  const { query, tablePagination } = usePage()
  const [username, setUsername] = useState('')
  const [status, setStatus] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)

  const filters: UserQuery = { username: username || undefined, status }
  const listQuery = useQuery({
    queryKey: ['users', { ...query, ...filters }],
    queryFn: () => listUsers({ ...query, ...filters }),
  })

  const saveMut = useMutation({
    mutationFn: (data: UserCreate | UserUpdate) =>
      editing ? updateUser(editing.id, data as UserUpdate) : createUser(data as UserCreate),
    onSuccess: () => {
      message.success('保存成功')
      setModalOpen(false)
      listQuery.refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteUser(id),
    onSuccess: () => {
      message.success('删除成功')
      listQuery.refetch()
    },
  })

  const createFields: FieldConfig[] = [
    { name: 'username', label: '用户名', type: 'text', required: true },
    { name: 'password', label: '密码', type: 'password', required: true },
    { name: 'role', label: '角色', type: 'select', required: true, options: ROLE_OPTIONS },
    { name: 'status', label: '状态', type: 'select', options: STATUS_OPTIONS },
  ]

  const editFields: FieldConfig[] = [
    { name: 'role', label: '角色', type: 'select', required: true, options: ROLE_OPTIONS },
    { name: 'status', label: '状态', type: 'select', options: STATUS_OPTIONS },
    { name: 'password', label: '重置密码(留空不修改)', type: 'password' },
  ]

  const columns: ColumnsType<User> = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '用户名', dataIndex: 'username' },
    { title: '角色', dataIndex: 'role', width: 120, render: (v: string) => <StatusTag value={v} map={ROLE_MAP} /> },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (v: number) => <StatusTag value={v} map={STATUS_MAP} />,
    },
    {
      title: '最近登录',
      dataIndex: 'last_login_at',
      width: 160,
      render: (v: string | null) => (v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-'),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 140,
      render: (_, record) =>
        canManageUsers ? (
          <Space>
            <Button type="link" size="small" onClick={() => { setEditing(record); setModalOpen(true) }}>
              编辑
            </Button>
            <Popconfirm title="确认删除该用户？" onConfirm={() => delMut.mutate(record.id)}>
              <Button type="link" size="small" danger>
                删除
              </Button>
            </Popconfirm>
          </Space>
        ) : null,
    },
  ]

  const handleSubmit = (values: Record<string, unknown>) => {
    if (editing) {
      const { password, ...rest } = values
      const payload: UserUpdate = { ...(rest as UserUpdate) }
      if (password) payload.password = password as string
      saveMut.mutate(payload)
    } else {
      saveMut.mutate(values as unknown as UserCreate)
    }
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="按用户名搜索"
          allowClear
          style={{ width: 220 }}
          onSearch={(v) => setUsername(v)}
        />
        <Select
          placeholder="状态"
          allowClear
          style={{ width: 120 }}
          options={STATUS_OPTIONS}
          onChange={(v) => setStatus(v as number | undefined)}
        />
        {canManageUsers && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); setModalOpen(true) }}>
            新增用户
          </Button>
        )}
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={listQuery.data?.list ?? []}
        loading={listQuery.isLoading}
        pagination={tablePagination(listQuery.data?.total ?? 0)}
      />
      <FormModal
        open={modalOpen}
        title={editing ? '编辑用户' : '新增用户'}
        fields={editing ? editFields : createFields}
        initialValues={
          editing
            ? { role: editing.role, status: editing.status }
            : { role: 'viewer', status: 1 }
        }
        confirmLoading={saveMut.isPending}
        onCancel={() => setModalOpen(false)}
        onSubmit={handleSubmit}
      />
    </div>
  )
}
