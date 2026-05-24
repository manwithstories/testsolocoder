<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <div class="profile-header">
        <el-avatar :size="80" :src="userStore.userInfo?.avatar">
          {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
        </el-avatar>
        <div class="profile-info">
          <h2>{{ userStore.userInfo?.nickname }}</h2>
          <div class="profile-meta">
            <el-tag :type="getRoleTagType(userStore.userInfo?.role)" size="small">
              {{ getRoleLabel(userStore.userInfo?.role) }}
            </el-tag>
            <el-tag v-if="userStore.userInfo?.status === 'verified'" type="success" size="small">
              已认证
            </el-tag>
            <el-tag v-else-if="userStore.userInfo?.status === 'frozen'" type="danger" size="small">
              已冻结
            </el-tag>
          </div>
        </div>
      </div>

      <div class="profile-stats">
        <div class="stat-item">
          <div class="stat-value">{{ userStore.userInfo?.order_count || 0 }}</div>
          <div class="stat-label">订单数</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">
            <el-rate :model-value="userStore.userInfo?.rating || 5" disabled size="small" />
          </div>
          <div class="stat-label">评分</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">¥{{ (userStore.userInfo?.balance || 0).toFixed(2) }}</div>
          <div class="stat-label">余额</div>
        </div>
      </div>
    </el-card>

    <el-card class="settings-card">
      <template #header>
        <div class="card-header">
          <el-icon><Setting /></el-icon>
          <span>账号设置</span>
        </div>
      </template>

      <el-form :model="profileForm" label-width="100px">
        <el-form-item label="昵称">
          <el-input v-model="profileForm.nickname" />
        </el-form-item>
        <el-form-item label="头像">
          <el-upload
            :auto-upload="false"
            :show-file-list="false"
            accept="image/*"
            @change="handleAvatarChange"
          >
            <el-avatar :size="60" :src="profileForm.avatar">
              {{ profileForm.nickname?.charAt(0) || 'U' }}
            </el-avatar>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleUpdate" :loading="updating">
            保存修改
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider />

      <div class="action-buttons">
        <el-button type="warning" @click="goToVerification">
          {{ userStore.userInfo?.status === 'verified' ? '查看认证' : '实名认证' }}
        </el-button>
        <el-button v-if="userStore.userInfo?.role === 'publisher'" type="success" @click="applyCourier">
          申请成为跑腿员
        </el-button>
        <el-button type="danger" @click="handleLogout">
          退出登录
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type UploadFile } from 'element-plus'
import { Setting } from '@element-plus/icons-vue'
import { userApi } from '@/api'
import { useUserStore } from '@/stores/user'
import type { UserRole } from '@/types'

const router = useRouter()
const userStore = useUserStore()
const updating = ref(false)

const profileForm = reactive({
  nickname: '',
  avatar: ''
})

const getRoleLabel = (role?: UserRole) => {
  const labels: Record<UserRole, string> = {
    publisher: '发布者',
    courier: '跑腿员',
    admin: '管理员'
  }
  return labels[role || 'publisher']
}

const getRoleTagType = (role?: UserRole): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  const types: Record<UserRole, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    publisher: 'primary',
    courier: 'success',
    admin: 'danger'
  }
  return types[role || 'publisher'] || 'primary'
}

const handleAvatarChange = (file: UploadFile) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      profileForm.avatar = e.target?.result as string
    }
    reader.readAsDataURL(file.raw)
  }
}

const handleUpdate = async () => {
  updating.value = true
  try {
    const res = await userApi.updateProfile({
      nickname: profileForm.nickname,
      avatar: profileForm.avatar
    })
    if (res.code === 200) {
      ElMessage.success('修改成功')
      userStore.updateUserInfo({
        nickname: profileForm.nickname,
        avatar: profileForm.avatar
      })
    }
  } catch (error) {
    console.error('Update failed:', error)
  } finally {
    updating.value = false
  }
}

const goToVerification = () => {
  router.push('/verification')
}

const applyCourier = async () => {
  router.push('/verification')
}

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/login')
}

onMounted(() => {
  if (userStore.userInfo) {
    profileForm.nickname = userStore.userInfo.nickname
    profileForm.avatar = userStore.userInfo.avatar
  }
})
</script>

<style lang="scss" scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  margin-bottom: 20px;

  .profile-header {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 20px;

    .profile-info {
      h2 {
        font-size: 24px;
        margin-bottom: 8px;
      }

      .profile-meta {
        display: flex;
        gap: 8px;
      }
    }
  }

  .profile-stats {
    display: flex;
    justify-content: space-around;
    padding-top: 20px;
    border-top: 1px solid #ebeef5;

    .stat-item {
      text-align: center;

      .stat-value {
        font-size: 24px;
        font-weight: bold;
        color: #667eea;
        margin-bottom: 4px;
      }

      .stat-label {
        color: #909399;
        font-size: 14px;
      }
    }
  }
}

.settings-card {
  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
  }

  .action-buttons {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
  }
}
</style>
