import { useState } from 'react'
import { Button, Descriptions, Drawer, Popconfirm, Select, Space, Table, message } from 'antd'
import { useMutation, useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import type { ColumnsType } from 'antd/es/table'
import { deleteDevice, listDevices, updateDeviceStatus } from '@/api/device'
import { createBan } from '@/api/ban'
import type { Device, DeviceQuery } from '@/types'
import FormModal, { type FieldConfig } from '@/components/FormModal'
import StatusTag from '@/components/StatusTag'
import { DEVICE_STATUS_MAP, DEVICE_STATUS_OPTIONS } from '@/constants/options'
import { usePage } from '@/hooks/usePage'
import { usePermission } from '@/hooks/usePermission'

export default function DevicesPage() {
  const { canWrite } = usePermission()
  const { query, tablePagination } = usePage()
  const [status, setStatus] = useState<string | undefined>()
  const [detail, setDetail] = useState<Device | null>(null)
  const [banTarget, setBanTarget] = useState<Device | null>(null)

  const filters: DeviceQuery = { status }
  const listQuery = useQuery({
    queryKey: ['devices', { ...query, ...filters }],
    queryFn: () => listDevices({ ...query, ...filters }),
  })

  const banMut = useMutation({
    mutationFn: async (reason?: string) => {
      if (!banTarget) return
      await createBan({
        device_id: banTarget.id,
        ban_type: 'device',
        target_value: banTarget.fingerprint_id,
        reason,
      })
      await updateDeviceStatus(banTarget.id, 'banned')
    },
    onSuccess: () => {
      message.success('已封禁')
      setBanTarget(null)
      listQuery.refetch()
    },
  })

  const unbanMut = useMutation({
    mutationFn: (id: number) => updateDeviceStatus(id, 'normal'),
    onSuccess: () => {
      message.success('已解封')
      listQuery.refetch()
    },
  })

  const delMut = useMutation({
    mutationFn: (id: number) => deleteDevice(id),
    onSuccess: () => {
      message.success('删除成功')
      listQuery.refetch()
    },
  })

  const banFields: FieldConfig[] = [{ name: 'reason', label: '封禁原因', type: 'text' }]

  const columns: ColumnsType<Device> = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '设备指纹', dataIndex: 'fingerprint_id', ellipsis: true },
    { title: 'User-Agent', dataIndex: 'user_agent', ellipsis: true, render: (v: string | null) => v ?? '-' },
    { title: '最近IP', dataIndex: 'last_seen_ip', width: 130, render: (v: string | null) => v ?? '-' },
    {
      title: '状态',
      dataIndex: 'device_status',
      width: 100,
      render: (v: string) => <StatusTag value={v} map={DEVICE_STATUS_MAP} />,
    },
    {
      title: '首次出现',
      dataIndex: 'first_seen_at',
      width: 160,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '最近出现',
      dataIndex: 'last_seen_at',
      width: 160,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) =>
        canWrite ? (
          <Space>
            <Button type="link" size="small" onClick={() => setDetail(record)}>
              详情
            </Button>
            {record.device_status !== 'banned' ? (
              <Button type="link" size="small" danger onClick={() => setBanTarget(record)}>
                封禁
              </Button>
            ) : (
              <Popconfirm title="确认解封该设备？" onConfirm={() => unbanMut.mutate(record.id)}>
                <Button type="link" size="small">
                  解封
                </Button>
              </Popconfirm>
            )}
            <Popconfirm title="确认删除该设备？" onConfirm={() => delMut.mutate(record.id)}>
              <Button type="link" size="small" danger>
                删除
              </Button>
            </Popconfirm>
          </Space>
        ) : (
          <Button type="link" size="small" onClick={() => setDetail(record)}>
            详情
          </Button>
        ),
    },
  ]

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Select
          placeholder="状态"
          allowClear
          style={{ width: 140 }}
          options={DEVICE_STATUS_OPTIONS}
          onChange={(v) => setStatus(v as string | undefined)}
        />
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={listQuery.data?.list ?? []}
        loading={listQuery.isLoading}
        pagination={tablePagination(listQuery.data?.total ?? 0)}
      />

      <Drawer title="设备详情" open={!!detail} onClose={() => setDetail(null)} width={520}>
        {detail && (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="ID">{detail.id}</Descriptions.Item>
            <Descriptions.Item label="综合指纹">{detail.fingerprint_id}</Descriptions.Item>
            <Descriptions.Item label="JA4">{detail.ja4_hash ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="JA3">{detail.ja3_hash ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="OS指纹">{detail.os_fingerprint ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="Header哈希">{detail.http_header_hash ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="User-Agent">{detail.user_agent ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="最近IP">{detail.last_seen_ip ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="状态">
              <StatusTag value={detail.device_status} map={DEVICE_STATUS_MAP} />
            </Descriptions.Item>
            <Descriptions.Item label="封禁原因">{detail.ban_reason ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="首次出现">
              {dayjs(detail.first_seen_at).format('YYYY-MM-DD HH:mm:ss')}
            </Descriptions.Item>
            <Descriptions.Item label="最近出现">
              {dayjs(detail.last_seen_at).format('YYYY-MM-DD HH:mm:ss')}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>

      <FormModal
        open={!!banTarget}
        title="封禁设备"
        fields={banFields}
        confirmLoading={banMut.isPending}
        onCancel={() => setBanTarget(null)}
        onSubmit={(values) => banMut.mutate(values.reason as string | undefined)}
      />
    </div>
  )
}
