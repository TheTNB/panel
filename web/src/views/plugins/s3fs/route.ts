import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 's3fs',
  path: '/plugins/s3fs',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-s3fs-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'S3fs',
        icon: 'mdi:server-network',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
