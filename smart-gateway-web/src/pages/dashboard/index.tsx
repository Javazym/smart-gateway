import { useMemo, useState } from 'react'
import { Card, Col, Empty, Row, Select, Space, Statistic, Typography } from 'antd'
import { useQuery } from '@tanstack/react-query'
import dayjs from 'dayjs'
import { listRequestLogs } from '@/api/log'
import type { RequestLog } from '@/types'
import LineChart from '@/components/LineChart'

type DimKey = 'source_ip' | 'device_fingerprint_id' | 'upstream_id'

const DIM_OPTIONS = [
  { label: '按来源IP', value: 'source_ip' },
  { label: '按设备指纹', value: 'device_fingerprint_id' },
  { label: '按上游节点', value: 'upstream_id' },
]

function aggregate(logs: RequestLog[], dim: DimKey) {
  const sorted = [...logs].sort(
    (a, b) => new Date(a.received_at).getTime() - new Date(b.received_at).getTime(),
  )
  if (sorted.length === 0) return { xAxis: [] as string[], series: [] as { name: string; data: number[] }[] }

  const labelOf = (l: RequestLog) => dayjs(l.received_at).format('HH:mm')
  const labels: string[] = []
  const labelIdx = new Map<string, number>()
  sorted.forEach((l) => {
    const lb = labelOf(l)
    if (!labelIdx.has(lb)) {
      labelIdx.set(lb, labels.length)
      labels.push(lb)
    }
  })

  const dimValue = (l: RequestLog): string => {
    if (dim === 'upstream_id') return l.upstream_id ? `节点 #${l.upstream_id}` : '未知'
    if (dim === 'source_ip') return l.source_ip || '未知'
    return l.device_fingerprint_id || '未知'
  }

  const counts = new Map<string, number[]>()
  sorted.forEach((l) => {
    const key = dimValue(l)
    const arr = counts.get(key) ?? new Array<number>(labels.length).fill(0)
    arr[labelIdx.get(labelOf(l))!] += 1
    counts.set(key, arr)
  })

  const ranked = [...counts.entries()]
    .map(([name, data]) => ({ name, data, total: data.reduce((s, v) => s + v, 0) }))
    .sort((a, b) => b.total - a.total)

  const series = ranked.slice(0, 6).map(({ name, data }) => ({ name, data }))
  if (ranked.length > 6) {
    const rest = new Array<number>(labels.length).fill(0)
    ranked.slice(6).forEach((s) => s.data.forEach((v, i) => (rest[i] += v)))
    series.push({ name: '其他', data: rest })
  }
  return { xAxis: labels, series }
}

function computeStats(logs: RequestLog[]) {
  const total = logs.length
  const ips = new Set(logs.map((l) => l.source_ip)).size
  const devices = new Set(logs.map((l) => l.device_fingerprint_id).filter(Boolean)).size
  const avgLatency = total ? Math.round(logs.reduce((s, l) => s + l.latency_ms, 0) / total) : 0
  return { total, ips, devices, avgLatency }
}

export default function DashboardPage() {
  const [dim, setDim] = useState<DimKey>('source_ip')

  const { data, isLoading, isRefetching } = useQuery({
    queryKey: ['dashboard-logs'],
    queryFn: () => listRequestLogs({ page: 1, size: 500 }),
    refetchInterval: 5000,
  })

  const logs = data?.list ?? []
  const chart = useMemo(() => aggregate(logs, dim), [logs, dim])
  const stats = useMemo(() => computeStats(logs), [logs])

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Typography.Text type="secondary">维度切换：</Typography.Text>
        <Select value={dim} options={DIM_OPTIONS} style={{ width: 160 }} onChange={(v) => setDim(v as DimKey)} />
        <Typography.Text type="secondary" style={{ fontSize: 12 }}>
          （每 5 秒轮询最近 500 条请求）
        </Typography.Text>
      </Space>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card size="small">
            <Statistic title="请求数(近500条)" value={stats.total} />
          </Card>
        </Col>
        <Col span={6}>
          <Card size="small">
            <Statistic title="独立来源IP" value={stats.ips} />
          </Card>
        </Col>
        <Col span={6}>
          <Card size="small">
            <Statistic title="独立设备指纹" value={stats.devices} />
          </Card>
        </Col>
        <Col span={6}>
          <Card size="small">
            <Statistic title="平均耗时(ms)" value={stats.avgLatency} />
          </Card>
        </Col>
      </Row>

      <Card title="请求量趋势">
        {logs.length === 0 ? (
          <Empty description="暂无请求数据" style={{ padding: 48 }} />
        ) : (
          <LineChart xAxis={chart.xAxis} series={chart.series} loading={isLoading || isRefetching} />
        )}
      </Card>
    </div>
  )
}
