import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'cron',
  path: '/cron',
  component: Layout,
  meta: {
    order: 70
  },
  children: [
    {
      name: 'cron-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'cronIndex.title',
        icon: 'mdi:timer-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
