<script setup lang="ts">
import app from '@/api/panel/app'
import cron from '@/api/panel/cron'
import dashboard from '@/api/panel/dashboard'
import website from '@/api/panel/website'
import Editor from '@guolao/vue-monaco-editor'
import { CronNaive } from '@vue-js-cron/naive-ui'
import { NInput } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const loading = ref(false)

const createModel = ref({
  name: '',
  type: 'shell',
  target: '',
  save: 1,
  backup_type: 'website',
  backup_path: '',
  script: '# 在此输入您要执行的脚本内容',
  time: '* * * * *'
})

const websites = ref<any>([])
const installedDbAndPhp = ref({
  php: [
    {
      label: '',
      value: ''
    }
  ],
  db: [
    {
      label: '',
      value: ''
    }
  ]
})

const mySQLInstalled = computed(() => {
  return installedDbAndPhp.value.db.find((item) => item.value === 'mysql')
})

const postgreSQLInstalled = computed(() => {
  return installedDbAndPhp.value.db.find((item) => item.value === 'postgresql')
})

const getWebsiteList = async (page: number, limit: number) => {
  const { data } = await website.list(page, limit)
  for (const item of data.items) {
    websites.value.push({
      label: item.name,
      value: item.name
    })
  }
  createModel.value.target = websites.value[0]?.value
}

const getPhpAndDb = async () => {
  const { data } = await dashboard.installedDbAndPhp()
  installedDbAndPhp.value = data
}

const handleSubmit = async () => {
  loading.value = true
  await cron
    .create(createModel.value)
    .then(() => {
      window.$message.success('创建成功')
      window.$bus.emit('task:refresh-cron')
      loading.value = false
      show.value = false
    })
    .catch(() => {
      loading.value = false
    })
}

watch(createModel, (value) => {
  if (value.backup_type === 'website') {
    createModel.value.target = websites.value[0]?.value
  } else {
    createModel.value.target = ''
  }
})

onMounted(() => {
  getPhpAndDb()
  app.isInstalled('nginx').then((res) => {
    if (res.data.installed) {
      getWebsiteList(1, 10000)
    }
  })
})
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建计划任务"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form>
      <n-form-item label="任务类型">
        <n-select
          v-model:value="createModel.type"
          :options="[
            { label: '运行脚本', value: 'shell' },
            { label: '备份数据', value: 'backup' },
            { label: '切割日志', value: 'cutoff' }
          ]"
        >
        </n-select>
      </n-form-item>
      <n-form-item label="任务名称">
        <n-input v-model:value="createModel.name" placeholder="任务名称" />
      </n-form-item>
      <n-form-item label="任务周期">
        <cron-naive v-model="createModel.time" locale="zh-cn"></cron-naive>
      </n-form-item>
      <div v-if="createModel.type === 'shell'">
        <n-text>脚本内容</n-text>
        <Editor
          v-model:value="createModel.script"
          language="shell"
          theme="vs-dark"
          height="40vh"
          mt-8
          :options="{
            automaticLayout: true,
            formatOnType: true,
            formatOnPaste: true
          }"
        />
      </div>
      <n-form-item v-if="createModel.type === 'backup'" label="备份类型">
        <n-radio-group v-model:value="createModel.backup_type">
          <n-radio value="website">网站目录</n-radio>
          <n-radio value="mysql" :disabled="!mySQLInstalled"> MySQL 数据库</n-radio>
          <n-radio value="postgres" :disabled="!postgreSQLInstalled"> PostgreSQL 数据库 </n-radio>
        </n-radio-group>
      </n-form-item>
      <n-form-item
        v-if="
          (createModel.backup_type === 'website' && createModel.type === 'backup') ||
          createModel.type === 'cutoff'
        "
        label="选择网站"
      >
        <n-select v-model:value="createModel.target" :options="websites" placeholder="选择网站" />
      </n-form-item>
      <n-form-item
        v-if="createModel.backup_type !== 'website' && createModel.type === 'backup'"
        label="数据库名"
      >
        <n-input v-model:value="createModel.target" placeholder="数据库名" />
      </n-form-item>
      <n-form-item v-if="createModel.type === 'backup'" label="保存目录">
        <n-input v-model:value="createModel.backup_path" placeholder="保存目录" />
      </n-form-item>
      <n-form-item v-if="createModel.type !== 'shell'" label="保留份数">
        <n-input-number v-model:value="createModel.save" />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]" pt-20>
      <n-col :span="24">
        <n-button type="info" block :loading="loading" @click="handleSubmit"> 提交 </n-button>
      </n-col>
    </n-row>
  </n-modal>
</template>

<style scoped lang="scss"></style>
