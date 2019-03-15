<template>
  <el-scrollbar wrap-class="scrollbar-wrapper">
    <el-menu
      v-loading="loading"
      :show-timeout="200"
      :default-active="$route.path"
      :collapse="isCollapse"
      mode="vertical"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
    >
      <sidebar-item v-for="route in constantRouterMap" :key="route.path" :item="route" :base-path="route.path"/>
    </el-menu>
  </el-scrollbar>
</template>
<script>
import { mapGetters } from 'vuex'
import SidebarItem from './SidebarItem'
import { getMenu } from '@/api/login'
import Layout from './../../Layout'
import user from './../../../../store/modules/user'
export default {
  components: { SidebarItem },
  computed: {
    ...mapGetters([
      'permission_routers',
      'sidebar'
    ]),
    isCollapse() {
      return !this.sidebar.opened
    }
  },
  data() {
    return {
      loading: false,
      constantRouterMap: ''
    }
  },
  created() {
    if (sessionStorage.getItem('userId')) {
      user.state.userId = sessionStorage.getItem('userId')
      sessionStorage.setItem('userId', user.state.userId)
    } else {
      sessionStorage.setItem('userId', user.state.userId)
    }
    this.initMenu()
    // this.constantRouterMap = [
    //   {
    //     path: '/redirect',
    //     component: Layout,
    //     hidden: true,
    //     children: [
    //       {
    //         path: '/redirect/:path*',
    //         component: () => import('@/views/redirect/index')
    //       }
    //     ]
    //   },
    //   {
    //     path: '/login',
    //     component: () => import('@/views/login/index'),
    //     hidden: true
    //   },
    //   {
    //     path: '/auth-redirect',
    //     component: () => import('@/views/login/authredirect'),
    //     hidden: true
    //   },
    //   {
    //     path: '/404',
    //     component: () => import('@/views/errorPage/404'),
    //     hidden: true
    //   },
    //   {
    //     path: '/401',
    //     component: () => import('@/views/errorPage/401'),
    //     hidden: true
    //   },
    //   {
    //     path: '/dashboard',
    //     component: Layout,
    //     redirect: '/dashboard/index',
    //     children: [
    //       {
    //         path: 'index',
    //         component: () => import('@/views/dashboard/index'),
    //         name: 'Dashboard',
    //         meta: { title: 'dashboard', icon: 'dashboard', noCache: true}
    //       }
    //     ]
    //   },
    //   {
    //     path: '/userManage',
    //     component: Layout,
    //     redirect: '/userManage/index',
    //     children: [
    //       {
    //         path: 'index',
    //         component: () => import('@/views/userManage/index'),
    //         name: 'userManage',
    //         meta: { title: 'userManage', icon: 'dashboard', noCache: true}
    //       }
    //     ]
    //   },
    // ]
  },
  methods: {
    initMenu() {
      this.loading = true
      getMenu(parseInt(user.state.userId)).then(res => {
        this.loading = false
        const menu = res.Data
        // user.state.button = JSON.stringify(menu)
        // let newMenuArr = menu.filter(function(val) {
        //   if(val.Menu === null){
        //     return false
        //   }else{
        //     return true
        //   }
        // })
        const newMenu = [
          {
            path: '/redirect',
            component: Layout,
            hidden: true,
            children: [
              {
                path: '/redirect/:path*',
                component: () => import('@/views/redirect/index')
              }
            ]
          },
          {
            path: '/login',
            component: () => import('@/views/login/index'),
            hidden: true
          },
          {
            path: '/auth-redirect',
            component: () => import('@/views/login/authredirect'),
            hidden: true
          },
          {
            path: '/401',
            component: () => import('@/views/errorPage/401'),
            hidden: true
          }]
        for (let i = 0; i < menu.length; i++) {
          if (menu[i].Menu) {
            const child = []
            newMenu.push({
              path: '/' + menu[i].MenuUrl,
              component: Layout,
              redirect: '/' + menu[i].MenuUrl,
              alwaysShow: false,
              children: child
            })
            for (let j = 0; j < menu[i].Menu.length; j++) {
              child.push({
                path: menu[i].Menu[j].MenuUrl,
                component: () => resolve => require(['@/views/' + menu[i].MenuUrl + '/' + menu[i].Menu[j].MenuUrl], resolve),
                name: menu[i].Menu[j].MenuUrl,
                meta: {
                  title: menu[i].Menu[j].MenuName,
                  icon: 'dashboard'
                }
              })
            }
          } else {
            newMenu.push({
              path: '/' + menu[i].MenuUrl,
              component: Layout,
              redirect: '/' + menu[i].MenuUrl + '/index',
              alwaysShow: false,
              children: [
                {
                  path: 'index',
                  component: () => resolve => require(['@/views/' + menu[i].MenuUrl + '/index'], resolve),
                  name: menu[i].MenuUrl,
                  meta: {
                    title: menu[i].MenuName,
                    icon: 'dashboard'
                  }
                }
              ]
            })
          }
        }
        this.constantRouterMap = newMenu
        console.log(this.constantRouterMap)
      }).catch(error => {
        console.log(error)
      })
    }
  }
}
</script>
