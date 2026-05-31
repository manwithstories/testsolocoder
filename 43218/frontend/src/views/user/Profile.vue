<template>
  <div class="profile-page">
    <div class="page-header">
      <h2 class="page-title">个人中心</h2>
    </div>

    <div class="profile-content">
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-button type="primary" link @click="showEditDialog = true">编辑</el-button>
          </div>
        </template>
        <div class="info-row">
          <span class="label">用户名：</span>
          <span>{{ userStore.userInfo?.username }}</span>
        </div>
        <div class="info-row">
          <span class="label">昵称：</span>
          <span>{{ userStore.userInfo?.nickname || '未设置' }}</span>
        </div>
        <div class="info-row">
          <span class="label">邮箱：</span>
          <span>{{ userStore.userInfo?.email || '未设置' }}</span>
        </div>
        <div class="info-row">
          <span class="label">手机号：</span>
          <span>{{ userStore.userInfo?.phone || '未设置' }}</span>
        </div>
        <div class="info-row">
          <span class="label">角色：</span>
          <el-tag>{{ getRoleText(userStore.userInfo?.role) }}</el-tag>
        </div>
        <div class="info-row">
          <span class="label">信用分：</span>
          <span class="credit-score">{{ userStore.userInfo?.creditScore }}</span>
          <el-rate :model-value="getRating(userStore.userInfo?.creditScore)" disabled size="small" />
        </div>
        <div class="info-row">
          <span class="label">实名认证：</span>
          <el-tag :type="userStore.userInfo?.isAuthenticated ? 'success' : 'warning'">
            {{ userStore.userInfo?.isAuthenticated ? '已认证' : '未认证' }}
          </el-tag>
          <el-button
            v-if="!userStore.userInfo?.isAuthenticated"
            type="primary"
            link
            @click="showAuthDialog = true"
          >
            去认证
          </el-button>
        </div>
      </el-card>

      <el-card class="stats-card">
        <template #header>
          <span>数据统计</span>
        </template>
        <div class="stats-grid" v-loading="loading">
          <div class="stat-item">
            <div class="stat-value">{{ userStats.sellOrders || 0 }}</div>
            <div class="stat-label">卖出订单</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ userStats.buyOrders || 0 }}</div>
            <div class="stat-label">购买订单</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ userStats.products || 0 }}</div>
            <div class="stat-label">发布商品</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ userStats.reviews || 0 }}</div>
            <div class="stat-label">评价数量</div>
          </div>
        </div>
      </el-card>
    </div>

    <el-dialog v-model="showEditDialog" title="编辑信息" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="editForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="editForm.phone" placeholder="请输入手机号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitEdit">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showAuthDialog" title="实名认证" width="500px">
      <el-form :model="authForm" label-width="80px">
        <el-form-item label="真实姓名">
          <el-input v-model="authForm.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="身份证号">
          <el-input v-model="authForm.idCard" placeholder="请输入身份证号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAuthDialog = false">取消</el-button>
        <el-button type="primary" :loading="submittingAuth" @click="submitAuth">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { userApi } from '@/api'
import type { UserStats } from '@/types'

const userStore = useUserStore()

const loading = ref(false)
const showEditDialog = ref(false)
const showAuthDialog = ref(false)
const submitting = ref(false)
const submittingAuth = ref(false)
const userStats = ref<UserStats>({
  sellOrders: 0,
  buyOrders: 0,
  products: 0,
  reviews: 0
})

const editForm = reactive({
  nickname: userStore.userInfo?.nickname || '',
  email: userStore.userInfo?.email || '',
  phone: userStore.userInfo?.phone || ''
})

const authForm = reactive({
  realName: '',
  idCard: ''
})

function getRoleText(role?: string): string {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    seller: '卖家',
    buyer: '买家',
    technician: '维修技师'
  }
  return roleMap[role || ''] || '用户'
}

function getRating(creditScore?: number): number {
  if (!creditScore) return 3
  if (creditScore >= 90) return 5
  if (creditScore >= 80) return 4
  if (creditScore >= 70) return 3
  if (creditScore >= 60) return 2
  return 1
}

async function fetchUserStats() {
  loading.value = true
  try {
    const res = await userApi.getUserStats()
    userStats.value = res.data
  } catch (error) {
    console.error('Failed to fetch user stats:', error)
  } finally {
    loading.value = false
  }
}

async function submitEdit() {
  submitting.value = true
  try {
    await userStore.updateProfile({
      nickname: editForm.nickname || undefined,
      email: editForm.email || undefined,
      phone: editForm.phone || undefined
    })
    ElMessage.success('更新成功')
    showEditDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    submitting.value = false
  }
}

async function submitAuth() {
  if (!authForm.realName || !authForm.idCard) {
    ElMessage.warning('请填写完整信息')
    return
  }

  submittingAuth.value = true
  try {
    await userApi.submitRealNameAuth(authForm)
    ElMessage.success('认证申请已提交，请等待审核')
    showAuthDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '提交失败')
  } finally {
    submittingAuth.value = false
  }
}

onMounted(() => {
  fetchUserStats()
})
</script>

<style lang="scss" scoped>
.profile-page {
  .profile-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }

  .info-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .info-row {
      display: flex;
      align-items: center;
      margin-bottom: 16px;

      .label {
        width: 100px;
        color: var(--text-lighter-color);
      }

      .credit-score {
        font-size: 18px;
        font-weight: 600;
        color: var(--warning-color);
        margin-right: 12px;
      }
    }
  }

  .stats-card {
    .stats-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 20px;

      .stat-item {
        text-align: center;
        padding: 20px;
        background: #f5f7fa;
        border-radius: 8px;

        .stat-value {
          font-size: 28px;
          font-weight: 600;
          color: var(--primary-color);
        }

        .stat-label {
          margin-top: 8px;
          color: var(--text-lighter-color);
          font-size: 14px;
        }
      }
    }
  }
}
</style>
