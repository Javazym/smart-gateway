import { useMemo, useState } from 'react'
import { Button, Input, Popconfirm, Select, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useMutation, useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import { createRoute, deleteRoute, listRoutes, updateRoute } from '@/api/route'
import { listGroups } from '@/api/upstream'
import { listRateLimits } from '@/api/rateLimit'
import { listCircuitBreakers } from '@/api/circuitBreaker'
import type { Route, RouteQuery } from '@/types'
import FormModal, { type FieldConfig } from '@/components/FormModal'
import StatusTag from '@/components/StatusTag'
import { MATCH_TYPE_OPTIONS, STATUS_MAP, STATUS_OPTIONS } from '@/constants/options'
import { usePage } from '@/hooks/usePage'
import { usePermission } from '@/hooks/usePermission'

const MATCH_TYPE_TEXT: Record<string, string> = {
  path: '路径',
  prefix: '前缀',
  domain: '域名',
  header: '请求头',
}

export default function RoutesPage() {
  const { canWrite } = usePermission()
  const { query, tablePagination } = usePage()
  const [name, setName] = useState('')
  const [status, setStatus] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Route | null>(null)

  const filters: RouteQuery = { name: name || undefined, status }
  const listQuery = useQuery({
    queryKey: ['routes', { ...query, ...filters }],
    queryFn: () => listRoutes({ ...query, ...filters }),
  })

  const groupsQuery = useQuery({ queryKey: ['upstream-groups', 'all'], queryFn: () => listGroups({ page: 1, size: 1000 }) })
  const rateQuery = useQuery({ queryKey: ['rate-limits', 'all'], queryFn: () => listRateLimits({ page: 1, size: 1000 }) })
  const cbQuery = useQuery({ queryKey: ['circuit-breakers', 'all'], queryFn: () => listCircuitBreakers({ page: 1, size: 1000 }) })

  const groupMap = useMemo(() => {
    const m: Record<number, string> = {}
    groupsQuery.data?.list.forEach((g) => {
      m[g.id] = g.name
    })
    return m
  }, [groupsQuery.data])

  const rateNameMap = useMemo(() => {
    const m: Record<number, string> = {}
    rateQuery.data?.list.forEach((r) => {
      m[r.id] = r.name
    })
    return m
  }, [rateQuery.data])

  const groupOptions = (groupsQuery.data?.list ?? []).map((g) => ({ label: g.name, value: g.id }))
  const rateOptions = [{ label: '无', value: 0 }, ...(rateQuery.data?.list ?? []).map((r) => ({ label: r.name, value: r.id }))]
  const cbOptions = [{ label: '无', value: 0 }, ...(cbQuery.data?.list ?? []).map((c) => ({ label: c.name, value: c.id }))]

  const saveMut = useMutation({
    mutationFn: (data: Partial<Route>) => (editing ? updateRoute(editing.id, data) : createRoute(data)),
    onSuccess: () => {
      message.success('保存成功')
      setModalOpen(false)
      listQuery.refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteRoute(id),
    onSuccess: () => {
      message.success('删除成功')
      listQuery.refetch()
    },
  })

  const fields: FieldConfig[] = [
    { name: 'name', label: '路由名', type: 'text', required: true },
    { name: 'match_type', label: '匹配类型', type: 'select', required: true, options: MATCH_TYPE_OPTIONS },
    { name: 'match_pattern', label: '匹配表达式', type: 'text', required: true },
    { name: 'upstream_group_id', label: '上游分组', type: 'select', required: true, options: groupOptions },
    { name: 'strip_prefix', label: '剥离前缀', type: 'text' },
    { name: 'priority', label: '优先级(小优先)', type: 'number', min: 0 },
    { name: 'rate_limit_id', label: '限流策略', type: 'select', options: rateOptions },
    { name: 'circuit_breaker_id', label: '熔断配置', type: 'select', options: cbOptions },
    { name: 'status', label: '状态', type: 'select', options: STATUS_OPTIONS },
  ]

  const columns: ColumnsType<Route> = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '路由名', dataIndex: 'name' },
    { title: '匹配类型', dataIndex: 'match_type', width: 100, render: (v: string) => MATCH_TYPE_TEXT[v] ?? v },
    { title: '匹配表达式', dataIndex: 'match_pattern' },
    { title: '上游分组', dataIndex: 'upstream_group_id', render: (v: number) => groupMap[v] ?? v },
    { title: '优先级', dataIndex: 'priority', width: 80 },
    { title: '限流策略', dataIndex: 'rate_limit_id', render: (v: number) => (v ? rateNameMap[v] ?? v : '-') },
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
            <Popconfirm title="确认删除该路由？" onConfirm={() => delMut.mutate(record.id)}>
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
        match_type: editing.match_type,
        match_pattern: editing.match_pattern,
        upstream_group_id: editing.upstream_group_id,
        strip_prefix: editing.strip_prefix,
        priority: editing.priority,
        rate_limit_id: editing.rate_limit_id,
        circuit_breaker_id: editing.circuit_breaker_id,
        status: editing.status,
      }
    : { match_type: 'path', priority: 0, status: 1, rate_limit_id: 0, circuit_breaker_id: 0 }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="按路由名搜索"
          allowClear
          style={{ width: 220 }}
          onSearch={(v) => setName(v)}
        />
        <Select
          placeholder="状态"
          allowClear
          style={{ width: 120 }}
          options={STATUS_OPTIONS}
          onChange={(v) => setStatus(v as number | undefined)}
        />
        {canWrite && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); setModalOpen(true) }}>
            新增路由
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
        title={editing ? '编辑路由' : '新增路由'}
        fields={fields}
        initialValues={initialValues}
        confirmLoading={saveMut.isPending}
        onCancel={() => setModalOpen(false)}
        onSubmit={(values) => saveMut.mutate(values as Partial<Route>)}
      />
    </div>
  )
}
