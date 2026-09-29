import { useState } from 'react'
import { Button, Popconfirm, Select, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useMutation, useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import { createBan, deleteBan, listBans, unban } from '@/api/ban'
import type { Ban, BanCreate, BanQuery } from '@/types'
import FormModal, { type FieldConfig } from '@/components/FormModal'
import StatusTag from '@/components/StatusTag'
import { BAN_STATUS_MAP, BAN_TYPE_OPTIONS, STATUS_OPTIONS } from '@/constants/options'
import { usePage } from '@/hooks/usePage'
import { usePermission } from '@/hooks/usePermission'

const BAN_TYPE_TEXT: Record<string, string> = {
  device: '设备',
  ip: 'IP',
}

export default function BansPage() {
  const { canWrite } = usePermission()
  const { query, tablePagination } = usePage()
  const [status, setStatus] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)

  const filters: BanQuery = { status }
  const listQuery = useQuery({
    queryKey: ['bans', { ...query, ...filters }],
    queryFn: () => listBans({ ...query, ...filters }),
  })

  const createMut = useMutation({
    mutationFn: (data: BanCreate) => createBan(data),
    onSuccess: () => {
      message.success('封禁成功')
      setModalOpen(false)
      listQuery.refetch()
    },
  })

  const unbanMut = useMutation({
    mutationFn: (id: number) => unban(id),
    onSuccess: () => {
      message.success('已解封')
      listQuery.refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteBan(id),
    onSuccess: () => {
      message.success('删除成功')
      listQuery.refetch()
    },
  })

  const fields: FieldConfig[] = [
    { name: 'ban_type', label: '封禁类型', type: 'select', required: true, options: BAN_TYPE_OPTIONS },
    { name: 'target_value', label: '封禁目标(指纹ID/IP)', type: 'text', required: true },
    { name: 'reason', label: '封禁原因', type: 'text' },
  ]

  const columns: ColumnsType<Ban> = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '类型', dataIndex: 'ban_type', width: 80, render: (v: string) => BAN_TYPE_TEXT[v] ?? v },
    { title: '目标', dataIndex: 'target_value', ellipsis: true },
    { title: '原因', dataIndex: 'reason', ellipsis: true, render: (v: string | null) => v ?? '-' },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (v: number) => <StatusTag value={v} map={BAN_STATUS_MAP} />,
    },
    {
      title: '封禁时间',
      dataIndex: 'banned_at',
      width: 160,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '解封时间',
      dataIndex: 'unbanned_at',
      width: 160,
      render: (v: string | null) => (v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 140,
      render: (_, record) =>
        canWrite ? (
          <Space>
            {record.status === 1 && (
              <Popconfirm title="确认解封？" onConfirm={() => unbanMut.mutate(record.id)}>
                <Button type="link" size="small">
                  解封
                </Button>
              </Popconfirm>
            )}
            <Popconfirm title="确认删除该记录？" onConfirm={() => delMut.mutate(record.id)}>
              <Button type="link" size="small" danger>
                删除
              </Button>
            </Popconfirm>
          </Space>
        ) : null,
    },
  ]

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Select
          placeholder="状态"
          allowClear
          style={{ width: 120 }}
          options={STATUS_OPTIONS}
          onChange={(v) => setStatus(v as number | undefined)}
        />
        {canWrite && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            新增封禁
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
        title="新增封禁"
        fields={fields}
        initialValues={{ ban_type: 'device' }}
        confirmLoading={createMut.isPending}
        onCancel={() => setModalOpen(false)}
        onSubmit={(values) => createMut.mutate(values as unknown as BanCreate)}
      />
    </div>
  )
}
