<script lang="ts" setup>
import info from '@/api/panel/info'
import type { CountInfo, HomePlugin, NowMonitor, SystemInfo } from './types'
import { router } from '@/router'
import { NButton, NPopconfirm } from 'naive-ui'
import { useAppStore } from '@/store'
import { useI18n } from 'vue-i18n'
import { formatBytes, formatPercent } from '@/utils/file'

const { t, locale } = useI18n()
const appStore = useAppStore()
const nowMonitor = ref<NowMonitor | null>(null)
const systemInfo = ref<SystemInfo | null>(null)
const homePlugins = ref<HomePlugin[] | null>(null)
const countInfo = ref<CountInfo>({
  website: 0,
  database: 0,
  ftp: 0,
  cron: 0
})

const cores = ref(0)
const netTotalSent = ref(0)
const netTotalRecv = ref(0)
const netCurrentSent = ref(0)
const netCurrentRecv = ref(0)
const diskTotalRead = ref(0)
const diskTotalWrite = ref(0)
const diskCurrentRead = ref(0)
const diskCurrentWrite = ref(0)

const getNowMonitor = async () => {
  info.nowMonitor().then((res) => {
    res.data.percent[0] = formatPercent(res.data.percent[0])
    res.data.mem.usedPercent = formatPercent(res.data.mem.usedPercent)
    // 计算 CPU 核心数
    if (cores.value == 0) {
      for (let i = 0; i < res.data.cpus.length; i++) {
        cores.value += res.data.cpus[i].cores
      }
    }
    // 计算网络流量
    let netTotalSentTemp = 0
    let netTotalRecvTemp = 0
    let netTotalSentOld = netTotalSent.value
    let netTotalRecvOld = netTotalRecv.value
    for (let i = 0; i < res.data.net.length; i++) {
      if (res.data.net[i].name === 'lo') {
        continue
      }
      netTotalSentTemp += res.data.net[i].bytesSent
      netTotalRecvTemp += res.data.net[i].bytesRecv
    }
    netTotalSent.value = netTotalSentTemp
    netTotalRecv.value = netTotalRecvTemp
    netCurrentSent.value = (netTotalSent.value - netTotalSentOld) / 3
    netCurrentRecv.value = (netTotalRecv.value - netTotalRecvOld) / 3
    // 计算磁盘读写
    let diskTotalReadTemp = 0
    let diskTotalWriteTemp = 0
    let diskTotalReadOld = diskTotalRead.value
    let diskTotalWriteOld = diskTotalWrite.value
    for (let i = 0; i < res.data.disk_io.length; i++) {
      diskTotalReadTemp += res.data.disk_io[i].readBytes
      diskTotalWriteTemp += res.data.disk_io[i].writeBytes
    }
    diskTotalRead.value = diskTotalReadTemp
    diskTotalWrite.value = diskTotalWriteTemp
    diskCurrentRead.value = (diskTotalRead.value - diskTotalReadOld) / 3
    diskCurrentWrite.value = (diskTotalWrite.value - diskTotalWriteOld) / 3

    nowMonitor.value = res.data
  })
}

const getSystemInfo = async () => {
  info.systemInfo().then((res) => {
    systemInfo.value = res.data
  })
}
const getCountInfo = async () => {
  info.countInfo().then((res) => {
    countInfo.value = res.data
  })
}

const getHomePlugins = async () => {
  info.homePlugins().then((res) => {
    homePlugins.value = res.data
  })
}

const handleRestartPanel = () => {
  clearInterval(homeInterval)
  window.$message.loading(t('homeIndex.system.restart.loading'))
  info.restart().then(() => {
    window.$message.success(t('homeIndex.system.restart.success'))
    setTimeout(() => {
      appStore.reloadPage()
    }, 3000)
  })
}

const handleUpdate = () => {
  info.checkUpdate().then((res) => {
    if (res.data.update) {
      router.push({ name: 'home-update' })
    } else {
      window.$message.success(t('homeIndex.system.update.success'))
    }
  })
}

