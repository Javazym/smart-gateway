import { useState } from 'react'
import { Button, Popconfirm, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useMutation, useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import {
  createGroup,
  createUpstream,
  deleteGroup,
  deleteUpstream,
  listGroups,
  listUpstreams,
  updateGroup,
  updateUpstream,
} from '@/api/upstream'
import type { Upstream, UpstreamGroup } from '@/types'
import FormModal, { type FieldConfig } from '@/components/FormModal'
import StatusTag from '@/components/StatusTag'
import { HEALTH_MAP, LB_ALGORITHM_OPTIONS, STATUS_MAP, STATUS_OPTIONS } from '@/constants/options'
import { usePage } from '@/hooks/usePage'
import { usePermission } from '@/hooks/usePermission'

const LB_TEXT: Record<string, string> = {
  round_robin: '轮询',
  weighted_rr: '加权轮询',
  ip_hash: 'IP哈希',
  least_conn: '最小连接数',
  consistent_hash: '一致性哈希',
}

const SCHEME_OPTIONS = [
  { label: 'http', value: 'http' },
  { label: 'https', value: 'https' },
]

function NodesTable({ groupId, canWrite }: { groupId: number; canWrite: boolean }) {
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['upstreams', groupId],
    queryFn: () => listUpstreams(groupId),
  })
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Upstream | null>(null)

  const saveMut = useMutation({
    mutationFn: (d: Partial<Upstream>) =>
      editing ? updateUpstream(editing.id, d) : createUpstream({ ...d, group_id: groupId }),
    onSuccess: () => {
      message.success('保存成功')
      setModalOpen(false)
      refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteUpstream(id),
    onSuccess: () => {
      message.success('删除成功')
      refetch()
    },
  })

  const fields: FieldConfig[] = [
    { name: 'name', label: '节点名', type: 'text', required: true },
    { name: 'scheme', label: '协议', type: 'select', required: true, options: SCHEME_OPTIONS },
    { name: 'host', label: '主机', type: 'text', required: true },
    { name: 'port', label: '端口', type: 'number', required: true, min: 1, max: 65535 },
    { name: 'weight', label: '权重', type: 'number', min: 0 },
    { name: 'max_conns', label: '最大连接数', type: 'number', min: 0 },
    { name: 'status', label: '状态', type: 'select', options: STATUS_OPTIONS },
  ]

  const columns: ColumnsType<Upstream> = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '节点名', dataIndex: 'name' },
    { title: '协议', dataIndex: 'scheme', width: 80 },
    {
      title: '地址',
      key: 'addr',
      render: (_, r) => `${r.host}:${r.port}`,
    },
    { title: '权重', dataIndex: 'weight', width: 70 },
    { title: '最大连接', dataIndex: 'max_conns', width: 90 },
    {
      title: '健康状态',
      dataIndex: 'health_status',
      width: 100,
      render: (v: string) => <StatusTag value={v} map={HEALTH_MAP} />,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 80,
      render: (v: number) => <StatusTag value={v} map={STATUS_MAP} />,
    },
    {
      title: '操作',
      key: 'action',
      width: 130,
      render: (_, record) =>
        canWrite ? (
          <Space>
            <Button type="link" size="small" onClick={() => { setEditing(record); setModalOpen(true) }}>
              编辑
            </Button>
            <Popconfirm title="确认删除该节点？" onConfirm={() => delMut.mutate(record.id)}>
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
        scheme: editing.scheme,
        host: editing.host,
        port: editing.port,
        weight: editing.weight,
        max_conns: editing.max_conns,
        status: editing.status,
      }
    : { scheme: 'http', weight: 1, max_conns: 0, status: 1 }

  return (
    <div style={{ padding: '8px 16px' }}>
      {canWrite && (
        <Button
          type="dashed"
          size="small"
          icon={<PlusOutlined />}
          style={{ marginBottom: 12 }}
          onClick={() => { setEditing(null); setModalOpen(true) }}
        >
          新增节点
        </Button>
      )}
      <Table
        rowKey="id"
        size="small"
        columns={columns}
        dataSource={data ?? []}
        loading={isLoading}
        pagination={false}
      />
      <FormModal
        open={modalOpen}
        title={editing ? '编辑节点' : '新增节点'}
        fields={fields}
        initialValues={initialValues}
        confirmLoading={saveMut.isPending}
        onCancel={() => setModalOpen(false)}
        onSubmit={(values) => saveMut.mutate(values as Partial<Upstream>)}
      />
    </div>
  )
}

export default function UpstreamsPage() {
  const { canWrite } = usePermission()
  const { query, tablePagination } = usePage()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<UpstreamGroup | null>(null)

  const listQuery = useQuery({
    queryKey: ['upstream-groups', query],
    queryFn: () => listGroups(query),
  })

  const saveMut = useMutation({
    mutationFn: (data: Partial<UpstreamGroup>) =>
      editing ? updateGroup(editing.id, data) : createGroup(data),
    onSuccess: () => {
      message.success('保存成功')
      setModalOpen(false)
      listQuery.refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteGroup(id),
    onSuccess: () => {
      message.success('删除成功')
      listQuery.refetch()
    },
  })

  const fields: FieldConfig[] = [
    { name: 'name', label: '分组名', type: 'text', required: true },
    { name: 'lb_algorithm', label: '负载均衡算法', type: 'select', required: true, options: LB_ALGORITHM_OPTIONS },
    { name: 'health_check_path', label: '健康检查路径', type: 'text' },
    { name: 'health_check_interval', label: '健康检查间隔', type: 'number', min: 1, addonAfter: '秒' },
    { name: 'status', label: '状态', type: 'select', options: STATUS_OPTIONS },
  ]

  const columns: ColumnsType<UpstreamGroup> = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '分组名', dataIndex: 'name' },
    { title: '负载均衡算法', dataIndex: 'lb_algorithm', width: 130, render: (v: string) => LB_TEXT[v] ?? v },
    { title: '健康检查路径', dataIndex: 'health_check_path', render: (v: string | null) => v ?? '-' },
    { title: '检查间隔(秒)', dataIndex: 'health_check_interval', width: 110 },
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
            <Popconfirm title="确认删除该分组？" onConfirm={() => delMut.mutate(record.id)}>
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
        lb_algorithm: editing.lb_algorithm,
        health_check_path: editing.health_check_path ?? undefined,
        health_check_interval: editing.health_check_interval,
        status: editing.status,
      }
    : { lb_algorithm: 'round_robin', health_check_interval: 10, status: 1 }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        {canWrite && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); setModalOpen(true) }}>
            新增分组
          </Button>
        )}
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={listQuery.data?.list ?? []}
        loading={listQuery.isLoading}
        pagination={tablePagination(listQuery.data?.total ?? 0)}
        expandable={{
          expandedRowRender: (record) => <NodesTable groupId={record.id} canWrite={canWrite} />,
          rowExpandable: () => true,
        }}
      />
      <FormModal
        open={modalOpen}
        title={editing ? '编辑分组' : '新增分组'}
        fields={fields}
        initialValues={initialValues}
        confirmLoading={saveMut.isPending}
        onCancel={() => setModalOpen(false)}
        onSubmit={(values) => saveMut.mutate(values as Partial<UpstreamGroup>)}
      />
    </div>
  )
}
