import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'mysql57',
  path: '/plugins/mysql57',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-mysql57-index',
      path: '',
      component: () => import('../mysql/IndexView.vue'),
      meta: {
        title: 'MySQL 5.7',
        icon: 'mdi:database',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
