<script setup lang="ts">
defineOptions({
  name: 'database-index'
})

import CreateDatabaseModal from '@/views/database/CreateDatabaseModal.vue'
import CreateDatabaseServerModal from '@/views/database/CreateServerModal.vue'
import DatabaseListView from '@/views/database/DatabaseList.vue'
import ServerListView from '@/views/database/ServerList.vue'
import { NButton } from 'naive-ui'

const currentTab = ref('database')

const createDatabaseModalShow = ref(false)
const createDatabaseServerModalShow = ref(false)
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button
        v-if="currentTab === 'database'"
        type="primary"
        @click="createDatabaseModalShow = true"
      >
        <TheIcon :size="18" icon="material-symbols:add" />
        创建数据库
      </n-button>
      <n-flex v-if="currentTab === 'server'">
        <n-button type="success">
          <TheIcon :size="18" icon="material-symbols:sync" />
          同步数据库
        </n-button>
        <n-button type="primary" @click="createDatabaseServerModalShow = true">
          <TheIcon :size="18" icon="material-symbols:add" />
          添加服务器
        </n-button>
      </n-flex>
    </template>
    <n-flex vertical>
      <n-tabs v-model:value="currentTab" type="line" animated>
        <n-tab-pane name="database" tab="数据库">
          <database-list-view v-model:type="currentTab" />
        </n-tab-pane>
        <n-tab-pane name="server" tab="服务器">
          <server-list-view v-model:type="currentTab" />
        </n-tab-pane>
      </n-tabs>
    </n-flex>
  </common-page>
  <create-database-modal v-model:show="createDatabaseModalShow" />
  <create-database-server-modal v-model:show="createDatabaseServerModalShow" />
</template>
