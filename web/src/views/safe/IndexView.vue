<script setup lang="ts">
import { NButton, NDataTable, NPopconfirm, NSpace, NTag } from 'naive-ui'

import safe from '@/api/panel/safe'
import { renderIcon } from '@/utils'
import CreateModal from '@/views/safe/CreateModal.vue'
import type { FirewallRule } from '@/views/safe/types'

const createModalShow = ref(false)

const model = ref({
  firewallStatus: false,
  sshStatus: false,
  pingStatus: false,
  sshPort: 22
})

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
    title: '网络协议',
    key: 'family',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any): any {
      return h(NTag, null, {
        default: () => {
          if (row.family !== '') {
            return row.family
          }
          return '无'
        }
      })
    }
  },
  {
    title: '端口',
    key: 'port',
    minWidth: 300,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any): any {
      if (row.port_start == row.port_end) {
        return row.port_start
      }
      return `${row.port_start}-${row.port_end}`
    }
  },
  {
    title: '策略',
    key: 'strategy',
    minWidth: 100,
    render(row: any): any {
      return h(
        NTag,
        {
          type:
            row.strategy === 'accept' ? 'success' : row.strategy === 'drop' ? 'warning' : 'error'
        },
        {
          default: () => {
            switch (row.strategy) {
              case 'accept':
                return '接受'
              case 'drop':
                return '丢弃'
              case 'reject':
                return '拒绝'
              default:
                return '未知'
            }
          }
        }
      )
    }
  },
  {
    title: '方向',
    key: 'direction',
    minWidth: 100,
    render(row: any): any {
      return h(NTag, null, {
        default: () => {
          switch (row.direction) {
            case 'in':
              return '传入'
            case 'out':
              return '传出'
            default:
              return '未知'
          }
        }
      })
    }
  },
  {
    title: '目标',
    key: 'address',
    minWidth: 100,
    render(row: any): any {
      return h(NTag, null, {
        default: () => {
          if (row.address === '') {
            return '所有'
          }
          return row.address
        }
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
  await safe.deleteFirewallRule(row).then(() => {
    window.$message.success('删除成功')
  })
  fetchFirewallRules(pagination.page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const fetchFirewallRules = async (page: number, limit: number) => {
  const { data } = await safe.firewallRules(page, limit)
  return data
}

const fetchSetting = async () => {
  safe.firewallStatus().then((res) => {
    model.value.firewallStatus = res.data
  })
  safe.ssh().then((res) => {
    model.value.sshStatus = res.data.status
    model.value.sshPort = res.data.port
  })
  safe.pingStatus().then((res) => {
    model.value.pingStatus = res.data
  })
}

const handleFirewallStatus = () => {
  safe.setFirewallStatus(model.value.firewallStatus).then(() => {
    window.$message.success('设置成功')
  })
}

const handleSsh = () => {
  safe.setSsh(model.value.sshStatus, model.value.sshPort).then(() => {
    window.$message.success('设置成功')
  })
}

const handlePingStatus = () => {
  safe.setPingStatus(model.value.pingStatus).then(() => {
    window.$message.success('设置成功')
  })
}

const batchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要删除的规则')
    return
  }

  for (const key of selectedRowKeys.value) {
    // 解析json
    const rule = JSON.parse(key)
    await safe.deleteFirewallRule(rule).then(() => {
      let port =
        rule.port_start == rule.port_end ? rule.port_start : `${rule.port_start}-${rule.port_end}`
      window.$message.success(`${rule.family} 规则 ${port}/${rule.protocol} 删除成功`)
    })
  }

  fetchFirewallRules(pagination.page, pagination.pageSize).then((res) => {
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
  fetchFirewallRules(page, pagination.pageSize).then((res) => {
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
  fetchSetting()
  onPageChange(1)
})
</script>

<template>
  <common-page show-footer>
    <n-space vertical>
      <n-card flex-1 rounded-10>
        <n-form inline>
          <n-form-item label="防火墙">
            <n-switch v-model:value="model.firewallStatus" @update:value="handleFirewallStatus" />
          </n-form-item>
          <n-form-item label="SSH">
            <n-switch v-model:value="model.sshStatus" @update:value="handleSsh" />
          </n-form-item>
          <n-form-item label="Ping">
            <n-switch v-model:value="model.pingStatus" @update:value="handlePingStatus" />
          </n-form-item>
          <n-form-item label="SSH端口">
            <n-input-number v-model:value="model.sshPort" @blur="handleSsh" />
          </n-form-item>
        </n-form>
      </n-card>
      <n-space flex items-center>
        <n-button type="primary" @click="createModalShow = true">
          <TheIcon :size="18" icon="material-symbols:add" />
          创建规则
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
      </n-space>

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
    </n-space>
  </common-page>
  <create-modal v-model:show="createModalShow" />
</template>

<style scoped lang="scss"></style>
