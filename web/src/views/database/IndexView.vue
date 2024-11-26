<script setup lang="ts">
defineOptions({
  name: 'database-index'
})

import CreateDatabaseModal from '@/views/database/CreateDatabaseModal.vue'
import CreateServerModal from '@/views/database/CreateServerModal.vue'
import CreateUserModal from '@/views/database/CreateUserModal.vue'
import DatabaseList from '@/views/database/DatabaseList.vue'
import ServerList from '@/views/database/ServerList.vue'
import UserList from '@/views/database/UserList.vue'
import { NButton } from 'naive-ui'

const currentTab = ref('database')

const createDatabaseModalShow = ref(false)
const createUserModalShow = ref(false)
const createServerModalShow = ref(false)
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
      <n-button v-if="currentTab === 'user'" type="primary" @click="createUserModalShow = true">
        <TheIcon :size="18" icon="material-symbols:add" />
        创建用户
      </n-button>
      <n-button v-if="currentTab === 'server'" type="primary" @click="createServerModalShow = true">
        <TheIcon :size="18" icon="material-symbols:add" />
        添加服务器
      </n-button>
    </template>
    <n-flex vertical>
      <n-tabs v-model:value="currentTab" type="line" animated>
        <n-tab-pane name="database" tab="数据库">
          <database-list v-model:type="currentTab" />
        </n-tab-pane>
        <n-tab-pane name="user" tab="用户">
          <user-list v-model:type="currentTab" />
        </n-tab-pane>
        <n-tab-pane name="server" tab="服务器">
          <server-list v-model:type="currentTab" />
        </n-tab-pane>
      </n-tabs>
    </n-flex>
  </common-page>
  <create-database-modal v-model:show="createDatabaseModalShow" />
  <create-user-modal v-model:show="createUserModalShow" />
  <create-server-modal v-model:show="createServerModalShow" />
</template>
