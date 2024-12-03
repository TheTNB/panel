<script lang="ts" setup>
import user from '@/api/panel/user'
import bgImg from '@/assets/images/login_bg.webp'
import logoImg from '@/assets/images/logo.png'
import { addDynamicRoutes } from '@/router'
import { useThemeStore, useUserStore } from '@/store'
import { getLocal, removeLocal, setLocal } from '@/utils'
import { rsaEncrypt } from '@/utils/encrypt'

const router = useRouter()
const route = useRoute()
const query = route.query
const { data: key, loading: isLoading } = useRequest(user.key, { initialData: '' })
const { data: isLogin } = useRequest(user.isLogin, { initialData: false })

interface LoginInfo {
  username: string
  password: string
  safe_login: boolean
}

const loginInfo = ref<LoginInfo>({
  username: '',
  password: '',
  safe_login: true
})

const localLoginInfo = getLocal('loginInfo') as LoginInfo
if (localLoginInfo) {
  loginInfo.value.username = localLoginInfo.username || ''
  loginInfo.value.password = localLoginInfo.password || ''
}

const userStore = useUserStore()
const themeStore = useThemeStore()
const loging = ref<boolean>(false)
const isRemember = useStorage('isRemember', false)

const logo = computed(() => themeStore.logo || logoImg)

async function handleLogin() {
  const { username, password, safe_login } = loginInfo.value
  if (!username || !password) {
    window.$message.warning('请输入用户名和密码')
    return
  }
  if (!key) {
    window.$message.warning('获取加密公钥失败，请刷新页面重试')
    return
  }
  try {
    user
      .login(
        rsaEncrypt(username, String(unref(key))),
        rsaEncrypt(password, String(unref(key))),
        safe_login
      )
      .then(async () => {
        loging.value = true
        window.$notification?.success({ title: '登录成功！', duration: 2500 })
        if (isRemember.value) {
          setLocal('loginInfo', { username, password })
        } else {
          removeLocal('loginInfo')
        }

        await addDynamicRoutes()
        await user.info().then((data: any) => {
          userStore.set(data)
        })
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

watch(
  () => isLogin,
  async () => {
    if (isLogin) {
      console.log(isLogin)
      await addDynamicRoutes()
      await user.info().then((data: any) => {
        userStore.set(data)
      })
      if (query.redirect) {
        const path = query.redirect as string
        Reflect.deleteProperty(query, 'redirect')
        await router.push({ path, query })
      } else {
        await router.push('/')
      }
    }
  }
)
</script>

<template>
  <AppPage :show-footer="true" :style="{ backgroundImage: `url(${bgImg})` }" bg-cover>
    <div m-auto min-w-345 f-c-c rounded-10 bg-white bg-opacity-60 p-15 card-shadow dark:bg-dark>
      <div w-480 flex-col px-20 py-35>
        <h5 color="#6a6a6a" f-c-c text-24 font-normal>
          <n-image :src="logo" height="50" preview-disabled mr-10 />{{ themeStore.name }}
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
            type="password"
            show-password-on="click"
            @keydown.enter="handleLogin"
          />
        </div>

        <div mt-20>
          <n-flex>
            <n-checkbox v-model:checked="loginInfo.safe_login" label="安全登录" />
            <n-checkbox v-model:checked="isRemember" label="记住我" />
          </n-flex>
        </div>

        <div mt-20>
          <n-button
            :loading="isLoading || loging"
            :disabled="isLoading || loging"
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
