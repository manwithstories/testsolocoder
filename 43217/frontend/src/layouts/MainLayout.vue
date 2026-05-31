<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElContainer, ElHeader, ElAside, ElMain, ElMenu, ElMenuItem, ElSubMenu, ElAvatar, ElDropdown, ElDropdownMenu, ElDropdownItem, ElIcon } from 'element-plus'
import { User, Setting, SwitchButton, DataAnalysis, OfficeBuilding, UserFilled, Document, Calendar, Money, Box, Histogram, Bell, TrendCharts, Warning } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const menuItems = computed(() => {
  const role = userStore.userRole
  if (role === 'hr') {
    return [
      { path: '/hr/employees', title: '员工管理', icon: User },
      { path: '/hr/departments', title: '部门管理', icon: OfficeBuilding },
      { path: '/hr/appointments', title: '预约管理', icon: Calendar },
      { path: '/hr/reports', title: '体检报告', icon: Document },
      { path: '/hr/budget', title: '预算管理', icon: Money },
      { path: '/hr/department-appointments', title: '部门预约分配', icon: Box },
      { path: '/hr/statistics', title: '数据统计', icon: Histogram },
      { path: '/hr/billings', title: '账单管理', icon: DataAnalysis },
      { path: '/hr/transactions', title: '交易记录', icon: TrendCharts },
      { path: '/hr/balance', title: '账户余额', icon: Money }
    ]
  } else if (role === 'agency') {
    return [
      { path: '/agency/packages', title: '套餐管理', icon: Box },
      { path: '/agency/timeslots', title: '时段管理', icon: Calendar },
      { path: '/agency/appointments', title: '预约管理', icon: UserFilled },
      { path: '/agency/reports', title: '报告上传', icon: Document },
      { path: '/agency/billings', title: '账单管理', icon: DataAnalysis }
    ]
  } else if (role === 'employee') {
    return [
      { path: '/employee/new-appointment', title: '新建预约', icon: Calendar },
      { path: '/employee/appointments', title: '我的预约', icon: UserFilled },
      { path: '/employee/reports', title: '我的报告', icon: Document },
      { path: '/employee/health-records', title: '健康档案', icon: User },
      { path: '/employee/health-trend', title: '趋势分析', icon: TrendCharts },
      { path: '/employee/abnormal-items', title: '异常指标', icon: Warning },
      { path: '/employee/reminders', title: '复查提醒', icon: Bell }
    ]
  }
  return []
})

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const activeMenu = computed(() => route.path)
</script>

<template>
  <ElContainer class="main-layout">
    <ElHeader class="layout-header">
      <div class="header-left">
        <h2 class="logo">健康管理平台</h2>
      </div>
      <div class="header-right">
        <ElDropdown>
          <div class="user-info">
            <ElAvatar :size="36" :icon="User" />
            <span class="username">{{ userStore.userInfo?.real_name || userStore.userInfo?.username }}</span>
          </div>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem @click="router.push('/dashboard')">
                <ElIcon><User /></ElIcon>
                个人中心
              </ElDropdownItem>
              <ElDropdownItem divided @click="handleLogout">
                <ElIcon><SwitchButton /></ElIcon>
                退出登录
              </ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </div>
    </ElHeader>
    <ElContainer>
      <ElAside class="layout-aside" width="220px">
        <ElMenu
          :default-active="activeMenu"
          class="side-menu"
          @select="(index: string) => router.push(index)"
        >
          <ElMenuItem index="/dashboard">
            <ElIcon><DataAnalysis /></ElIcon>
            <span>首页</span>
          </ElMenuItem>
          <ElMenuItem
            v-for="item in menuItems"
            :key="item.path"
            :index="item.path"
          >
            <ElIcon><component :is="item.icon" /></ElIcon>
            <span>{{ item.title }}</span>
          </ElMenuItem>
        </ElMenu>
      </ElAside>
      <ElMain class="layout-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </ElMain>
    </ElContainer>
  </ElContainer>
</template>

<style scoped lang="scss">
.main-layout {
  height: 100vh;
}

.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(90deg, #409eff 0%, #66b1ff 100%);
  color: #fff;
  padding: 0 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

  .logo {
    margin: 0;
    font-size: 20px;
    font-weight: bold;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;

    .username {
      color: #fff;
      font-size: 14px;
    }
  }
}

.layout-aside {
  background-color: #304156;
  overflow-x: hidden;

  .side-menu {
    border-right: none;
    background-color: #304156;

    :deep(.el-menu-item) {
      color: #bfcbd9;
      background-color: #304156;

      &:hover {
        background-color: #263445;
        color: #fff;
      }

      &.is-active {
        background-color: #409eff;
        color: #fff;
      }
    }
  }
}

.layout-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
