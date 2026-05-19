<template>
  <el-container class="app-container">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <el-icon :size="28" color="#409eff"><Reading /></el-icon>
        <span>图书管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        class="sidebar-menu"
        :collapse="false"
        background-color="transparent"
        text-color="#606266"
        active-text-color="#409eff"
      >
        <template v-for="route in menuRoutes" :key="route.path">
          <el-menu-item :index="route.path">
            <el-icon><component :is="route.meta.icon" /></el-icon>
            <template #title>{{ route.meta.title }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-title">{{ currentTitle }}</div>
        <div class="header-right">
          <el-button type="primary" round @click="showAddBook">
            <el-icon><Plus /></el-icon>
            添加图书
          </el-button>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
    <AddBookDialog v-model="showAddDialog" @success="onBookAdded" />
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { routes } from '@/router'
import AddBookDialog from '@/components/AddBookDialog.vue'

const route = useRoute()
const showAddDialog = ref(false)

const activeMenu = computed(() => route.path)

const menuRoutes = computed(() => {
  return routes.filter((r: RouteRecordRaw) => !r.meta?.hidden && r.path !== '/')
})

const currentTitle = computed(() => {
  return route.meta.title || '图书管理系统'
})

const showAddBook = () => {
  showAddDialog.value = true
}

const onBookAdded = () => {
  showAddDialog.value = false
}
</script>

<style scoped lang="scss">
.app-container {
  height: 100vh;
}

.sidebar {
  background-color: #fff;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  border-bottom: 1px solid #ebeef5;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.header-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.main-content {
  padding: 24px;
  background-color: #f5f7fa;
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
