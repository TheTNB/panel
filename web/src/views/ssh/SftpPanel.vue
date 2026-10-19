<script setup lang="ts">
import { darkTheme } from 'naive-ui'
import { useGettext } from 'vue3-gettext'

import ws from '@/api/ws'
import { useConfirm } from '@/components/system/composables/useConfirm'
import { useResponsive } from '@/components/system/composables/useResponsive'
import { darkThemeOverrides } from '@/design/naive-overrides'
import { formatBytes, joinPath } from '@/utils'
import SftpBrowser from '@/views/ssh/SftpBrowser.vue'

const { $gettext } = useGettext()
const { confirmAction } = useConfirm()
const { isMobile } = useResponsive()

// 手机单栏显示，分段切换左右两栏；两栏都保持挂载，切换不丢选中项和路径
const activePane = ref<'left' | 'right'>('left')

const props = defineProps<{
  hosts: { label: string; value: number }[]
}>()

const leftRef = ref<InstanceType<typeof SftpBrowser> | null>(null)
const rightRef = ref<InstanceType<typeof SftpBrowser> | null>(null)

const leftHost = ref(-1)
const leftPath = ref('/')
const rightHost = ref(props.hosts.find((h) => h.value !== -1)?.value ?? -1)
const rightPath = ref('/')

interface TransferItem {
  id: string
  name: string
  srcId: number
  srcPath: string
  dstId: number
  dstPath: string
  status: 'waiting' | 'running' | 'success' | 'error'
  transferred: number
  total: number
  speed: number
  errorMsg: string
  ws: WebSocket | null
  lastTransferred: number
  lastTime: number
}

const transfers = ref<TransferItem[]>([])

// 面板本机在终端界面约定为 -1,API 侧用 0 表示
const toApi = (id: number) => (id === -1 ? 0 : id)

const hostLabel = (id: number) => props.hosts.find((h) => h.value === id)?.label || `#${id}`

const handleTransfer = async (direction: 'ltr' | 'rtl') => {
  const from =
    direction === 'ltr'
      ? { browser: leftRef.value, host: leftHost.value, path: leftPath.value }
      : { browser: rightRef.value, host: rightHost.value, path: rightPath.value }
  const to =
    direction === 'ltr'
      ? { host: rightHost.value, path: rightPath.value }
      : { host: leftHost.value, path: leftPath.value }

  const files = from.browser?.selectedFiles || []
  if (!files.length) {
    window.$message.warning($gettext('Please select files to transfer'))
    return
  }
  if (from.host === to.host && from.path === to.path) {
    window.$message.warning($gettext('Source and destination are the same directory'))
    return
  }

  // 目录递归传输可能耗时较长,先二次确认
  const dirs = files.filter((f) => f.is_dir)
  if (dirs.length) {
    const ok = await confirmAction({
      title: $gettext('Transfer Directories'),
      content: $gettext(
        'The selection contains %{ count } directory(ies), which will be transferred recursively and may take a long time. For large directories, compressing before transferring is recommended. Continue?',
        { count: dirs.length },
      ),
    })
    if (!ok) return
  }

  for (const file of files) {
    transfers.value.push({
      id: `transfer-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`,
      name: file.name,
      srcId: from.host,
      srcPath: joinPath(from.path, file.name),
      dstId: to.host,
      dstPath: joinPath(to.path, file.name),
      status: 'waiting',
      transferred: 0,
      // 目录总量未知,由后端进度回报填充
      total: file.is_dir ? 0 : file.size,
      speed: 0,
      errorMsg: '',
      ws: null,
      lastTransferred: 0,
      lastTime: 0,
    })
  }
  processQueue()
}

const processQueue = () => {
  if (transfers.value.some((t) => t.status === 'running')) return
  const next = transfers.value.find((t) => t.status === 'waiting')
  if (next) startTransfer(next)
}

