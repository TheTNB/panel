import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'mysql84',
  path: '/plugins/mysql84',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-mysql84-index',
      path: '',
      component: () => import('../mysql/IndexView.vue'),
      meta: {
        title: 'MySQL 8.4',
        icon: 'mdi:database',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
