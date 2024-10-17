<script lang="ts" setup>
import { LineChart } from 'echarts/charts'
import {
  DataZoomComponent,
  GridComponent,
  LegendComponent,
  TitleComponent,
  TooltipComponent
} from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { NButton, NPopconfirm } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import dashboard from '@/api/panel/dashboard'
import { router } from '@/router'
import { useAppStore } from '@/store'
import { formatDateTime, formatDuration, toTimestamp } from '@/utils/common'
import { formatBytes, formatPercent } from '@/utils/file'
import VChart from 'vue-echarts'
import type { CountInfo, HomeApp, Realtime, SystemInfo } from './types'

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent
])

const { locale } = useI18n()
const appStore = useAppStore()
const realtime = ref<Realtime | null>(null)
const systemInfo = ref<SystemInfo | null>(null)
const homeApps = ref<HomeApp[] | null>(null)
const homeAppsLoading = ref(false)
const countInfo = ref<CountInfo>({
  website: 0,
  database: 0,
  ftp: 0,
  cron: 0
})

const nets = ref<Array<string>>([]) // 选择的网卡
const disks = ref<Array<string>>([]) // 选择的硬盘
const chartType = ref('net')
const unitType = ref('KB')
const units = [
  { label: 'B', value: 'B' },
  { label: 'KB', value: 'KB' },
  { label: 'MB', value: 'MB' },
  { label: 'GB', value: 'GB' }
]

const cores = ref(0)
const diskReadBytes = ref<Array<number>>([])
const diskWriteBytes = ref<Array<number>>([])
const netBytesSent = ref<Array<number>>([])
const netBytesRecv = ref<Array<number>>([])
const timeDiskData = ref<Array<string>>([])
const timeNetData = ref<Array<string>>([])
const total = reactive({
  diskReadBytes: 0,
  diskWriteBytes: 0,
  diskRWBytes: 0,
  diskRWTime: 0,
  netBytesSent: 0,
  netBytesRecv: 0
})

const current = reactive({
  diskReadBytes: 0,
  diskWriteBytes: 0,
  diskRWBytes: 0,
  diskRWTime: 0,
  netBytesSent: 0,
  netBytesRecv: 0,
  time: 0
})

const statusColor = (percentage: number) => {
  if (percentage >= 90) {
    return 'var(--error-color)'
  } else if (percentage >= 80) {
    return 'var(--warning-color)'
  } else if (percentage >= 70) {
    return 'var(--info-color)'
  }
  return 'var(--success-color)'
}

const statusText = (percentage: number) => {
  if (percentage >= 90) {
    return '运行堵塞'
  } else if (percentage >= 80) {
    return '运行缓慢'
  } else if (percentage >= 70) {
    return '运行正常'
  }
  return '运行流畅'
}

