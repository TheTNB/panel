import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'postgresql15',
  path: '/plugins/postgresql15',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-postgresql15-index',
      path: '',
      component: () => import('../postgresql/IndexView.vue'),
      meta: {
        title: 'PostgreSQL 15',
        icon: 'mdi:database',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
