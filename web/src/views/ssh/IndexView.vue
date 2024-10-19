<script setup lang="ts">
defineOptions({
  name: 'ssh-index'
})

import { AttachAddon } from '@xterm/addon-attach'
import { ClipboardAddon } from '@xterm/addon-clipboard'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import { WebglAddon } from '@xterm/addon-webgl'
import { Terminal } from '@xterm/xterm'
import '@xterm/xterm/css/xterm.css'
import { useI18n } from 'vue-i18n'

import ssh from '@/api/panel/ssh'

const { t } = useI18n()

const model = ref({
  host: '',
  port: 22,
  user: '',
  password: ''
})

const terminal = ref<HTMLElement | null>(null)
const term = ref()
const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
const ws = new WebSocket(`${protocol}://${window.location.host}/api/ssh/session`)
const attachAddon = new AttachAddon(ws)
const fitAddon = new FitAddon()
const clipboardAddon = new ClipboardAddon()
const webLinksAddon = new WebLinksAddon()
const webglAddon = new WebglAddon()

const handleSave = () => {
  ssh
    .saveInfo(model.value.host, model.value.port, model.value.user, model.value.password)
    .then(() => {
      window.$message.success(t('sshIndex.alerts.save'))
    })
}

const getInfo = () => {
  ssh.info().then((res) => {
    model.value.host = res.data.host
    model.value.port = res.data.port
    model.value.user = res.data.user
    model.value.password = res.data.password
  })
}

const openSession = () => {
  term.value = new Terminal({
    lineHeight: 1.2,
    fontSize: 14,
    fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
    cursorBlink: true,
    cursorStyle: 'underline',
    tabStopWidth: 4,
    theme: { background: '#111', foreground: '#fff' }
  })

  term.value.loadAddon(attachAddon)
  term.value.loadAddon(fitAddon)
  term.value.loadAddon(clipboardAddon)
  term.value.loadAddon(webLinksAddon)
  term.value.loadAddon(webglAddon)
  webglAddon.onContextLoss(() => {
    webglAddon.dispose()
  })

  ws.onopen = () => {
    term.value.open(terminal.value!)
    fitAddon.fit()
    term.value.focus()
    window.addEventListener(
      'resize',
      () => {
        fitAddon.fit()
      },
      false
    )
  }

  ws.onclose = () => {
    term.value.write('\r\n连接已关闭，请刷新重试。')
    term.value.write('\r\nConnection closed. Please refresh.')
    window.removeEventListener('resize', () => {
      fitAddon.fit()
    })
  }

  ws.onerror = (event) => {
    term.value.write('\r\n连接发生错误，请刷新重试。')
    term.value.write('\r\nConnection error. Please refresh .')
    console.error(event)
    ws.close()
  }

  term.value.onResize(({ cols, rows }: { cols: number; rows: number }) => {
    if (ws.readyState === 1) {
      ws.send(
        JSON.stringify({
          resize: true,
          columns: cols,
          rows: rows
        })
      )
    }
  })
}

const onTermWheel = (event: WheelEvent) => {
  if (event.ctrlKey) {
    event.preventDefault()
    if (event.deltaY > 0) {
      if (term.value.options.fontSize > 12) {
        term.value.options.fontSize = term.value.options.fontSize - 1
      }
    } else {
      term.value.options.fontSize = term.value.options.fontSize + 1
    }
    fitAddon.fit()
  }
}

onMounted(() => {
  getInfo()
  openSession()
})
</script>

<template>
  <common-page show-footer>
    <n-space vertical>
      <n-form inline>
        <n-form-item :label="$t('sshIndex.save.fields.host.label')">
          <n-input
            v-model:value="model.host"
            :placeholder="$t('sshIndex.save.fields.host.placeholder')"
            clearable
            size="small"
          />
        </n-form-item>
        <n-form-item :label="$t('sshIndex.save.fields.port.label')">
          <n-input-number
            v-model:value="model.port"
            :placeholder="$t('sshIndex.save.fields.port.placeholder')"
            clearable
            size="small"
          />
        </n-form-item>
        <n-form-item :label="$t('sshIndex.save.fields.username.label')">
          <n-input
            v-model:value="model.user"
            :placeholder="$t('sshIndex.save.fields.username.placeholder')"
            clearable
            size="small"
          />
        </n-form-item>
        <n-form-item :label="$t('sshIndex.save.fields.password.label')">
          <n-input
            v-model:value="model.password"
            :placeholder="$t('sshIndex.save.fields.password.placeholder')"
            clearable
            size="small"
          />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" size="small" @click="handleSave">
            {{ $t('sshIndex.save.actions.submit') }}
          </n-button>
        </n-form-item>
      </n-form>
      <n-card>
        <div ref="terminal" @wheel="onTermWheel" h-600></div>
      </n-card>
    </n-space>
  </common-page>
</template>
