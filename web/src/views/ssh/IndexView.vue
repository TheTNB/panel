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
  const term = new Terminal({
    lineHeight: 1.2,
    fontSize: 14,
    fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
    cursorBlink: true,
    cursorStyle: 'underline',
    scrollback: 1000,
    scrollSensitivity: 15,
    tabStopWidth: 4,
    theme: { background: '#111', foreground: '#fff' }
  })

  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const ws = new WebSocket(`${protocol}://${window.location.host}/api/ssh/session`)

  const attachAddon = new AttachAddon(ws)
  term.loadAddon(attachAddon)
  const fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  const clipboardAddon = new ClipboardAddon()
  term.loadAddon(clipboardAddon)
  const webLinksAddon = new WebLinksAddon()
  term.loadAddon(webLinksAddon)
  const webglAddon = new WebglAddon()
  term.loadAddon(webglAddon)
  webglAddon.onContextLoss(() => {
    webglAddon.dispose()
  })

  ws.onopen = () => {
    term.open(document.getElementById('terminal') as HTMLElement)
    fitAddon.fit()
    term.write('\r\n欢迎来到耗子面板 SSH，连接成功。')
    term.write('\r\nWelcome to Rat Panel SSH. Connection success.\r\n')
    term.focus()
  }

  ws.onclose = () => {
    term.write('\r\nSSH连接已关闭，请刷新页面。')
    term.write('\r\nSSH connection closed. Please refresh the page.\r\n')
  }

  ws.onerror = (event) => {
    term.write('\r\nSSH连接发生错误，请刷新页面。')
    term.write('\r\nSSH connection error. Please refresh the page.\r\n')
    console.error(event)
    ws.close()
  }

  term.onResize(({ cols, rows }) => {
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

  window.addEventListener(
    'resize',
    () => {
      fitAddon.fit()
    },
    false
  )
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
      <div
        id="terminal"
        style="width: 100%; height: 70vh; background-color: #000000; margin-top: 20px"
      ></div>
    </n-space>
  </common-page>
</template>
