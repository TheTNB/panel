<script setup lang="ts">
import type { App } from '@/views/app/types'
import { useI18n } from 'vue-i18n'
import app from '../../api/panel/app'

const { t } = useI18n()

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const operation = defineModel<string>('operation', { type: String, required: true })
const info = defineModel<App>('info', { type: Object, required: true })

const doSubmit = ref(false)

const model = reactive({
  channel: '',
  version: ''
})

const options = computed(() => {
  return info.value.channels.map((channel) => {
    return {
      label: channel.name,
      value: channel.slug
    }
  })
})

const handleSubmit = () => {
  app.install(info.value.slug, model.channel).then(() => {
    window.$message.success(t('appIndex.alerts.install'))
  })
}

const handleChannelUpdate = (value: string) => {
  const channel = info.value.channels.find((channel) => channel.slug === value)
  if (channel) {
    model.version = channel.subs[0].version
  }
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="operation + ' ' + info.name"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form :model="model">
      <n-form-item path="channel" label="渠道">
        <n-select
          v-model:value="model.channel"
          :options="options"
          @update-value="handleChannelUpdate"
        />
      </n-form-item>
      <n-form-item path="channel" label="版本号">
        <n-input v-model:value="model.version" placeholder="请选择渠道" readonly disabled />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button type="info" block :loading="doSubmit" :disabled="doSubmit" @click="handleSubmit">
          提交
        </n-button>
      </n-col>
    </n-row>
  </n-modal>
</template>

<style scoped lang="scss"></style>
