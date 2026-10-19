<script setup lang="ts">
defineOptions({
  name: 'frpc-config-tune',
})

import { useGettext } from 'vue3-gettext'

import frp from '@/api/apps/frp'
import { useResponsive } from '@/components/system/composables/useResponsive'

import { authMethodOptions, logLevelOptions, toOptions, useBoolOptions } from './options'

const { $gettext } = useGettext()
const { isMobile } = useResponsive()
const currentTab = ref('general')
const saveLoading = ref(false)

const model = ref<Record<string, string>>({
  user: '',
  server_addr: '',
  server_port: '',
  login_fail_exit: '',
  nat_hole_stun_server: '',
  dns_server: '',
  udp_packet_size: '',
  auth_method: '',
  auth_token: '',
  transport_protocol: '',
  transport_pool_count: '',
  transport_tcp_mux: '',
  transport_tcp_mux_keepalive_interval: '',
  transport_dial_server_timeout: '',
  transport_dial_server_keepalive: '',
  transport_heartbeat_interval: '',
  transport_heartbeat_timeout: '',
  transport_connect_server_local_ip: '',
  transport_proxy_url: '',
  transport_tls_enable: '',
  transport_tls_cert_file: '',
  transport_tls_key_file: '',
  transport_tls_trusted_ca_file: '',
  transport_tls_server_name: '',
  web_server_addr: '',
  web_server_port: '',
  web_server_user: '',
  web_server_password: '',
  web_server_pprof_enable: '',
  log_to: '',
  log_level: '',
  log_max_days: '',
  log_disable_print_color: '',
})

const boolOptions = useBoolOptions()

const protocolOptions = toOptions(['tcp', 'kcp', 'quic', 'websocket', 'wss'])

useRequest(frp.client()).onSuccess(({ data }: any) => {
  model.value = { ...model.value, ...data }
})

const handleSave = () => {
  saveLoading.value = true
  useRequest(frp.saveClient(model.value))
    .onSuccess(() => {
      window.$message.success($gettext('Saved successfully'))
    })
    .onComplete(() => {
      saveLoading.value = false
    })
}
</script>

