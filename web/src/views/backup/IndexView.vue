<script setup lang="ts">
import backup from '@/api/panel/backup'
import ListView from '@/views/backup/ListView.vue'
import { NButton, NInput } from 'naive-ui'

const currentTab = ref('website')
const createModal = ref(false)
const createModel = ref({
  target: '',
  path: ''
})
const oldTab = ref('')

const handleCreate = () => {
  backup.create(currentTab.value, createModel.value.target, createModel.value.path).then(() => {
    createModal.value = false
    window.$message.success('创建成功')
    // 有点low，但是没找到更好的办法
    oldTab.value = currentTab.value
    currentTab.value = ''
    setTimeout(() => {
      currentTab.value = oldTab.value
    }, 0)
  })
}
</script>

<template>
  <common-page show-footer>
    <template #action>
      <div flex items-center>
        <n-button class="ml-16" type="primary" @click="createModal = true">
          <TheIcon :size="18" icon="material-symbols:add" />
          创建备份
        </n-button>
      </div>
    </template>
    <n-flex vertical>
      <n-alert type="info">此处仅显示面板默认备份目录的文件。</n-alert>
      <n-tabs v-model:value="currentTab" type="line" animated>
        <n-tab-pane name="website" tab="网站">
          <list-view v-model:type="currentTab" />
        </n-tab-pane>
        <n-tab-pane name="mysql" tab="MySQL">
          <list-view v-model:type="currentTab" />
        </n-tab-pane>
        <n-tab-pane name="postgres" tab="PostgreSQL">
          <list-view v-model:type="currentTab" />
        </n-tab-pane>
      </n-tabs>
    </n-flex>
  </common-page>
  <n-modal
    v-model:show="createModal"
    preset="card"
    title="创建备份"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="createModal = false"
  >
    <n-form :model="createModel">
      <n-form-item path="name" label="名称">
        <n-input
          v-model:value="createModel.target"
          type="text"
          @keydown.enter.prevent
          placeholder="输入网站/数据库名称"
        />
      </n-form-item>
      <n-form-item path="path" label="目录">
        <n-input
          v-model:value="createModel.path"
          type="text"
          @keydown.enter.prevent
          placeholder="留空使用默认路径"
        />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button type="info" block @click="handleCreate">提交</n-button>
      </n-col>
    </n-row>
  </n-modal>
</template>
