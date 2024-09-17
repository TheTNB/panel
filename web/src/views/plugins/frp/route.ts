import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'frp',
  path: '/plugins/frp',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-frp-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Frp',
        icon: 'mdi:virtual-private-network',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
