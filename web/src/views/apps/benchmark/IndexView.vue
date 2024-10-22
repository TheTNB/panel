<script setup lang="ts">
import benchmark from '@/api/apps/benchmark'
import TheIcon from '@/components/custom/TheIcon.vue'

defineOptions({
  name: 'apps-benchmark-index'
})

const inTest = ref(false)
const current = ref('CPU')
const progress = ref(0)

const tests = [
  'image',
  'machine',
  'compile',
  'encryption',
  'compression',
  'physics',
  'json',
  'memory',
  'disk'
]

const cpu = ref({
  image: {
    single: 0,
    multi: 0
  },
  machine: {
    single: 0,
    multi: 0
  },
  compile: {
    single: 0,
    multi: 0
  },
  encryption: {
    single: 0,
    multi: 0
  },
  compression: {
    single: 0,
    multi: 0
  },
  physics: {
    single: 0,
    multi: 0
  },
  json: {
    single: 0,
    multi: 0
  }
})

const cpuTotal = computed(() => {
  return {
    single: Object.values(cpu.value).reduce((a, b) => a + b.single, 0),
    multi: Object.values(cpu.value).reduce((a, b) => a + b.multi, 0)
  }
})

const memory = ref({
  score: 0,
  bandwidth: '待跑分',
  latency: '待跑分'
})

const disk = ref({
  score: 0,
  1024: {
    read_iops: '待跑分',
    read_speed: '待跑分',
    write_iops: '待跑分',
    write_speed: '待跑分'
  },
  4: {
    read_iops: '待跑分',
    read_speed: '待跑分',
    write_iops: '待跑分',
    write_speed: '待跑分'
  },
  512: {
    read_iops: '待跑分',
    read_speed: '待跑分',
    write_iops: '待跑分',
    write_speed: '待跑分'
  },
  64: {
    read_iops: '待跑分',
    read_speed: '待跑分',
    write_iops: '待跑分',
    write_speed: '待跑分'
  }
})

const handleTest = async () => {
  inTest.value = true
  progress.value = 0
  for (let i = 0; i < tests.length; i++) {
    const test = tests[i]
    current.value = test
    if (test != 'memory' && test != 'disk') {
      for (let j = 0; j < 2; j++) {
        const { data } = await benchmark.test(test, j === 1)
        cpu.value[test as keyof typeof cpu.value][j === 1 ? 'multi' : 'single'] = data
      }
    } else {
      const { data } = await benchmark.test(test, false)
      if (test === 'memory') {
        memory.value = data
      } else {
        disk.value = data
      }
    }
    progress.value = Math.round(((i + 1) / tests.length) * 100)
  }
  inTest.value = false
}
</script>

