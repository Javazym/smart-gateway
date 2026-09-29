import { useState } from 'react'
import type { TablePaginationConfig } from 'antd'

// 分页状态与 Antd Table 分页配置生成
export function usePage(defaultSize = 10) {
  const [page, setPage] = useState(1)
  const [size, setSize] = useState(defaultSize)

  return {
    page,
    size,
    query: { page, size },
    tablePagination: (total: number): TablePaginationConfig => ({
      current: page,
      pageSize: size,
      total,
      showSizeChanger: true,
      showTotal: (t) => `共 ${t} 条`,
      onChange: (p, s) => {
        setPage(p)
        setSize(s)
      },
    }),
  }
}
