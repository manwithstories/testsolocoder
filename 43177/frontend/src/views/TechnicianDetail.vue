<template>
  <div class="technician-detail-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">技师详情</h2>
        <el-button @click="router.back()">返回</el-button>
      </div>

      <el-card v-if="profile" class="profile-card">
        <div class="profile-header">
          <el-avatar :size="100" :src="profile.user?.avatar">
            {{ profile.user?.username?.charAt(0) }}
          </el-avatar>
          <div class="profile-info">
            <h2>{{ profile.user?.real_name || profile.user?.username }}</h2>
            <div class="profile-specialty">{{ profile.specialty }}</div>
            <div class="profile-rating">
              <el-rate :model-value="profile.rating" disabled />
              <span class="rating-value">{{ profile.rating }}</span>
            </div>
            <div class="profile-stats">
              <div class="stat-item">
                <div class="stat-value">{{ profile.completed_orders }}</div>
                <div class="stat-label">完成订单</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ profile.experience_years }}</div>
                <div class="stat-label">从业年限</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ profile.service_radius }}km</div>
                <div class="stat-label">服务范围</div>
              </div>
            </div>
          </div>
        </div>

        <el-descriptions :column="2" border class="mt-20">
          <el-descriptions-item label="手机号">{{ profile.user?.phone }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ profile.user?.email || '-' }}</el-descriptions-item>
          <el-descriptions-item label="认证状态">
            <el-tag :type="profile.is_verified ? 'success' : 'warning'">
              {{ profile.is_verified ? '已认证' : '未认证' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="接单上限">{{ profile.max_active_orders }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card v-if="recentOrders.length > 0" class="mt-20">
        <template #header>
          <div class="card-header">最近订单</div>
        </template>
        <el-table :data="recentOrders" style="width: 100%">
          <el-table-column prop="order_no" label="工单号" width="200" />
          <el-table-column prop="title" label="服务内容" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :class="`status-${row.status}`" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="final_price" label="费用" width="100">
            <template #default="{ row }">
              ¥{{ row.final_price || row.quoted_price }}
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-card v-if="recentReviews.length > 0" class="mt-20">
        <template #header>
          <div class="card-header">最近评价</div>
        </template>
        <div v-for="review in recentReviews" :key="review.id" class="review-item">
          <div class="review-header">
            <el-rate :model-value="review.rating" disabled size="small" />
            <span class="review-time">{{ formatTime(review.created_at) }}</span>
          </div>
          <div class="review-content">{{ review.content }}</div>
          <div v-if="review.reply" class="review-reply">
            <strong>技师回复：</strong>{{ review.reply }}
          </div>
        </div>
      </el-card>

      <div v-if="userStore.isCustomer" class="action-bar">
        <el-button type="primary" size="large" @click="createOrder">
          预约该技师
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { userApi } from '@/api/user'
import { useUserStore } from '@/store/user'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const profile = ref<any>(null)
const recentOrders = ref<any[]>([])
const recentReviews = ref<any[]>([])

onMounted(() => {
  loadTechnicianDetail()
})

async function loadTechnicianDetail() {
  const id = route.params.id
  try {
    const res = await userApi.getTechnicianDetail(Number(id))
    if (res.data) {
      profile.value = res.data.profile
      recentOrders.value = res.data.recent_orders || []
      recentReviews.value = res.data.recent_reviews || []
    }
  } catch (error) {
    console.error('Failed to load technician detail:', error)
  }
}

function getStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    pending: '待分配',
    assigned: '待接单',
    accepted: '已接单',
    on_site: '已到达',
    repairing: '维修中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function createOrder() {
  router.push('/orders/create')
}
</script>

<style scoped>
.technician-detail-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.profile-header {
  display: flex;
  gap: 30px;
  align-items: flex-start;
}

.profile-info {
  flex: 1;
}

.profile-info h2 {
  margin: 0 0 10px 0;
  color: #303133;
}

.profile-specialty {
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
}

.profile-rating {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.rating-value {
  color: #fba934;
  font-weight: 600;
  font-size: 18px;
}

.profile-stats {
  display: flex;
  gap: 30px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  font-weight: 600;
}

.review-item {
  padding: 15px 0;
  border-bottom: 1px solid #e4e7ed;
}

.review-item:last-child {
  border-bottom: none;
}

.review-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.review-time {
  font-size: 12px;
  color: #909399;
}

.review-content {
  color: #606266;
  margin-bottom: 10px;
}

.review-reply {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  color: #606266;
}

.action-bar {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding: 20px;
}
</style>
