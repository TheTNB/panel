<script setup lang="ts">
import ws from '@/api/ws'
import type { LogInst } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const props = defineProps({
  path: String
})

const log = ref('')
const logRef = ref<LogInst | null>(null)
let logWs: WebSocket | null = null

const init = async () => {
  const cmd = `tail -n 40 -f '${props.path}'`
  ws.exec(cmd)
    .then((ws: WebSocket) => {
      logWs = ws
      ws.onmessage = (event) => {
        log.value += event.data + '\n'
        const lines = log.value.split('\n')
        if (lines.length > 2000) {
          log.value = lines.slice(lines.length - 2000).join('\n')
        }
      }
    })
    .catch(() => {
      window.$message.error('获取日志流失败')
    })
}

const handleClose = () => {
  if (logWs) {
    logWs.close()
  }
  log.value = ''
}

watch(
  () => props.path,
  () => {
    handleClose()
    init()
  }
)

watchEffect(() => {
  if (log.value) {
    nextTick(() => {
      logRef.value?.scrollTo({ position: 'bottom', silent: true })
    })
  }
})

defineExpose({
  init
})
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="日志"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="handleClose"
    @mask-click="handleClose"
  >
    <n-log ref="logRef" :log="log" trim :rows="40" />
  </n-modal>
</template>

<style scoped lang="scss"></style>
