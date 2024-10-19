<script setup lang="ts">
defineOptions({
  name: 'container-index'
})

import ContainerView from '@/views/container/ContainerView.vue'
import ImageView from '@/views/container/ImageView.vue'
import NetworkView from '@/views/container/NetworkView.vue'
import VolumeView from '@/views/container/VolumeView.vue'

const currentTab = ref('container')
</script>

<template>
  <common-page show-footer>
    <n-tabs v-model:value="currentTab" type="line" animated size="large">
      <n-tab-pane name="container" tab="容器">
        <container-view />
      </n-tab-pane>
      <n-tab-pane name="image" tab="镜像">
        <image-view />
      </n-tab-pane>
      <n-tab-pane name="network" tab="网络">
        <network-view />
      </n-tab-pane>
      <n-tab-pane name="volume" tab="卷">
        <volume-view />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
