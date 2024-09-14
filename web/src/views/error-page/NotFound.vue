<script lang="ts" setup>
const { replace } = useRouter()
</script>

<template>
  <AppPage>
    <n-result description="抱歉，您访问的页面不存在。" m-auto status="404">
      <template #icon>
        <img src="@/assets/images/404.webp" width="500" />
      </template>
      <template #footer>
        <n-button @click="replace('/')"> 返回首页</n-button>
      </template>
    </n-result>
  </AppPage>
</template>
