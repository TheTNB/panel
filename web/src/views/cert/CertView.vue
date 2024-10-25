<script setup lang="ts">
import Editor from '@guolao/vue-monaco-editor'
import type { MessageReactive } from 'naive-ui'
import { NButton, NDataTable, NFlex, NPopconfirm, NSpace, NSwitch, NTable, NTag } from 'naive-ui'

import cert from '@/api/panel/cert'
import type { Cert } from '@/views/cert/types'

const props = defineProps({
  algorithms: Array<any>,
  websites: Array<any>,
  accounts: Array<any>,
  dns: Array<any>,
  caProviders: Array<any>
})

const { algorithms, websites, accounts, dns, caProviders } = toRefs(props)

let messageReactive: MessageReactive | null = null

const updateCertModel = ref<any>({
  domains: [],
  dns_id: 0,
  type: 'P256',
  account_id: null,
  website_id: 0,
  auto_renew: true
})
const updateCertModal = ref(false)
const updateCert = ref<any>()
const showModal = ref(false)
const showCertModel = ref<any>({
  cert: '',
  key: ''
})
const deployCertModal = ref(false)
const deployCertModel = ref<any>({
  id: 0,
  website_id: 0
})

const columns: any = [
  {
    title: '域名',
    key: 'domains',
    minWidth: 200,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(
        'span',
        {
          type: row.status == 'active' ? 'success' : 'error'
        },
        {
          default: () => row.domains.join(', ')
        }
      )
    }
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(
        NTag,
        {
          type: 'info',
          bordered: false
        },
        {
          default: () => {
            switch (row.type) {
              case 'P256':
                return 'EC 256'
              case 'P384':
                return 'EC 384'
              case '2048':
                return 'RSA 2048'
              case '4096':
                return 'RSA 4096'
              default:
                return '未知'
            }
          }
        }
      )
    }
  },
  {
    title: '关联账号',
    key: 'account_id',
    minWidth: 400,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(NFlex, null, {
        default: () => [
          h(NTag, null, { default: () => (row.account?.email == null ? '无' : row.account.email) }),
          h(NTag, null, {
            default: () =>
              caProviders?.value?.find((item: any) => item.value === row.account?.ca)?.label
          })
        ]
      })
    }
  },
  {
    title: '关联网站',
    key: 'website_id',
    minWidth: 150,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(
        NTag,
        {
          type: row.website == null ? 'error' : 'success',
          bordered: false
        },
        {
          default: () => (row.website?.name == null ? '无' : row.website.name)
        }
      )
    }
  },
  {
    title: '关联DNS',
    key: 'dns_id',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(
        NTag,
        {
          type: row.dns == null ? 'error' : 'success',
          bordered: false
        },
        {
          default: () => (row.dns?.name == null ? '无' : row.dns.name)
        }
      )
    }
  },
  {
    title: '自动续签',
    key: 'auto_renew',
    width: 100,
    align: 'center',
    resizable: true,
    render(row: any) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.auto_renew
        //onUpdateValue: () => handleAutoRenewUpdate(row)
      })
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 350,
    align: 'center',
    hideInExcel: true,
    resizable: true,
    render(row: any) {
      return [
        row.cert_url == ''
          ? h(
              NButton,
              {
                size: 'small',
                type: 'info',
                style: 'margin-left: 15px;',
                onClick: async () => {
                  messageReactive = window.$message.loading('请稍后...', {
                    duration: 0
                  })
                  // 没有设置 DNS 接口和网站则获取解析记录
                  if (row.dns_id == 0 && row.website_id == 0) {
                    const { data } = await cert.manualDNS(row.id)
                    messageReactive.destroy()
                    window.$message.info('请先前往域名处设置 DNS 解析，再继续签发')
                    const d = window.$dialog.info({
                      style: 'width: 60vw',
                      title: '待设置DNS 记录列表',
                      content: () => {
                        return h(NTable, [
                          h('thead', [
                            h('tr', [
                              h('th', '域名'),
                              h('th', '类型'),
                              h('th', '主机记录'),
                              h('th', '记录值')
                            ])
                          ]),
                          h(
                            'tbody',
                            data.map((item: any) =>
                              h('tr', [
                                h('td', item?.domain),
                                h('td', 'TXT'),
                                h('td', item?.name),
                                h('td', item?.value)
                              ])
                            )
                          )
                        ])
                      },
                      positiveText: '签发',
                      onPositiveClick: async () => {
                        d.loading = true
                        messageReactive = window.$message.loading('请稍后...', {
                          duration: 0
                        })
                        cert
                          .obtain(row.id)
                          .then(() => {
                            window.$message.success('签发成功')
                            onPageChange(1)
                          })
                          .finally(() => {
                            d.loading = false
                            messageReactive?.destroy()
                          })
                      }
                    })
                  } else {
                    cert
                      .obtain(row.id)
                      .then(() => {
                        window.$message.success('签发成功')
                        onPageChange(1)
                      })
                      .finally(() => {
                        messageReactive?.destroy()
                      })
                  }
                }
              },
              {
                default: () => '签发'
              }
            )
          : null,
        row.cert != '' && row.key != ''
          ? h(
              NButton,
              {
                size: 'small',
                type: 'info',
                onClick: () => {
                  if (row.website_id != 0) {
                    deployCertModel.value.website_id = row.website_id
                  } else {
                    deployCertModel.value.website_id = 0
                  }
                  deployCertModel.value.id = row.id
                  deployCertModal.value = true
                }
              },
              {
                default: () => '部署'
              }
            )
          : null,
        row.cert != '' && row.key != ''
          ? h(
              NButton,
              {
                size: 'small',
                type: 'success',
                style: 'margin-left: 15px;',
                onClick: async () => {
                  messageReactive = window.$message.loading('请稍后...', {
                    duration: 0
                  })
                  await cert.renew(row.id)
                  messageReactive.destroy()
                  window.$message.success('续签成功')
                  onPageChange(1)
                }
              },
              {
                default: () => '续签'
              }
            )
          : null,
        row.cert != '' && row.key != ''
          ? h(
              NButton,
              {
                size: 'small',
                type: 'tertiary',
                style: 'margin-left: 15px;',
                onClick: () => {
                  showCertModel.value.cert = row.cert
                  showCertModel.value.key = row.key
                  showModal.value = true
                }
              },
              {
                default: () => '查看'
              }
            )
          : null,
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 15px;',
            onClick: () => {
              updateCert.value = row.id
              updateCertModel.value.domains = row.domains
              updateCertModel.value.dns_id = row.dns_id
              updateCertModel.value.type = row.type
              updateCertModel.value.account_id = row.account_id
              updateCertModel.value.website_id = row.website_id
              updateCertModel.value.auto_renew = row.auto_renew
              updateCertModal.value = true
            }
          },
          {
            default: () => '修改'
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: async () => {
              await cert.certDelete(row.id)
              window.$message.success('删除成功')
              onPageChange(1)
            }
          },
          {
            default: () => {
              return '确定删除证书吗？'
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
                  default: () => '删除'
                }
              )
            }
          }
        )
      ]
    }
  }
]
const data = ref<Cert[]>([] as Cert[])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200]
})

