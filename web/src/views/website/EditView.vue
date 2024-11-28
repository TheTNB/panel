<script setup lang="ts">
defineOptions({
  name: 'website-edit'
})

import Editor from '@guolao/vue-monaco-editor'
import type { MessageReactive } from 'naive-ui'
import { NButton } from 'naive-ui'

import cert from '@/api/panel/cert'
import dashboard from '@/api/panel/dashboard'
import website from '@/api/panel/website'
import type { Cert } from '@/views/cert/types'
import type { WebsiteListen, WebsiteSetting } from '@/views/website/types'

let messageReactive: MessageReactive | null = null

const current = ref('listen')
const route = useRoute()
const { id } = route.params

const setting = ref<WebsiteSetting>({
  id: 0,
  name: '',
  listens: [] as WebsiteListen[],
  domains: [],
  root: '',
  path: '',
  index: [],
  php: 0,
  open_basedir: false,
  https: false,
  ssl_certificate: '',
  ssl_certificate_key: '',
  ssl_not_before: '',
  ssl_not_after: '',
  ssl_dns_names: [],
  ssl_issuer: '',
  ssl_ocsp_server: [],
  http_redirect: false,
  hsts: false,
  ocsp: false,
  rewrite: '',
  raw: '',
  log: ''
})
const installedDbAndPhp = ref({
  php: [
    {
      label: '不使用',
      value: 0
    }
  ],
  db: [
    {
      label: '',
      value: ''
    }
  ]
})
const certs = ref<Cert[]>([] as Cert[])
const { data: rewrites }: { data: any } = useRequest(website.rewrites, {
  initialData: {}
})
const rewriteOptions = computed(() => {
  return Object.keys(rewrites.value).map((key) => ({
    label: key,
    value: key
  }))
})
const rewriteValue = ref(null)

const title = computed(() => {
  if (setting.value) {
    return `编辑网站 - ${setting.value.name}`
  }
  return '编辑网站 - 加载中...'
})
const certOptions = computed(() => {
  return certs.value.map((item) => ({
    label: item.domains.join(', '),
    value: item.id
  }))
})
const selectedCert = ref(null)

const fetchPhpAndDb = async () => {
  const { data } = await dashboard.installedDbAndPhp()
  installedDbAndPhp.value = data
}

const fetchWebsiteSetting = async () => {
  await website.config(Number(id)).then((res) => {
    setting.value = res.data
  })
}

const fetchCertList = async () => {
  const { data } = await cert.certs(1, 10000)
  certs.value = data.items
}

const handleSave = async () => {
  // 如果没有任何监听地址设置了https，则自动添加443
  if (setting.value.https && !setting.value.listens.some((item) => item.https)) {
    setting.value.listens.push({
      address: '443',
      https: true,
      quic: true
    })
  }
  // 如果关闭了https，自动禁用所有https和quic
  if (!setting.value.https) {
    setting.value.listens = setting.value.listens.filter((item) => item.address !== '443') // 443直接删掉
    setting.value.listens.forEach((item) => {
      item.https = false
      item.quic = false
    })
  }

  await website.saveConfig(Number(id), setting.value).then(() => {
    fetchWebsiteSetting()
    window.$message.success('保存成功')
  })
}

const handleReset = async () => {
  await website.resetConfig(Number(id)).then(() => {
    fetchWebsiteSetting()
    window.$message.success('重置成功')
  })
}

const handleRewrite = (value: string) => {
  setting.value.rewrite = rewrites.value[value] || ''
}

const handleObtainCert = async () => {
  messageReactive = window.$message.loading('请稍后...', {
    duration: 0
  })
  await website
    .obtainCert(Number(id))
    .then(() => {
      fetchWebsiteSetting()
      window.$message.success('签发成功')
    })
    .finally(() => {
      messageReactive?.destroy()
    })
}

const handleSelectCert = (value: number) => {
  const cert = certs.value.find((item) => item.id === value)
  if (cert) {
    setting.value.ssl_certificate = cert.cert
    setting.value.ssl_certificate_key = cert.key
  }
}

