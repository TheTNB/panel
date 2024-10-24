<script setup lang="ts">
defineOptions({
  name: 'ssh-index'
})

import ssh from '@/api/panel/ssh'
import ws from '@/api/ws'
import TheIcon from '@/components/custom/TheIcon.vue'
import CreateModal from '@/views/ssh/CreateModal.vue'
import UpdateModal from '@/views/ssh/UpdateModal.vue'
import '@fontsource-variable/jetbrains-mono/wght-italic.css'
import '@fontsource-variable/jetbrains-mono/wght.css'
import { AttachAddon } from '@xterm/addon-attach'
import { ClipboardAddon } from '@xterm/addon-clipboard'
import { FitAddon } from '@xterm/addon-fit'
import { Unicode11Addon } from '@xterm/addon-unicode11'
import { WebLinksAddon } from '@xterm/addon-web-links'
import { WebglAddon } from '@xterm/addon-webgl'
import { Terminal } from '@xterm/xterm'
import '@xterm/xterm/css/xterm.css'
import { NButton, NFlex, NPopconfirm } from 'naive-ui'

const terminal = ref<HTMLElement | null>(null)
const term = ref()
let sshWs: WebSocket | null = null
const fitAddon = new FitAddon()
const webglAddon = new WebglAddon()

const current = ref(0)
const collapsed = ref(true)
const create = ref(false)
const update = ref(false)
const updateId = ref(0)

const list = ref<any[]>([])

const fetchData = async () => {
  list.value = []
  const { data } = await ssh.list(1, 10000)
  if (data.items.length === 0) {
    window.$message.info('请先创建主机')
    return
  }
  data.items.forEach((item: any) => {
    list.value.push({
      label: item.name === '' ? item.host : item.name,
      key: item.id,
      extra: () => {
        return h(
          NFlex,
          {
            size: 'small',
            style: 'float: right'
          },
          {
            default: () => [
              h(
                NButton,
                {
                  type: 'primary',
                  size: 'small',
                  onClick: () => {
                    update.value = true
                    updateId.value = item.id
                  }
                },
                {
                  default: () => {
                    return '编辑'
                  }
                }
              ),
              h(
                NPopconfirm,
                {
                  onPositiveClick: () => handleDelete(item.id)
                },
                {
                  default: () => {
                    return '确定删除主机吗？'
                  },
                  trigger: () => {
                    return h(
                      NButton,
                      {
                        size: 'small',
                        type: 'error'
                      },
                      {
                        default: () => {
                          return '删除'
                        }
                      }
                    )
                  }
                }
              )
            ]
          }
        )
      }
    })
  })
  await openSession(updateId.value === 0 ? Number(list.value[0].key) : updateId.value)
}

const handleDelete = async (id: number) => {
  await ssh.delete(id)
  list.value = list.value.filter((item: any) => item.key !== id)
  if (current.value === id) {
    if (list.value.length > 0) {
      await openSession(Number(list.value[0].key))
    } else {
      term.value.dispose()
    }
    if (list.value.length === 0) {
      create.value = true
    }
  }
}

const handleChange = (key: number) => {
  openSession(key)
}

const openSession = async (id: number) => {
  closeSession()
  await ws.ssh(id).then((ws) => {
    sshWs = ws
    term.value = new Terminal({
      allowProposedApi: true,
      lineHeight: 1.2,
      fontSize: 14,
      fontFamily: `'JetBrains Mono Variable', monospace`,
      cursorBlink: true,
      cursorStyle: 'underline',
      tabStopWidth: 4,
      theme: { background: '#111', foreground: '#fff' }
    })

    term.value.loadAddon(new AttachAddon(ws))
    term.value.loadAddon(fitAddon)
    term.value.loadAddon(new ClipboardAddon())
    term.value.loadAddon(new WebLinksAddon())
    term.value.loadAddon(new Unicode11Addon())
    term.value.unicode.activeVersion = '11'
    term.value.loadAddon(webglAddon)
    webglAddon.onContextLoss(() => {
      webglAddon.dispose()
    })
    term.value.open(terminal.value!)

    fitAddon.fit()
    term.value.focus()
    window.addEventListener('resize', onResize, false)
    current.value = id

    ws.onclose = () => {
      term.value.write('\r\n连接已关闭，请刷新重试。')
      term.value.write('\r\nConnection closed. Please refresh.')
      window.removeEventListener('resize', onResize)
    }

    ws.onerror = (event) => {
      term.value.write('\r\n连接发生错误，请刷新重试。')
      term.value.write('\r\nConnection error. Please refresh .')
      console.error(event)
      ws.close()
    }
  })
}

const closeSession = () => {
  try {
    term.value.dispose()
    sshWs?.close()
    terminal.value!.innerHTML = ''
  } catch {
    /* empty */
  }
}

const onResize = () => {
  fitAddon.fit()
  if (sshWs != null && sshWs.readyState === 1) {
    const { cols, rows } = term.value
    sshWs.send(
      JSON.stringify({
        resize: true,
        columns: cols,
        rows: rows
      })
    )
  }
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
  // https://github.com/xtermjs/xterm.js/pull/5178
  document.fonts.ready.then((fontFaceSet: any) =>
    Promise.all(Array.from(fontFaceSet).map((el: any) => el.load())).then(fetchData)
  )
  window.$bus.on('ssh:refresh', fetchData)
})

onUnmounted(() => {
  closeSession()
  window.$bus.off('ssh:refresh')
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button type="primary" @click="create = true">
        <TheIcon :size="18" icon="material-symbols:add" />
        创建主机
      </n-button>
    </template>
    <n-layout has-sider sider-placement="right">
      <n-layout content-style="overflow: visible" bg-hex-111>
        <div ref="terminal" @wheel="onTermWheel" h-75vh></div>
      </n-layout>
      <n-layout-sider
        bordered
        :collapsed-width="0"
        :collapsed="collapsed"
        show-trigger
        :native-scrollbar="false"
        @collapse="collapsed = true"
        @expand="collapsed = false"
        @after-enter="onResize"
        @after-leave="onResize"
      >
        <n-menu
          v-model:value="current"
          :collapsed="collapsed"
          :collapsed-width="0"
          :collapsed-icon-size="0"
          :options="list"
          @update-value="handleChange"
        />
      </n-layout-sider>
    </n-layout>
  </common-page>
  <create-modal v-model:show="create" />
  <update-modal v-model:show="update" v-model:id="updateId" />
</template>

<style scoped lang="scss">
:deep(.xterm) {
  padding: 4rem !important;
}

:deep(.xterm .xterm-viewport::-webkit-scrollbar) {
  border-radius: 0.4rem;
  height: 6px;
  width: 8px;
}

:deep(.xterm .xterm-viewport::-webkit-scrollbar-thumb) {
  background-color: #666;
  border-radius: 0.4rem;
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
  transition: all 1s;
}

:deep(.xterm .xterm-viewport:hover::-webkit-scrollbar-thumb) {
  background-color: #aaa;
}

:deep(.xterm .xterm-viewport::-webkit-scrollbar-track) {
  background-color: #111;
  border-radius: 0.4rem;
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
  transition: all 1s;
}

:deep(.xterm .xterm-viewport:hover::-webkit-scrollbar-track) {
  background-color: #444;
}
</style>
