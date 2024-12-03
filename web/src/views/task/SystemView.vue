<script setup lang="ts">
import { NButton, NDataTable, NPopconfirm, NTag } from 'naive-ui'

import process from '@/api/panel/process'
import { formatBytes, formatDateTime, formatPercent, renderIcon } from '@/utils'

const columns: any = [
  {
    title: 'PID',
    key: 'pid',
    width: 120,
    ellipsis: { tooltip: true }
  },
  {
    title: '名称',
    key: 'name',
    minWidth: 250,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '父进程 ID',
    key: 'ppid',
    width: 120,
    ellipsis: { tooltip: true }
  },
  {
    title: '线程数',
    key: 'num_threads',
    width: 100,
    ellipsis: { tooltip: true }
  },
  {
    title: '用户',
    key: 'username',
    minWidth: 100,
    ellipsis: { tooltip: true }
  },
  {
    title: '状态',
    key: 'status',
    minWidth: 150,
    ellipsis: { tooltip: true },
    render(row: any) {
      switch (row.status) {
        case 'R':
          return h(NTag, { type: 'success' }, { default: () => '运行' })
        case 'S':
          return h(NTag, { type: 'warning' }, { default: () => '睡眠' })
        case 'T':
          return h(NTag, { type: 'error' }, { default: () => '停止' })
        case 'I':
          return h(NTag, { type: 'primary' }, { default: () => '空闲' })
        case 'Z':
          return h(NTag, { type: 'error' }, { default: () => '僵尸' })
        case 'W':
          return h(NTag, { type: 'warning' }, { default: () => '等待' })
        case 'L':
          return h(NTag, { type: 'info' }, { default: () => '锁定' })
        default:
          return h(NTag, { type: 'default' }, { default: () => row.status })
      }
    }
  },
  {
    title: 'CPU',
    key: 'cpu',
    minWidth: 100,
    ellipsis: { tooltip: true },
    render(row: any): string {
      return formatPercent(row.cpu) + '%'
    }
  },
  {
    title: '内存',
    key: 'rss',
    minWidth: 100,
    ellipsis: { tooltip: true },
    render(row: any): string {
      return formatBytes(row.rss)
    }
  },
  {
    title: '启动时间',
    key: 'start_time',
    width: 160,
    ellipsis: { tooltip: true },
    render(row: any): string {
      return formatDateTime(row.start_time)
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    align: 'center',
    hideInExcel: true,
    render(row: any) {
      return h(
        NPopconfirm,
        {
          onPositiveClick: async () => {
            await process.kill(row.pid)
            await refresh()
            window.$message.success(`进程 ${row.pid} 已终止`)
          }
        },
        {
          default: () => {
            return '确定终止进程 ' + row.pid + ' ?'
          },
          trigger: () => {
            return h(
              NButton,
              {
                size: 'small',
                type: 'error'
              },
              {
                default: () => '终止',
                icon: renderIcon('material-symbols:stop-circle-outline-rounded', { size: 14 })
              }
            )
          }
        }
      )
    }
  }
]

const { loading, data, page, total, pageSize, pageCount, refresh } = usePagination(
  (page, pageSize) => process.list(page, pageSize),
  {
    initialData: { total: 0, list: [] },
    total: (res: any) => res.total,
    data: (res: any) => res.items
  }
)
</script>

<template>
  <n-flex vertical>
    <n-data-table
      striped
      remote
      :scroll-x="1400"
      :loading="loading"
      :columns="columns"
      :data="data"
      :row-key="(row: any) => row.pid"
      v-model:page="page"
      v-model:pageSize="pageSize"
      :pagination="{
        page: page,
        pageCount: pageCount,
        pageSize: pageSize,
        itemCount: total,
        showQuickJumper: true,
        showSizePicker: true,
        pageSizes: [20, 50, 100, 200]
      }"
    />
  </n-flex>
</template>
