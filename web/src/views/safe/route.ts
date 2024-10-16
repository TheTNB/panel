import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'safe',
  path: '/safe',
  component: Layout,
  meta: {
    order: 30
  },
  children: [
    {
      name: 'safe-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: '系统安全',
        icon: 'mdi:server-security',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
