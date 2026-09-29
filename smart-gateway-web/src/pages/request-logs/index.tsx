import { useState } from 'react'
import { Button, Form, Input, InputNumber, Space, Table, Tag } from 'antd'
import { ReloadOutlined, SearchOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import { listRequestLogs } from '@/api/log'
import type { RequestLog, RequestLogQuery } from '@/types'
import { usePage } from '@/hooks/usePage'

function statusColor(code: number) {
  if (code >= 200 && code < 300) return 'green'
  if (code >= 300 && code < 400) return 'orange'
  if (code >= 400) return 'red'
  return 'default'
}

export default function RequestLogsPage() {
  const { query, tablePagination } = usePage(20)
  const [filters, setFilters] = useState<RequestLogQuery>({})
  const [form] = Form.useForm()

  const listQuery = useQuery({
    queryKey: ['request-logs', { ...query, ...filters }],
    queryFn: () => listRequestLogs({ ...query, ...filters }),
  })

  const onSearch = (values: Record<string, unknown>) => {
    setFilters({
      source_ip: (values.source_ip as string) || undefined,
      device_fingerprint_id: (values.device_fingerprint_id as string) || undefined,
      upstream_id: (values.upstream_id as number) || undefined,
      route_id: (values.route_id as number) || undefined,
    })
  }

  const onReset = () => {
    form.resetFields()
    setFilters({})
  }

  const columns: ColumnsType<RequestLog> = [
    {
      title: '时间',
      dataIndex: 'received_at',
      width: 160,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm:ss'),
    },
    { title: '方法', dataIndex: 'method', width: 80 },
    { title: '路径', dataIndex: 'path', ellipsis: true },
    { title: '来源IP', dataIndex: 'source_ip', width: 130 },
    { title: '上游节点', dataIndex: 'upstream_id', width: 90, render: (v: number) => (v ? v : '-') },
    { title: '路由', dataIndex: 'route_id', width: 80, render: (v: number) => (v ? v : '-') },
    {
      title: '状态码',
      dataIndex: 'status_code',
      width: 90,
      render: (v: number) => <Tag color={statusColor(v)}>{v}</Tag>,
    },
    { title: '耗时(ms)', dataIndex: 'latency_ms', width: 90 },
    { title: '大小', dataIndex: 'response_size', width: 90 },
    { title: '设备指纹', dataIndex: 'device_fingerprint_id', ellipsis: true, render: (v: string) => v || '-' },
  ]

  return (
    <div>
      <Form form={form} layout="inline" style={{ marginBottom: 16 }} onFinish={onSearch}>
        <Form.Item name="source_ip" label="来源IP">
          <Input placeholder="来源IP" allowClear style={{ width: 140 }} />
        </Form.Item>
        <Form.Item name="device_fingerprint_id" label="设备指纹">
          <Input placeholder="设备指纹" allowClear style={{ width: 160 }} />
        </Form.Item>
        <Form.Item name="upstream_id" label="上游节点ID">
          <InputNumber placeholder="节点ID" min={1} style={{ width: 120 }} />
        </Form.Item>
        <Form.Item name="route_id" label="路由ID">
          <InputNumber placeholder="路由ID" min={1} style={{ width: 120 }} />
        </Form.Item>
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
              查询
            </Button>
            <Button icon={<ReloadOutlined />} onClick={onReset}>
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={listQuery.data?.list ?? []}
        loading={listQuery.isLoading}
        pagination={tablePagination(listQuery.data?.total ?? 0)}
      />
    </div>
  )
}
