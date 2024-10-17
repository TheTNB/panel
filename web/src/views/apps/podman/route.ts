import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'podman',
  path: '/apps/podman',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-podman-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Podman',
        icon: 'devicon:podman',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
