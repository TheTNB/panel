import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'openresty',
  path: '/plugins/openresty',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-openresty-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'OpenResty',
        icon: 'mdi:server-network',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
