<script setup lang="ts">
import PhpView from '@/views/apps/php/PhpView.vue'

defineOptions({
  name: 'apps-php80-index'
})
</script>

<template>
  <php-view :version="80" />
</template>
