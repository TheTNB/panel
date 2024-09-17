<script setup lang="ts">
const show = defineModel<boolean>('show', { type: Boolean, required: true })
const file = defineModel<string>('file', { type: String, required: true })
const editor = ref<any>(null)

const handleRefresh = () => {
  editor.value.get()
}

const handleSave = () => {
  editor.value.save()
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="'编辑 - ' + file"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <template #header-extra>
      <n-flex>
        <n-button @click="handleRefresh"> 刷新 </n-button>
        <n-button type="primary" @click="handleSave"> 保存 </n-button>
      </n-flex>
    </template>
    <code-editor ref="editor" :path="file" />
  </n-modal>
</template>

<style scoped lang="scss"></style>
