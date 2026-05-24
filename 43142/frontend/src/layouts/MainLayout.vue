<template>
  <el-container class="main-layout">
    <el-header class="header">
      <div class="header-content">
        <div class="logo" @click="goHome">
          <el-icon><Briefcase /></el-icon>
          <span>招聘平台</span>
        </div>
        <el-menu
          mode="horizontal" :default-active="activeMenu" class="nav-menu" @select="handleMenuSelect">
          <el-menu-item index="jobs">
            <el-icon><Opportunity /></el-icon>
            <span>职位</span>
          </el-menu-item>
          <el-menu-item v-if="isCompany" index="my-jobs">
            <el-icon><Management /></el-icon>
            <span>我的职位</span>
          </el-menu-item>
          <el-menu-item v-if="isApplicant" index="resumes">
            <el-icon><Document /></el-icon>
            <span>我的简历</span>
          </el-menu-item>
          <el-menu-item v-if="isApplicant" index="applications">
            <el-icon><Tickets /></el-icon>
            <span>我的投递</span>
          </el-menu-item>
          <el-menu-item v-if="isCompany" index="company-applications">
            <el-icon><Tickets /></el-icon>
            <span>收到简历</span>
          </el-menu-item>
          <el-menu-item v-if="isApplicant" index="interviews">
            <el-icon><Calendar /></el-icon>
            <span>我的面试</span>
          </el-menu-item>
          <el-menu-item v-if="isCompany" index="company-interviews">
            <el-icon><Calendar /></el-icon>
            <span>面试管理</span>
          </el-menu-item>
          <el-menu-item v-if="isCompany || isAdmin" index="stats">
            <el-icon><DataAnalysis /></el-icon>
            <span>统计分析</span>
          </el-menu-item>
        </el-menu>
        <div class="user-info">
          <el-dropdown v-if="userStore.isLoggedIn" @command="handleUserCommand">
            <span class="user-name">
              <el-avatar :size="28">{{ userNameFirstChar }}</el-avatar>
              {{ userStore.userName }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button v-else type="primary" @click="goLogin">登录</el-button>
          <el-button v-if="!userStore.isLoggedIn" type="success" @click="goRegister">注册</el-button>
        </div>
      </div>
    </el-header>
    <el-main class="main-content">
      <slot />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCompany = computed(() => userStore.hasRole('company'))
const isApplicant = computed(() => userStore.hasRole('applicant'))
const isAdmin = computed(() => userStore.hasRole('admin'))

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/my/jobs')) return 'my-jobs'
  if (path.startsWith('/resumes')) return 'resumes'
  if (path.startsWith('/applications') && isApplicant.value) return 'applications'
  if (path.startsWith('/company/applications')) return 'company-applications'
  if (path.startsWith('/interviews')) return 'interviews'
  if (path.startsWith('/company/interviews')) return 'company-interviews'
  if (path.startsWith('/stats')) return 'stats'
  if (path.startsWith('/jobs')) return 'jobs'
  return 'jobs'
})

const userNameFirstChar = computed(() => {
  return userStore.userName.charAt(0).toUpperCase()
})

function handleMenuSelect(index: string) {
  switch (index) {
    case 'jobs':
      router.push('/jobs')
      break
    case 'my-jobs':
      router.push('/my/jobs')
      break
    case 'resumes':
      router.push('/resumes')
      break
    case 'applications':
      router.push('/applications')
      break
    case 'company-applications':
      router.push('/company/applications')
      break
    case 'interviews':
      router.push('/interviews')
      break
    case 'company-interviews':
      router.push('/company/interviews')
      break
    case 'stats':
      router.push('/stats')
      break
  }
}

function handleUserCommand(command: string) {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
      userStore.logout()
      router.push('/login')
    }
  }

function goHome() {
    router.push('/jobs')
  }

function goLogin() {
    router.push('/login')
  }

function goRegister() {
    router.push('/register')
  }
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
  height: 60px;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: 600;
  color: #409eff;
  cursor: pointer;
  margin-right: 40px;
}

.logo .el-icon {
  font-size: 28px;
  margin-right: 8px;
}

.nav-menu {
  flex: 1;
  border-bottom: none;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-name {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.user-name .el-avatar {
  margin-right: 8px;
}

.main-content {
  padding: 0;
}
</style>
