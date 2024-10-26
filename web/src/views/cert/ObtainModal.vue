<script setup lang="ts">
import cert from '@/api/panel/cert'
import type { MessageReactive } from 'naive-ui'
import { NButton, NTable } from 'naive-ui'

let messageReactive: MessageReactive | null = null

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const id = defineModel<number>('id', { type: Number, required: true })

const model = ref({
  type: 'auto'
})

const options = [
  { label: '自动', value: 'auto' },
  { label: '手动', value: 'manual' },
  { label: '自签名', value: 'self-signed' }
]

const handleSubmit = async () => {
  messageReactive = window.$message.loading('请稍后...', {
    duration: 0
  })
  if (model.value.type == 'auto') {
    await cert
      .obtainAuto(id.value)
      .then(() => {
        window.$message.success('签发成功')
        show.value = false
      })
      .finally(() => {
        messageReactive?.destroy()
        window.$bus.emit('cert:refresh-cert')
        window.$bus.emit('cert:refresh-async')
      })
  } else if (model.value.type == 'manual') {
    const { data } = await cert.manualDNS(id.value)
    messageReactive.destroy()
    window.$message.info('请先前往域名处设置 DNS 解析，再继续签发')
    const d = window.$dialog.info({
      style: 'width: 60vw',
      title: '待设置DNS 记录列表',
      content: () => {
        return h(NTable, [
          h('thead', [
            h('tr', [h('th', '域名'), h('th', '类型'), h('th', '主机记录'), h('th', '记录值')])
          ]),
          h(
            'tbody',
            data.map((item: any) =>
              h('tr', [
                h('td', item?.domain),
                h('td', 'TXT'),
                h('td', item?.name),
                h('td', item?.value)
              ])
            )
          )
        ])
      },
      positiveText: '签发',
      onPositiveClick: async () => {
        d.loading = true
        messageReactive = window.$message.loading('请稍后...', {
          duration: 0
        })
        await cert
          .obtainManual(id.value)
          .then(() => {
            window.$message.success('签发成功')
            show.value = false
          })
          .finally(() => {
            d.loading = false
            messageReactive?.destroy()
            window.$bus.emit('cert:refresh-cert')
            window.$bus.emit('cert:refresh-async')
          })
      }
    })
  } else {
    await cert
      .obtainSelfSigned(id.value)
      .then(() => {
        window.$message.success('签发成功')
        show.value = false
      })
      .finally(() => {
        messageReactive?.destroy()
        window.$bus.emit('cert:refresh-cert')
        window.$bus.emit('cert:refresh-async')
      })
  }
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="签发证书"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form :model="model">
      <n-form-item path="type" label="签发模式">
        <n-select v-model:value="model.type" :options="options" />
      </n-form-item>
      <n-button type="info" block @click="handleSubmit">提交</n-button>
    </n-form>
  </n-modal>
</template>

<style scoped lang="scss"></style>
