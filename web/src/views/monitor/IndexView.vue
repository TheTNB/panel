<script setup lang="ts">
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent
} from 'echarts/components'
import monitor from '@/api/panel/monitor'
import type { MonitorData } from '@/views/monitor/types'
import { NButton } from 'naive-ui'

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent
])

const data = ref<MonitorData>({
  times: [],
  load: {
    load1: [],
    load5: [],
    load15: []
  },
  cpu: {
    percent: []
  },
  mem: {
    total: '',
    used: [],
    available: []
  },
  swap: {
    total: '',
    used: [],
    free: []
  },
  net: {
    sent: [],
    recv: [],
    tx: [],
    rx: []
  }
})

const start = ref(Math.floor(new Date(new Date().setHours(0, 0, 0, 0)).getTime()))
const end = ref(Math.floor(Date.now()))

const monitorSwitch = ref(false)
const saveDay = ref(30)

const load = ref<any>({
  title: {
    text: '负载',
    textAlign: 'left',
    textStyle: {
      fontSize: 20
    }
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    align: 'left',
    data: ['1分钟', '5分钟', '15分钟']
  },
  xAxis: [{ type: 'category', boundaryGap: false, data: data.value.times }],
  yAxis: [
    {
      type: 'value'
    }
  ],
  dataZoom: {
    show: true,
    realtime: true,
    start: 0,
    end: 100
  },
  series: [
    {
      name: '1分钟',
      type: 'line',
      smooth: true,
      data: data.value.load.load1,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    },
    {
      name: '5分钟',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.load.load5,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    },
    {
      name: '15分钟',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.load.load15,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    }
  ]
})

const cpu = ref<any>({
  title: {
    text: 'CPU',
    textAlign: 'left',
    textStyle: {
      fontSize: 20
    }
  },
  tooltip: {
    trigger: 'axis'
  },
  xAxis: [{ type: 'category', boundaryGap: false, data: data.value.times }],
  yAxis: [
    {
      name: '单位 %',
      min: 0,
      max: 100,
      type: 'value',
      axisLabel: {
        formatter: '{value} %'
      }
    }
  ],
  dataZoom: {
    show: true,
    realtime: true,
    start: 0,
    end: 100
  },
  series: [
    {
      name: '使用率',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.cpu.percent,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    }
  ]
})

const mem = ref<any>({
  title: {
    text: '内存',
    textAlign: 'left',
    textStyle: {
      fontSize: 20
    }
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    align: 'left',
    data: ['内存', 'Swap']
  },
  xAxis: [{ type: 'category', boundaryGap: false, data: data.value.times }],
  yAxis: [
    {
      name: '单位 MB',
      min: 0,
      max: data.value.mem.total,
      type: 'value',
      axisLabel: {
        formatter: '{value} M'
      }
    }
  ],
  dataZoom: {
    show: true,
    realtime: true,
    start: 0,
    end: 100
  },
  series: [
    {
      name: '内存',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.mem.used,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    },
    {
      name: 'Swap',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.swap.used,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    }
  ]
})

const net = ref<any>({
  title: {
    text: '网络',
    textAlign: 'left',
    textStyle: {
      fontSize: 20
    }
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    align: 'left',
    data: ['总计出', '总计入', '每秒出', '每秒入']
  },
  xAxis: [{ type: 'category', boundaryGap: false, data: data.value.times }],
  yAxis: [
    {
      name: '单位 Mb',
      type: 'value',
      axisLabel: {
        formatter: '{value} Mb'
      }
    }
  ],
  dataZoom: {
    show: true,
    realtime: true,
    start: 0,
    end: 100
  },
  series: [
    {
      name: '总计出',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.net.sent,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    },
    {
      name: '总计入',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.net.recv,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    },
    {
      name: '每秒出',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.net.tx,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    },
    {
      name: '每秒入',
      type: 'line',
      smooth: true,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      },
      data: data.value.net.rx,
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      },
      markLine: {
        data: [{ type: 'average', name: '平均值' }]
      }
    }
  ]
})

const getData = async () => {
  monitor.list(start.value, end.value).then((res) => {
    data.value = res.data
  })
}

