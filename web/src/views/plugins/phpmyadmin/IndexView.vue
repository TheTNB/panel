<script setup lang="ts">
import { NButton } from 'naive-ui'
import phpmyadmin from '@/api/plugins/phpmyadmin'
import Editor from '@guolao/vue-monaco-editor'
import { themeConfig, themeDarkConfig, tokenConf } from 'monaco-editor-nginx/cjs/conf'
import suggestions from 'monaco-editor-nginx/cjs/suggestions'
import { directives } from 'monaco-editor-nginx/cjs/directives'

const currentTab = ref('status')
const hostname = ref(window.location.hostname)
const port = ref(0)
const path = ref('')
const newPort = ref(0)
const url = computed(() => {
  return `http://${hostname.value}:${port.value}/${path.value}`
})
const config = ref('')

const getInfo = async () => {
  phpmyadmin.info().then((res: any) => {
    path.value = res.data.path
    port.value = res.data.port
    newPort.value = res.data.port
  })
}

const handleSave = async () => {
  await phpmyadmin.port(newPort.value)
  window.$message.success('保存成功')
  await getInfo()
}

const getConfig = async () => {
  const { data } = await phpmyadmin.getConfig()
  return data
}

const handleSaveConfig = async () => {
  await phpmyadmin.saveConfig(config.value)
  window.$message.success('保存成功')
}

const editorOnBeforeMount = (monaco: any) => {
  monaco.languages.register({
    id: 'nginx'
  })

  monaco.languages.setMonarchTokensProvider('nginx', tokenConf)
  monaco.editor.defineTheme('nginx-theme', themeConfig)
  monaco.editor.defineTheme('nginx-theme-dark', themeDarkConfig)

  monaco.languages.registerCompletionItemProvider('nginx', {
    provideCompletionItems: (model: any, position: any) => {
      const word = model.getWordUntilPosition(position)
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn
      }
      return { suggestions: suggestions(range) }
    }
  })

  monaco.languages.registerHoverProvider('nginx', {
    provideHover: (model: any, position: any) => {
      const word = model.getWordAtPosition(position)
      if (!word) return
      const data = directives.find((item) => item.n === word.word || item.n === `$${word.word}`)
      if (!data) return
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn
      }
      const contents = [{ value: `**\`${data.n}\`** | ${data.m} | ${data.c || ''}` }]
      if (data.s) {
        contents.push({ value: `**syntax:** ${data.s || ''}` })
      }
      if (data.v) {
        contents.push({ value: `**default:** ${data.v || ''}` })
      }
      if (data.d) {
        contents.push({ value: `${data.d}` })
      }
      return {
        contents: [...contents],
        range: range
      }
    }
  })
}

onMounted(() => {
  getInfo()
  getConfig().then((res) => {
    config.value = res
  })
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button v-if="currentTab == 'status'" class="ml-16" type="primary" @click="handleSave">
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'config'"
        class="ml-16"
        type="primary"
        @click="handleSaveConfig"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="status" tab="状态">
        <n-space vertical>
          <n-card title="访问信息" rounded-10>
            <n-alert type="info">
              访问地址: <a :href="url" target="_blank">{{ url }}</a>
            </n-alert>
          </n-card>
          <n-card title="修改端口" rounded-10>
            <n-input-number v-model:value="newPort" min="1" />
            修改 phpMyAdmin 访问端口
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="config" tab="修改配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 phpMyAdmin 的 OpenResty
            配置文件，如果你不了解各参数的含义，请不要随意修改！
          </n-alert>
          <Editor
            v-model:value="config"
            language="nginx"
            theme="nginx-theme-dark"
            height="60vh"
            mt-8
            @before-mount="editorOnBeforeMount"
            :options="{
              automaticLayout: true,
              formatOnType: true,
              formatOnPaste: true
            }"
          />
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
