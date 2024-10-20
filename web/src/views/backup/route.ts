import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'backup',
  path: '/backup',
  component: Layout,
  meta: {
    order: 60
  },
  children: [
    {
      name: 'backup-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: '备份',
        icon: 'mdi:backup-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
