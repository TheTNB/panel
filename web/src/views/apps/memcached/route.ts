import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'memcached',
  path: '/apps/memcached',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-memcached-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Memcached',
        icon: 'logos:memcached',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
