<script setup lang="ts">
import PathInput from '@/views/file/PathInput.vue'
import ToolBar from '@/views/file/ToolBar.vue'
import ListTable from '@/views/file/ListTable.vue'
import type { Marked } from '@/views/file/types'
import ArchiveModal from '@/views/file/ArchiveModal.vue'
import PermissionModal from '@/views/file/PermissionModal.vue'

const path = ref('/www')
const selected = ref<string[]>([])
const marked = ref<Marked[]>([])

const archive = ref(false)
const permission = ref(false)
</script>

<template>
  <common-page show-footer>
    <n-flex vertical :size="20">
      <path-input v-model:path="path" />
      <tool-bar
        v-model:path="path"
        v-model:selected="selected"
        v-model:marked="marked"
        v-model:archive="archive"
        v-model:permission="permission"
      />
      <list-table
        v-model:path="path"
        v-model:selected="selected"
        v-model:marked="marked"
        v-model:archive="archive"
        v-model:permission="permission"
      />
      <archive-modal v-model:show="archive" v-model:path="path" v-model:selected="selected" />
      <permission-modal v-model:show="permission" v-model:selected="selected" />
    </n-flex>
  </common-page>
</template>
