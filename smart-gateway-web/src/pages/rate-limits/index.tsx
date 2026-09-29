import { useState } from 'react'
import { Button, Popconfirm, Select, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useMutation, useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import {
  createRateLimit,
  deleteRateLimit,
  listRateLimits,
  updateRateLimit,
} from '@/api/rateLimit'
import type { RateLimitPolicy, RateLimitQuery } from '@/types'
import FormModal, { type FieldConfig } from '@/components/FormModal'
import StatusTag from '@/components/StatusTag'
import { KEY_TYPE_OPTIONS, RATE_ALGORITHM_OPTIONS, STATUS_MAP, STATUS_OPTIONS } from '@/constants/options'
import { usePage } from '@/hooks/usePage'
import { usePermission } from '@/hooks/usePermission'

const ALGORITHM_TEXT: Record<string, string> = {
  token_bucket: '令牌桶',
  leaky_bucket: '漏桶',
  sliding_window: '滑动窗口',
  fixed_window: '固定窗口',
}

const KEY_TYPE_TEXT: Record<string, string> = {
  ip: 'IP',
  device: '设备指纹',
  route: '路由',
}

export default function RateLimitsPage() {
  const { canWrite } = usePermission()
  const { query, tablePagination } = usePage()
  const [status, setStatus] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<RateLimitPolicy | null>(null)

  const filters: RateLimitQuery = { status }
  const listQuery = useQuery({
    queryKey: ['rate-limits', { ...query, ...filters }],
    queryFn: () => listRateLimits({ ...query, ...filters }),
  })

  const saveMut = useMutation({
    mutationFn: (data: Partial<RateLimitPolicy>) =>
      editing ? updateRateLimit(editing.id, data) : createRateLimit(data),
    onSuccess: () => {
      message.success('保存成功')
      setModalOpen(false)
      listQuery.refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteRateLimit(id),
    onSuccess: () => {
      message.success('删除成功')
      listQuery.refetch()
    },
  })

  const fields: FieldConfig[] = [
    { name: 'name', label: '策略名', type: 'text', required: true },
    { name: 'algorithm', label: '算法', type: 'select', required: true, options: RATE_ALGORITHM_OPTIONS },
    { name: 'rate', label: '阈值(次数)', type: 'number', required: true, min: 1 },
    { name: 'burst', label: '突发容量', type: 'number', min: 0 },
    { name: 'window_seconds', label: '时间窗口', type: 'number', required: true, min: 1, addonAfter: '秒' },
    { name: 'key_type', label: '限流维度', type: 'select', required: true, options: KEY_TYPE_OPTIONS },
    { name: 'status', label: '状态', type: 'select', options: STATUS_OPTIONS },
  ]

  const columns: ColumnsType<RateLimitPolicy> = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '策略名', dataIndex: 'name' },
    { title: '算法', dataIndex: 'algorithm', width: 110, render: (v: string) => ALGORITHM_TEXT[v] ?? v },
    { title: '阈值', dataIndex: 'rate', width: 80 },
    { title: '突发', dataIndex: 'burst', width: 80 },
    { title: '窗口(秒)', dataIndex: 'window_seconds', width: 90 },
    { title: '维度', dataIndex: 'key_type', width: 100, render: (v: string) => KEY_TYPE_TEXT[v] ?? v },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (v: number) => <StatusTag value={v} map={STATUS_MAP} />,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      width: 160,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 140,
      render: (_, record) =>
        canWrite ? (
          <Space>
            <Button type="link" size="small" onClick={() => { setEditing(record); setModalOpen(true) }}>
              编辑
            </Button>
            <Popconfirm title="确认删除该策略？" onConfirm={() => delMut.mutate(record.id)}>
              <Button type="link" size="small" danger>
                删除
              </Button>
            </Popconfirm>
          </Space>
        ) : null,
    },
  ]

  const initialValues = editing
    ? {
        name: editing.name,
        algorithm: editing.algorithm,
        rate: editing.rate,
        burst: editing.burst,
        window_seconds: editing.window_seconds,
        key_type: editing.key_type,
        status: editing.status,
      }
    : { algorithm: 'token_bucket', burst: 0, window_seconds: 1, key_type: 'ip', status: 1 }

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
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); setModalOpen(true) }}>
            新增策略
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
        title={editing ? '编辑限流策略' : '新增限流策略'}
        fields={fields}
        initialValues={initialValues}
        confirmLoading={saveMut.isPending}
        onCancel={() => setModalOpen(false)}
        onSubmit={(values) => saveMut.mutate(values as Partial<RateLimitPolicy>)}
      />
    </div>
  )
}
