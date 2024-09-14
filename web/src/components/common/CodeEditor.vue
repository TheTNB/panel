<script setup lang="ts">
import Editor from '@guolao/vue-monaco-editor'
import { themeConfig, themeDarkConfig, tokenConf } from 'monaco-editor-nginx/cjs/conf'
import suggestions from 'monaco-editor-nginx/cjs/suggestions'
import { directives } from 'monaco-editor-nginx/cjs/directives'
import file from '@/api/panel/file'
import { languageByPath } from '@/utils/file'

const props = defineProps({
  path: String,
  readOnly: Boolean
})

const disabled = ref(false) // 在出现错误的情况下禁用保存
const data = ref('')

const get = async () => {
  await file
    .content(props.path!)
    .then((res) => {
      data.value = res.data
      window.$message.success('获取成功')
    })
    .catch(() => {
      disabled.value = true
    })
}

const save = async () => {
  if (disabled.value) {
    window.$message.error('当前状态下不可保存')
    return
  }
  await file.save(props.path!, data.value)
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
  get()
})

defineExpose({
  get,
  save
})
</script>

<template>
  <Editor
    v-model:value="data"
    :language="languageByPath(props.path!)"
    theme="vs-dark"
    height="60vh"
    @before-mount="editorOnBeforeMount"
    :options="{
      automaticLayout: true,
      formatOnType: true,
      formatOnPaste: true,
      wordWrap: 'on'
    }"
  />
</template>

<style scoped lang="scss"></style>