<template>
  <n-tabs v-model:value="currentTab" type="line" :placement="isMobile ? 'top' : 'left'" animated>
    <n-tab-pane name="general" :tab="$gettext('General')">
      <n-flex vertical>
        <n-alert type="info">
          {{ $gettext('Leave a field empty to comment out the option and use the frpc default.') }}
        </n-alert>
        <n-form>
          <n-form-item :label="$gettext('Server Address (serverAddr)')">
            <n-input v-model:value="model.server_addr" placeholder="0.0.0.0" />
          </n-form-item>
          <n-form-item :label="$gettext('Server Port (serverPort)')">
            <n-input v-model:value="model.server_port" placeholder="7000" />
          </n-form-item>
          <n-form-item :label="$gettext('User (user)')">
            <n-input
              v-model:value="model.user"
              :placeholder="$gettext('Proxy names will be prefixed with this user')"
            />
          </n-form-item>
          <n-form-item :label="$gettext('Exit On Login Failure (loginFailExit)')">
            <n-select v-model:value="model.login_fail_exit" :options="boolOptions" clearable />
          </n-form-item>
          <n-form-item :label="$gettext('STUN Server (natHoleStunServer)')">
            <n-input
              v-model:value="model.nat_hole_stun_server"
              placeholder="stun.easyvoip.com:3478"
            />
          </n-form-item>
          <n-form-item :label="$gettext('DNS Server (dnsServer)')">
            <n-input v-model:value="model.dns_server" placeholder="8.8.8.8" />
          </n-form-item>
          <n-form-item :label="$gettext('UDP Packet Size (udpPacketSize)')">
            <n-input v-model:value="model.udp_packet_size" placeholder="1500" />
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
        </n-flex>
      </n-flex>
    </n-tab-pane>
    <n-tab-pane name="auth" :tab="$gettext('Authentication')">
      <n-flex vertical>
        <n-alert type="warning">
          {{ $gettext('The token must be identical on frps and frpc, otherwise login will fail.') }}
        </n-alert>
        <n-form>
          <n-form-item :label="$gettext('Method (auth.method)')">
            <n-select v-model:value="model.auth_method" :options="authMethodOptions" clearable />
          </n-form-item>
          <n-form-item :label="$gettext('Token (auth.token)')">
            <n-input
              v-model:value="model.auth_token"
              type="password"
              show-password-on="click"
              :placeholder="$gettext('Leave empty for no authentication')"
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
        </n-flex>
      </n-flex>
    </n-tab-pane>
    <n-tab-pane name="transport" :tab="$gettext('Transport')">
      <n-flex vertical>
        <n-alert type="warning">
          {{ $gettext('Do not modify these options unless you know what they do.') }}
        </n-alert>
        <n-form>
          <n-form-item :label="$gettext('Protocol (transport.protocol)')">
            <n-select
              v-model:value="model.transport_protocol"
              :options="protocolOptions"
              clearable
            />
          </n-form-item>
          <n-form-item :label="$gettext('Pool Count (transport.poolCount)')">
            <n-input v-model:value="model.transport_pool_count" placeholder="5" />
          </n-form-item>
          <n-form-item :label="$gettext('TCP Mux (transport.tcpMux)')">
            <n-select v-model:value="model.transport_tcp_mux" :options="boolOptions" clearable />
          </n-form-item>
          <n-form-item
            :label="$gettext('TCP Mux Keepalive Interval (transport.tcpMuxKeepaliveInterval)')"
          >
            <n-input v-model:value="model.transport_tcp_mux_keepalive_interval" placeholder="30" />
          </n-form-item>
          <n-form-item :label="$gettext('Dial Timeout (transport.dialServerTimeout)')">
            <n-input v-model:value="model.transport_dial_server_timeout" placeholder="10" />
          </n-form-item>
          <n-form-item :label="$gettext('Dial Keepalive (transport.dialServerKeepalive)')">
            <n-input v-model:value="model.transport_dial_server_keepalive" placeholder="7200" />
          </n-form-item>
          <n-form-item :label="$gettext('Heartbeat Interval (transport.heartbeatInterval)')">
            <n-input v-model:value="model.transport_heartbeat_interval" placeholder="30" />
          </n-form-item>
          <n-form-item :label="$gettext('Heartbeat Timeout (transport.heartbeatTimeout)')">
            <n-input v-model:value="model.transport_heartbeat_timeout" placeholder="90" />
          </n-form-item>
          <n-form-item :label="$gettext('Local IP (transport.connectServerLocalIP)')">
            <n-input
              v-model:value="model.transport_connect_server_local_ip"
              placeholder="0.0.0.0"
            />
          </n-form-item>
          <n-form-item :label="$gettext('Proxy URL (transport.proxyURL)')">
            <n-input
              v-model:value="model.transport_proxy_url"
              placeholder="socks5://user:passwd@192.168.1.128:1080"
            />
          </n-form-item>
          <n-form-item :label="$gettext('Enable TLS (transport.tls.enable)')">
            <n-select v-model:value="model.transport_tls_enable" :options="boolOptions" clearable />
          </n-form-item>
          <n-form-item :label="$gettext('TLS Certificate (transport.tls.certFile)')">
            <n-input v-model:value="model.transport_tls_cert_file" placeholder="client.crt" />
          </n-form-item>
          <n-form-item :label="$gettext('TLS Private Key (transport.tls.keyFile)')">
            <n-input v-model:value="model.transport_tls_key_file" placeholder="client.key" />
          </n-form-item>
          <n-form-item :label="$gettext('TLS CA (transport.tls.trustedCaFile)')">
            <n-input v-model:value="model.transport_tls_trusted_ca_file" placeholder="ca.crt" />
          </n-form-item>
          <n-form-item :label="$gettext('TLS Server Name (transport.tls.serverName)')">
            <n-input v-model:value="model.transport_tls_server_name" placeholder="example.com" />
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
        </n-flex>
      </n-flex>
    </n-tab-pane>
    <n-tab-pane name="admin" :tab="$gettext('Admin API')">
      <n-flex vertical>
        <n-alert type="info">
          {{ $gettext('The frpc admin API is available only when the port is set.') }}
        </n-alert>
        <n-form>
          <n-form-item :label="$gettext('Listen Address (webServer.addr)')">
            <n-input v-model:value="model.web_server_addr" placeholder="127.0.0.1" />
          </n-form-item>
          <n-form-item :label="$gettext('Listen Port (webServer.port)')">
            <n-input v-model:value="model.web_server_port" placeholder="7400" />
          </n-form-item>
          <n-form-item :label="$gettext('Username (webServer.user)')">
            <n-input v-model:value="model.web_server_user" placeholder="admin" />
          </n-form-item>
          <n-form-item :label="$gettext('Password (webServer.password)')">
            <n-input
              v-model:value="model.web_server_password"
              type="password"
              show-password-on="click"
            />
          </n-form-item>
          <n-form-item :label="$gettext('Enable pprof (webServer.pprofEnable)')">
            <n-select
              v-model:value="model.web_server_pprof_enable"
              :options="boolOptions"
              clearable
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
        </n-flex>
      </n-flex>
    </n-tab-pane>
    <n-tab-pane name="log" :tab="$gettext('Logs')">
      <n-flex vertical>
        <n-form>
          <n-form-item :label="$gettext('Output (log.to)')">
            <n-input v-model:value="model.log_to" placeholder="./frpc.log" />
          </n-form-item>
          <n-form-item :label="$gettext('Level (log.level)')">
            <n-select v-model:value="model.log_level" :options="logLevelOptions" clearable />
          </n-form-item>
          <n-form-item :label="$gettext('Max Days (log.maxDays)')">
            <n-input v-model:value="model.log_max_days" placeholder="3" />
          </n-form-item>
          <n-form-item :label="$gettext('Disable Colors (log.disablePrintColor)')">
            <n-select
              v-model:value="model.log_disable_print_color"
              :options="boolOptions"
              clearable
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
        </n-flex>
      </n-flex>
    </n-tab-pane>
  </n-tabs>
</template>
