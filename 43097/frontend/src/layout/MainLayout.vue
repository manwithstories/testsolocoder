<template>
  <el-container class="main-layout">
    <el-aside :width="appStore.sidebarWidth" class="sidebar">
      <div class="logo" v-show="!appStore.sidebarCollapsed">
        <span class="logo-text">HMS</span>
      </div>
      <div class="logo-collapsed" v-show="appStore.sidebarCollapsed">
        <span class="logo-text">H</span>
      </div>
      <Sidebar />
    </el-aside>
    <el-container>
      <el-header class="header">
        <Header />
      </el-header>
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/store/app'
import { useUserStore } from '@/store/user'
import Sidebar from '@/components/Sidebar.vue'
import Header from '@/components/Header.vue'

const appStore = useAppStore()
const userStore = useUserStore()

onMounted(() => {
  userStore.initFromStorage()
})
</script>

<style scoped lang="scss">
.main-layout {
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  overflow: hidden;
  display: flex;
  flex-direction: column;

  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #2b2f3a;
    border-bottom: 1px solid #1f2d3d;

    .logo-text {
      font-size: 24px;
      font-weight: bold;
      color: #fff;
      letter-spacing: 2px;
    }
  }

  .logo-collapsed {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #2b2f3a;
    border-bottom: 1px solid #1f2d3d;

    .logo-text {
      font-size: 24px;
      font-weight: bold;
      color: #fff;
    }
  }
}

.header {
  padding: 0;
  height: 60px;
  background: #fff;
}

.main-content {
  background-color: #f0f2f5;
  padding: 24px;
  overflow-y: auto;
  height: calc(100vh - 60px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .main-content {
    padding: 16px;
  }
}
</style>
