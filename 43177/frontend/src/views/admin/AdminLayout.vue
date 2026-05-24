<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="220px" class="admin-aside">
        <div class="logo">
          <h2>管理后台</h2>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="admin-menu"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
        >
          <el-menu-item index="/admin/dashboard">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据概览</span>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/technicians/verify">
            <el-icon><CircleCheck /></el-icon>
            <span>技师审核</span>
          </el-menu-item>
          <el-menu-item index="/admin/orders">
            <el-icon><Document /></el-icon>
            <span>工单管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/refunds">
            <el-icon><RefreshLeft /></el-icon>
            <span>退款审核</span>
          </el-menu-item>
          <el-menu-item index="/admin/categories">
            <el-icon><Menu /></el-icon>
            <span>分类管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/service-items">
            <el-icon><Setting /></el-icon>
            <span>服务项目</span>
          </el-menu-item>
          <el-menu-item index="/admin/parts">
            <el-icon><Tools /></el-icon>
            <span>配件管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/part-requests">
            <el-icon><Box /></el-icon>
            <span>配件申请</span>
          </el-menu-item>
          <el-menu-item index="/admin/withdraws">
            <el-icon><Money /></el-icon>
            <span>提现审核</span>
          </el-menu-item>
          <el-menu-item index="/admin/reports">
            <el-icon><TrendCharts /></el-icon>
            <span>财务报表</span>
          </el-menu-item>
          <el-menu-item index="/admin/reviews">
            <el-icon><ChatDotRound /></el-icon>
            <span>差评处理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="admin-header">
          <div class="header-left">
            <router-link to="/home">返回前台</router-link>
          </div>
          <div class="header-right">
            <el-dropdown>
              <span class="user-info">
                <el-avatar :size="32">{{ userStore.userInfo?.username?.charAt(0) }}</el-avatar>
                <span class="username">{{ userStore.userInfo?.username }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="router.push('/profile')">个人中心</el-dropdown-item>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  DataAnalysis,
  User,
  CircleCheck,
  Document,
  RefreshLeft,
  Menu,
  Setting,
  Tools,
  Box,
  Money,
  TrendCharts,
  ChatDotRound
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

function handleLogout() {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.admin-aside {
  background-color: #304156;
  overflow: auto;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.admin-menu {
  border-right: none;
}

.admin-menu :deep(.el-menu-item) {
  height: 50px;
  line-height: 50px;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}

.header-left a {
  color: #409eff;
  text-decoration: none;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.username {
  font-size: 14px;
}

.admin-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
