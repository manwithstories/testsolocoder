<template>
  <div class="dashboard-layout">
    <el-container>
      <el-aside width="240px" class="sidebar">
        <div class="logo">
          <el-icon size="28" color="#409EFF"><Museum /></el-icon>
          <span>博物馆管理系统</span>
        </div>
        <el-menu
          :default-active="$route.path"
          class="sidebar-menu"
          router
          background-color="#1f2937"
          text-color="#9ca3af"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataAnalysis /></el-icon>
            <span>控制面板</span>
          </el-menu-item>

          <el-sub-menu index="collection" v-if="userStore.isAdmin">
            <template #title>
              <el-icon><Picture /></el-icon>
              <span>藏品管理</span>
            </template>
            <el-menu-item index="/dashboard/collections">藏品列表</el-menu-item>
            <el-menu-item index="/dashboard/collections/categories">分类管理</el-menu-item>
          </el-sub-menu>

          <el-sub-menu index="exhibition" v-if="userStore.isAdmin">
            <template #title>
              <el-icon><Collection /></el-icon>
              <span>展览管理</span>
            </template>
            <el-menu-item index="/dashboard/exhibitions">展览列表</el-menu-item>
          </el-sub-menu>

          <el-menu-item index="/dashboard/reservations" v-if="userStore.isAdmin">
            <el-icon><Calendar /></el-icon>
            <span>预约管理</span>
          </el-menu-item>

          <el-sub-menu index="guide">
            <template #title>
              <el-icon><Microphone /></el-icon>
              <span>导览管理</span>
            </template>
            <el-menu-item index="/dashboard/guide/schedules" v-if="userStore.isGuide">导览排班</el-menu-item>
            <el-menu-item index="/dashboard/guide/contents" v-if="userStore.isAdmin || userStore.isGuide">导览内容</el-menu-item>
          </el-sub-menu>

          <el-menu-item index="/dashboard/research" v-if="userStore.isAdmin">
            <el-icon><Document /></el-icon>
            <span>学术申请</span>
          </el-menu-item>

          <el-menu-item index="/dashboard/statistics" v-if="userStore.isAdmin">
            <el-icon><TrendCharts /></el-icon>
            <span>统计分析</span>
          </el-menu-item>

          <el-menu-item index="/dashboard/museums" v-if="userStore.isAdmin">
            <el-icon><OfficeBuilding /></el-icon>
            <span>博物馆管理</span>
          </el-menu-item>

          <el-menu-item index="/dashboard/users" v-if="userStore.isAdmin">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>

          <el-sub-menu index="my">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>个人中心</span>
            </template>
            <el-menu-item index="/dashboard/profile">基本资料</el-menu-item>
            <el-menu-item index="/dashboard/my-reservations">我的预约</el-menu-item>
            <el-menu-item index="/dashboard/my-visits">参观记录</el-menu-item>
            <el-menu-item index="/dashboard/my-research">我的申请</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="dashboard-header">
          <div class="breadcrumb">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index">
                {{ item }}
              </el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-avatar :size="32" :src="userStore.user?.avatar">
                  {{ userStore.user?.nickname?.charAt(0) }}
                </el-avatar>
                <span>{{ userStore.user?.nickname }}</span>
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="home">返回前台</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <el-main class="dashboard-main">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const breadcrumbs = computed(() => {
  return route.matched
    .filter(r => r.meta?.title)
    .map(r => r.meta.title as string)
})

const handleCommand = (command: string) => {
  if (command === 'home') {
    router.push('/')
  } else if (command === 'logout') {
    userStore.logout()
    ElMessage.success('退出成功')
    router.push('/login')
  }
}
</script>

<style scoped lang="scss">
.dashboard-layout {
  min-height: 100vh;

  :deep(.el-container) {
    height: 100vh;
  }
}

.sidebar {
  background: #1f2937;

  .logo {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    color: #fff;
    font-size: 18px;
    font-weight: 600;
    border-bottom: 1px solid #374151;
  }

  .sidebar-menu {
    border: none;
    height: calc(100vh - 64px);
  }

  :deep(.el-menu-item), :deep(.el-sub-menu__title) {
    height: 50px;
    line-height: 50px;
  }
}

.dashboard-header {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;

  .header-right {
    .user-info {
      display: flex;
      align-items: center;
      gap: 10px;
      cursor: pointer;
    }
  }
}

.dashboard-main {
  background: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
