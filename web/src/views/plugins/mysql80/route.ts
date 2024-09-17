import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'mysql80',
  path: '/plugins/mysql80',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-mysql80-index',
      path: '',
      component: () => import('../mysql/IndexView.vue'),
      meta: {
        title: 'MySQL 8.0',
        icon: 'mdi:database',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
