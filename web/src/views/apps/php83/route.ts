import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php83',
  path: '/apps/php83',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-php83-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'PHP 8.3',
        icon: 'logos:php',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
