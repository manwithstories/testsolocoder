<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="240px" class="admin-aside">
        <div class="admin-logo">
          <el-icon size="32"><Tickets /></el-icon>
          <span>票务管理后台</span>
        </div>
        <el-menu
          :default-active="$route.path"
          class="admin-menu"
          router
        >
          <el-menu-item index="/admin/shows">
            <el-icon><Film /></el-icon>
            <span>演出管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/orders">
            <el-icon><Document /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/statistics">
            <el-icon><DataAnalysis /></el-icon>
            <span>统计分析</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="admin-header">
          <div class="header-left">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/admin/shows' }">管理后台</el-breadcrumb-item>
              <el-breadcrumb-item>{{ $route.name }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="header-right">
            <el-button type="text" @click="$router.push('/')">
              <el-icon><House /></el-icon>
              前台首页
            </el-button>
          </div>
        </el-header>

        <el-main class="admin-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { Tickets, Film, Document, DataAnalysis, House } from '@element-plus/icons-vue'
</script>

<style lang="scss" scoped>
.admin-layout {
  min-height: 100vh;

  .admin-aside {
    background: #001529;
    color: white;

    .admin-logo {
      display: flex;
      align-items: center;
      gap: 12px;
      height: 64px;
      padding: 0 24px;
      font-size: 18px;
      font-weight: bold;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);

      el-icon {
        color: #409eff;
      }
    }

    .admin-menu {
      background: transparent;
      border-right: none;

      :deep(.el-menu-item) {
        color: rgba(255, 255, 255, 0.65);

        &:hover {
          background: rgba(255, 255, 255, 0.05);
          color: white;
        }

        &.is-active {
          background: #409eff;
          color: white;
        }
      }
    }
  }

  .admin-header {
    background: white;
    border-bottom: 1px solid #f0f0f0;
    padding: 0 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .admin-main {
    background: #f5f7fa;
    padding: 24px;
  }
}
</style>
