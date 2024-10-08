<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import setting from '@/api/panel/setting'
import { useThemeStore } from '@/store'
import type { Setting } from '@/views/setting/types'

const { t } = useI18n()
const themeStore = useThemeStore()

const model = ref<Setting>({
  name: '',
  locale: '',
  username: '',
  password: '',
  email: '',
  port: 8888,
  entrance: '',
  website_path: '',
  backup_path: '',
  https: false,
  cert: '',
  key: ''
})

const locales = [
  { label: '简体中文', value: 'zh_CN' },
  { label: 'English', value: 'en' }
]

const getSetting = () => {
  setting.list().then((res) => {
    model.value = res.data
  })
}

const handleSave = () => {
  setting.update(model.value).then(() => {
    window.$message.success(t('settingIndex.edit.toasts.success'))
    setTimeout(() => {
      maybeHardReload()
    }, 1000)
  })
}

const maybeHardReload = () => {
  if (model.value.locale !== themeStore.locale) {
    window.location.reload()
  }
}

onMounted(() => {
  getSetting()
})
</script>

<template>
  <n-space vertical>
    <n-alert type="info">
      {{ $t('settingIndex.info') }}
    </n-alert>
    <n-form>
      <n-form-item :label="$t('settingIndex.edit.fields.name.label')">
        <n-input
          v-model:value="model.name"
          :placeholder="$t('settingIndex.edit.fields.name.placeholder')"
        />
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.locale.label')">
        <n-select v-model:value="model.locale" :options="locales"> </n-select>
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.username.label')">
        <n-input
          v-model:value="model.username"
          :placeholder="$t('settingIndex.edit.fields.username.placeholder')"
        />
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.password.label')">
        <n-input
          v-model:value="model.password"
          :placeholder="$t('settingIndex.edit.fields.password.placeholder')"
        />
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.email.label')">
        <n-input
          v-model:value="model.email"
          :placeholder="$t('settingIndex.edit.fields.email.placeholder')"
        />
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.port.label')">
        <n-input-number
          v-model:value="model.port"
          :placeholder="$t('settingIndex.edit.fields.port.placeholder')"
        />
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.entrance.label')">
        <n-input
          v-model:value="model.entrance"
          :placeholder="$t('settingIndex.edit.fields.entrance.placeholder')"
        />
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.path.label')">
        <n-input
          v-model:value="model.website_path"
          :placeholder="$t('settingIndex.edit.fields.path.placeholder')"
        />
      </n-form-item>
      <n-form-item :label="$t('settingIndex.edit.fields.backup.label')">
        <n-input
          v-model:value="model.backup_path"
          :placeholder="$t('settingIndex.edit.fields.backup.placeholder')"
        />
      </n-form-item>
    </n-form>
  </n-space>
  <n-button type="primary" @click="handleSave">
    {{ $t('settingIndex.edit.actions.submit') }}
  </n-button>
</template>

<style scoped lang="scss"></style>
