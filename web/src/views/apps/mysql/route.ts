import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'mysql',
  path: '/apps/mysql',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-mysql-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Percona（MySQL）',
        icon: 'logos:percona',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
