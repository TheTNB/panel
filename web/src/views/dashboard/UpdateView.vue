<script setup lang="ts">
defineOptions({
  name: 'home-update'
})

import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import type { MessageReactive } from 'naive-ui'
import { NButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import dashboard from '@/api/panel/dashboard'
import { router } from '@/router'
import { formatDateTime } from '@/utils'
import type { Version } from '@/views/home/types'

const { t } = useI18n()
const versions = ref<Version[] | null>(null)
let messageReactive: MessageReactive | null = null

const getVersions = () => {
  dashboard.updateInfo().then((res: any) => {
    versions.value = res.data
  })
}

const handleUpdate = () => {
  window.$dialog.warning({
    title: t('homeUpdate.confirm.update.title'),
    content: t('homeUpdate.confirm.update.content'),
    positiveText: t('homeUpdate.confirm.update.positiveText'),
    negativeText: t('homeUpdate.confirm.update.negativeText'),
    onPositiveClick: () => {
      messageReactive = window.$message.loading(t('homeUpdate.confirm.update.loading'), {
        duration: 0
      })
      dashboard
        .update()
        .then(() => {
          messageReactive?.destroy()
          window.$message.success(t('homeUpdate.alerts.success'))
          setTimeout(() => {
            setTimeout(() => {
              window.location.reload()
            }, 400)
            router.push({ name: 'home-index' })
          }, 2500)
        })
        .catch(() => {
          messageReactive?.destroy()
        })
    },
    onNegativeClick: () => {
      window.$message.info(t('homeUpdate.alerts.info'))
    }
  })
}

onMounted(() => {
  getVersions()
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <div>
        <n-button v-if="versions" class="ml-16" type="primary" @click="handleUpdate">
          <TheIcon :size="18" icon="material-symbols:arrow-circle-up-outline-rounded" />
          {{ $t('homeUpdate.button.update') }}
        </n-button>
      </div>
    </template>
    <n-timeline v-if="versions" pt-10>
      <n-timeline-item
        v-for="(item, index) in versions"
        :type="Number(index) == 0 ? 'info' : 'default'"
        :key="index"
        :title="item.version"
        :time="formatDateTime(item.updated_at)"
      >
        <MdPreview
          v-model="item.description"
          noMermaid
          noKatex
          noIconfont
          noHighlight
          noImgZoomIn
        />
      </n-timeline-item>
    </n-timeline>
    <div v-else pt-40>
      <n-result status="418" title="Loading..." :description="$t('homeUpdate.loading')"> </n-result>
    </div>
  </common-page>
</template>
