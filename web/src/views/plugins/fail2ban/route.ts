import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'fail2ban',
  path: '/plugins/fail2ban',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-fail2ban-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Fail2ban',
        icon: 'mdi:wall-fire',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