let eggCount = 0
const getEgg = () => {
  eggCount++
  if (eggCount > 10) {
    return t('homeIndex.eggs.count.gt10')
  } else if (eggCount > 4) {
    return t('homeIndex.eggs.count.gt4')
  } else {
    return t('homeIndex.eggs.count.gt0')
  }
}

const toSponsor = () => {
  if (locale.value === 'en') {
    window.open('https://opencollective.com/tnb')
  } else {
    window.open('https://afdian.com/a/TheTNB')
  }
}

const toGit = () => {
  window.open('https://github.com/TheTNB/panel')
}

const handleManagePlugin = (slug: string) => {
  router.push({ name: 'plugins-' + slug + '-index' })
}

const quantifier = computed(() => {
  return locale.value === 'en' ? '' : ' 个'
})

let homeInterval: any = null

onMounted(() => {
  getNowMonitor()
  getSystemInfo()
  getCountInfo()
  getHomePlugins()
  homeInterval = setInterval(() => {
    getNowMonitor()
  }, 3000)
})

onUnmounted(() => {
  clearInterval(homeInterval)
})
</script>

<template>
  <AppPage :show-footer="true" min-w-375>
    <div flex-1>
      <n-space vertical>
        <div>
          <n-card :segmented="true" rounded-10 size="small">
            <n-page-header :subtitle="systemInfo?.panel_version">
              <n-grid :cols="4">
                <n-gi>
                  <n-statistic
                    :label="$t('homeIndex.website')"
                    :value="countInfo.website + quantifier"
                  />
                </n-gi>
                <n-gi>
                  <n-statistic
                    :label="$t('homeIndex.database')"
                    :value="countInfo.database + quantifier"
                  />
                </n-gi>
                <n-gi>
                  <n-statistic label="FTP" :value="countInfo.ftp + quantifier" />
                </n-gi>
                <n-gi>
                  <n-statistic :label="$t('homeIndex.cron')" :value="countInfo.cron + quantifier" />
                </n-gi>
              </n-grid>
              <template #title>{{ $t('name') }}</template>
              <template #extra>
                <n-space>
                  <n-button type="primary" @click="toSponsor">
                    {{ $t('homeIndex.sponsor') }}
                  </n-button>
                  <n-button @click="toGit">{{ $t('homeIndex.git') }}</n-button>
                </n-space>
              </template>
            </n-page-header>
          </n-card>
        </div>
        <n-grid
          x-gap="12"
          y-gap="12"
          cols="1 s:1 m:1 l:3 xl:3 2xl:3"
          item-responsive
          responsive="screen"
        >
          <n-gi>
            <n-card
              :segmented="true"
              rounded-10
              size="small"
              :title="$t('homeIndex.resources.title')"
            >
              <n-space v-if="nowMonitor" vertical :size="30">
                <n-thing>
                  <template #avatar>
                    <n-avatar>
                      <n-icon>
                        <icon-mdi:cpu-64-bit />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #header> CPU</template>
                  <template #description>
                    <n-progress
                      type="line"
                      :percentage="nowMonitor.percent[0]"
                      :indicator-placement="'inside'"
                    />
                  </template>
                  <p>
                    {{
                      $t('homeIndex.resources.cpu.used', {
                        used: nowMonitor.cpus.length,
                        total: cores
                      })
                    }}
                  </p>
                  <p>{{ nowMonitor.cpus[0].modelName }}</p>
                </n-thing>
                <n-thing v-if="nowMonitor">
                  <template #avatar>
                    <n-avatar>
                      <n-icon>
                        <icon-mdi:memory />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #header> {{ $t('homeIndex.resources.memory.title') }}</template>
                  <template #description>
                    <n-progress
                      type="line"
                      status="info"
                      :percentage="nowMonitor.mem.usedPercent"
                      :indicator-placement="'inside'"
                    />
                  </template>
                  <p>
                    {{
                      $t('homeIndex.resources.memory.physical.used', {
                        used: formatBytes(nowMonitor.mem.used),
                        total: formatBytes(nowMonitor.mem.total)
                      })
                    }}
                  </p>
                  <p>
                    {{
                      $t('homeIndex.resources.memory.swap.used', {
                        used: formatBytes(nowMonitor.swap.used),
                        total: formatBytes(nowMonitor.swap.total)
                      })
                    }}
                  </p>
                </n-thing>
              </n-space>
              <n-skeleton v-else text :repeat="10" />
            </n-card>
          </n-gi>
          <n-gi>
            <n-card :segmented="true" rounded-10 size="small" :title="$t('homeIndex.loads.title')">
              <n-space v-if="nowMonitor" vertical size="large">
                <n-thing>
                  <template #avatar>
                    <n-avatar>
                      <n-icon>
                        <icon-mdi:gauge-empty />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #header>
                    {{ $t('homeIndex.loads.time', { time: '1' }) }}
                  </template>
                  <n-popover trigger="hover" placement="top-end">
                    <template #trigger>
                      <n-progress
                        type="line"
                        :percentage="formatPercent((nowMonitor.load.load1 / cores) * 100)"
                        :indicator-placement="'inside'"
                      />
                    </template>
                    <span>
                      {{ $t('homeIndex.loads.load', { load: '1' }) }}
                      <n-tag type="primary">{{ nowMonitor.load.load1 }}</n-tag>
                    </span>
                  </n-popover>
                </n-thing>
                <n-thing>
                  <template #avatar>
                    <n-avatar>
                      <n-icon>
                        <!--系统负载-->
                        <icon-mdi:gauge />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #header>
                    {{ $t('homeIndex.loads.time', { time: '5' }) }}
                  </template>
                  <n-popover trigger="hover" placement="top-end">
                    <template #trigger>
                      <n-progress
                        type="line"
                        :percentage="formatPercent((nowMonitor.load.load5 / cores) * 100)"
                        :indicator-placement="'inside'"
                      />
                    </template>
                    <span>
                      {{ $t('homeIndex.loads.load', { load: '5' }) }}
                      <n-tag type="primary">{{ nowMonitor.load.load5 }}</n-tag>
                    </span>
                  </n-popover>
                </n-thing>
                <n-thing>
                  <template #avatar>
                    <n-avatar>
                      <n-icon>
                        <icon-mdi:gauge-full />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #header>
                    {{ $t('homeIndex.loads.time', { time: '15' }) }}
                  </template>
                  <n-popover trigger="hover" placement="top-end">
                    <template #trigger>
                      <n-progress
                        type="line"
                        :percentage="formatPercent((nowMonitor.load.load15 / cores) * 100)"
                        :indicator-placement="'inside'"
                      />
                    </template>
                    <span>
                      {{ $t('homeIndex.loads.load', { load: '15' }) }}
                      <n-tag type="primary">{{ nowMonitor.load.load15 }}</n-tag>
                    </span>
                  </n-popover>
                </n-thing>
              </n-space>
              <n-skeleton v-else text :repeat="10" />
            </n-card>
          </n-gi>
          <n-gi>
            <n-card
              :segmented="true"
              rounded-10
              size="small"
              :title="$t('homeIndex.traffic.title')"
            >
              <n-space v-if="nowMonitor" vertical :size="36">
                <n-thing>
                  <template #avatar>
                    <n-avatar>
                      <n-icon>
                        <icon-mdi:server-network />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #header> {{ $t('homeIndex.traffic.network.title') }}</template>
                  <p>
                    {{
                      $t('homeIndex.traffic.network.current', {
                        sent: formatBytes(netCurrentSent),
                        received: formatBytes(netCurrentRecv)
                      })
                    }}
                  </p>
                  <p>
                    {{
                      $t('homeIndex.traffic.network.total', {
                        sent: formatBytes(netTotalSent),
                        received: formatBytes(netTotalRecv)
                      })
                    }}
                  </p>
                </n-thing>
                <n-thing>
                  <template #avatar>
                    <n-avatar>
                      <n-icon>
                        <icon-mdi:harddisk />
                      </n-icon>
                    </n-avatar>
                  </template>
                  <template #header> {{ $t('homeIndex.traffic.disk.title') }}</template>
                  <p>
                    {{
                      $t('homeIndex.traffic.disk.current', {
                        read: formatBytes(diskCurrentRead),
                        write: formatBytes(diskCurrentWrite)
                      })
                    }}
                  </p>
                  <p>
                    {{
                      $t('homeIndex.traffic.disk.total', {
                        read: formatBytes(diskTotalRead),
                        write: formatBytes(diskTotalWrite)
                      })
                    }}
                  </p>
                </n-thing>
              </n-space>
              <n-skeleton v-else text :repeat="10" />
            </n-card>
          </n-gi>
        </n-grid>
        <n-grid
          x-gap="12"
          y-gap="12"
          cols="1 s:1 m:2 l:3 xl:3 2xl:3"
          item-responsive
          responsive="screen"
        >
          <n-gi span="2 s:1 m:1 l:2">
            <div min-w-375 flex-1>
              <n-card
                :segmented="true"
                rounded-10
                size="small"
                :title="$t('homeIndex.store.title')"
              >
                <n-space v-if="nowMonitor" class="pb-10 pt-10">
                  <div v-for="item in nowMonitor?.disk_usage" :key="item.path">
                    <n-popover trigger="hover">
                      <template #trigger>
                        <n-space vertical class="flex items-center">
                          <p>{{ item.path }}</p>
                          <n-progress :percentage="formatPercent(item.usedPercent)" type="circle">
                          </n-progress>
                          <p>{{ formatBytes(item.used) }} / {{ formatBytes(item.total) }}</p>
                        </n-space>
                      </template>
                      <n-table :single-line="false">
                        <tr>
                          <th>{{ $t('homeIndex.store.columns.path') }}</th>
                          <td>{{ item.path }}</td>
                        </tr>
                        <tr>
                          <th>{{ $t('homeIndex.store.columns.type') }}</th>
                          <td>{{ item.fstype }}</td>
                        </tr>
                        <tr>
                          <th>Inodes {{ $t('homeIndex.store.columns.used') }}</th>
                          <td>{{ formatPercent(item.inodesUsedPercent) }}%</td>
                        </tr>
                        <tr>
                          <th>Inodes {{ $t('homeIndex.store.columns.total') }}</th>
                          <td>{{ item.inodesTotal }}</td>
                        </tr>
                        <tr>
                          <th>Inodes {{ $t('homeIndex.store.columns.used') }}</th>
                          <td>{{ item.inodesUsed }}</td>
                        </tr>
                        <tr>
                          <th>Inodes {{ $t('homeIndex.store.columns.free') }}</th>
                          <td>{{ item.inodesFree }}</td>
                        </tr>
                      </n-table>
                    </n-popover>
                  </div>
                </n-space>
                <n-skeleton v-else text :repeat="9" />
              </n-card>
            </div>
          </n-gi>
          <n-gi>
            <div min-w-375 flex-1>
              <n-card
                :segmented="true"
                rounded-10
                size="small"
                :title="$t('homeIndex.system.title')"
              >
                <n-table :single-line="false">
                  <tr>
                    <th>{{ $t('homeIndex.system.columns.os') }}</th>
                    <td>
                      {{ systemInfo?.os_name || $t('homeIndex.system.columns.loading') }}
                    </td>
                  </tr>
                  <tr>
                    <th>{{ $t('homeIndex.system.columns.panel') }}</th>
                    <td>
                      {{ systemInfo?.panel_version || $t('homeIndex.system.columns.loading') }}
                    </td>
                  </tr>
                  <tr>
                    <th>{{ $t('homeIndex.system.columns.uptime') }}</th>
                    <td>{{ systemInfo?.uptime || $t('homeIndex.system.columns.loading') }} 天</td>
                  </tr>
                  <tr>
                    <th>{{ $t('homeIndex.system.columns.operate') }}</th>
                    <td>
                      <n-space>
                        <n-popconfirm @positive-click="handleRestartPanel">
                          <template #trigger>
                            <n-button type="warning">
                              <n-icon size="20">
                                <icon-mdi:restart />
                              </n-icon>
                              {{ $t('homeIndex.system.restart.label') }}
                            </n-button>
                          </template>
                          {{ $t('homeIndex.system.restart.confirm') }}
                        </n-popconfirm>
                        <n-button type="success" @click="handleUpdate">
                          <n-icon size="20">
                            <icon-mdi:arrow-up-bold-circle-outline />
                          </n-icon>
                          {{ $t('homeIndex.system.update.label') }}
                        </n-button>
                      </n-space>
                    </td>
                  </tr>
                </n-table>
              </n-card>
            </div>
          </n-gi>
        </n-grid>
        <n-grid
          x-gap="12"
          y-gap="12"
          cols="1 s:1 m:2 l:3 xl:3 2xl:3"
          item-responsive
          responsive="screen"
        >
          <n-gi span="2 s:1 m:1 l:2">
            <div min-w-375 flex-1>
              <n-card
                :segmented="true"
                rounded-10
                size="small"
                :title="$t('homeIndex.plugins.title')"
              >
                <n-grid
                  v-if="homePlugins"
                  x-gap="12"
                  y-gap="12"
                  cols="3 s:1 m:2 l:3"
                  item-responsive
                  responsive="screen"
                >
                  <n-gi v-for="item in homePlugins" :key="item.name">
                    <n-card
                      :segmented="true"
                      size="small"
                      cursor-pointer
                      rounded-10
                      hover:card-shadow
                      @click="handleManagePlugin(item.slug)"
                    >
                      <n-space>
                        <n-thing>
                          <template #avatar>
                            <n-avatar class="mt-4">
                              <n-icon>
                                <icon-mdi:package-variant-closed />
                              </n-icon>
                            </n-avatar>
                          </template>
                          <template #header>
                            {{ item.name }}
                          </template>
                          <template #description>
                            {{ item.version }}
                          </template>
                        </n-thing>
                      </n-space>
                    </n-card>
                  </n-gi>
                </n-grid>
                <n-skeleton v-else text :repeat="9" />
              </n-card>
            </div>
          </n-gi>
          <n-gi>
            <div min-w-375 flex-1>
              <n-card
                :segmented="true"
                rounded-10
                size="small"
                :title="$t('homeIndex.about.title')"
              >
                <template #header-extra>
                  <n-popover trigger="hover">
                    <template #trigger>
                      <n-icon size="20">
                        <icon-mdi:about-circle-outline />
                      </n-icon>
                    </template>
                    <span>{{ getEgg() }}</span>
                  </n-popover>
                </template>
                <n-space vertical :size="12">
                  <n-alert type="success">
                    {{ $t('homeIndex.about.tnb') }}
                  </n-alert>
                  <n-alert type="info">
                    <span
                      v-html="
                        $t('homeIndex.about.specialThanks', {
                          supporter: `<a target='_blank' href='https://www.weixiaoduo.com/'>「薇晓朵」<\/a>`
                        })
                      "
                    >
                    </span>
                  </n-alert>
                  <n-image
                    src="https://mirror.ghproxy.com/https://raw.githubusercontent.com/TheTNB/sponsor/main/sponsors.svg"
                    width="100%"
                    preview-disabled
                    lazy
                    @click="toSponsor"
                  />
                </n-space>
              </n-card>
            </div>
          </n-gi>
        </n-grid>
      </n-space>
    </div>
  </AppPage>
</template>
