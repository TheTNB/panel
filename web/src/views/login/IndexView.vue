<script lang="ts" setup>
import { useStorage } from '@vueuse/core'
import { useUserStore } from '@/store'
import { getLocal, removeLocal, setLocal } from '@/utils'
import bgImg from '@/assets/images/login_bg.webp'
import { addDynamicRoutes } from '@/router'
import user from '@/api/panel/user'
import { title } from '@/main'

const router = useRouter()
const route = useRoute()
const query = route.query

interface LoginInfo {
  username: string
  password: string
}

const loginInfo = ref<LoginInfo>({
  username: '',
  password: ''
})

const localLoginInfo = getLocal('loginInfo') as LoginInfo
if (localLoginInfo) {
  loginInfo.value.username = localLoginInfo.username || ''
  loginInfo.value.password = localLoginInfo.password || ''
}

const loging = ref<boolean>(false)
const isRemember = useStorage('isRemember', false)
const userStore = useUserStore()

async function handleLogin() {
  const { username, password } = loginInfo.value
  if (!username || !password) {
    window.$message.warning('请输入用户名和密码')
    return
  }
  try {
    user.login(username, password).then(async () => {
      loging.value = true
      window.$notification?.success({ title: '登录成功！', duration: 2500 })
      if (isRemember.value) {
        setLocal('loginInfo', { username, password })
      } else {
        removeLocal('loginInfo')
      }

      await addDynamicRoutes()
      await userStore.getUserInfo()
      if (query.redirect) {
        const path = query.redirect as string
        Reflect.deleteProperty(query, 'redirect')
        await router.push({ path, query })
      } else {
        await router.push('/')
      }
    })
  } catch (error) {
    console.error(error)
  }
  loging.value = false
}

onMounted(async () => {
  // 已登录自动跳转
  await user.isLogin().then(async (res) => {
    if (res.data) {
      await addDynamicRoutes()
      await userStore.getUserInfo()
      if (query.redirect) {
        const path = query.redirect as string
        Reflect.deleteProperty(query, 'redirect')
        await router.push({ path, query })
      } else {
        await router.push('/')
      }
    }
  })
})
</script>

<template>
  <AppPage :show-footer="true" :style="{ backgroundImage: `url(${bgImg})` }" bg-cover>
    <div m-auto min-w-345 f-c-c rounded-10 bg-white bg-opacity-60 p-15 card-shadow dark:bg-dark>
      <div hidden w-380 px-20 py-35 md:block>
        <img alt="login_banner" src="@/assets/images/login_banner.png" w-full />
      </div>

      <div w-320 flex-col px-20 py-35>
        <h5 color="#6a6a6a" f-c-c text-24 font-normal>
          <img class="mr-10" height="50" src="@/assets/images/logo.png" />{{ title }}
        </h5>
        <div mt-30>
          <n-input
            v-model:value="loginInfo.username"
            :maxlength="32"
            autofocus
            class="h-50 items-center pl-10 text-16"
            placeholder="用户名"
          />
        </div>
        <div mt-30>
          <n-input
            v-model:value="loginInfo.password"
            :maxlength="32"
            class="h-50 items-center pl-10 text-16"
            placeholder="密码"
            show-password-on="mousedown"
            type="password"
            @keydown.enter="handleLogin"
          />
        </div>

        <div mt-20>
          <n-checkbox
            :checked="isRemember"
            :on-update:checked="(val: boolean) => (isRemember = val)"
            label="记住我"
          />
        </div>

        <div mt-20>
          <n-button
            :loading="loging"
            type="primary"
            h-50
            w-full
            rounded-5
            text-16
            @click="handleLogin"
          >
            登录
          </n-button>
        </div>
      </div>
    </div>
  </AppPage>
</template>
