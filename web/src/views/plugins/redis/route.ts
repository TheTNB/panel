import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'redis',
  path: '/plugins/redis',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-redis-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Redis',
        icon: 'mdi:database',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
