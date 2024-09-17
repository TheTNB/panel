import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'pureftpd',
  path: '/plugins/pureftpd',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-pureftpd-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Pure-FTPd',
        icon: 'mdi:server-network',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
