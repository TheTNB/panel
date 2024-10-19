import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'benchmark',
  path: '/apps/benchmark',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-benchmark-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: '耗子跑分',
        icon: 'dashicons:performance',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