const chartDisk = computed(() => {
  return {
    title: {
      text: chartType.value == 'net' ? '网络' : '硬盘',
      textAlign: 'left',
      textStyle: {
        fontSize: 20
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      formatter: function (params: any) {
        let res = params[0].name + '<br/>'
        params.forEach(function (item: any) {
          res += `${item.marker} ${item.seriesName}: ${item.value} ${unitType.value}<br/>`
        })
        return res
      }
    },
    legend: {
      align: 'left',
      data: chartType.value == 'net' ? ['发送', '接收'] : ['读取', '写入']
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: timeDiskData.value
    },
    yAxis: {
      name: `单位 ${unitType.value}`,
      type: 'value',
      axisLabel: {
        formatter: `{value} ${unitType.value}`
      }
    },
    series: [
      {
        name: chartType.value == 'net' ? '发送' : '读取',
        type: 'line',
        smooth: true,
        data: chartType.value == 'net' ? netBytesSent.value : diskReadBytes.value,
        markPoint: {
          data: [
            { type: 'max', name: '最大值' },
            { type: 'min', name: '最小值' }
          ]
        },
        markLine: {
          data: [{ type: 'average', name: '平均值' }]
        },
        lineStyle: {
          color: 'rgb(247, 184, 81)'
        },
        itemStyle: {
          color: 'rgb(247, 184, 81)'
        },
        areaStyle: {
          color: 'rgb(247, 184, 81)'
        }
      },
      {
        name: chartType.value == 'net' ? '接收' : '写入',
        type: 'line',
        smooth: true,
        data: chartType.value == 'net' ? netBytesRecv.value : diskWriteBytes.value,
        markPoint: {
          data: [
            { type: 'max', name: '最大值' },
            { type: 'min', name: '最小值' }
          ]
        },
        markLine: {
          data: [{ type: 'average', name: '平均值' }]
        },
        lineStyle: {
          color: 'rgb(82, 169, 255)'
        },
        itemStyle: {
          color: 'rgb(82, 169, 255)'
        },
        areaStyle: {
          color: 'rgb(82, 169, 255)'
        }
      }
    ]
  }
})

let isFetching = false

const fetchCurrent = async () => {
  if (isFetching) return
  isFetching = true
  dashboard
    .current(nets.value, disks.value)
    .then(({ data }) => {
      data.percent = formatPercent(data.percent)
      data.mem.usedPercent = formatPercent(data.mem.usedPercent)
      // 计算 CPU 核心数
      if (cores.value == 0) {
        for (let i = 0; i < data.cpus.length; i++) {
          cores.value += data.cpus[i].cores
        }
      }
      // 计算实时数据
      let time = current.time == 0 ? 3 : toTimestamp(data.time) - current.time
      let netTotalSentTemp = 0
      let netTotalRecvTemp = 0
      for (let i = 0; i < data.net.length; i++) {
        if (data.net[i].name === 'lo') {
          continue
        }
        netTotalSentTemp += data.net[i].bytesSent
        netTotalRecvTemp += data.net[i].bytesRecv
      }
      current.netBytesSent =
        total.netBytesSent != 0 ? (netTotalSentTemp - total.netBytesSent) / time : 0
      current.netBytesRecv =
        total.netBytesRecv != 0 ? (netTotalRecvTemp - total.netBytesRecv) / time : 0
      total.netBytesSent = netTotalSentTemp
      total.netBytesRecv = netTotalRecvTemp
      // 计算硬盘读写
      let diskTotalReadTemp = 0
      let diskTotalWriteTemp = 0
      let diskRWTimeTemp = 0
      for (let i = 0; i < data.disk_io.length; i++) {
        diskTotalReadTemp += data.disk_io[i].readBytes
        diskTotalWriteTemp += data.disk_io[i].writeBytes
        diskRWTimeTemp += data.disk_io[i].readTime + data.disk_io[i].writeTime
      }
      current.diskReadBytes =
        total.diskReadBytes != 0 ? (diskTotalReadTemp - total.diskReadBytes) / time : 0
      current.diskWriteBytes =
        total.diskWriteBytes != 0 ? (diskTotalWriteTemp - total.diskWriteBytes) / time : 0
      current.diskRWBytes =
        total.diskRWBytes != 0
          ? (diskTotalReadTemp + diskTotalWriteTemp - total.diskRWBytes) / time
          : 0
      current.diskRWTime =
        total.diskRWTime != 0 ? Number(((diskRWTimeTemp - total.diskRWTime) / time).toFixed(2)) : 0
      current.time = toTimestamp(data.time)
      total.diskReadBytes = diskTotalReadTemp
      total.diskWriteBytes = diskTotalWriteTemp
      total.diskRWBytes = diskTotalReadTemp + diskTotalWriteTemp
      total.diskRWTime = diskRWTimeTemp

      // 图表数据填充
      netBytesSent.value.push(calculateSize(current.netBytesSent))
      if (netBytesSent.value.length > 10) {
        netBytesSent.value.splice(0, 1)
      }
      netBytesRecv.value.push(calculateSize(current.netBytesRecv))
      if (netBytesRecv.value.length > 10) {
        netBytesRecv.value.splice(0, 1)
      }
      diskReadBytes.value.push(calculateSize(current.diskReadBytes))
      if (diskReadBytes.value.length > 10) {
        diskReadBytes.value.splice(0, 1)
      }
      diskWriteBytes.value.push(calculateSize(current.diskWriteBytes))
      if (diskWriteBytes.value.length > 10) {
        diskWriteBytes.value.splice(0, 1)
      }
      timeDiskData.value.push(formatDateTime(data.time))
      if (timeDiskData.value.length > 10) {
        timeDiskData.value.splice(0, 1)
      }
      timeNetData.value.push(formatDateTime(data.time))
      if (timeNetData.value.length > 10) {
        timeNetData.value.splice(0, 1)
      }

      realtime.value = data
    })
    .finally(() => {
      isFetching = false
    })
}

const fetchSystemInfo = async () => {
  dashboard.systemInfo().then((res) => {
    systemInfo.value = res.data
  })
}

const fetchCountInfo = async () => {
  dashboard.countInfo().then((res) => {
    countInfo.value = res.data
  })
}

const fetchHomeApps = async () => {
  homeAppsLoading.value = true
  dashboard.homeApps().then((res) => {
    homeApps.value = res.data
    homeAppsLoading.value = false
  })
}

const handleRestartPanel = () => {
  clearInterval(homeInterval)
  window.$message.loading('面板重启中...')
  dashboard.restart().then(() => {
    window.$message.success('面板重启成功')
    setTimeout(() => {
      appStore.reloadPage()
    }, 3000)
  })
}

const handleUpdate = () => {
  dashboard.checkUpdate().then((res) => {
    if (res.data.update) {
      router.push({ name: 'home-update' })
    } else {
      window.$message.success('当前已是最新版本')
    }
  })
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

const handleManageApp = (slug: string) => {
  router.push({ name: 'apps-' + slug + '-index' })
}

const calculateSize = (bytes: any) => {
  switch (unitType.value) {
    case 'B':
      return Number(bytes.toFixed(2))
    case 'KB':
      return Number((bytes / 1024).toFixed(2))
    case 'MB':
      return Number((bytes / 1024 / 1024).toFixed(2))
    case 'GB':
      return Number((bytes / 1024 / 1024 / 1024).toFixed(2))
    default:
      return 0
  }
}

const clearCurrent = () => {
  total.netBytesSent = 0
  total.netBytesRecv = 0
  total.diskReadBytes = 0
  total.diskWriteBytes = 0
  total.diskRWBytes = 0
  total.diskRWTime = 0
  netBytesSent.value = []
  netBytesRecv.value = []
  timeNetData.value = []
  diskReadBytes.value = []
  diskWriteBytes.value = []
  timeDiskData.value = []
}

const quantifier = computed(() => {
  return locale.value === 'en' ? '' : ' 个'
})

let homeInterval: any = null

onMounted(() => {
  fetchCurrent()
  fetchSystemInfo()
  fetchCountInfo()
  fetchHomeApps()
  homeInterval = setInterval(() => {
    fetchCurrent()
  }, 3000)
})

onUnmounted(() => {
  clearInterval(homeInterval)
})

if (import.meta.hot) {
  import.meta.hot.accept()
  import.meta.hot.dispose(() => {
    clearInterval(homeInterval)
  })
}
</script>

<template>
  <AppPage :show-footer="true" min-w-375>
    <div flex-1>
      <n-space vertical>
        <n-card :segmented="true" rounded-10 size="small">
          <n-page-header :subtitle="systemInfo?.panel_version">
            <n-grid :cols="4" pb-10>
              <n-gi>
                <n-statistic label="网站" :value="countInfo.website + quantifier" />
              </n-gi>
              <n-gi>
                <n-statistic label="数据库" :value="countInfo.database + quantifier" />
              </n-gi>
              <n-gi>
                <n-statistic label="FTP" :value="countInfo.ftp + quantifier" />
              </n-gi>
              <n-gi>
                <n-statistic label="计划任务" :value="countInfo.cron + quantifier" />
              </n-gi>
            </n-grid>
            <template #title>耗子面板</template>
            <template #extra>
              <n-space>
                <n-button type="primary" @click="toSponsor"> 赞助支持 </n-button>
                <n-button @click="toGit">开源地址</n-button>
              </n-space>
            </template>
          </n-page-header>
        </n-card>

        <n-card :segmented="true" rounded-10 size="small" title="资源总览">
          <n-flex v-if="realtime" size="large">
            <n-popover trigger="hover">
              <template #trigger>
                <n-flex vertical flex items-center p-20 pl-40 pr-40>
                  <p>负载状态</p>
                  <n-progress
                    type="dashboard"
                    :percentage="formatPercent((realtime.load.load1 / cores) * 100)"
                    :color="statusColor((realtime.load.load1 / cores) * 100)"
                  >
                  </n-progress>
                  <p>{{ statusText((realtime.load.load1 / cores) * 100) }}</p>
                </n-flex>
              </template>
              <n-table :single-line="false" striped>
                <tr>
                  <th>最近 1 分钟</th>
                  <td>
                    {{ formatPercent((realtime.load.load1 / cores) * 100) }}% /
                    {{ realtime.load.load1 }}
                  </td>
                </tr>
                <tr>
                  <th>最近 5 分钟</th>
                  <td>
                    {{ formatPercent((realtime.load.load5 / cores) * 100) }}% /
                    {{ realtime.load.load5 }}
                  </td>
                </tr>
                <tr>
                  <th>最近 15 分钟</th>
                  <td>
                    {{ formatPercent((realtime.load.load15 / cores) * 100) }}% /
                    {{ realtime.load.load15 }}
                  </td>
                </tr>
              </n-table>
            </n-popover>
            <n-popover trigger="hover">
              <template #trigger>
                <n-flex vertical flex items-center p-20 pl-40 pr-40>
                  <p>CPU</p>
                  <n-progress
                    type="dashboard"
                    :percentage="realtime.percent"
                    :color="statusColor(realtime.percent)"
                  >
                  </n-progress>
                  <p>{{ cores }} 核心</p>
                </n-flex>
              </template>
              <n-table :single-line="false" striped>
                <tr>
                  <th>型号</th>
                  <td>{{ realtime.cpus[0].modelName }}</td>
                </tr>
                <tr>
                  <th>参数</th>
                  <td>
                    {{ realtime.cpus.length }} CPU {{ cores }} 核心
                    {{ formatBytes(realtime.cpus[0].cacheSize * 1024) }} 缓存
                  </td>
                </tr>
                <tr v-for="item in realtime.cpus" :key="item.modelName">
                  <th>CPU-{{ item.cpu }}</th>
                  <td>
                    使用率 {{ formatPercent(realtime.percents[item.cpu]) }}% 频率 {{ item.mhz }} MHz
                  </td>
                </tr>
              </n-table>
            </n-popover>
            <n-popover trigger="hover">
              <template #trigger>
                <n-flex vertical flex items-center p-20 pl-40 pr-40>
                  <p>内存</p>
                  <n-progress
                    type="dashboard"
                    :percentage="realtime.mem.usedPercent"
                    :color="statusColor(realtime.mem.usedPercent)"
                  >
                  </n-progress>
                  <p>{{ formatBytes(realtime.mem.total) }}</p>
                </n-flex>
              </template>
              <n-table :single-line="false" striped>
                <tr>
                  <th>活跃</th>
                  <td>
                    {{ formatBytes(realtime.mem.active) }}
                  </td>
                </tr>
                <tr>
                  <th>不活跃</th>
                  <td>
                    {{ formatBytes(realtime.mem.inactive) }}
                  </td>
                </tr>
                <tr>
                  <th>空闲</th>
                  <td>
                    {{ formatBytes(realtime.mem.free) }}
                  </td>
                </tr>
                <tr>
                  <th>共享</th>
                  <td>
                    {{ formatBytes(realtime.mem.shared) }}
                  </td>
                </tr>
                <tr>
                  <th>已提交</th>
                  <td>
                    {{ formatBytes(realtime.mem.committedas) }}
                  </td>
                </tr>
                <tr>
                  <th>提交限制</th>
                  <td>
                    {{ formatBytes(realtime.mem.commitlimit) }}
                  </td>
                </tr>
                <tr>
                  <th>SWAP大小</th>
                  <td>
                    {{ formatBytes(realtime.mem.swaptotal) }}
                  </td>
                </tr>
                <tr>
                  <th>SWAP已用</th>
                  <td>
                    {{ formatBytes(realtime.mem.swapcached) }}
                  </td>
                </tr>
                <tr>
                  <th>SWAP可用</th>
                  <td>
                    {{ formatBytes(realtime.mem.swapfree) }}
                  </td>
                </tr>
                <tr>
                  <th>物理内存大小</th>
                  <td>
                    {{ formatBytes(realtime.mem.total) }}
                  </td>
                </tr>
                <tr>
                  <th>物理内存已用</th>
                  <td>
                    {{ formatBytes(realtime.mem.used) }}
                  </td>
                </tr>
                <tr>
                  <th>物理内存可用</th>
                  <td>
                    {{ formatBytes(realtime.mem.available) }}
                  </td>
                </tr>
                <tr>
                  <th>buffers/cached</th>
                  <td>
                    {{ formatBytes(realtime.mem.buffers) }} / {{ formatBytes(realtime.mem.cached) }}
                  </td>
                </tr>
              </n-table>
            </n-popover>
            <n-popover v-for="item in realtime?.disk_usage" :key="item.path" trigger="hover">
              <template #trigger>
                <n-flex vertical flex items-center p-20 pl-40 pr-40>
                  <p>{{ item.path }}</p>
                  <n-progress
                    type="dashboard"
                    :percentage="formatPercent(item.usedPercent)"
                    :color="statusColor(item.usedPercent)"
                  >
                  </n-progress>
                  <p>{{ formatBytes(item.used) }} / {{ formatBytes(item.total) }}</p>
                </n-flex>
              </template>
              <n-table :single-line="false">
                <tr>
                  <th>挂载点</th>
                  <td>{{ item.path }}</td>
                </tr>
                <tr>
                  <th>文件系统</th>
                  <td>{{ item.fstype }}</td>
                </tr>
                <tr>
                  <th>Inodes 使用率</th>
                  <td>{{ formatPercent(item.inodesUsedPercent) }}%</td>
                </tr>
                <tr>
                  <th>Inodes 总数</th>
                  <td>{{ item.inodesTotal }}</td>
                </tr>
                <tr>
                  <th>Inodes 已用</th>
                  <td>{{ item.inodesUsed }}</td>
                </tr>
                <tr>
                  <th>Inodes 可用</th>
                  <td>{{ item.inodesFree }}</td>
                </tr>
              </n-table>
            </n-popover>
          </n-flex>
          <n-skeleton v-else text :repeat="10" />
        </n-card>
        <n-grid
          x-gap="12"
          y-gap="12"
          cols="1 s:1 m:1 l:2 xl:2 2xl:2"
          item-responsive
          responsive="screen"
        >
          <n-gi>
            <n-flex vertical>
              <n-card :segmented="true" size="small" title="快捷应用" min-h-340 rounded-10>
                <n-scrollbar max-h-270>
                  <n-grid
                    v-if="!homeAppsLoading"
                    x-gap="12"
                    y-gap="12"
                    cols="3 s:1 m:2 l:3"
                    item-responsive
                    responsive="screen"
                  >
                    <n-gi v-for="item in homeApps" :key="item.name">
                      <n-card
                        :segmented="true"
                        size="small"
                        cursor-pointer
                        rounded-10
                        hover:card-shadow
                        @click="handleManageApp(item.slug)"
                      >
                        <n-space>
                          <n-thing>
                            <template #avatar>
                              <div class="mt-8">
                                <TheIcon :size="30" :icon="item.icon" />
                              </div>
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
                </n-scrollbar>
                <n-text v-if="!homeAppsLoading && !homeApps">
                  您还没有设置任何应用在此显示！
                </n-text>
                <n-skeleton v-if="homeAppsLoading" text :repeat="12" />
              </n-card>
              <n-card :segmented="true" rounded-10 size="small" title="系统信息">
                <n-table v-if="systemInfo" :single-line="false">
                  <tr>
                    <th>主机名</th>
                    <td>
                      {{ systemInfo?.hostname || '加载中...' }}
                    </td>
                  </tr>
                  <tr>
                    <th>系统版本</th>
                    <td>
                      {{ `${systemInfo?.os_name} ${systemInfo?.kernel_arch}` || '加载中...' }}
                    </td>
                  </tr>
                  <tr>
                    <th>内核版本</th>
                    <td>
                      {{ systemInfo?.kernel_version || '加载中...' }}
                    </td>
                  </tr>
                  <tr>
                    <th>运行时间</th>
                    <td>
                      {{ formatDuration(Number(systemInfo?.uptime)) || '加载中...' }}
                    </td>
                  </tr>
                  <tr>
                    <th>操作</th>
                    <td>
                      <n-space>
                        <n-popconfirm @positive-click="handleRestartPanel">
                          <template #trigger>
                            <n-button type="warning" size="small">
                              <TheIcon :size="20" icon="mdi:restart" />
                              重启面板
                            </n-button>
                          </template>
                          确定要重启面板吗？
                        </n-popconfirm>
                        <n-button type="success" @click="handleUpdate" size="small">
                          <TheIcon :size="20" icon="mdi:arrow-up-bold-circle-outline" />
                          检查更新
                        </n-button>
                      </n-space>
                    </td>
                  </tr>
                </n-table>
                <n-skeleton v-else text :repeat="9" />
              </n-card>
            </n-flex>
          </n-gi>
          <n-gi>
            <n-card :segmented="true" rounded-10 size="small" title="实时监控">
              <n-flex vertical v-if="systemInfo">
                <n-form
                  inline
                  label-placement="left"
                  label-width="auto"
                  require-mark-placement="right-hanging"
                >
                  <n-form-item>
                    <n-radio-group v-model:value="chartType">
                      <n-radio-button value="net" label="网络" />
                      <n-radio-button value="disk" label="硬盘" />
                    </n-radio-group>
                  </n-form-item>
                  <n-form-item label="单位" ml-auto>
                    <n-select
                      v-model:value="unitType"
                      :options="units"
                      @update-value="clearCurrent"
                      w-80
                    ></n-select>
                  </n-form-item>
                  <n-form-item v-if="chartType == 'net'" label="网卡">
                    <n-select
                      multiple
                      v-model:value="nets"
                      :options="systemInfo.nets"
                      @update-value="clearCurrent"
                      w-200
                    ></n-select>
                  </n-form-item>
                  <n-form-item v-if="chartType == 'disk'" label="硬盘">
                    <n-select
                      multiple
                      v-model:value="disks"
                      :options="systemInfo.disks"
                      @update-value="clearCurrent"
                      w-200
                    ></n-select>
                  </n-form-item>
                </n-form>
                <n-flex v-if="chartType == 'net'">
                  <n-tag>总发送 {{ formatBytes(total.netBytesSent) }}</n-tag>
                  <n-tag>总接收 {{ formatBytes(total.netBytesRecv) }}</n-tag>
                  <n-tag>实时发送 {{ formatBytes(current.netBytesSent) }}/s</n-tag>
                  <n-tag>实时接收 {{ formatBytes(current.netBytesRecv) }}/s</n-tag>
                </n-flex>
                <n-flex v-if="chartType == 'disk'">
                  <n-tag>读取 {{ formatBytes(total.diskReadBytes) }}</n-tag>
                  <n-tag>写入 {{ formatBytes(total.diskWriteBytes) }}</n-tag>
                  <n-tag>实时读写 {{ formatBytes(current.diskRWBytes) }}/s</n-tag>
                  <n-tag>读写延迟 {{ current.diskRWTime }}ms</n-tag>
                </n-flex>
                <n-card :bordered="false" h-497>
                  <v-chart class="chart" :option="chartDisk" autoresize />
                </n-card>
              </n-flex>
              <n-skeleton v-else text :repeat="24" />
            </n-card>
          </n-gi>
        </n-grid>
      </n-space>
    </div>
  </AppPage>
</template>
