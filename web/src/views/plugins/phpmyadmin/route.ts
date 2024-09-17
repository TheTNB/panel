import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'phpmyadmin',
  path: '/plugins/phpmyadmin',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-phpmyadmin-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'phpMyAdmin',
        icon: 'mdi:database',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
