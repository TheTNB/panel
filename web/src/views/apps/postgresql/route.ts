import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'postgresql',
  path: '/apps/postgresql',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-postgresql-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'PostgreSQL',
        icon: 'logos:postgresql',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