const startTransfer = async (item: TransferItem) => {
  item.status = 'running'
  item.lastTime = performance.now()
  try {
    const socket = await ws.sshTransfer({
      src_id: toApi(item.srcId),
      src_path: item.srcPath,
      dst_id: toApi(item.dstId),
      dst_path: item.dstPath,
    })
    item.ws = socket
    socket.onmessage = (ev) => {
      let msg: any
      try {
        msg = JSON.parse(ev.data)
      } catch {
        return
      }
      if (msg.status === 'progress') {
        const now = performance.now()
        const dt = (now - item.lastTime) / 1000
        if (dt > 0) item.speed = (msg.transferred - item.lastTransferred) / dt
        item.lastTransferred = msg.transferred
        item.lastTime = now
        item.transferred = msg.transferred
        item.total = msg.total
      } else if (msg.status === 'success') {
        item.status = 'success'
        item.transferred = item.total
        finishTransfer(item)
      } else if (msg.status === 'error') {
        item.status = 'error'
        item.errorMsg = msg.msg
        finishTransfer(item)
      }
    }
    socket.onclose = () => {
      if (item.status === 'running') {
        item.status = 'error'
        item.errorMsg = $gettext('Connection closed')
        finishTransfer(item)
      }
    }
  } catch {
    item.status = 'error'
    item.errorMsg = $gettext('Failed to connect')
    processQueue()
  }
}

const finishTransfer = (item: TransferItem) => {
  item.ws?.close()
  item.ws = null
  if (item.status === 'success') {
    refreshDestination(item)
  }
  processQueue()
}

// 传输完成后刷新正在浏览目标目录的一侧
const refreshDestination = (item: TransferItem) => {
  const dstDir = item.dstPath.slice(0, item.dstPath.lastIndexOf('/')) || '/'
  if (leftHost.value === item.dstId && leftPath.value === dstDir) leftRef.value?.refresh()
  if (rightHost.value === item.dstId && rightPath.value === dstDir) rightRef.value?.refresh()
}

const handleCancel = (item: TransferItem) => {
  if (item.status !== 'waiting' && item.status !== 'running') return
  const socket = item.ws
  item.status = 'error'
  item.errorMsg = $gettext('Cancelled')
  item.ws = null
  socket?.close()
  processQueue()
}

const clearFinished = () => {
  transfers.value = transfers.value.filter((t) => t.status === 'waiting' || t.status === 'running')
}

const percent = (item: TransferItem) => {
  if (item.total <= 0) return item.status === 'success' ? 100 : 0
  return Math.min(100, Math.round((item.transferred / item.total) * 100))
}

onUnmounted(() => {
  transfers.value.forEach((item) => {
    if (item.ws) {
      item.ws.onmessage = null
      item.ws.onclose = null
      item.ws.close()
    }
  })
})
</script>

