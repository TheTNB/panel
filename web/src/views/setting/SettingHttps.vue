<script setup lang="ts">
import setting from '@/api/panel/setting'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = ref({
  https: false,
  cert: '',
  key: ''
})

const handleSave = () => {
  setting.updateHttps(model.value).then(() => {
    window.$message.success(t('settingIndex.edit.toasts.success'))
    setTimeout(() => {}, 1000)
  })
}

onMounted(() => {
  setting.getHttps().then((res) => {
    model.value = res.data
  })
})
</script>

<template>
  <n-space vertical>
    <n-alert type="warning"> 错误的证书会导致面板无法访问，请谨慎操作！ </n-alert>
    <n-form>
      <n-form-item :label="$t('settingIndex.edit.fields.https.label')">
        <n-switch v-model:value="model.https" />
      </n-form-item>
      <n-form-item v-if="model.https" :label="$t('settingIndex.edit.fields.cert.label')">
        <n-input v-model:value="model.cert" type="textarea" />
      </n-form-item>
      <n-form-item v-if="model.https" :label="$t('settingIndex.edit.fields.key.label')">
        <n-input v-model:value="model.key" type="textarea" />
      </n-form-item>
    </n-form>
  </n-space>
  <n-button type="primary" @click="handleSave">
    {{ $t('settingIndex.edit.actions.submit') }}
  </n-button>
</template>

<style scoped lang="scss"></style>
