<template>
  <div class="header-container">
    <div class="header-left">
      <div class="toggle-btn" @click="appStore.toggleSidebar()">
        <el-icon :size="20">
          <component :is="appStore.sidebarCollapsed ? 'Expand' : 'Fold'" />
        </el-icon>
      </div>
      <span class="system-title">酒店管理系统</span>
    </div>
    <div class="header-right">
      <el-dropdown @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="32" class="user-avatar">
            {{ userStore.userName?.charAt(0) || 'U' }}
          </el-avatar>
          <span class="user-name">{{ userStore.userName || '用户' }}</span>
          <span class="user-role">({{ roleText }})</span>
          <el-icon><CaretBottom /></el-icon>
        </div>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { CaretBottom, User, SwitchButton } from '@element-plus/icons-vue'
import { useAppStore } from '@/store/app'
import { useUserStore } from '@/store/user'
import { UserRole } from '@/types'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const roleText = computed(() => {
  const roleMap: Record<UserRole, string> = {
    [UserRole.ADMIN]: '超级管理员',
    [UserRole.MANAGER]: '经理',
    [UserRole.RECEPTIONIST]: '前台',
    [UserRole.STAFF]: '员工'
  }
  return roleMap[userStore.userRole as UserRole] || '用户'
})

const handleCommand = async (command: string) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
      .then(async () => {
        await userStore.logout()
        ElMessage.success('退出登录成功')
        router.push('/login')
      })
      .catch(() => {
        // 用户取消
      })
  }
}
</script>

<style scoped lang="scss">
.header-container {
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: relative;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;

  .toggle-btn {
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 6px;
    transition: all 0.3s;

    &:hover {
      background: #f5f7fa;
    }
  }

  .system-title {
    font-size: 18px;
    font-weight: 600;
    color: #303133;
  }
}

.header-right {
  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    padding: 8px 12px;
    border-radius: 6px;
    transition: all 0.3s;

    &:hover {
      background: #f5f7fa;
    }

    .user-avatar {
      background: #409EFF;
    }

    .user-name {
      font-size: 14px;
      color: #303133;
    }

    .user-role {
      font-size: 12px;
      color: #909399;
    }
  }
}

@media (max-width: 768px) {
  .header-container {
    padding: 0 16px;
  }

  .system-title {
    font-size: 16px;
  }

  .user-role {
    display: none;
  }
}
</style>
