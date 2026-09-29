import ReactECharts from 'echarts-for-react'

interface LineSeries {
  name: string
  data: number[]
}

interface LineChartProps {
  xAxis: string[]
  series: LineSeries[]
  loading?: boolean
}

export default function LineChart({ xAxis, series, loading }: LineChartProps) {
  const option = {
    tooltip: { trigger: 'axis' },
    legend: { data: series.map((s) => s.name), top: 0 },
    grid: { left: 48, right: 24, top: 48, bottom: 40 },
    xAxis: {
      type: 'category' as const,
      boundaryGap: false,
      data: xAxis,
    },
    yAxis: { type: 'value' as const, name: '请求数', minInterval: 1 },
    series: series.map((s) => ({
      name: s.name,
      type: 'line' as const,
      smooth: true,
      showSymbol: false,
      data: s.data,
    })),
  }

  return (
    <ReactECharts
      option={option}
      notMerge
      style={{ height: 420 }}
      showLoading={loading}
    />
  )
}
