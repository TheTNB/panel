import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php80',
  path: '/apps/php80',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-php80-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'PHP 8.0',
        icon: 'logos:php',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
