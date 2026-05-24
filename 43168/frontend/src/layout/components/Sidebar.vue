<template>
  <el-aside :width="collapsed ? '64px' : '220px'" class="sidebar">
    <div class="logo">
      <img src="/vite.svg" alt="logo" class="logo-img" v-if="!collapsed" />
      <span class="logo-text" v-if="!collapsed">管理后台</span>
      <span class="logo-text-mini" v-else>管</span>
    </div>
    <el-menu
      :default-active="activeMenu"
      :collapse="collapsed"
      :collapse-transition="false"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
      router
    >
      <el-menu-item
        v-for="menu in menus"
        :key="menu.path"
        :index="menu.path"
      >
        <el-icon><component :is="menu.icon" /></el-icon>
        <template #title>{{ menu.title }}</template>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { roleMenus } from '@/router'

const appStore = useAppStore()
const userStore = useUserStore()
const route = useRoute()

const collapsed = computed(() => appStore.sidebarCollapsed)

const activeMenu = computed(() => route.path)

const menus = computed(() => {
  if (!userStore.userInfo) return []
  return roleMenus[userStore.userInfo.role] || []
})
</script>

<style lang="scss" scoped>
.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b3648;
  overflow: hidden;

  .logo-img {
    width: 32px;
    height: 32px;
    margin-right: 8px;
  }

  .logo-text {
    color: #fff;
    font-size: 16px;
    font-weight: bold;
    white-space: nowrap;
  }

  .logo-text-mini {
    color: #fff;
    font-size: 20px;
    font-weight: bold;
  }
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-menu-item) {
  background-color: #304156 !important;

  &:hover {
    background-color: #263445 !important;
  }

  &.is-active {
    background-color: #409EFF !important;
    color: #fff !important;
  }
}
</style>