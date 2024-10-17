<script setup lang="ts">
import { NButton, NDataTable, NPopconfirm, NTag } from 'naive-ui'

import firewall from '@/api/panel/firewall'
import { renderIcon } from '@/utils'
import CreateForwardModal from '@/views/safe/CreateForwardModal.vue'
import type { FirewallRule } from '@/views/safe/types'

const createModalShow = ref(false)

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: '传输协议',
    key: 'protocol',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any): any {
      return h(NTag, null, {
        default: () => {
          if (row.protocol !== '') {
            return row.protocol
          }
          return '无'
        }
      })
    }
  },
  {
    title: '端口',
    key: 'port',
    width: 150,
    render(row: any): any {
      return h(NTag, null, {
        default: () => {
          return row.port
        }
      })
    }
  },
  {
    title: '目标 IP',
    key: 'target_ip',
    minWidth: 200,
    render(row: any): any {
      return h(
        NTag,
        {
          type: 'info'
        },
        {
          default: () => {
            return row.target_ip
          }
        }
      )
    }
  },
  {
    title: '目标端口',
    key: 'target_port',
    width: 150,
    render(row: any): any {
      return h(
        NTag,
        {
          type: 'info'
        },
        {
          default: () => {
            return row.target_port
          }
        }
      )
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
            onPositiveClick: () => handleDelete(row)
          },
          {
            default: () => {
              return '确定要删除吗？'
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

const data = ref<FirewallRule[]>([] as FirewallRule[])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200]
})

const selectedRowKeys = ref<any>([])

const handleDelete = async (row: any) => {
  await firewall.deleteForward(row).then(() => {
    window.$message.success('删除成功')
  })
  fetchFirewallForwards(pagination.page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const fetchFirewallForwards = async (page: number, limit: number) => {
  const { data } = await firewall.forwards(page, limit)
  return data
}

const batchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要删除的规则')
    return
  }

  for (const key of selectedRowKeys.value) {
    // 解析json
    const rule = JSON.parse(key)
    await firewall.deleteForward(rule).then(() => {
      window.$message.success(`${rule.protocol} ${rule.target_ip}:${rule.target_port} 删除成功`)
    })
  }

  fetchFirewallForwards(pagination.page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const onPageChange = (page: number) => {
  pagination.page = page
  fetchFirewallForwards(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

watch(createModalShow, () => {
  onPageChange(1)
})

onMounted(() => {
  onPageChange(1)
})
</script>

<template>
  <n-flex vertical>
    <n-card flex-1 rounded-10>
      <n-flex items-center>
        <n-button type="primary" @click="createModalShow = true">
          <TheIcon :size="18" icon="material-symbols:add" />
          创建转发
        </n-button>
        <n-popconfirm @positive-click="batchDelete">
          <template #trigger>
            <n-button type="warning">
              <TheIcon :size="18" icon="material-symbols:delete-outline" />
              批量删除
            </n-button>
          </template>
          确定要批量删除吗？
        </n-popconfirm>
      </n-flex>
    </n-card>
    <n-data-table
      striped
      remote
      :scroll-x="1000"
      :loading="false"
      :columns="columns"
      :data="data"
      :row-key="(row: any) => JSON.stringify(row)"
      :pagination="pagination"
      @update:checked-row-keys="onChecked"
      @update:page="onPageChange"
      @update:page-size="onPageSizeChange"
    />
  </n-flex>
  <create-forward-modal v-model:show="createModalShow" />
</template>

<style scoped lang="scss"></style>
