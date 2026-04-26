import { AppRouteRecord } from '@/types/router'

export const dataManageRoutes: AppRouteRecord = {
  path: '/data-manage',
  name: 'DataManage',
  component: '/index/index',
  meta: {
    title: 'menus.dataManage.title',
    icon: 'ri:database-2-line',
    roles: ['R_SUPER', 'R_ADMIN']
  },
  children: [
    {
      path: 'training-plan',
      name: 'TrainingPlan',
      component: '/data-manage/training-plan',
      meta: {
        title: 'menus.dataManage.trainingPlan',
        keepAlive: true,
        roles: ['R_SUPER', 'R_ADMIN']
      }
    },
    {
      path: 'classroom-manage',
      name: 'ClassroomManage',
      component: '/data-manage/classroom-manage',
      meta: {
        title: 'menus.dataManage.classroomManage',
        keepAlive: true,
        roles: ['R_SUPER', 'R_ADMIN']
      }
    },
    {
      path: 'class-time',
      name: 'ClassTime',
      component: '/data-manage/class-time',
      meta: {
        title: 'menus.dataManage.classTime',
        keepAlive: true,
        roles: ['R_SUPER', 'R_ADMIN']
      }
    },
    {
      path: 'class-manage',
      name: 'ClassManage',
      component: '/data-manage/class-manage',
      meta: {
        title: 'menus.dataManage.classManage',
        keepAlive: true,
        roles: ['R_SUPER', 'R_ADMIN']
      }
    },
    {
      path: 'course-library',
      name: 'CourseLibrary',
      component: '/data-manage/course-library',
      meta: {
        title: 'menus.dataManage.courseLibrary',
        keepAlive: true,
        roles: ['R_SUPER', 'R_ADMIN']
      }
    },
    {
      path: 'common-material',
      name: 'CommonMaterial',
      component: '/data-manage/common-material',
      meta: {
        title: 'menus.dataManage.commonMaterial',
        keepAlive: true,
        roles: ['R_SUPER', 'R_ADMIN']
      }
    }
  ]
}
