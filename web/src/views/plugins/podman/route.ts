import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'podman',
  path: '/plugins/podman',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-podman-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Podman',
        icon: 'mdi:cup-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
