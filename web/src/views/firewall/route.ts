import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'firewall',
  path: '/firewall',
  component: Layout,
  meta: {
    order: 30
  },
  children: [
    {
      name: 'firewall-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: '防火墙',
        icon: 'mdi:shield-check-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
