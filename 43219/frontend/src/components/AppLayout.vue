<template>
  <el-container class="layout">
    <el-header class="header">
      <div class="logo">家政服务平台</div>
      <el-menu mode="horizontal" :default-active="activeMenu" @select="onSelect" class="menu">
        <el-menu-item index="/services">服务项目</el-menu-item>
        <el-menu-item v-if="userStore.role==='customer'" index="/booking/new">预约服务</el-menu-item>
        <el-menu-item v-if="userStore.isLoggedIn" index="/bookings">我的预约</el-menu-item>
        <el-menu-item v-if="userStore.isLoggedIn" index="/orders">订单管理</el-menu-item>
        <el-menu-item v-if="userStore.role==='customer'" index="/tickets">我的工单</el-menu-item>
        <el-menu-item v-if="userStore.role==='company'" index="/company/services">公司管理</el-menu-item>
        <el-menu-item v-if="userStore.role==='company'" index="/company/bookings">预约审核</el-menu-item>
        <el-menu-item v-if="userStore.role==='company'" index="/company/finance">财务</el-menu-item>
        <el-menu-item v-if="userStore.role==='company'" index="/company/stats">统计</el-menu-item>
        <el-menu-item v-if="userStore.role==='staff'" index="/staff/profile">个人中心</el-menu-item>
        <el-menu-item v-if="userStore.role==='staff'" index="/staff/schedule">我的档期</el-menu-item>
        <el-menu-item v-if="userStore.role==='staff'" index="/staff/orders">接单处理</el-menu-item>
        <el-menu-item v-if="userStore.role==='staff'" index="/staff/earnings">我的收益</el-menu-item>
        <el-menu-item v-if="userStore.role==='admin'" index="/admin/stats">运营统计</el-menu-item>
        <el-menu-item v-if="userStore.role==='admin'" index="/admin/tickets">工单处理</el-menu-item>
        <el-menu-item v-if="userStore.role==='admin'" index="/admin/staff">人员管理</el-menu-item>
      </el-menu>
      <div class="user-area">
        <template v-if="userStore.isLoggedIn">
          <span class="muted">{{ userStore.user?.real_name || userStore.user?.username }}</span>
          <el-tag size="small" style="margin-left:8px">{{ roleLabel }}</el-tag>
          <el-button size="small" text @click="onLogout">退出</el-button>
        </template>
        <template v-else>
          <el-button size="small" @click="$router.push('/login')">登录</el-button>
        </template>
      </div>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const route = useRoute()
const router = useRouter()
const activeMenu = computed(() => route.path)

const roleLabel = computed(() => ({
  company: '家政公司',
  staff: '家政人员',
  customer: '服务客户',
  admin: '管理员',
} as Record<string, string>)[userStore.role] || '')

function onSelect(idx: string) {
  router.push(idx)
}

function onLogout() {
  userStore.logout()
  router.replace('/login')
}
</script>

<style scoped>
.layout { min-height: 100vh; }
.header {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #eee;
  padding: 0 24px;
}
.logo {
  font-weight: 700;
  font-size: 18px;
  color: #409eff;
  margin-right: 24px;
  white-space: nowrap;
}
.menu { flex: 1; border-bottom: none; }
.user-area { display: flex; align-items: center; white-space: nowrap; }
</style>
