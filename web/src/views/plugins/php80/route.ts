import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'php80',
  path: '/plugins/php80',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-php80-index',
      path: '',
      component: () => import('../php/IndexView.vue'),
      meta: {
        title: 'PHP 8.0',
        icon: 'mdi:language-php',
        role: ['admin'],
        requireAuth: true,
        php: 80
      }
    }
  ]
} as RouteType
