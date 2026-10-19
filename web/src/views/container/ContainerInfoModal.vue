<script setup lang="ts">
import { useGettext } from 'vue3-gettext'

import container from '@/api/panel/container'
import { useResponsive } from '@/components/system/composables/useResponsive'

const { $gettext } = useGettext()
const { isMobile } = useResponsive()

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const props = defineProps<{
  containerId: string
}>()

const loading = ref(false)
const info = ref<any>(null)
const currentTab = ref('basic')

// 挂载列表
const mounts = computed(() => info.value?.Mounts || [])

// 端口映射（NetworkSettings.Ports: "80/tcp" -> [{HostIp, HostPort}]）
const ports = computed(() => {
  const result: string[] = []
  const portMap = info.value?.NetworkSettings?.Ports || {}
  for (const [containerPort, bindings] of Object.entries(portMap)) {
    if (!bindings) {
      result.push(containerPort)
      continue
    }
    for (const binding of bindings as any[]) {
      result.push(
        `${binding.HostIp ? binding.HostIp + ':' : ''}${binding.HostPort}->${containerPort}`,
      )
    }
  }
  return result
})

// 网络与 IP
const networks = computed(() => {
  const result: { name: string; ip: string }[] = []
  const networkMap = info.value?.NetworkSettings?.Networks || {}
  for (const [name, settings] of Object.entries(networkMap)) {
    result.push({ name, ip: (settings as any).IPAddress || '-' })
  }
  return result
})

// 启动命令
const command = computed(() => {
  const entrypoint = info.value?.Config?.Entrypoint || []
  const cmd = info.value?.Config?.Cmd || []
  return [...entrypoint, ...cmd].join(' ') || '-'
})

// 标签
const labels = computed(() => {
  const labelMap = info.value?.Config?.Labels || {}
  return Object.entries(labelMap).map(([key, value]) => `${key}=${value}`)
})

// 原始 inspect 输出
const rawJson = computed(() => (info.value ? JSON.stringify(info.value, null, 2) : ''))

const formatTime = (time: string) => {
  if (!time || time.startsWith('0001-')) return '-'
  return new Date(time).toLocaleString()
}

watch(show, (val) => {
  if (val) {
    info.value = null
    currentTab.value = 'basic'
    loading.value = true
    useRequest(container.containerInspect(props.containerId))
      .onSuccess(({ data }: any) => {
        info.value = data
      })
      .onComplete(() => {
        loading.value = false
      })
  }
})
</script>

<template>
  <app-modal
    v-model:show="show"
    size="lg"
    :title="
      $gettext('Container Info') + (info ? ' - ' + String(info.Name || '').replace(/^\//, '') : '')
    "
  >
    <n-spin :show="loading">
      <n-tabs v-model:value="currentTab" type="line" animated>
        <!-- 基本信息 -->
        <n-tab-pane name="basic" :tab="$gettext('Basic Info')">
          <n-flex v-if="info" vertical :size="16">
            <n-descriptions :column="isMobile ? 1 : 2" label-placement="left" bordered size="small">
              <n-descriptions-item :label="$gettext('Container Name')">
                {{ String(info.Name || '').replace(/^\//, '') || '-' }}
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Container ID')">
                {{ String(info.Id || '').slice(0, 12) }}
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Image')">
                {{ info.Config?.Image || '-' }}
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Status')">
                <n-tag
                  :type="info.State?.Status === 'running' ? 'success' : 'default'"
                  size="small"
                >
                  {{ info.State?.Status || '-' }}
                </n-tag>
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Created At')">
                {{ formatTime(info.Created) }}
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Started At')">
                {{ formatTime(info.State?.StartedAt) }}
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Restart Policy')">
                {{ info.HostConfig?.RestartPolicy?.Name || 'no' }}
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Privileged')">
                {{ info.HostConfig?.Privileged ? $gettext('Yes') : $gettext('No') }}
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Network')" :span="isMobile ? 1 : 2">
                <n-flex v-if="networks.length > 0" :size="4">
                  <n-tag v-for="net in networks" :key="net.name" size="small">
                    {{ net.name }}: {{ net.ip }}
                  </n-tag>
                </n-flex>
                <template v-else>{{ info.HostConfig?.NetworkMode || '-' }}</template>
              </n-descriptions-item>
              <n-descriptions-item
                :label="$gettext('Ports (Host->Container)')"
                :span="isMobile ? 1 : 2"
              >
                <n-flex v-if="ports.length > 0" :size="4">
                  <n-tag v-for="port in ports" :key="port" size="small">{{ port }}</n-tag>
                </n-flex>
                <template v-else>-</template>
              </n-descriptions-item>
              <n-descriptions-item :label="$gettext('Command')" :span="isMobile ? 1 : 2">
                <span style="font-family: monospace; word-break: break-all">{{ command }}</span>
              </n-descriptions-item>
            </n-descriptions>

            <!-- 挂载 -->
            <n-card :title="$gettext('Mounts')" size="small" embedded>
              <n-table v-if="mounts.length > 0" :bordered="true" :single-line="false" size="small">
                <thead>
                  <tr>
                    <th>{{ $gettext('Type') }}</th>
                    <th>{{ $gettext('Host Path / Volume') }}</th>
                    <th>{{ $gettext('Container Path') }}</th>
                    <th>{{ $gettext('Mode') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(mount, index) in mounts" :key="index">
                    <td>
                      <n-tag size="small">{{ mount.Type }}</n-tag>
                    </td>
                    <td style="font-family: monospace; word-break: break-all">
                      {{ mount.Type === 'volume' ? mount.Name : mount.Source }}
                    </td>
                    <td style="font-family: monospace; word-break: break-all">
                      {{ mount.Destination }}
                    </td>
                    <td>
                      <n-tag :type="mount.RW ? 'success' : 'warning'" size="small">
                        {{ mount.RW ? $gettext('Read-Write') : $gettext('Read-Only') }}
                      </n-tag>
                    </td>
                  </tr>
                </tbody>
              </n-table>
              <n-empty v-else :description="$gettext('No mounts')" size="small" />
            </n-card>

            <!-- 环境变量 -->
            <n-card :title="$gettext('Environment Variables')" size="small" embedded>
              <n-flex v-if="(info.Config?.Env || []).length > 0" vertical :size="2">
                <span
                  v-for="env in info.Config.Env"
                  :key="env"
                  style="font-family: monospace; font-size: 13px; word-break: break-all"
                >
                  {{ env }}
                </span>
              </n-flex>
              <n-empty v-else :description="$gettext('No environment variables')" size="small" />
            </n-card>

            <!-- 标签 -->
            <n-card v-if="labels.length > 0" :title="$gettext('Labels')" size="small" embedded>
              <n-flex vertical :size="2">
                <span
                  v-for="label in labels"
                  :key="label"
                  style="font-family: monospace; font-size: 13px; word-break: break-all"
                >
                  {{ label }}
                </span>
              </n-flex>
            </n-card>
          </n-flex>
        </n-tab-pane>

        <!-- 原始输出 -->
        <n-tab-pane name="raw" :tab="$gettext('Raw Output')">
          <pre
            style="
              max-height: 60vh;
              overflow: auto;
              margin: 0;
              padding: 12px;
              font-family: monospace;
              font-size: 13px;
              line-height: 1.6;
              background: var(--n-color-embedded, rgba(128, 128, 128, 0.08));
              border-radius: 4px;
              white-space: pre-wrap;
              word-break: break-all;
            "
            >{{ rawJson }}</pre>
        </n-tab-pane>
      </n-tabs>
    </n-spin>
  </app-modal>
</template>
