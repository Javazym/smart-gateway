import { useState } from 'react'
import { Button, Form, Input, Space, Table } from 'antd'
import { ReloadOutlined, SearchOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import { listAuditLogs } from '@/api/log'
import type { AuditLog, AuditLogQuery } from '@/types'
import { usePage } from '@/hooks/usePage'

export default function AuditLogsPage() {
  const { query, tablePagination } = usePage(20)
  const [filters, setFilters] = useState<AuditLogQuery>({})
  const [form] = Form.useForm()

  const listQuery = useQuery({
    queryKey: ['audit-logs', { ...query, ...filters }],
    queryFn: () => listAuditLogs({ ...query, ...filters }),
  })

  const onSearch = (values: Record<string, unknown>) => {
    setFilters({
      username: (values.username as string) || undefined,
      action: (values.action as string) || undefined,
    })
  }

  const onReset = () => {
    form.resetFields()
    setFilters({})
  }

  const columns: ColumnsType<AuditLog> = [
    {
      title: '时间',
      dataIndex: 'created_at',
      width: 160,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm:ss'),
    },
    { title: '操作人', dataIndex: 'username', width: 120, render: (v: string | null) => v ?? '-' },
    { title: '操作', dataIndex: 'action', width: 160 },
    { title: '目标类型', dataIndex: 'target_type', width: 120, render: (v: string | null) => v ?? '-' },
    { title: '目标ID', dataIndex: 'target_id', width: 120, render: (v: string | null) => v ?? '-' },
    { title: '详情', dataIndex: 'detail', ellipsis: true, render: (v: string | null) => v ?? '-' },
    { title: '来源IP', dataIndex: 'source_ip', width: 130, render: (v: string | null) => v ?? '-' },
  ]

  return (
    <div>
      <Form form={form} layout="inline" style={{ marginBottom: 16 }} onFinish={onSearch}>
        <Form.Item name="username" label="操作人">
          <Input placeholder="操作人" allowClear style={{ width: 140 }} />
        </Form.Item>
        <Form.Item name="action" label="操作">
          <Input placeholder="操作类型" allowClear style={{ width: 160 }} />
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
