import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php84',
  path: '/apps/php84',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-php84-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'PHP 8.4',
        icon: 'logos:php',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
