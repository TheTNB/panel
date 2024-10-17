import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'phpmyadmin',
  path: '/apps/phpmyadmin',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-phpmyadmin-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'phpMyAdmin',
        icon: 'simple-icons:phpmyadmin',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
