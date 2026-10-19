<script setup lang="ts">
defineOptions({
  name: 'apache-config-tune',
})

import { useGettext } from 'vue3-gettext'

import apache from '@/api/apps/apache'
import { useResponsive } from '@/components/system/composables/useResponsive'
import { apachePreset, type TuneProfile } from '@/utils/tunepreset'

const { $gettext } = useGettext()
const { isMobile } = useResponsive()
const currentTab = ref('mpm')

// MPM 事件模型
const startServers = ref<number | null>(null)
const minSpareThreads = ref<number | null>(null)
const maxSpareThreads = ref<number | null>(null)
const threadsPerChild = ref<number | null>(null)
const maxRequestWorkers = ref<number | null>(null)
const maxConnectionsPerChild = ref<number | null>(null)

// 连接设置
const timeout = ref<number | null>(null)
const keepAlive = ref('')
const maxKeepAliveRequests = ref<number | null>(null)
const keepAliveTimeout = ref<number | null>(null)

const saveLoading = ref(false)
const showPresetModal = ref(false)

const handlePreset = (profile: TuneProfile) => {
  const r = apachePreset(profile)
  startServers.value = r.startServers
  threadsPerChild.value = r.threadsPerChild
  maxRequestWorkers.value = r.maxRequestWorkers
  minSpareThreads.value = r.minSpareThreads
  maxSpareThreads.value = r.maxSpareThreads
  maxConnectionsPerChild.value = r.maxConnectionsPerChild
  timeout.value = r.timeout
  keepAlive.value = r.keepAlive
  maxKeepAliveRequests.value = r.maxKeepAliveRequests
  keepAliveTimeout.value = r.keepAliveTimeout
  window.$message.success($gettext('Generated, review the values and save'))
}

const onOffOptions = [
  { label: 'On', value: 'On' },
  { label: 'Off', value: 'Off' },
]

useRequest(apache.configTune()).onSuccess(({ data }: any) => {
  startServers.value = Number(data.start_servers) || null
  minSpareThreads.value = Number(data.min_spare_threads) || null
  maxSpareThreads.value = Number(data.max_spare_threads) || null
  threadsPerChild.value = Number(data.threads_per_child) || null
  maxRequestWorkers.value = Number(data.max_request_workers) || null
  maxConnectionsPerChild.value = data.max_connections_per_child
    ? Number(data.max_connections_per_child)
    : null
  timeout.value = Number(data.timeout) || null
  keepAlive.value = data.keep_alive || null
  maxKeepAliveRequests.value = data.max_keep_alive_requests
    ? Number(data.max_keep_alive_requests)
    : null
  keepAliveTimeout.value = Number(data.keep_alive_timeout) || null
})

const getConfigData = () => ({
  start_servers: String(startServers.value ?? ''),
  min_spare_threads: String(minSpareThreads.value ?? ''),
  max_spare_threads: String(maxSpareThreads.value ?? ''),
  threads_per_child: String(threadsPerChild.value ?? ''),
  max_request_workers: String(maxRequestWorkers.value ?? ''),
  max_connections_per_child: String(maxConnectionsPerChild.value ?? ''),
  timeout: String(timeout.value ?? ''),
  keep_alive: keepAlive.value ?? '',
  max_keep_alive_requests: String(maxKeepAliveRequests.value ?? ''),
  keep_alive_timeout: String(keepAliveTimeout.value ?? ''),
})

const handleSave = () => {
  saveLoading.value = true
  useRequest(apache.saveConfigTune(getConfigData()))
    .onSuccess(() => {
      window.$message.success($gettext('Saved successfully'))
    })
    .onComplete(() => {
      saveLoading.value = false
    })
}
</script>

