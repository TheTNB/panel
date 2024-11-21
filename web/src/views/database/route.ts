import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'database',
  path: '/database',
  component: Layout,
  meta: {
    order: 2
  },
  children: [
    {
      name: 'database-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: '数据库',
        icon: 'mdi:database',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
