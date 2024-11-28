<script setup lang="ts">
import file from '@/api/panel/file'
import { decodeBase64 } from '@/utils'
import { languageByPath } from '@/utils/file'
import Editor from '@guolao/vue-monaco-editor'

const props = defineProps({
  path: {
    type: String,
    required: true
  },
  readOnly: {
    type: Boolean,
    required: true
  }
})

const disabled = ref(false) // 在出现错误的情况下禁用保存
const data = ref('')

const get = async () => {
  await file
    .content(props.path)
    .then((res) => {
      data.value = decodeBase64(res.data.content)
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
  await file.save(props.path, data.value)
  window.$message.success('保存成功')
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
    :language="languageByPath(props.path)"
    theme="vs-dark"
    height="60vh"
    :options="{
      automaticLayout: true,
      formatOnType: true,
      formatOnPaste: true,
      wordWrap: 'on'
    }"
  />
</template>

<style scoped lang="scss"></style>
