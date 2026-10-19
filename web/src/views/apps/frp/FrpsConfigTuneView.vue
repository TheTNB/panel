<script setup lang="ts">
defineOptions({
  name: 'frps-config-tune',
})

import { useGettext } from 'vue3-gettext'

import frp from '@/api/apps/frp'
import { useResponsive } from '@/components/system/composables/useResponsive'

import { authMethodOptions, logLevelOptions, useBoolOptions } from './options'

const { $gettext } = useGettext()
const { isMobile } = useResponsive()
const currentTab = ref('general')
const saveLoading = ref(false)

const model = ref<Record<string, string>>({
  bind_addr: '',
  bind_port: '',
  kcp_bind_port: '',
  quic_bind_port: '',
  proxy_bind_addr: '',
  sub_domain_host: '',
  max_ports_per_client: '',
  udp_packet_size: '',
  detailed_errors_to_client: '',
  auth_method: '',
  auth_token: '',
  vhost_http_port: '',
  vhost_https_port: '',
  vhost_http_timeout: '',
  tcpmux_http_connect_port: '',
  custom_404_page: '',
  web_server_addr: '',
  web_server_port: '',
  web_server_user: '',
  web_server_password: '',
  web_server_pprof_enable: '',
  enable_prometheus: '',
  transport_max_pool_count: '',
  transport_tcp_mux: '',
  transport_tcp_mux_keepalive_interval: '',
  transport_tcp_keepalive: '',
  transport_heartbeat_timeout: '',
  transport_tls_force: '',
  transport_tls_cert_file: '',
  transport_tls_key_file: '',
  transport_tls_trusted_ca_file: '',
  log_to: '',
  log_level: '',
  log_max_days: '',
  log_disable_print_color: '',
})

const boolOptions = useBoolOptions()

useRequest(frp.server()).onSuccess(({ data }: any) => {
  model.value = { ...model.value, ...data }
})

