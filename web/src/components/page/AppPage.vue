<script lang="ts" setup>
interface Props {
  showFooter?: boolean
}

withDefaults(defineProps<Props>(), {
  showFooter: false
})
</script>

<template>
  <transition appear mode="out-in" name="fade-slide">
    <section class="cus-scroll-y wh-full flex-col bg-[#f5f6fb] p-15" dark:bg-hex-121212>
      <slot />
      <AppFooter v-if="showFooter" mt-15 />
    </section>
  </transition>
</template>
