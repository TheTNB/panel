<script setup lang="ts">
import { renderIcon } from '@/utils'
import { NButton, NPopconfirm, NTag } from 'naive-ui'

import database from '@/api/panel/database'

const columns: any = [
  {
    title: '数据库名',
    key: 'name',
    minWidth: 100,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '服务器',
    key: 'server',
    width: 300
  },
  {
    title: '编码',
    key: 'encoding',
    width: 200,
    render(row: any) {
      return h(NTag, null, {
        default: () => row.encoding
      })
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    align: 'center',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row.server_id, row.name)
          },
          {
            default: () => {
              return '确定删除数据库吗？'
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  style: 'margin-left: 15px;'
                },
                {
                  default: () => '删除',
                  icon: renderIcon('material-symbols:delete-outline', { size: 14 })
                }
              )
            }
          }
        )
      ]
    }
  }
]

const { loading, data, page, total, pageSize, pageCount, refresh } = usePagination(
  (page, pageSize) => database.list(page, pageSize),
  {
    initialData: { total: 0, list: [] },
    total: (res: any) => res.total,
    data: (res: any) => res.items
  }
)

const handleDelete = async (serverID: number, name: string) => {
  await database.delete(serverID, name).then(() => {
    window.$message.success('删除成功')
    refresh()
  })
}

onMounted(() => {
  window.$bus.on('database:refresh', () => {
    refresh()
  })
})

onUnmounted(() => {
  window.$bus.off('database:refresh')
})
</script>

<template>
  <n-data-table
    striped
    remote
    :scroll-x="1200"
    :loading="loading"
    :columns="columns"
    :data="data"
    :row-key="(row: any) => row.name"
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
</template>

<style scoped lang="scss"></style>
