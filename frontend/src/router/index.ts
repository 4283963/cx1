import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/mobile/rules'
  },
  {
    path: '/mobile',
    redirect: '/mobile/rules'
  },
  {
    path: '/mobile/rules',
    name: 'MobileRules',
    component: () => import('@/views/mobile/LinkageRules.vue'),
    meta: { title: '联动规则配置' }
  },
  {
    path: '/mobile/rules/new',
    name: 'MobileRuleNew',
    component: () => import('@/views/mobile/LinkageRuleForm.vue'),
    meta: { title: '新建联动规则' }
  },
  {
    path: '/mobile/rules/:id/edit',
    name: 'MobileRuleEdit',
    component: () => import('@/views/mobile/LinkageRuleForm.vue'),
    meta: { title: '编辑联动规则' }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/EnvironmentMonitor.vue'),
    meta: { title: '环境数据监控大屏' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) || '智能家居控制系统'
  next()
})

export default router
