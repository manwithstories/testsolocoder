<template>
  <div id="app">
    <el-container>
      <el-header v-if="showHeader">
        <div class="header-container">
          <div class="logo" @click="goHome">
            <el-icon><Camera /></el-icon>
            <span>摄影器材租赁平台</span>
          </div>
          <div class="nav-menu">
            <el-menu
              :default-active="activeMenu"
              mode="horizontal"
              @select="handleMenuSelect"
              :ellipsis="false"
            >
              <el-menu-item index="/">首页</el-menu-item>
              <el-menu-item index="/equipments">设备列表</el-menu-item>
              <el-menu-item v-if="userStore.isLoggedIn" index="/orders">我的订单</el-menu-item>
              <el-menu-item v-if="userStore.isOwner()" index="/my-equipments">我的设备</el-menu-item>
              <el-menu-item v-if="userStore.isOwner()" index="/export">数据导出</el-menu-item>
              <el-menu-item v-if="userStore.isAdmin()" index="/admin/users">用户管理</el-menu-item>
            </el-menu>
          </div>
          <div class="user-actions">
            <template v-if="userStore.isLoggedIn">
              <el-dropdown @command="handleUserCommand">
                <span class="user-info">
                  <el-avatar :size="32" :src="userStore.user?.avatar">
                    {{ userStore.user?.username?.charAt(0)?.toUpperCase() }}
                  </el-avatar>
                  <span class="username">{{ userStore.user?.username }}</span>
                  <el-icon><ArrowDown /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                    <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
            <template v-else>
              <el-button type="primary" @click="goLogin">登录</el-button>
              <el-button @click="goRegister">注册</el-button>
            </template>
          </div>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const showHeader = computed(() => {
  return !['/login', '/register'].includes(route.path)
})

const activeMenu = computed(() => route.path)

onMounted(() => {
  if (userStore.token) {
    userStore.loadUser()
  }
})

function goHome() {
  router.push('/')
}

function goLogin() {
  router.push('/login')
}

function goRegister() {
  router.push('/register')
}

function handleMenuSelect(index: string) {
  router.push(index)
}

function handleUserCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      userStore.logout()
      router.push('/login')
      break
  }
}
</script>

<style scoped>
.el-header {
  padding: 0;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  color: #409eff;
}

.nav-menu {
  flex: 1;
  margin: 0 20px;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 0 10px;
}

.username {
  font-size: 14px;
}

.el-main {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
</style>
