<script setup lang="ts">
import container from '@/api/panel/container'

const props = defineProps({
  show: Boolean
})

const { show } = toRefs(props)
const doSubmit = ref(false)

const createModel = reactive({
  name: '',
  publish_all_ports: false,
  ports: [
    {
      container_start: 80,
      container_end: 80,
      host_start: 80,
      host_end: 80,
      host: '',
      protocol: 'tcp'
    }
  ],
  network: '',
  volumes: [
    {
      host: '/www',
      container: '/www',
      mode: 'rw'
    }
  ],
  cpus: 0,
  memory: 0,
  env: [],
  command: [],
  tty: false,
  restart_policy: 'no',
  labels: [],
  entrypoint: [],
  auto_remove: false,
  image: '',
  cpu_shares: 1024,
  privileged: false,
  open_stdin: false
})
const networks = ref<any>({})

const restartPolicyOptions = [
  { label: '无', value: 'no' },
  { label: '始终', value: 'always' },
  { label: '失败时（默认重启 5 次）', value: 'on-failure' },
  { label: '未手动停止则重启', value: 'unless-stopped' }
]

const addPortRow = () => {
  createModel.ports.push({
    container_start: 80,
    container_end: 80,
    host_start: 80,
    host_end: 80,
    host: '',
    protocol: 'tcp'
  })
}

const removePortRow = (index: number) => {
  createModel.ports.splice(index, 1)
}

const addVolumeRow = () => {
  createModel.volumes.push({
    host: '/www',
    container: '/www',
    mode: 'rw'
  })
}

const removeVolumeRow = (index: number) => {
  createModel.volumes.splice(index, 1)
}

const getNetworks = async () => {
  const { data } = await container.networkList(1, 1000)
  networks.value = data.items.map((item: any) => {
    return {
      label: item.name,
      value: item.id
    }
  })
  if (networks.value.length > 0) {
    createModel.network = networks.value[0].value
  }
}

const handleSubmit = () => {
  doSubmit.value = true
  container
    .containerCreate(createModel)
    .then(() => {
      window.$message.success('创建成功')
      handleClose()
    })
    .catch(() => {
      window.$message.error('创建失败')
    })
    .finally(() => {
      doSubmit.value = false
    })
}

const emit = defineEmits(['close'])

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  getNetworks()
})
</script>