<template>
  <!-- 终端环境固定暗色,不随面板主题切换 -->
  <n-config-provider :theme="darkTheme" :theme-overrides="darkThemeOverrides" abstract>
    <div class="sftp-panel" :class="{ single: isMobile }">
      <n-tabs v-if="isMobile" v-model:value="activePane" type="segment" size="small">
        <n-tab name="left" :tab="$gettext('Left Pane')" />
        <n-tab name="right" :tab="$gettext('Right Pane')" />
      </n-tabs>
      <div class="sftp-panes">
        <SftpBrowser
          v-show="!isMobile || activePane === 'left'"
          ref="leftRef"
          v-model:host-id="leftHost"
          v-model:path="leftPath"
          :hosts="props.hosts"
        />
        <div class="sftp-actions">
          <n-button
            secondary
            :circle="!isMobile"
            :title="$gettext('Transfer to right')"
            @click="handleTransfer('ltr')"
          >
            <template #icon>
              <i-mdi-arrow-right />
            </template>
            <template v-if="isMobile">{{ $gettext('Transfer to right') }}</template>
          </n-button>
          <n-button
            secondary
            :circle="!isMobile"
            :title="$gettext('Transfer to left')"
            @click="handleTransfer('rtl')"
          >
            <template #icon>
              <i-mdi-arrow-left />
            </template>
            <template v-if="isMobile">{{ $gettext('Transfer to left') }}</template>
          </n-button>
        </div>
        <SftpBrowser
          v-show="!isMobile || activePane === 'right'"
          ref="rightRef"
          v-model:host-id="rightHost"
          v-model:path="rightPath"
          :hosts="props.hosts"
        />
      </div>

      <div v-if="transfers.length" class="sftp-queue">
        <div class="sftp-queue-header">
          <span>{{ $gettext('Transfer Queue') }} ({{ transfers.length }})</span>
          <n-button text size="tiny" @click="clearFinished">
            {{ $gettext('Clear Finished') }}
          </n-button>
        </div>
        <div class="sftp-queue-list">
          <div v-for="item in transfers" :key="item.id" class="sftp-queue-item">
            <span class="queue-icon">
              <i-mdi-check-circle v-if="item.status === 'success'" class="text-green-500" />
              <i-mdi-close-circle v-else-if="item.status === 'error'" class="text-red-500" />
              <i-mdi-progress-clock v-else-if="item.status === 'waiting'" />
              <i-mdi-swap-horizontal v-else />
            </span>
            <span class="queue-name" :title="`${item.srcPath} -> ${item.dstPath}`">
              {{ item.name }}
            </span>
            <span class="queue-route">
              {{ hostLabel(item.srcId) }}
              <i-mdi-arrow-right class="text-xs inline" />
              {{ hostLabel(item.dstId) }}
            </span>
            <template v-if="item.status === 'running'">
              <n-progress
                type="line"
                :percentage="percent(item)"
                :show-indicator="false"
                class="queue-progress"
              />
              <span class="queue-meta">
                {{ percent(item) }}% · {{ formatBytes(item.speed) }}/s
              </span>
            </template>
            <span
              v-else-if="item.status === 'error'"
              class="queue-meta error"
              :title="item.errorMsg"
            >
              {{ item.errorMsg }}
            </span>
            <span v-else-if="item.status === 'success'" class="queue-meta">
              {{ formatBytes(item.total) }}
            </span>
            <span v-else class="queue-meta">{{ $gettext('Waiting') }}</span>
            <n-button
              v-if="item.status === 'waiting' || item.status === 'running'"
              text
              size="tiny"
              :title="$gettext('Cancel')"
              @click="handleCancel(item)"
            >
              <template #icon>
                <i-mdi-close />
              </template>
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </n-config-provider>
</template>

<style scoped lang="scss">
.sftp-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
  gap: 10px;
  background: transparent;

  /* 终端环境固定暗色变量,值取自 variables.scss 的 html.dark */
  --color-text-primary: #f8fafc;
  --color-text-secondary: #cbd5e1;
  --color-bg-elevated: #151b27;
  --color-bg-subtle: #1e2533;
  --color-border-default: #1f2937;
  --color-brand-subtle: rgba(54, 173, 106, 0.14);
  --color-error-fg: #f87171;
}

.sftp-panes {
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: stretch;
  gap: 10px;

  > .sftp-browser {
    flex: 1;
  }
}

.sftp-actions {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
  flex-shrink: 0;
}

.sftp-queue {
  flex-shrink: 0;
  max-height: 180px;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--color-border-default);
  border-radius: var(--radius-md);
  background: var(--color-bg-elevated);
  overflow: hidden;
}

.sftp-queue-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  border-bottom: 1px solid var(--color-border-default);
  font-size: 12px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.sftp-queue-list {
  overflow-y: auto;
  padding: 4px;
}

.sftp-queue-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 8px;
  font-size: 12px;
  color: var(--color-text-primary);
}

.queue-icon {
  display: inline-flex;
  align-items: center;
  flex-shrink: 0;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.queue-name {
  flex-shrink: 0;
  max-width: 180px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.queue-route {
  flex-shrink: 0;
  color: var(--color-text-secondary);
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.queue-progress {
  flex: 1;
  min-width: 60px;
}

.queue-meta {
  margin-left: auto;
  flex-shrink: 0;
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
  max-width: 260px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

  &.error {
    color: var(--color-error-fg, #ef4444);
  }
}

// 单栏（手机）：传输按钮横排提到栏上方；队列项的源与目标换到第二行，避免横向溢出
.sftp-panel.single {
  .sftp-panes {
    flex-direction: column;
  }

  .sftp-actions {
    order: -1;
    flex-direction: row;

    > .n-button {
      flex: 1;
    }
  }

  .sftp-queue-item {
    flex-wrap: wrap;
    row-gap: 2px;
  }

  .queue-name {
    flex: 1;
    min-width: 0;
    max-width: none;
  }

  .queue-route {
    order: 1;
    width: 100%;
  }

  .queue-meta {
    max-width: 40%;
  }
}
</style>