<template>
  <common-page show-footer>
    <n-flex vertical>
      <n-alert type="warning">
        跑分结果仅供参考，受系统资源调度和缓存等因素影响可能与实际性能有所偏差！
      </n-alert>
      <n-alert v-if="inTest" title="跑分中，可能需要较长时间..." type="info">
        当前项目：{{ current }}
      </n-alert>
      <n-progress v-if="inTest" :percentage="progress" color="var(--primary-color)" processing />
    </n-flex>
    <n-flex vertical items-center pt-40>
      <div w-800>
        <n-grid :cols="3">
          <n-gi>
            <n-popover trigger="hover">
              <template #trigger>
                <n-flex vertical items-center>
                  <div v-if="cpuTotal.single !== 0 && cpuTotal.multi !== 0">
                    单核
                    <n-number-animation :from="0" :to="cpuTotal.single" show-separator />
                    / 多核
                    <n-number-animation :from="0" :to="cpuTotal.multi" show-separator />
                  </div>
                  <div v-else>待跑分</div>
                  <n-progress
                    type="circle"
                    :percentage="100"
                    :stroke-width="3"
                    color="var(--primary-color)"
                  >
                    <TheIcon :size="50" icon="bi:cpu" color="var(--primary-color)" />
                  </n-progress>
                  CPU
                </n-flex>
              </template>
              <n-table :single-line="false" striped>
                <tr>
                  <th>图像处理</th>
                  <td>单核 {{ cpu.image.single }} / 多核 {{ cpu.image.multi }}</td>
                </tr>
                <tr>
                  <th>机器学习</th>
                  <td>单核 {{ cpu.machine.single }} / 多核 {{ cpu.machine.multi }}</td>
                </tr>
                <tr>
                  <th>程序编译</th>
                  <td>单核 {{ cpu.compile.single }} / 多核 {{ cpu.compile.multi }}</td>
                </tr>
                <tr>
                  <th>AES 加密</th>
                  <td>单核 {{ cpu.encryption.single }} / 多核 {{ cpu.encryption.multi }}</td>
                </tr>
                <tr>
                  <th>压缩/解压缩</th>
                  <td>单核 {{ cpu.compression.single }} / 多核 {{ cpu.compression.multi }}</td>
                </tr>
                <tr>
                  <th>物理仿真</th>
                  <td>单核 {{ cpu.physics.single }} / 多核 {{ cpu.physics.multi }}</td>
                </tr>
                <tr>
                  <th>JSON 解析</th>
                  <td>单核 {{ cpu.json.single }} / 多核 {{ cpu.json.multi }}</td>
                </tr>
              </n-table>
            </n-popover>
          </n-gi>
          <n-gi>
            <n-popover trigger="hover">
              <template #trigger>
                <n-flex vertical items-center>
                  <div v-if="memory.score !== 0">
                    <n-number-animation :from="0" :to="memory.score" show-separator />
                  </div>
                  <div v-else>待跑分</div>
                  <n-progress
                    type="circle"
                    :percentage="100"
                    :stroke-width="3"
                    color="var(--primary-color)"
                  >
                    <TheIcon :size="50" icon="bi:memory" color="var(--primary-color)" />
                  </n-progress>
                  内存
                </n-flex>
              </template>
              <n-table :single-line="false" striped>
                <tr>
                  <th>内存带宽</th>
                  <td>{{ memory.bandwidth }}</td>
                </tr>
                <tr>
                  <th>内存延迟</th>
                  <td>{{ memory.latency }}</td>
                </tr>
              </n-table>
            </n-popover>
          </n-gi>
          <n-gi>
            <n-popover trigger="hover">
              <template #trigger>
                <n-flex vertical items-center>
                  <div v-if="disk.score !== 0">
                    <n-number-animation :from="0" :to="disk.score" show-separator />
                  </div>
                  <div v-else>待跑分</div>
                  <n-progress
                    type="circle"
                    :percentage="100"
                    :stroke-width="3"
                    color="var(--primary-color)"
                  >
                    <TheIcon :size="50" icon="bi:hdd-stack" color="var(--primary-color)" />
                  </n-progress>
                  硬盘
                </n-flex>
              </template>
              <n-table :single-line="false" striped>
                <tr>
                  <th>4KB 读取</th>
                  <td>速度 {{ disk['4'].read_speed }} / {{ disk['4'].read_iops }} IOPS</td>
                </tr>
                <tr>
                  <th>4KB 写入</th>
                  <td>速度 {{ disk['4'].write_speed }} / {{ disk['4'].write_iops }} IOPS</td>
                </tr>
                <tr>
                  <th>64KB 读取</th>
                  <td>速度 {{ disk['64'].read_speed }} / {{ disk['64'].read_iops }} IOPS</td>
                </tr>
                <tr>
                  <th>64KB 写入</th>
                  <td>速度 {{ disk['64'].write_speed }} / {{ disk['64'].write_iops }} IOPS</td>
                </tr>
                <tr>
                  <th>512KB 读取</th>
                  <td>速度 {{ disk['512'].read_speed }} / {{ disk['512'].read_iops }} IOPS</td>
                </tr>
                <tr>
                  <th>512KB 写入</th>
                  <td>速度 {{ disk['512'].write_speed }} / {{ disk['512'].write_iops }} IOPS</td>
                </tr>
                <tr>
                  <th>1MB 读取</th>
                  <td>速度 {{ disk['1024'].read_speed }} / {{ disk['1024'].read_iops }} IOPS</td>
                </tr>
                <tr>
                  <th>1MB 写入</th>
                  <td>速度 {{ disk['1024'].write_speed }} / {{ disk['1024'].write_iops }} IOPS</td>
                </tr>
              </n-table>
            </n-popover>
          </n-gi>
        </n-grid>
      </div>
      <n-button
        type="primary"
        size="large"
        :disabled="inTest"
        :loading="inTest"
        @click="handleTest"
        mt-40
        w-200
      >
        {{ inTest ? '跑分中...' : '开始跑分' }}
      </n-button>
    </n-flex>
  </common-page>
</template>
