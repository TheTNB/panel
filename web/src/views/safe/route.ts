import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'safe',
  path: '/safe',
  component: Layout,
  meta: {
    order: 4
  },
  children: [
    {
      name: 'safe-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'safeIndex.title',
        icon: 'mdi:server-security',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