const clearLog = async () => {
  await website.clearLog(Number(id)).then(() => {
    fetchWebsiteSetting()
    window.$message.success('清空成功')
  })
}

const onCreateListen = () => {
  return {
    address: '',
    https: false,
    quic: false
  }
}

onMounted(async () => {
  await fetchWebsiteSetting()
  await fetchPhpAndDb()
  await fetchCertList()
})
</script>

<template>
  <common-page show-footer :title="title">
    <template #action>
      <n-flex>
        <n-tag v-if="current === 'config'" type="warning">
          如果您修改了原文，那么点击保存后，其余的修改将不会生效！
        </n-tag>
        <n-popconfirm v-if="current === 'config'" @positive-click="handleReset">
          <template #trigger>
            <n-button type="success">
              <TheIcon :size="18" icon="material-symbols:refresh" />
              重置配置
            </n-button>
          </template>
          确定要重置配置吗？
        </n-popconfirm>
        <n-button v-if="current === 'https'" class="ml-16" type="success" @click="handleObtainCert">
          <TheIcon :size="18" icon="material-symbols:done-rounded" />
          一键签发证书
        </n-button>
        <n-button v-if="current !== 'log'" class="ml-16" type="primary" @click="handleSave">
          <TheIcon :size="18" icon="material-symbols:save-outline" />
          保存
        </n-button>
        <n-popconfirm v-if="current === 'log'" @positive-click="clearLog">
          <template #trigger>
            <n-button type="primary">
              <TheIcon :size="18" icon="material-symbols:delete-outline" />
              清空日志
            </n-button>
          </template>
          确定要清空吗？
        </n-popconfirm>
      </n-flex>
    </template>

    <n-tabs v-model:value="current" type="line" animated>
      <n-tab-pane name="listen" tab="域名监听">
        <n-form v-if="setting">
          <n-form-item label="域名">
            <n-dynamic-input
              v-model:value="setting.domains"
              placeholder="example.com"
              :min="1"
              show-sort-button
            />
          </n-form-item>
          <n-form-item label="监听地址">
            <n-dynamic-input
              v-model:value="setting.listens"
              show-sort-button
              :on-create="onCreateListen"
            >
              <template #default="{ value }">
                <div w-full flex items-center>
                  <n-input v-model:value="value.address" clearable />
                  <n-checkbox v-model:checked="value.https" ml-20 mr-20 w-120> HTTPS </n-checkbox>
                  <n-checkbox v-model:checked="value.quic" w-200> QUIC(HTTP3) </n-checkbox>
                </div>
              </template>
            </n-dynamic-input>
          </n-form-item>
        </n-form>
        <n-skeleton v-else text :repeat="10" />
      </n-tab-pane>
      <n-tab-pane name="basic" tab="基本设置">
        <n-form v-if="setting">
          <n-form-item label="网站目录">
            <n-input v-model:value="setting.path" placeholder="输入网站目录（绝对路径）" />
          </n-form-item>
          <n-form-item label="运行目录">
            <n-input
              v-model:value="setting.root"
              placeholder="输入运行目录（Laravel 等程序需要）（绝对路径）"
            />
          </n-form-item>
          <n-form-item label="默认文档">
            <n-dynamic-tags v-model:value="setting.index" />
          </n-form-item>
          <n-form-item label="PHP 版本">
            <n-select
              v-model:value="setting.php"
              :default-value="0"
              :options="installedDbAndPhp.php"
              placeholder="选择PHP版本"
              @keydown.enter.prevent
            >
            </n-select>
          </n-form-item>
          <n-form-item label="防跨站攻击（PHP）">
            <n-switch v-model:value="setting.open_basedir" />
          </n-form-item>
        </n-form>
        <n-skeleton v-else text :repeat="10" />
      </n-tab-pane>
      <n-tab-pane name="https" tab="HTTPS">
        <n-flex vertical v-if="setting">
          <n-card v-if="setting.https && setting.ssl_issuer != ''">
            <n-descriptions title="证书信息" :column="2">
              <n-descriptions-item>
                <template #label>证书有效期</template>
                <n-flex>
                  <n-tag>{{ setting.ssl_not_before }}</n-tag>
                  -
                  <n-tag>{{ setting.ssl_not_after }}</n-tag>
                </n-flex>
              </n-descriptions-item>
              <n-descriptions-item>
                <template #label>颁发者</template>
                <n-flex>
                  <n-tag>{{ setting.ssl_issuer }}</n-tag>
                </n-flex>
              </n-descriptions-item>
              <n-descriptions-item>
                <template #label>域名</template>
                <n-flex>
                  <n-tag v-for="item in setting.ssl_dns_names" :key="item">{{ item }}</n-tag>
                </n-flex>
              </n-descriptions-item>
              <n-descriptions-item>
                <template #label>OCSP</template>
                <n-flex>
                  <n-tag v-for="item in setting.ssl_ocsp_server" :key="item">{{ item }}</n-tag>
                </n-flex>
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
          <n-form>
            <n-grid :cols="24" :x-gap="24">
              <n-form-item-gi :span="12" label="总开关（只有打开了总开关，下面的设置才会生效！）">
                <n-switch v-model:value="setting.https" />
              </n-form-item-gi>
              <n-form-item-gi v-if="setting.https" :span="12" label="使用已有证书">
                <n-select
                  v-model:value="selectedCert"
                  :options="certOptions"
                  @update-value="handleSelectCert"
                />
              </n-form-item-gi>
            </n-grid>
          </n-form>
          <n-form inline>
            <n-form-item label="HSTS">
              <n-switch v-model:value="setting.hsts" />
            </n-form-item>
            <n-form-item label="HTTP 跳转">
              <n-switch v-model:value="setting.http_redirect" />
            </n-form-item>
            <n-form-item label="OCSP 装订">
              <n-switch v-model:value="setting.ocsp" />
            </n-form-item>
          </n-form>
          <n-form>
            <n-form-item label="证书">
              <n-input
                v-model:value="setting.ssl_certificate"
                type="textarea"
                placeholder="输入 PEM 证书文件的内容"
              />
            </n-form-item>
            <n-form-item label="私钥">
              <n-input
                v-model:value="setting.ssl_certificate_key"
                type="textarea"
                placeholder="输入 KEY 私钥文件的内容"
              />
            </n-form-item>
          </n-form>
        </n-flex>
        <n-skeleton v-else text :repeat="10" />
      </n-tab-pane>
      <n-tab-pane name="rewrite" tab="伪静态">
        <n-flex vertical>
          <n-form label-placement="left" label-width="auto">
            <n-form-item label="预设">
              <n-select
                v-model:value="rewriteValue"
                clearable
                :options="rewriteOptions"
                @update-value="handleRewrite"
              />
            </n-form-item>
          </n-form>
          <Editor
            v-if="setting"
            v-model:value="setting.rewrite"
            language="ini"
            theme="vs-dark"
            height="60vh"
            :options="{
              automaticLayout: true,
              formatOnType: true,
              formatOnPaste: true
            }"
          />
        </n-flex>
      </n-tab-pane>
      <n-tab-pane name="config" tab="配置原文">
        <n-flex vertical>
          <n-alert type="warning" w-full>
            如果您不了解配置规则，请勿随意修改，否则可能会导致网站无法访问或面板功能异常！如果已经遇到问题，可尝试重置配置！
          </n-alert>
          <Editor
            v-if="setting"
            v-model:value="setting.raw"
            language="ini"
            theme="vs-dark"
            height="60vh"
            :options="{
              automaticLayout: true,
              formatOnType: true,
              formatOnPaste: true
            }"
          />
        </n-flex>
      </n-tab-pane>
      <n-tab-pane name="log" tab="访问日志">
        <n-flex vertical>
          <n-flex flex items-center>
            <n-alert type="warning" w-full>
              全部日志可通过下载文件
              <n-tag>{{ setting.log }}</n-tag>
              查看。
            </n-alert>
          </n-flex>
          <realtime-log :path="setting.log" />
        </n-flex>
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
