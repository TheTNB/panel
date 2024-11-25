<script setup lang="ts">
import { renderIcon } from '@/utils'
import { NButton, NInput, NInputGroup, NPopconfirm, NTag } from 'naive-ui'

import database from '@/api/panel/database'
import { formatDateTime } from '@/utils'

const columns: any = [
  {
    title: '用户名',
    key: 'username',
    minWidth: 100,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return row.username || '无'
    }
  },
  {
    title: '密码',
    key: 'password',
    width: 200,
    render(row: any) {
      return h(NInputGroup, null, {
        default: () => [
          h(NInput, {
            value: row.password,
            type: 'password',
            showPasswordOn: 'click',
            readonly: true
          })
        ]
      })
    }
  },
  {
    title: '主机',
    key: 'host',
    width: 200,
    render(row: any) {
      return h(NTag, null, {
        default: () => row.host
      })
    }
  },
  {
    title: '备注',
    key: 'remark',
    minWidth: 250,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(NInput, {
        size: 'small',
        value: row.remark
        /*onBlur: () => handleRemark(row),
        onUpdateValue(v) {
          row.remark = v
        }*/
      })
    }
  },
  {
    title: '更新日期',
    key: 'updated_at',
    width: 200,
    ellipsis: { tooltip: true },
    render(row: any) {
      return formatDateTime(row.updated_at)
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
            onPositiveClick: () => handleDelete(row.id)
          },
          {
            default: () => {
              return '确定删除服务器吗？'
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
  (page, pageSize) => database.userList(page, pageSize),
  {
    initialData: { total: 0, list: [] },
    total: (res: any) => res.total,
    data: (res: any) => res.items
  }
)

const handleDelete = async (id: number) => {
  await database.userDelete(id).then(() => {
    window.$message.success('删除成功')
    refresh()
  })
}

onMounted(() => {
  window.$bus.on('database-user:refresh', () => {
    refresh()
  })
})

onUnmounted(() => {
  window.$bus.off('database-user:refresh')
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
