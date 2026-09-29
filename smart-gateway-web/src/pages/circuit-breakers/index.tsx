import { useState } from 'react'
import { Button, Popconfirm, Select, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useMutation, useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import {
  createCircuitBreaker,
  deleteCircuitBreaker,
  listCircuitBreakers,
  updateCircuitBreaker,
} from '@/api/circuitBreaker'
import type { CircuitBreaker, CircuitBreakerQuery } from '@/types'
import FormModal, { type FieldConfig } from '@/components/FormModal'
import StatusTag from '@/components/StatusTag'
import { STATUS_MAP, STATUS_OPTIONS } from '@/constants/options'
import { usePage } from '@/hooks/usePage'
import { usePermission } from '@/hooks/usePermission'

export default function CircuitBreakersPage() {
  const { canWrite } = usePermission()
  const { query, tablePagination } = usePage()
  const [status, setStatus] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<CircuitBreaker | null>(null)

  const filters: CircuitBreakerQuery = { status }
  const listQuery = useQuery({
    queryKey: ['circuit-breakers', { ...query, ...filters }],
    queryFn: () => listCircuitBreakers({ ...query, ...filters }),
  })

  const saveMut = useMutation({
    mutationFn: (data: Partial<CircuitBreaker>) =>
      editing ? updateCircuitBreaker(editing.id, data) : createCircuitBreaker(data),
    onSuccess: () => {
      message.success('保存成功')
      setModalOpen(false)
      listQuery.refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteCircuitBreaker(id),
    onSuccess: () => {
      message.success('删除成功')
      listQuery.refetch()
    },
  })

  const fields: FieldConfig[] = [
    { name: 'name', label: '熔断器名', type: 'text', required: true },
    { name: 'failure_threshold', label: '连续失败阈值', type: 'number', required: true, min: 1 },
    { name: 'success_threshold', label: '半开成功阈值', type: 'number', required: true, min: 1 },
    { name: 'timeout_ms', label: '调用超时', type: 'number', required: true, min: 1, addonAfter: '毫秒' },
    { name: 'open_duration_seconds', label: '熔断打开时长', type: 'number', required: true, min: 1, addonAfter: '秒' },
    { name: 'half_open_max_requests', label: '半开探测数', type: 'number', required: true, min: 1 },
    { name: 'status', label: '状态', type: 'select', options: STATUS_OPTIONS },
  ]

  const columns: ColumnsType<CircuitBreaker> = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '熔断器名', dataIndex: 'name' },
    { title: '失败阈值', dataIndex: 'failure_threshold', width: 90 },
    { title: '成功阈值', dataIndex: 'success_threshold', width: 90 },
    { title: '超时(ms)', dataIndex: 'timeout_ms', width: 100 },
    { title: '打开时长(s)', dataIndex: 'open_duration_seconds', width: 110 },
    { title: '半开探测', dataIndex: 'half_open_max_requests', width: 90 },
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
            <Popconfirm title="确认删除该熔断器？" onConfirm={() => delMut.mutate(record.id)}>
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
        failure_threshold: editing.failure_threshold,
        success_threshold: editing.success_threshold,
        timeout_ms: editing.timeout_ms,
        open_duration_seconds: editing.open_duration_seconds,
        half_open_max_requests: editing.half_open_max_requests,
        status: editing.status,
      }
    : { failure_threshold: 5, success_threshold: 3, timeout_ms: 3000, open_duration_seconds: 30, half_open_max_requests: 1, status: 1 }

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
            新增熔断器
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
        title={editing ? '编辑熔断配置' : '新增熔断配置'}
        fields={fields}
        initialValues={initialValues}
        confirmLoading={saveMut.isPending}
        onCancel={() => setModalOpen(false)}
        onSubmit={(values) => saveMut.mutate(values as Partial<CircuitBreaker>)}
      />
    </div>
  )
}
