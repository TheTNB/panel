import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php82',
  path: '/apps/php82',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-php82-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'PHP 8.2',
        icon: 'logos:php',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
