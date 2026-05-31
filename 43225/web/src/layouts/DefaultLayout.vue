<template>
  <div class="layout-wrapper">
    <div class="sidebar" :class="{ collapsed: isCollapsed }">
      <div class="logo">
        <el-icon :size="28"><Ship /></el-icon>
        <span v-if="!isCollapsed">船舶租赁平台</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapsed"
        router
        background-color="#001529"
        text-color="#ffffffa6"
        active-text-color="#1890ff"
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <template #title>首页</template>
        </el-menu-item>
        <el-menu-item index="/ships">
          <el-icon><Ship /></el-icon>
          <template #title>船舶列表</template>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasRole(['owner', 'admin'])" index="/my-ships">
          <el-icon><Sailboat /></el-icon>
          <template #title>我的船只</template>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasRole(['owner', 'admin'])" index="/ship-create">
          <el-icon><Plus /></el-icon>
          <template #title>发布船只</template>
        </el-menu-item>
        <el-menu-item index="/docks">
          <el-icon><OfficeBuilding /></el-icon>
          <template #title>码头列表</template>
        </el-menu-item>
        <el-menu-item index="/berth-reservations">
          <el-icon><Calendar /></el-icon>
          <template #title>泊位预约</template>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasRole(['admin'])" index="/berths">
          <el-icon><Grid /></el-icon>
          <template #title>泊位管理</template>
        </el-menu-item>
        <el-menu-item index="/rentals">
          <el-icon><Tickets /></el-icon>
          <template #title>租赁订单</template>
        </el-menu-item>
        <el-menu-item index="/my-rentals">
          <el-icon><Document /></el-icon>
          <template #title>我的租赁</template>
        </el-menu-item>
        <el-menu-item index="/voyage-logs">
          <el-icon><Notebook /></el-icon>
          <template #title>航海日志</template>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasRole(['owner', 'admin'])" index="/maintenance">
          <el-icon><Tools /></el-icon>
          <template #title>维修保养</template>
        </el-menu-item>
        <el-menu-item index="/finance">
          <el-icon><Money /></el-icon>
          <template #title>财务结算</template>
        </el-menu-item>
        <el-menu-item index="/reviews">
          <el-icon><Star /></el-icon>
          <template #title>评价管理</template>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasRole(['admin'])" index="/admin/users">
          <el-icon><UserFilled /></el-icon>
          <template #title>用户管理</template>
        </el-menu-item>
      </el-menu>
    </div>
    <div class="main-container">
      <div class="header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="toggleCollapse">
            <Fold v-if="!isCollapsed" />
            <Expand v-else />
          </el-icon>
        </div>
        <div class="header-right">
          <el-dropdown trigger="click" @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="32" :src="userStore.user?.avatar_url">
                {{ userStore.user?.full_name?.[0] }}
              </el-avatar>
              <span class="username">{{ userStore.user?.full_name }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      <div class="content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapsed = ref(false)
const activeMenu = computed(() => route.path)

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

const handleCommand = (command: string) => {
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

<style lang="scss" scoped>
@use '@/styles/variables.scss' as *;

.layout-wrapper {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: $sidebar-width;
  background: #001529;
  transition: width $transition-duration;
  overflow: hidden;

  &.collapsed {
    width: $sidebar-collapsed-width;
  }

  .logo {
    height: $header-height;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-size: 18px;
    font-weight: 600;
    gap: 8px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  :deep(.el-menu) {
    border-right: none;
  }
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.header {
  height: $header-height;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  z-index: 10;

  .header-left {
    display: flex;
    align-items: center;

    .collapse-btn {
      font-size: 20px;
      cursor: pointer;
      color: rgba(0, 0, 0, 0.65);

      &:hover {
        color: #1890ff;
      }
    }
  }

  .header-right {
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 8px 12px;
      border-radius: 4px;

      &:hover {
        background: rgba(0, 0, 0, 0.04);
      }

      .username {
        font-size: 14px;
      }
    }
  }
}

.content {
  flex: 1;
  padding: 24px;
  background: #f0f2f5;
  overflow-y: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
