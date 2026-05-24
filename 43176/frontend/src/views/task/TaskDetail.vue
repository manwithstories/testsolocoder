<template>
  <div class="task-detail-container">
    <div v-if="task" class="task-detail">
      <el-card class="task-card">
        <div class="task-header">
          <span :class="['task-type-tag', task.type]">
            {{ getTaskTypeLabel(task.type) }}
          </span>
          <span :class="['status-tag', task.status]">
            {{ getStatusLabel(task.status) }}
          </span>
        </div>

        <h2 class="task-title">{{ task.title }}</h2>
        <p class="task-desc">{{ task.description }}</p>

        <div class="task-meta">
          <div class="meta-item">
            <el-icon><Location /></el-icon>
            <span>起点: {{ task.start_addr }}</span>
          </div>
          <div class="meta-item">
            <el-icon><Location /></el-icon>
            <span>终点: {{ task.end_addr }}</span>
          </div>
          <div class="meta-item">
            <el-icon><Clock /></el-icon>
            <span>截止: {{ formatTime(task.deadline) }}</span>
          </div>
          <div class="meta-item">
            <el-icon><Money /></el-icon>
            <span class="reward">报酬: ¥{{ task.reward }}</span>
          </div>
        </div>

        <div v-if="task.images && task.images.length > 0" class="task-images">
          <el-image
            v-for="img in task.images"
            :key="img.id"
            :src="img.image_url"
            :preview-src-list="task.images.map(i => i.image_url)"
            fit="cover"
            class="task-image"
          />
        </div>

        <div class="publisher-info">
          <el-avatar :size="40" :src="task.publisher?.avatar">
            {{ task.publisher?.nickname?.charAt(0) }}
          </el-avatar>
          <div class="publisher-detail">
            <span class="publisher-name">{{ task.publisher?.nickname }}</span>
            <span class="publisher-rating">
              <el-rate :model-value="task.publisher?.rating" disabled size="small" />
            </span>
          </div>
        </div>

        <div v-if="task.courier" class="courier-info">
          <el-avatar :size="40" :src="task.courier?.avatar">
            {{ task.courier?.nickname?.charAt(0) }}
          </el-avatar>
          <div class="courier-detail">
            <span class="courier-name">{{ task.courier?.nickname }}</span>
            <span class="courier-rating">
              <el-rate :model-value="task.courier?.rating" disabled size="small" />
            </span>
          </div>
        </div>
      </el-card>

      <div v-if="task.order" class="order-card">
        <el-card>
          <template #header>
            <div class="card-header">
              <el-icon><List /></el-icon>
              <span>订单信息</span>
            </div>
          </template>
          <div class="order-info">
            <div class="info-item">
              <span>订单状态:</span>
              <span :class="['status-tag', task.order.status]">
                {{ getOrderStatusLabel(task.order.status) }}
              </span>
            </div>
            <div class="info-item">
              <span>开始时间:</span>
              <span>{{ formatTime(task.order.start_time) }}</span>
            </div>
            <div class="info-item">
              <span>完成时间:</span>
              <span>{{ formatTime(task.order.end_time) }}</span>
            </div>
            <div class="info-item">
              <span>服务费:</span>
              <span>¥{{ task.order.service_fee }}</span>
            </div>
            <div class="info-item">
              <span>实际收入:</span>
              <span class="reward">¥{{ task.order.actual_payment }}</span>
            </div>
          </div>

          <div v-if="task.order.proof_images && task.order.proof_images.length > 0" class="proof-images">
            <h4>送达凭证</h4>
            <el-image
              v-for="(img, index) in task.order.proof_images"
              :key="index"
              :src="img.image_url"
              fit="cover"
              class="proof-image"
              :preview-src-list="task.order.proof_images.map(i => i.image_url)"
            />
          </div>
        </el-card>
      </div>

      <div class="action-bar">
        <template v-if="task.status === 'pending' && canAccept">
          <el-button type="primary" size="large" @click="handleAccept">
            立即接单
          </el-button>
        </template>

        <template v-else-if="task.status === 'accepted' && isCourier">
          <el-button type="primary" size="large" @click="handleStart">
            开始任务
          </el-button>
        </template>

        <template v-else-if="task.status === 'in_progress' && isCourier">
          <el-upload
            :auto-upload="false"
            :show-file-list="false"
            @change="handleUploadProof"
          >
            <el-button type="primary" size="large">
              上传凭证完成任务
            </el-button>
          </el-upload>
        </template>

        <template v-else-if="task.status === 'completed' && isPublisher">
          <el-button type="primary" size="large" @click="goToReview">
            评价跑腿员
          </el-button>
        </template>

        <template v-if="task.status === 'pending' && isPublisher">
          <el-button type="danger" size="large" @click="handleCancel">
            取消任务
          </el-button>
        </template>

        <el-button v-if="task.order" size="large" @click="goToOrder">
          查看订单详情
        </el-button>
      </div>
    </div>

    <el-empty v-else description="任务不存在" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type UploadFile } from 'element-plus'