const onPageChange = (page: number) => {
  pagination.page = page
  getCertList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const getCertList = async (page: number, limit: number) => {
  const { data } = await cert.certs(page, limit)
  return data
}

const handleUpdateCert = async () => {
  await cert.certUpdate(updateCert.value, updateCertModel.value)
  window.$message.success('更新成功')
  updateCertModal.value = false
  onPageChange(1)
  updateCertModel.value.domains = []
  updateCertModel.value.dns_id = 0
  updateCertModel.value.type = 'P256'
  updateCertModel.value.account_id = 0
  updateCertModel.value.website_id = 0
  updateCertModel.value.auto_renew = true
}

const handleDeployCert = async () => {
  await cert.deploy(deployCertModel.value.id, deployCertModel.value.website_id)
  window.$message.success('部署成功')
  deployCertModal.value = false
  deployCertModel.value.id = 0
  deployCertModel.value.website_id = 0
  onPageChange(1)
}

const handleShowModalClose = () => {
  showCertModel.value.cert = ''
  showCertModel.value.key = ''
}

onMounted(() => {
  onPageChange(pagination.page)
  window.$bus.on('cert:refresh-cert', () => {
    onPageChange(pagination.page)
  })
})

onUnmounted(() => {
  window.$bus.off('cert:refresh-cert')
})
</script>

<template>
  <n-space vertical size="large">
    <n-data-table
      striped
      remote
      :scroll-x="1400"
      :loading="false"
      :columns="columns"
      :data="data"
      :row-key="(row: any) => row.id"
      :pagination="pagination"
      @update:page="onPageChange"
      @update:page-size="onPageSizeChange"
    />
  </n-space>
  <n-modal
    v-model:show="updateCertModal"
    preset="card"
    title="修改证书"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-alert type="info">
        可以通过选择网站 / DNS 中的任意一项来自动签发和部署证书，也可以手动输入域名并设置 DNS
        解析来签发证书
      </n-alert>
      <n-form :model="updateCertModel">
        <n-form-item label="域名">
          <n-dynamic-input
            v-model:value="updateCertModel.domains"
            placeholder="example.com"
            :min="1"
            show-sort-button
          />
        </n-form-item>
        <n-form-item path="type" label="密钥类型">
          <n-select
            v-model:value="updateCertModel.type"
            placeholder="选择密钥类型"
            clearable
            :options="algorithms"
          />
        </n-form-item>
        <n-form-item path="website_id" label="网站">
          <n-select
            v-model:value="updateCertModel.website_id"
            placeholder="选择用于部署证书的网站"
            clearable
            :options="websites"
          />
        </n-form-item>
        <n-form-item path="account_id" label="账号">
          <n-select
            v-model:value="updateCertModel.account_id"
            placeholder="选择用于签发证书的账号"
            clearable
            :options="accounts"
          />
        </n-form-item>
        <n-form-item path="account_id" label="DNS">
          <n-select
            v-model:value="updateCertModel.dns_id"
            placeholder="选择用于签发证书的DNS"
            clearable
            :options="dns"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleUpdateCert">提交</n-button>
    </n-space>
  </n-modal>
  <n-modal
    v-model:show="deployCertModal"
    preset="card"
    title="部署证书"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="deployCertModel">
        <n-form-item path="website_id" label="网站">
          <n-select
            v-model:value="deployCertModel.website_id"
            placeholder="选择需要部署证书的网站"
            clearable
            :options="websites"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleDeployCert">提交</n-button>
    </n-space>
  </n-modal>
  <n-modal
    v-model:show="showModal"
    preset="card"
    title="查看证书"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="handleShowModalClose"
  >
    <n-tabs type="line" animated>
      <n-tab-pane name="cert" tab="证书">
        <Editor
          v-model:value="showCertModel.cert"
          theme="vs-dark"
          height="60vh"
          mt-8
          :options="{
            readOnly: true,
            automaticLayout: true
          }"
        />
      </n-tab-pane>
      <n-tab-pane name="key" tab="密钥">
        <Editor
          v-model:value="showCertModel.key"
          theme="vs-dark"
          height="60vh"
          mt-8
          :options="{
            readOnly: true,
            automaticLayout: true
          }"
        />
      </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>

<style scoped lang="scss"></style>
