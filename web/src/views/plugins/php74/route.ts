import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php74',
  path: '/plugins/php74',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-php74-index',
      path: '',
      component: () => import('../php/IndexView.vue'),
      meta: {
        title: 'PHP 7.4',
        icon: 'mdi:language-php',
        role: ['admin'],
        requireAuth: true,
        php: 74
      }
    }
  ]
} as RouteType