const getSwitchAndDays = async () => {
  monitor.switchAndDays().then((res) => {
    monitorSwitch.value = res.data.switch
    saveDay.value = res.data.days
  })
}

const handleSwitch = async () => {
  monitor.switch(monitorSwitch.value).then(() => {
    window.$message.success('操作成功')
  })
}

const handleDays = async () => {
  monitor.saveDays(saveDay.value).then(() => {
    window.$message.success('操作成功')
  })
}

const handleClear = async () => {
  monitor.clear().then(() => {
    window.$message.success('操作成功')
  })
}

// 监听 data 的变化
watch(data, () => {
  load.value.xAxis[0].data = data.value.times
  load.value.series[0].data = data.value.load.load1
  load.value.series[1].data = data.value.load.load5
  load.value.series[2].data = data.value.load.load15
  cpu.value.xAxis[0].data = data.value.times
  cpu.value.series[0].data = data.value.cpu.percent
  mem.value.xAxis[0].data = data.value.times
  mem.value.yAxis[0].max = data.value.mem.total
  mem.value.series[0].data = data.value.mem.used
  mem.value.series[1].data = data.value.swap.used
  net.value.xAxis[0].data = data.value.times
  net.value.series[0].data = data.value.net.sent
  net.value.series[1].data = data.value.net.recv
  net.value.series[2].data = data.value.net.tx
  net.value.series[3].data = data.value.net.rx
})

// 监听时间选择的变化
watch([start, end], () => {
  // 开始时间不能大于结束时间
  if (start.value > end.value) {
    window.$message.error('开始时间不能大于结束时间')
    return
  }
  getData()
})

onMounted(() => {
  getSwitchAndDays()
  getData()
})
</script>

<template>
  <common-page show-footer>
    <n-card :segmented="true" size="small" flex items-center rounded-10>
      <n-form
        inline
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-grid cols="1 s:1 m:1 l:24 xl:24 2xl:24" item-responsive responsive="screen">
          <n-form-item-gi :span="3" label="开启监控">
            <n-switch v-model:value="monitorSwitch" @update-value="handleSwitch" />
          </n-form-item-gi>
          <n-form-item-gi :span="6" label="保存天数">
            <n-input-number v-model:value="saveDay">
              <template #suffix> 天 </template>
            </n-input-number>
          </n-form-item-gi>
          <n-form-item-gi :span="2">
            <n-button type="primary" @click="handleDays">确定</n-button>
          </n-form-item-gi>
          <n-form-item-gi :span="9" label="时间选择">
            <n-date-picker v-model:value="start" type="datetime" />
            -
            <n-date-picker v-model:value="end" type="datetime" />
          </n-form-item-gi>
          <n-form-item-gi :span="4" label="操作">
            <n-popconfirm @positive-click="handleClear">
              <template #trigger>
                <n-button type="error">
                  <TheIcon :size="18" class="mr-5" icon="material-symbols:delete-outline" />
                  清除监控记录
                </n-button>
              </template>
              确定要清空吗？
            </n-popconfirm>
          </n-form-item-gi>
        </n-grid>
      </n-form>
    </n-card>
    <n-grid cols="1 s:1 m:1 l:2 xl:2 2xl:2" item-responsive responsive="screen" pt-20>
      <n-gi m-10>
        <n-card :segmented="true" rounded-10 style="height: 40vh">
          <v-chart class="chart" :option="load" autoresize />
        </n-card>
      </n-gi>
      <n-gi m-10>
        <n-card :segmented="true" rounded-10 style="height: 40vh">
          <v-chart class="chart" :option="cpu" autoresize />
        </n-card>
      </n-gi>
      <n-gi m-10>
        <n-card :segmented="true" rounded-10 style="height: 40vh">
          <v-chart class="chart" :option="mem" autoresize />
        </n-card>
      </n-gi>
      <n-gi m-10>
        <n-card :segmented="true" rounded-10 style="height: 40vh">
          <v-chart class="chart" :option="net" autoresize />
        </n-card>
      </n-gi>
    </n-grid>
  </common-page>
</template>

<style scoped lang="scss"></style>
