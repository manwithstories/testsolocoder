<template>
  <div class="dashboard">
    <div class="page-header">
      <h2 class="page-title">工作台</h2>
      <p class="text-muted">欢迎回来，{{ userStore.userInfo?.realName }}！</p>
    </div>

    <el-row :gutter="20" class="mb-24">
      <el-col :span="6">
        <div class="stat-card primary">
          <div class="stat-label">我的申请</div>
          <div class="stat-value">{{ stats.totalApplications }}</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card warning">
          <div class="stat-label">处理中</div>
          <div class="stat-value">{{ stats.processingApps }}</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card success">
          <div class="stat-label">已完成</div>
          <div class="stat-value">{{ stats.completedApps }}</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card danger">
          <div class="stat-label">已驳回</div>
          <div class="stat-value">{{ stats.rejectedApps }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="card">
          <div class="flex-between mb-16">
            <h3 class="font-weight-600">最近申请</h3>
            <el-button type="primary" link @click="goToApplications">查看全部</el-button>
          </div>
          <el-table :data="recentApplications" style="width: 100%" v-loading="loading">
            <el-table-column prop="applicationNo" label="申请编号" width="140" />
            <el-table-column prop="companyName" label="公司名称" />
            <el-table-column prop="companyType" label="公司类型" width="120">
              <template #default="{ row }">
                {{ getCompanyTypeText(row.companyType) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <span :class="['status-tag', getStatusClass(row.status)]">
                  {{ getStatusText(row.status) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button type="primary" link @click="viewDetail(row.id)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="card">
          <div class="flex-between mb-16">
            <h3 class="font-weight-600">最新通知</h3>
            <el-button type="primary" link @click="goToNotifications">查看全部</el-button>
          </div>
          <div class="notification-list">
            <div
              v-for="notification in recentNotifications"
              :key="notification.id"
              class="notification-item"
              :class="{ unread: !notification.isRead }"
              @click="viewNotification(notification.id)"
            >
              <div class="notification-title truncate">{{ notification.title }}</div>
              <div class="notification-time text-muted">{{ formatTime(notification.createdAt) }}</div>
            </div>
            <div v-if="recentNotifications.length === 0" class="empty-state">
              <el-icon :size="48"><Bell /></el-icon>
              <p>暂无通知</p>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { applicationApi } from '@/api/application'
import { notificationApi } from '@/api/notification'
import { Application, Notification, ApplicationStatus, CompanyType } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const recentApplications = ref<Application[]>([])
const recentNotifications = ref<Notification[]>([])

const stats = computed(() => ({
  totalApplications: recentApplications.value.length,
  processingApps: recentApplications.value.filter(a => a.status === 'processing').length,
  completedApps: recentApplications.value.filter(a => a.status === 'completed').length,
  rejectedApps: recentApplications.value.filter(a => a.status === 'rejected').length
}))

const fetchData = async () => {
  loading.value = true
  try {
    const [appRes, notifRes] = await Promise.all([
      applicationApi.getList({ page: 1, pageSize: 5 }),
      notificationApi.getList({ page: 1, pageSize: 5 })
    ])
    recentApplications.value = appRes?.list || []
    recentNotifications.value = notifRes?.list || []
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

const getStatusText = (status: ApplicationStatus) => {
  const map: Record<ApplicationStatus, string> = {
    draft: '草稿',
    pending_review: '待审核',
    reviewing: '审核中',
    processing: '处理中',
    completed: '已完成',
    rejected: '已驳回',
    cancelled: '已取消',
    payment_pending: '待支付'
  }
  return map[status] || status
}

const getStatusClass = (status: ApplicationStatus) => {
  const map: Record<ApplicationStatus, string> = {
    draft: 'draft',
    pending_review: 'pending',
    reviewing: 'pending',
    processing: 'processing',
    completed: 'completed',
    rejected: 'rejected',
    cancelled: 'cancelled',
    payment_pending: 'pending'
  }
  return map[status] || ''
}

const getCompanyTypeText = (type: CompanyType) => {
  const map: Record<CompanyType, string> = {
    llc: '有限责任公司',
    joint_stock: '股份有限公司',
    sole: '个人独资',
    partnership: '合伙企业'
  }
  return map[type] || type
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

const goToApplications = () => {
  router.push('/applications')
}

const goToNotifications = () => {
  router.push('/notifications')
}

const viewDetail = (id: number) => {
  router.push(`/applications/${id}`)
}

const viewNotification = async (id: number) => {
  try {
    await notificationApi.markAsRead(id)
    router.push('/notifications')
  } catch (error) {
    console.error('标记已读失败:', error)
  }
}

onMounted(fetchData)
</script>

<style scoped>
.notification-list {
  max-height: 300px;
  overflow-y: auto;
}

.notification-item {
  padding: 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s;
}

.notification-item:hover {
  background: #f5f7fa;
}

.notification-item.unread {
  background: #ecf5ff;
}

.notification-title {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}

.notification-time {
  font-size: 12px;
}
</style>
