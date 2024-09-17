<script setup lang="ts">
import info from '@/api/panel/info'
import { NButton } from 'naive-ui'
import type { MessageReactive } from 'naive-ui'
import type { PanelInfo } from '@/views/home/types'
import { formatDateTime } from '@/utils'
import { router } from '@/router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const versions = ref<PanelInfo[] | null>(null)
let messageReactive: MessageReactive | null = null

const getVersions = () => {
  info.updateInfo().then((res: any) => {
    versions.value = res.data.reverse()
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
      info
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
          <TheIcon
            :size="18"
            class="mr-5"
            icon="material-symbols:arrow-circle-up-outline-rounded"
          />
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
        :time="formatDateTime(item.date)"
      >
        <pre v-html="item.body" />
      </n-timeline-item>
    </n-timeline>
    <div v-else pt-40>
      <n-result status="418" title="Loading..." :description="$t('homeUpdate.loading')"> </n-result>
    </div>
  </common-page>
</template>
