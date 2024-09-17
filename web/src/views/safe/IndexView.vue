<script setup lang="ts">
import { NButton, NDataTable, NPopconfirm, NSpace } from 'naive-ui'
import type { FirewallRule } from '@/views/safe/types'
import safe from '@/api/panel/safe'
import { renderIcon } from '@/utils'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const model = ref({
  firewallStatus: false,
  sshStatus: false,
  pingStatus: false,
  sshPort: 22
})

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: t('safeIndex.columns.port'),
    key: 'port',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('safeIndex.columns.protocol'),
    key: 'protocol',
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('safeIndex.columns.actions'),
    key: 'actions',
    width: 140,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row),
            onNegativeClick: () => {
              window.$message.info(t('safeIndex.alerts.undelete'))
            }
          },
          {
            default: () => {
              return t('safeIndex.confirm.delete')
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
                  default: () => t('safeIndex.buttons.delete'),
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
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
})

const selectedRowKeys = ref<any>([])

const addModel = ref({
  port: '',
  protocol: 'tcp'
})

const handleDelete = async (row: any) => {
  await safe.deleteFirewallRule(row.port, row.protocol).then(() => {
    window.$message.success(t('safeIndex.alerts.delete'))
  })
  getFirewallRules(pagination.page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const handleAdd = async () => {
  await safe.addFirewallRule(addModel.value.port, addModel.value.protocol).then(() => {
    window.$message.success(t('safeIndex.alerts.add'))
    addModel.value.port = ''
    addModel.value.protocol = 'tcp'
  })
  getFirewallRules(pagination.page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const getFirewallRules = async (page: number, limit: number) => {
  const { data } = await safe.firewallRules(page, limit)
  return data
}

const getSetting = async () => {
  safe.firewallStatus().then((res) => {
    model.value.firewallStatus = res.data
  })
  safe.sshStatus().then((res) => {
    model.value.sshStatus = res.data
  })
  safe.pingStatus().then((res) => {
    model.value.pingStatus = res.data
  })
  safe.sshPort().then((res) => {
    model.value.sshPort = res.data
  })
}

const handleFirewallStatus = () => {
  safe.setFirewallStatus(model.value.firewallStatus).then(() => {
    window.$message.success(t('safeIndex.alerts.setup'))
  })
}

const handleSshStatus = () => {
  safe.setSshStatus(model.value.sshStatus).then(() => {
    window.$message.success(t('safeIndex.alerts.setup'))
  })
}

const handlePingStatus = () => {
  safe.setPingStatus(model.value.pingStatus).then(() => {
    window.$message.success(t('safeIndex.alerts.setup'))
  })
}

const handleSshPort = () => {
  safe.setSshPort(model.value.sshPort).then(() => {
    window.$message.success(t('safeIndex.alerts.setup'))
  })
}

const batchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info(t('safeIndex.alerts.select'))
    return
  }

  for (const key of selectedRowKeys.value) {
    // 通过 / 分割端口和协议
    const [port, protocol] = key.split('/')
    if (!port || !protocol) {
      continue
    }

    await safe.deleteFirewallRule(port, protocol).then(() => {
      let rule = data.value.find((item) => item.port === port && item.protocol === protocol)
      window.$message.success(
        t('safeIndex.alerts.ruleDelete', { rule: rule?.port + '/' + rule?.protocol })
      )
    })
  }

  getFirewallRules(pagination.page, pagination.pageSize).then((res) => {
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
  getFirewallRules(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

onMounted(() => {
  getSetting()
  getFirewallRules(pagination.page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
})
</script>

<template>
  <common-page show-footer>
    <n-space vertical>
      <n-card flex-1 rounded-10>
        <n-form inline>
          <n-form-item :label="$t('safeIndex.filter.fields.firewall.label')">
            <n-switch
              v-model:value="model.firewallStatus"
              @update:value="handleFirewallStatus"
              :checkedChildren="$t('safeIndex.filter.fields.firewall.checked')"
              :unCheckedChildren="$t('safeIndex.filter.fields.firewall.unchecked')"
            />
          </n-form-item>
          <n-form-item :label="$t('safeIndex.filter.fields.ssh.label')">
            <n-switch
              v-model:value="model.sshStatus"
              @update:value="handleSshStatus"
              :checkedChildren="$t('safeIndex.filter.fields.ssh.checked')"
              :unCheckedChildren="$t('safeIndex.filter.fields.ssh.unchecked')"
            />
          </n-form-item>
          <n-form-item :label="$t('safeIndex.filter.fields.ping.label')">
            <n-switch
              v-model:value="model.pingStatus"
              @update:value="handlePingStatus"
              :checkedChildren="$t('safeIndex.filter.fields.ping.checked')"
              :unCheckedChildren="$t('safeIndex.filter.fields.ping.unchecked')"
            />
          </n-form-item>
          <n-form-item :label="$t('safeIndex.filter.fields.port.label')">
            <n-input-number v-model:value="model.sshPort" @blur="handleSshPort" />
          </n-form-item>
        </n-form>
      </n-card>
      <n-space flex items-center>
        <n-popconfirm @positive-click="batchDelete">
          <template #trigger>
            <n-button type="warning"> {{ $t('safeIndex.buttons.batchDelete') }} </n-button>
          </template>
          {{ $t('safeIndex.confirm.batchDelete') }}
        </n-popconfirm>
        <n-text>{{ $t('safeIndex.portControl.title') }}</n-text>
        <n-input
          v-model:value="addModel.port"
          :placeholder="$t('safeIndex.portControl.fields.port.placeholder')"
        />
        <n-select
          v-model:value="addModel.protocol"
          :placeholder="$t('safeIndex.portControl.fields.protocol.placeholder')"
          style="width: 120px"
          :options="[
            { label: 'TCP', value: 'tcp' },
            { label: 'UDP', value: 'udp' }
          ]"
        />
        <n-button type="primary" @click="handleAdd">
          {{ $t('safeIndex.buttons.add') }}
        </n-button>
      </n-space>

      <n-data-table
        striped
        remote
        :loading="false"
        :scroll-x="1200"
        :columns="columns"
        :data="data"
        :row-key="(row: any) => row.port + '/' + row.protocol"
        :pagination="pagination"
        @update:checked-row-keys="onChecked"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </n-space>
  </common-page>
</template>

<style scoped lang="scss"></style>
