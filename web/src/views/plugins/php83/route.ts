import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php83',
  path: '/plugins/php83',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-php83-index',
      path: '',
      component: () => import('../php/IndexView.vue'),
      meta: {
        title: 'PHP 8.3',
        icon: 'mdi:language-php',
        role: ['admin'],
        requireAuth: true,
        php: 83
      }
    }
  ]
} as RouteType
