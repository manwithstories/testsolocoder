<template>
  <el-container class="main-layout">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <el-icon size="24"><Drone /></el-icon>
        <span>无人机租赁平台</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#001529"
        text-color="#fff"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/drones">
          <el-icon><Box /></el-icon>
          <span>设备列表</span>
        </el-menu-item>
        <el-menu-item v-if="role === 'owner'" index="/my-drones">
          <el-icon><Tools /></el-icon>
          <span>我的设备</span>
        </el-menu-item>
        <el-menu-item v-if="role === 'owner'" index="/drone/create">
          <el-icon><Plus /></el-icon>
          <span>添加设备</span>
        </el-menu-item>
        <el-menu-item index="/orders">
          <el-icon><List /></el-icon>
          <span>租赁订单</span>
        </el-menu-item>
        <el-menu-item index="/services">
          <el-icon><Service /></el-icon>
          <span>航拍服务</span>
        </el-menu-item>
        <el-menu-item v-if="role === 'client'" index="/service/create">
          <el-icon><Edit /></el-icon>
          <span>发布需求</span>
        </el-menu-item>
        <el-menu-item v-if="role === 'pilot'" index="/flights">
          <el-icon><Location /></el-icon>
          <span>飞行记录</span>
        </el-menu-item>
        <el-menu-item v-if="role === 'pilot'" index="/flight/create">
          <el-icon><AddLocation /></el-icon>
          <span>添加记录</span>
        </el-menu-item>
        <el-menu-item index="/insurance">
          <el-icon><Shield /></el-icon>
          <span>保险理赔</span>
        </el-menu-item>
        <el-menu-item index="/reviews">
          <el-icon><Star /></el-icon>
          <span>评价中心</span>
        </el-menu-item>
        <el-menu-item v-if="role === 'owner' || role === 'pilot'" index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据统计</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <span class="page-title">{{ $route.meta.title }}</span>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="userInfo?.avatar">
                {{ userInfo?.nickname?.[0] || 'U' }}
              </el-avatar>
              <span class="username">{{ userInfo?.nickname || userInfo?.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const userInfo = computed(() => userStore.userInfo)
const role = computed(() => userStore.role)

function handleCommand(command: string) {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      router.push('/login')
    }).catch(() => {})
  }
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
}
.sidebar {
  background-color: #001529;
  overflow-y: auto;
}
.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 60px;
  color: #fff;
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid #1f2d3d;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 20px;
}
.page-title {
  font-size: 16px;
  font-weight: 500;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}
.username {
  font-size: 14px;
}
.main-content {
  background: #f5f7fa;
  padding: 20px;
}
</style>
