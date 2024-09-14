import type { RouteModule, RoutesType, RouteType } from '~/types/router'

export const basicRoutes: RoutesType = [
  {
    name: '404',
    path: '/404',
    component: () => import('@/views/error-page/NotFound.vue'),
    isHidden: true
  },

  {
    name: 'Login',
    path: '/login',
    component: () => import('@/views/login/IndexView.vue'),
    isHidden: true,
    meta: {
      title: '登录页'
    }
  }
]

export const NOT_FOUND_ROUTE: RouteType = {
  name: 'NotFound',
  path: '/:pathMatch(.*)*',
  redirect: '/404',
  isHidden: true
}

export const EMPTY_ROUTE: RouteType = {
  name: 'Empty',
  path: '/:pathMatch(.*)*',
  component: () => {}
}

const modules = import.meta.glob('@/views/**/route.ts', {
  eager: true
}) as RouteModule
const asyncRoutes: RoutesType = []
Object.keys(modules).forEach((key) => {
  asyncRoutes.push(modules[key].default)
})

export { asyncRoutes }