<template>
  <n-flex vertical>
    <n-tabs v-model:value="currentTab" type="line" :placement="isMobile ? 'top' : 'left'" animated>
      <n-tab-pane name="mpm" :tab="$gettext('MPM Event')">
        <n-flex vertical>
          <n-alert type="info">
            {{ $gettext('Worker thread pool settings for the event MPM.') }}
          </n-alert>
          <n-form>
            <n-form-item :label="$gettext('Start Servers (StartServers)')">
              <n-input-number
                class="w-full"
                v-model:value="startServers"
                :placeholder="$gettext('e.g. 3')"
                :min="1"
              />
            </n-form-item>
            <n-form-item :label="$gettext('Min Spare Threads (MinSpareThreads)')">
              <n-input-number
                class="w-full"
                v-model:value="minSpareThreads"
                :placeholder="$gettext('e.g. 75')"
                :min="1"
              />
            </n-form-item>
            <n-form-item :label="$gettext('Max Spare Threads (MaxSpareThreads)')">
              <n-input-number
                class="w-full"
                v-model:value="maxSpareThreads"
                :placeholder="$gettext('e.g. 250')"
                :min="1"
              />
            </n-form-item>
            <n-form-item :label="$gettext('Threads Per Child (ThreadsPerChild)')">
              <n-input-number
                class="w-full"
                v-model:value="threadsPerChild"
                :placeholder="$gettext('e.g. 25')"
                :min="1"
              />
            </n-form-item>
            <n-form-item :label="$gettext('Max Request Workers (MaxRequestWorkers)')">
              <n-input-number
                class="w-full"
                v-model:value="maxRequestWorkers"
                :placeholder="$gettext('e.g. 400')"
                :min="1"
              />
            </n-form-item>
            <n-form-item :label="$gettext('Max Connections Per Child (MaxConnectionsPerChild)')">
              <n-input-number
                class="w-full"
                v-model:value="maxConnectionsPerChild"
                :placeholder="$gettext('0 means unlimited')"
                :min="0"
              />
            </n-form-item>
          </n-form>
          <n-flex>
            <n-button
              type="primary"
              :loading="saveLoading"
              :disabled="saveLoading"
              @click="handleSave"
            >
              {{ $gettext('Save') }}
            </n-button>
            <n-button type="info" @click="showPresetModal = true">
              {{ $gettext('Generate Recommended Configuration') }}
            </n-button>
          </n-flex>
        </n-flex>
      </n-tab-pane>
      <n-tab-pane name="connection" :tab="$gettext('Connection')">
        <n-flex vertical>
          <n-alert type="info">
            {{ $gettext('Connection and keep-alive settings.') }}
          </n-alert>
          <n-form>
            <n-form-item :label="$gettext('Timeout (Timeout)')">
              <n-input-number
                class="w-full"
                v-model:value="timeout"
                :placeholder="$gettext('e.g. 60')"
                :min="1"
              />
            </n-form-item>
            <n-form-item :label="$gettext('Keep Alive (KeepAlive)')">
              <n-select v-model:value="keepAlive" :options="onOffOptions" clearable />
            </n-form-item>
            <n-form-item :label="$gettext('Max Keep Alive Requests (MaxKeepAliveRequests)')">
              <n-input-number
                class="w-full"
                v-model:value="maxKeepAliveRequests"
                :placeholder="$gettext('0 means unlimited')"
                :min="0"
              />
            </n-form-item>
            <n-form-item :label="$gettext('Keep Alive Timeout (KeepAliveTimeout)')">
              <n-input-number
                class="w-full"
                v-model:value="keepAliveTimeout"
                :placeholder="$gettext('e.g. 5')"
                :min="1"
              />
            </n-form-item>
          </n-form>
          <n-flex>
            <n-button
              type="primary"
              :loading="saveLoading"
              :disabled="saveLoading"
              @click="handleSave"
            >
              {{ $gettext('Save') }}
            </n-button>
            <n-button type="info" @click="showPresetModal = true">
              {{ $gettext('Generate Recommended Configuration') }}
            </n-button>
          </n-flex>
        </n-flex>
      </n-tab-pane>
    </n-tabs>
    <tune-preset-modal
      v-model:show="showPresetModal"
      :fields="['memory', 'cpu']"
      :memory-ratio="0.5"
      @generate="handlePreset"
    />
  </n-flex>
</template>
