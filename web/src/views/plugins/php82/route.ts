import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php82',
  path: '/plugins/php82',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-php82-index',
      path: '',
      component: () => import('../php/IndexView.vue'),
      meta: {
        title: 'PHP 8.2',
        icon: 'mdi:language-php',
        role: ['admin'],
        requireAuth: true,
        php: 82
      }
    }
  ]
} as RouteType
