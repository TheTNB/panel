import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php74',
  path: '/apps/php74',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-php74-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'PHP 7.4',
        icon: 'logos:php',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
