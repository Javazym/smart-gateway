import { Tag } from 'antd'
import type { TagMeta } from '@/constants/options'

export default function StatusTag({
  value,
  map,
}: {
  value: string | number | null | undefined
  map: Record<string, TagMeta>
}) {
  if (value === null || value === undefined || value === '') {
    return <Tag>-</Tag>
  }
  const item = map[String(value)]
  if (!item) {
    return <Tag>{String(value)}</Tag>
  }
  return <Tag color={item.color}>{item.text}</Tag>
}
