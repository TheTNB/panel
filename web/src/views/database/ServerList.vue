<script setup lang="ts">
import { renderIcon } from '@/utils'
import { NButton, NInput, NInputGroup, NPopconfirm, NTag } from 'naive-ui'

import database from '@/api/panel/database'
import { formatDateTime } from '@/utils'
import UpdateServerModal from '@/views/database/UpdateServerModal.vue'

const updateModal = ref(false)
const updateID = ref(0)

const columns: any = [
  {
    title: '类型',
    key: 'type',
    width: 150,
    render(row: any) {
      return h(
        NTag,
        { type: 'info' },
        {
          default: () => {
            switch (row.type) {
              case 'mysql':
                return 'MySQL'
              case 'postgresql':
                return 'PostgreSQL'
              default:
                return row.type
            }
          }
        }
      )
    }
  },
  {
    title: '名称',
    key: 'name',
    minWidth: 100,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '用户名',
    key: 'username',
    width: 150,
    ellipsis: { tooltip: true },
    render(row: any) {
      return row.username || '无'
    }
  },
  {
    title: '密码',
    key: 'password',
    width: 250,
    render(row: any) {
      return h(NInputGroup, null, {
        default: () => [
          h(NInput, {
            value: row.password,
            type: 'password',
            showPasswordOn: 'click',
            readonly: true,
            placeholder: '无'
          }),
          h(
            NButton,
            {
              type: 'primary',
              ghost: true,
              onClick: () => {
                navigator.clipboard.writeText(row.password)
                window.$message.success('复制成功')
              }
            },
            { default: () => '复制' }
          )
        ]
      })
    }
  },
  {
    title: '主机',
    key: 'host',
    width: 150,
    render(row: any) {
      return h(NTag, null, {
        default: () => `${row.host}:${row.port}`
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
        value: row.remark,
        onBlur: () => handleRemark(row),
        onUpdateValue(v) {
          row.remark = v
        }
      })
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row: any) {
      return h(
        NTag,
        { type: row.status === 'valid' ? 'success' : 'error' },
        { default: () => (row.status === 'valid' ? '有效' : '无效') }
      )
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
    width: 300,
    align: 'center',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NPopconfirm,
          {
            onPositiveClick: async () => {
              await database.serverSync(row.id).then(() => {
                window.$message.success('同步成功')
                refresh()
              })
            }
          },
          {
            default: () => {
              return '确定同步数据库用户（不包括密码）到面板？'
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'success'
                },
                {
                  default: () => '同步',
                  icon: renderIcon('material-symbols:sync', { size: 14 })
                }
              )
            }
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 15px;',
            onClick: () => {
              updateID.value = row.id
              updateModal.value = true
            }
          },
          {
            default: () => '修改',
            icon: renderIcon('material-symbols:edit-outline', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => {
              // 防手贱
              if (['local_mysql', 'local_postgresql'].includes(row.name)) {
                window.$message.error('内置服务器不能删除，如需删除请卸载对应应用')
                return
              }
              handleDelete(row.id)
            }
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
  (page, pageSize) => database.serverList(page, pageSize),
  {
    initialData: { total: 0, list: [] },
    total: (res: any) => res.total,
    data: (res: any) => res.items
  }
)

const handleDelete = async (id: number) => {
  await database.serverDelete(id).then(() => {
    window.$message.success('删除成功')
    refresh()
  })
}

const handleRemark = (row: any) => {
  database.serverRemark(row.id, row.remark).then(() => {
    window.$message.success('修改成功')
  })
}

onMounted(() => {
  window.$bus.on('database-server:refresh', () => {
    refresh()
  })
})

onUnmounted(() => {
  window.$bus.off('database-server:refresh')
})
</script>

<template>
  <n-data-table
    striped
    remote
    :scroll-x="1700"
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
  <update-server-modal v-model:id="updateID" v-model:show="updateModal" />
</template>

<style scoped lang="scss"></style>
