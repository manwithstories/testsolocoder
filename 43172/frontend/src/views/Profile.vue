<template>
  <div class="page-container">
    <el-row :gutter="20">
      <el-col :span="8">
        <div class="card profile-card">
          <div class="avatar-section">
            <el-avatar :size="80" :src="userStore.user?.avatar">
              {{ userStore.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <h3>{{ userStore.user?.real_name || userStore.username }}</h3>
            <p class="user-role">{{ getRoleLabel(userStore.userRole) }}</p>
          </div>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">{{ userStore.username }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ userStore.user?.email }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ userStore.user?.phone || '未填写' }}</el-descriptions-item>
            <el-descriptions-item label="信用评分">
              <el-tag :type="getCreditTagType(userStore.user?.credit_score || 100)">
                {{ userStore.user?.credit_score }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="注册时间">
              {{ formatDate(userStore.user?.created_at || '') }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </el-col>
      <el-col :span="16">
        <div class="card">
          <h2 class="section-title">编辑资料</h2>
          <el-form :model="form" label-width="100px">
            <el-form-item label="头像">
              <el-upload
                :auto-upload="false"
                :limit="1"
                list-type="picture-card"
                :on-change="handleAvatarChange"
                accept="image/*"
              >
                <el-icon><Plus /></el-icon>
              </el-upload>
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="form.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item label="收货地址">
              <el-input
                v-model="form.address"
                type="textarea"
                :rows="3"
                placeholder="请输入收货地址"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="handleSave">
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <div v-if="userStore.userRole === 'authenticator' && userStore.user?.authenticator_profile" class="card">
          <h2 class="section-title">鉴定师资质</h2>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="资质编号">
              {{ userStore.user.authenticator_profile.license_number }}
            </el-descriptions-item>
            <el-descriptions-item label="审核状态">
              <el-tag :type="getAuthStatusType(userStore.user.authenticator_profile.status)">
                {{ getAuthStatusLabel(userStore.user.authenticator_profile.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="专业领域">
              {{ userStore.user.authenticator_profile.specialties || '未填写' }}
            </el-descriptions-item>
            <el-descriptions-item label="累计鉴定">
              {{ userStore.user.authenticator_profile.completed_count }} 单
            </el-descriptions-item>
            <el-descriptions-item label="评分">
              {{ userStore.user.authenticator_profile.rating.toFixed(1) }} 分
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div v-if="userStore.userRole === 'seller' && userStore.user?.seller_profile" class="card">
          <h2 class="section-title">卖家信息</h2>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="店铺名称">
              {{ userStore.user.seller_profile.store_name || '未设置' }}
            </el-descriptions-item>
            <el-descriptions-item label="累计销量">
              {{ userStore.user.seller_profile.total_sales }} 单
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import type { UploadFile } from 'element-plus'
import dayjs from 'dayjs'
import { Plus } from '@element-plus/icons-vue'

const userStore = useUserStore()
const loading = ref(false)

const form = reactive({
  avatar: '',
  phone: '',
  address: ''
})

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    buyer: '买家',
    seller: '卖家',
    authenticator: '鉴定师',
    admin: '管理员'
  }
  return labels[role] || role
}

const getCreditTagType = (score: number) => {
  if (score >= 90) return 'success'
  if (score >= 70) return 'warning'
  return 'danger'
}

const getAuthStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return types[status] || 'info'
}

const getAuthStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return labels[status] || status
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const handleAvatarChange = (file: UploadFile) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      form.avatar = e.target?.result as string
    }
    reader.readAsDataURL(file.raw)
  }
}

const handleSave = async () => {
  loading.value = true
  try {
    const success = await userStore.updateProfile({
      avatar: form.avatar || undefined,
      phone: form.phone || undefined,
      address: form.address || undefined
    })
    if (success) {
      ElMessage.success('保存成功')
    }
  } catch (error) {
    console.error('Save error:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (userStore.user) {
    form.phone = userStore.user.phone || ''
    form.address = userStore.user.address || ''
  }
})
</script>

<style lang="scss" scoped>
.profile-card {
  .avatar-section {
    text-align: center;
    margin-bottom: 20px;
    
    .el-avatar {
      margin-bottom: 12px;
    }
    
    h3 {
      font-size: 20px;
      font-weight: 600;
      margin-bottom: 4px;
    }
    
    .user-role {
      font-size: 14px;
      color: var(--text-secondary);
    }
  }
}
</style>
