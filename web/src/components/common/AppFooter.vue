<template>
  <footer color="#6a6a6a" f-c-c flex-col text-14>
    <p>
      © 2022 - {{ year }}
      <a
        hover="decoration-primary color-primary"
        target="__blank"
        href="https://github.com/TheTNB/panel"
      >
        耗子面板
      </a>
      All Rights Reserved.
    </p>
    <p>
      <a
        hover="decoration-primary color-primary"
        target="_blank"
        href="https://jq.qq.com/?_wv=1027&k=I1oJKSTH"
      >
        Q群 12370907
      </a>
      <n-divider vertical />
      <a hover="decoration-primary color-primary" target="_blank" href="https://tom.moe/docs">
        使用文档
      </a>
      <n-divider vertical />
      <a
        hover="decoration-primary color-primary"
        target="_blank"
        href="https://tom.moe/c/technical/panel"
      >
        交流社区
      </a>
      <n-divider vertical />
      <a
        hover="decoration-primary color-primary"
        target="_blank"
        href="https://afdian.com/a/TheTNB"
      >
        赞助支持
      </a>
    </p>
  </footer>
</template>

<script lang="ts" setup>
const year = new Date().getFullYear()
</script>
