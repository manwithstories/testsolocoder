<template>
  <div class="layout-container">
    <el-aside class="layout-sidebar">
      <div style="height: 60px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 18px; font-weight: bold; border-bottom: 1px solid #1f3a5f;">
        问卷调查平台
      </div>
      <el-menu
        :default-active="activeMenu"
        class="el-menu-vertical"
        background-color="#001529"
        text-color="#fff"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/surveys">
          <el-icon><Document /></el-icon>
          <span>问卷列表</span>
        </el-menu-item>
        <el-menu-item index="/surveys/create">
          <el-icon><EditPen /></el-icon>
          <span>创建问卷</span>
        </el-menu-item>
        <el-menu-item v-if="isAdmin" index="/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="layout-main">
      <el-header class="layout-header">
        <el-dropdown @command="handleCommand">
          <span class="el-dropdown-link" style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
            <el-avatar :size="32" :src="user?.avatar">
              {{ user?.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <span>{{ user?.nickname || '用户' }}</span>
            <el-icon><ArrowDown /></el-icon>
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
      </el-header>

      <el-main class="layout-content">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store'
import { authApi } from '@/api/auth'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const user = computed(() => userStore.user)
const isAdmin = computed(() => user.value?.role === 'admin')
const activeMenu = computed(() => route.path)

onMounted(async () => {
  if (!user.value) {
    try {
      const profile = await authApi.getProfile()
      userStore.setUserInfo(profile)
    } catch (e) {
      console.error('Failed to load profile')
    }
  }
})

const handleCommand = (command: string) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.clearUser()
    router.push('/login')
  }
}
</script>

<style scoped>
.layout-sidebar {
  width: 220px;
}

.layout-main {
  display: flex;
  flex-direction: column;
}

.layout-header {
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  padding: 0 20px;
  justify-content: flex-end;
}

.layout-content {
  flex: 1;
  background: #f5f7fa;
}
</style>