import { Location, Clock, Money, List } from '@element-plus/icons-vue'
import { taskApi } from '@/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import type { Task, TaskType, TaskStatus, OrderStatus } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const task = ref<Task | null>(null)
const proofImages = ref<string[]>([])

const taskId = computed(() => Number(route.params.id))
const isPublisher = computed(() => task.value?.publisher_id === userStore.userInfo?.id)
const isCourier = computed(() => task.value?.courier_id === userStore.userInfo?.id)
const canAccept = computed(() =>
  userStore.isCourier && !isPublisher.value && userStore.userInfo?.status === 'verified'
)

const taskTypeLabels: Record<TaskType, string> = {
  buy: '代购',
  pickup: '代取',
  deliver: '代送',
  queue: '排队代办',
  errand: '其他代办'
}

const statusLabels: Record<TaskStatus, string> = {
  pending: '待接单',
  accepted: '已接单',
  in_progress: '进行中',
  completed: '已完成',
  cancelled: '已取消',
  timeout: '已超时'
}

const orderStatusLabels: Record<OrderStatus, string> = {
  pending: '待处理',
  accepted: '已接单',
  in_progress: '进行中',
  delivered: '已送达',
  completed: '已完成',
  cancelled: '已取消'
}

const getTaskTypeLabel = (type: TaskType) => taskTypeLabels[type] || '代办'
const getStatusLabel = (status: TaskStatus) => statusLabels[status] || status
const getOrderStatusLabel = (status: OrderStatus) => orderStatusLabels[status] || status

const formatTime = (time?: string) => {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

const fetchTask = async () => {
  try {
    const res = await taskApi.get(taskId.value)
    if (res.code === 200) {
      task.value = res.data as Task
    }
  } catch (error) {
    console.error('Failed to fetch task:', error)
  }
}

const handleAccept = async () => {
  try {
    const res = await taskApi.accept(taskId.value)
    if (res.code === 200) {
      ElMessage.success('接单成功')
      fetchTask()
    }
  } catch (error) {
    console.error('Failed to accept task:', error)
  }
}

const handleStart = async () => {
  try {
    ElMessage.info('请在订单详情中上报位置和进度')
    router.push(`/orders/${task.value?.order?.id}`)
  } catch (error) {
    console.error('Failed to start task:', error)
  }
}

const handleUploadProof = (file: UploadFile) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      proofImages.value.push(e.target?.result as string)
      completeTask()
    }
    reader.readAsDataURL(file.raw)
  }
}

const completeTask = async () => {
  try {
    const res = await taskApi.complete(taskId.value, { proof_images: proofImages.value })
    if (res.code === 200) {
      ElMessage.success('任务已完成')
      fetchTask()
    }
  } catch (error) {
    console.error('Failed to complete task:', error)
  }
}

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该任务吗？', '提示', {
      type: 'warning'
    })
    const res = await taskApi.cancel(taskId.value, { reason: '用户取消' })
    if (res.code === 200) {
      ElMessage.success('任务已取消')
      router.push('/my-tasks')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to cancel task:', error)
    }
  }
}

const goToReview = () => {
  router.push(`/reviews?order_id=${task.value?.order?.id}`)
}

const goToOrder = () => {
  router.push(`/orders/${task.value?.order?.id}`)
}

onMounted(() => {
  fetchTask()
})
</script>

<style lang="scss" scoped>
.task-detail-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.task-card,
.order-card {
  margin-bottom: 20px;
}

.task-header {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.task-title {
  font-size: 22px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.task-desc {
  color: #606266;
  font-size: 15px;
  line-height: 1.6;
  margin-bottom: 20px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.task-meta {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 20px;

  .meta-item {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #606266;

    .el-icon {
      color: #667eea;
    }

    .reward {
      color: #ff5722;
      font-weight: bold;
      font-size: 18px;
    }
  }
}

.task-images {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;

  .task-image {
    width: 100px;
    height: 100px;
    border-radius: 8px;
  }
}

.publisher-info,
.courier-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-top: 16px;
}

.publisher-detail,
.courier-detail {
  display: flex;
  flex-direction: column;

  .publisher-name,
  .courier-name {
    font-weight: 500;
    color: #303133;
  }
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.order-info {
  .info-item {
    display: flex;
    justify-content: space-between;
    padding: 8px 0;
    border-bottom: 1px solid #ebeef5;

    &:last-child {
      border-bottom: none;
    }

    .reward {
      color: #ff5722;
      font-weight: bold;
    }
  }
}

.proof-images {
  margin-top: 20px;

  h4 {
    margin-bottom: 12px;
  }

  .proof-image {
    width: 120px;
    height: 120px;
    border-radius: 8px;
    margin-right: 12px;
  }
}

.action-bar {
  display: flex;
  gap: 12px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: sticky;
  bottom: 20px;
}

@media (max-width: 768px) {
  .task-meta {
    grid-template-columns: 1fr;
  }

  .action-bar {
    flex-direction: column;
  }
}
</style>
