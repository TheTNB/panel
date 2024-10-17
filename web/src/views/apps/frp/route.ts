import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'frp',
  path: '/apps/frp',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-frp-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Frp',
        icon: 'icon-park-outline:connection-box',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