<template>
  <n-modal
    title="创建容器"
    preset="card"
    style="width: 60vw"
    size="huge"
    :show="show"
    :bordered="false"
    :segmented="false"
    @close="handleClose"
  >
    <n-form :model="createModel">
      <n-form-item path="name" label="容器名">
        <n-input v-model:value="createModel.name" type="text" @keydown.enter.prevent />
      </n-form-item>
      <n-form-item path="name" label="镜像">
        <n-input v-model:value="createModel.image" type="text" @keydown.enter.prevent />
      </n-form-item>
      <n-form-item path="exposedAll" label="端口">
        <n-radio
          :checked="!createModel.publish_all_ports"
          :value="false"
          @change="createModel.publish_all_ports = !$event.target.value"
        >
          映射端口
        </n-radio>
        <n-radio
          :checked="createModel.publish_all_ports"
          :value="true"
          @change="createModel.publish_all_ports = !!$event.target.value"
        >
          暴露所有
        </n-radio>
      </n-form-item>
      <n-form-item path="ports" label="端口映射" v-if="!createModel.publish_all_ports">
        <n-space vertical>
          <n-table striped>
            <thead>
              <tr>
                <th>IP</th>
                <th>主机（起始）</th>
                <th>主机（结束）</th>
                <th>容器（起始）</th>
                <th>容器（结束）</th>
                <th>协议</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in createModel.ports" :key="index">
                <td>
                  <n-input
                    v-model:value="item.host"
                    type="text"
                    @keydown.enter.prevent
                    placeholder="可留空"
                  />
                </td>
                <td>
                  <n-input-number
                    v-model:value="item.host_start"
                    type="text"
                    @keydown.enter.prevent
                  />
                </td>
                <td>
                  <n-input-number
                    v-model:value="item.host_end"
                    type="text"
                    @keydown.enter.prevent
                  />
                </td>
                <td>
                  <n-input-number
                    v-model:value="item.container_start"
                    type="text"
                    @keydown.enter.prevent
                  />
                </td>
                <td>
                  <n-input-number
                    v-model:value="item.container_end"
                    type="text"
                    @keydown.enter.prevent
                  />
                </td>
                <td>
                  <n-radio
                    :checked="item.protocol === 'tcp'"
                    value="tcp"
                    name="protocol"
                    @change="item.protocol = $event.target.value"
                  >
                    TCP
                  </n-radio>
                  <n-radio
                    :checked="item.protocol === 'udp'"
                    value="udp"
                    name="protocol"
                    @change="item.protocol = $event.target.value"
                  >
                    UDP
                  </n-radio>
                </td>
                <td><n-button @click="removePortRow(index)" size="small">删除</n-button></td>
              </tr>
            </tbody>
          </n-table>
          <n-button @click="addPortRow">添加</n-button>
        </n-space>
      </n-form-item>
      <n-form-item path="network" label="网络">
        <n-select v-model:value="createModel.network" :options="networks" />
      </n-form-item>
      <n-form-item path="mount" label="挂载">
        <n-space vertical>
          <n-table striped>
            <thead>
              <tr>
                <th>主机目录</th>
                <th>容器目录</th>
                <th>权限</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in createModel.volumes" :key="index">
                <td>
                  <n-input v-model:value="item.host" type="text" @keydown.enter.prevent />
                </td>
                <td>
                  <n-input v-model:value="item.container" type="text" @keydown.enter.prevent />
                </td>
                <td>
                  <n-radio
                    :checked="item.mode === 'rw'"
                    value="rw"
                    name="mode"
                    @change="item.mode = $event.target.value"
                  >
                    读写
                  </n-radio>
                  <n-radio
                    :checked="item.mode === 'ro'"
                    value="ro"
                    name="mode"
                    @change="item.mode = $event.target.value"
                  >
                    只读
                  </n-radio>
                </td>
                <td><n-button @click="removeVolumeRow(index)" size="small">删除</n-button></td>
              </tr>
            </tbody>
          </n-table>
          <n-button @click="addVolumeRow">添加</n-button>
        </n-space>
      </n-form-item>
      <n-form-item path="command" label="启动命令">
        <n-dynamic-input v-model:value="createModel.command" placeholder="命令" />
      </n-form-item>
      <n-form-item path="entrypoint" label="入口点">
        <n-dynamic-input v-model:value="createModel.entrypoint" placeholder="入口点" />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="8">
          <n-form-item path="memory" label="内存">
            <n-input-number v-model:value="createModel.memory" />
          </n-form-item>
        </n-col>
        <n-col :span="8">
          <n-form-item path="cpus" label="CPU">
            <n-input-number v-model:value="createModel.cpus" />
          </n-form-item>
        </n-col>
        <n-col :span="8">
          <n-form-item path="cpu_shares" label="CPU 权重">
            <n-input-number v-model:value="createModel.cpu_shares" />
          </n-form-item>
        </n-col>
      </n-row>
      <n-row :gutter="[0, 24]">
        <n-col :span="6">
          <n-form-item path="tty" label="伪终端（-t）">
            <n-switch v-model:value="createModel.tty" />
          </n-form-item>
        </n-col>
        <n-col :span="6">
          <n-form-item path="open_stdin" label="标准输入（-i）">
            <n-switch v-model:value="createModel.open_stdin" />
          </n-form-item>
        </n-col>
        <n-col :span="6">
          <n-form-item path="auto_remove" label="退出后自动删除">
            <n-switch v-model:value="createModel.auto_remove" />
          </n-form-item>
        </n-col>
        <n-col :span="6">
          <n-form-item path="privileged" label="特权模式">
            <n-switch v-model:value="createModel.privileged" />
          </n-form-item>
        </n-col>
      </n-row>
      <n-form-item path="restart_policy" label="重启策略">
        <n-select
          v-model:value="createModel.restart_policy"
          placeholder="选择重启策略"
          :options="restartPolicyOptions"
        >
          {{ createModel.restart_policy || '选择重启策略' }}
        </n-select>
      </n-form-item>
      <n-form-item path="env" label="环境变量">
        <n-dynamic-input
          v-model:value="createModel.env"
          preset="pair"
          key-placeholder="环境变量名"
          value-placeholder="环境变量值"
        />
      </n-form-item>
      <n-form-item path="env" label="标签">
        <n-dynamic-input
          v-model:value="createModel.labels"
          preset="pair"
          key-placeholder="标签名"
          value-placeholder="标签值"
        />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button type="info" block :loading="doSubmit" :disabled="doSubmit" @click="handleSubmit">
          提交
        </n-button>
      </n-col>
    </n-row>
  </n-modal>
</template>

<style scoped lang="scss"></style>