const handleSave = () => {
  saveLoading.value = true
  useRequest(frp.saveServer(model.value))
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
          {{ $gettext('Leave a field empty to comment out the option and use the frps default.') }}
        </n-alert>
        <n-form>
          <n-form-item :label="$gettext('Bind Address (bindAddr)')">
            <n-input v-model:value="model.bind_addr" placeholder="0.0.0.0" />
          </n-form-item>
          <n-form-item :label="$gettext('Bind Port (bindPort)')">
            <n-input v-model:value="model.bind_port" placeholder="7000" />
          </n-form-item>
          <n-form-item :label="$gettext('KCP Bind Port (kcpBindPort)')">
            <n-input
              v-model:value="model.kcp_bind_port"
              :placeholder="$gettext('Leave empty to disable KCP')"
            />
          </n-form-item>
          <n-form-item :label="$gettext('QUIC Bind Port (quicBindPort)')">
            <n-input
              v-model:value="model.quic_bind_port"
              :placeholder="$gettext('Leave empty to disable QUIC')"
            />
          </n-form-item>
          <n-form-item :label="$gettext('Proxy Bind Address (proxyBindAddr)')">
            <n-input
              v-model:value="model.proxy_bind_addr"
              :placeholder="$gettext('Defaults to bindAddr')"
            />
          </n-form-item>
          <n-form-item :label="$gettext('Subdomain Host (subDomainHost)')">
            <n-input v-model:value="model.sub_domain_host" placeholder="frps.com" />
          </n-form-item>
          <n-form-item :label="$gettext('Max Ports Per Client (maxPortsPerClient)')">
            <n-input
              v-model:value="model.max_ports_per_client"
              :placeholder="$gettext('0 means no limit')"
            />
          </n-form-item>
          <n-form-item :label="$gettext('UDP Packet Size (udpPacketSize)')">
            <n-input v-model:value="model.udp_packet_size" placeholder="1500" />
          </n-form-item>
          <n-form-item :label="$gettext('Detailed Errors To Client (detailedErrorsToClient)')">
            <n-select
              v-model:value="model.detailed_errors_to_client"
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
    <n-tab-pane name="vhost" :tab="$gettext('HTTP(S)')">
      <n-flex vertical>
        <n-alert type="info">
          {{ $gettext('Required for http and https type proxies.') }}
        </n-alert>
        <n-form>
          <n-form-item :label="$gettext('HTTP Port (vhostHTTPPort)')">
            <n-input v-model:value="model.vhost_http_port" placeholder="80" />
          </n-form-item>
          <n-form-item :label="$gettext('HTTPS Port (vhostHTTPSPort)')">
            <n-input v-model:value="model.vhost_https_port" placeholder="443" />
          </n-form-item>
          <n-form-item :label="$gettext('HTTP Response Timeout (vhostHTTPTimeout)')">
            <n-input v-model:value="model.vhost_http_timeout" placeholder="60" />
          </n-form-item>
          <n-form-item :label="$gettext('TCPMux HTTP Connect Port (tcpmuxHTTPConnectPort)')">
            <n-input
              v-model:value="model.tcpmux_http_connect_port"
              :placeholder="$gettext('Required for tcpmux type proxies')"
            />
          </n-form-item>
          <n-form-item :label="$gettext('Custom 404 Page (custom404Page)')">
            <n-input v-model:value="model.custom_404_page" placeholder="/path/to/404.html" />
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
    <n-tab-pane name="dashboard" :tab="$gettext('Dashboard')">
      <n-flex vertical>
        <n-alert type="info">
          {{ $gettext('The frps dashboard is available only when the port is set.') }}
        </n-alert>
        <n-form>
          <n-form-item :label="$gettext('Listen Address (webServer.addr)')">
            <n-input v-model:value="model.web_server_addr" placeholder="127.0.0.1" />
          </n-form-item>
          <n-form-item :label="$gettext('Listen Port (webServer.port)')">
            <n-input v-model:value="model.web_server_port" placeholder="7500" />
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
          <n-form-item :label="$gettext('Enable Prometheus (enablePrometheus)')">
            <n-select v-model:value="model.enable_prometheus" :options="boolOptions" clearable />
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
          <n-form-item :label="$gettext('Max Pool Count (transport.maxPoolCount)')">
            <n-input v-model:value="model.transport_max_pool_count" placeholder="5" />
          </n-form-item>
          <n-form-item :label="$gettext('TCP Mux (transport.tcpMux)')">
            <n-select v-model:value="model.transport_tcp_mux" :options="boolOptions" clearable />
          </n-form-item>
          <n-form-item
            :label="$gettext('TCP Mux Keepalive Interval (transport.tcpMuxKeepaliveInterval)')"
          >
            <n-input v-model:value="model.transport_tcp_mux_keepalive_interval" placeholder="30" />
          </n-form-item>
          <n-form-item :label="$gettext('TCP Keepalive (transport.tcpKeepalive)')">
            <n-input v-model:value="model.transport_tcp_keepalive" placeholder="7200" />
          </n-form-item>
          <n-form-item :label="$gettext('Heartbeat Timeout (transport.heartbeatTimeout)')">
            <n-input v-model:value="model.transport_heartbeat_timeout" placeholder="90" />
          </n-form-item>
          <n-form-item :label="$gettext('Force TLS (transport.tls.force)')">
            <n-select v-model:value="model.transport_tls_force" :options="boolOptions" clearable />
          </n-form-item>
          <n-form-item :label="$gettext('TLS Certificate (transport.tls.certFile)')">
            <n-input v-model:value="model.transport_tls_cert_file" placeholder="server.crt" />
          </n-form-item>
          <n-form-item :label="$gettext('TLS Private Key (transport.tls.keyFile)')">
            <n-input v-model:value="model.transport_tls_key_file" placeholder="server.key" />
          </n-form-item>
          <n-form-item :label="$gettext('TLS CA (transport.tls.trustedCaFile)')">
            <n-input v-model:value="model.transport_tls_trusted_ca_file" placeholder="ca.crt" />
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
            <n-input v-model:value="model.log_to" placeholder="./frps.log" />
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
