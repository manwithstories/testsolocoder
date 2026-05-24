<template>
  <div class="order-detail-container">
    <div v-if="order" class="order-detail">
      <el-card class="order-card">
        <div class="order-header">
          <span :class="['status-tag', order.status]">
            {{ getStatusLabel(order.status) }}
          </span>
          <span class="order-no">订单号: {{ order.id }}</span>
        </div>

        <div class="order-steps">
          <el-steps :active="getStepActive()" finish-status="success">
            <el-step title="已下单" />
            <el-step title="已接单" />
            <el-step title="进行中" />
            <el-step title="已完成" />
          </el-steps>
        </div>

        <div class="task-info">
          <h3>{{ order.task?.title }}</h3>
          <div class="task-meta">
            <span :class="['task-type-tag', order.task?.type]">
              {{ getTaskTypeLabel(order.task?.type) }}
            </span>
            <span class="reward">¥{{ order.reward }}</span>
          </div>
          <p class="task-desc">{{ order.task?.description }}</p>
        </div>

        <div class="location-info">
          <div class="location-item">
            <el-icon><Location /></el-icon>
            <div>
              <span class="label">起点</span>
              <span class="address">{{ order.task?.start_addr }}</span>
            </div>
          </div>
          <div class="location-arrow">
            <el-icon><ArrowDown /></el-icon>
          </div>
          <div class="location-item">
            <el-icon><Location /></el-icon>
            <div>
              <span class="label">终点</span>
              <span class="address">{{ order.task?.end_addr }}</span>
            </div>
          </div>
        </div>

        <div class="time-info">
          <div class="time-item">
            <span>开始时间:</span>
            <span>{{ formatTime(order.start_time) }}</span>
          </div>
          <div class="time-item">
            <span>完成时间:</span>
            <span>{{ formatTime(order.end_time) }}</span>
          </div>
          <div class="time-item">
            <span>服务费:</span>
            <span>¥{{ order.service_fee }}</span>
          </div>
          <div class="time-item">
            <span>实际收入:</span>
            <span class="reward">¥{{ order.actual_payment }}</span>
          </div>
        </div>
      </el-card>

      <el-card class="tracking-card">
        <template #header>
          <div class="card-header">
            <el-icon><LocationInformation /></el-icon>
            <span>实时追踪</span>
          </div>
        </template>

        <div v-if="order.tracks && order.tracks.length > 0" class="track-list">
          <div
            v-for="track in order.tracks"
            :key="track.id"
            class="track-item"
          >
            <div class="track-time">{{ formatTime(track.created_at) }}</div>
            <div class="track-content">
              <span class="track-event">{{ getEventLabel(track.event_type) }}</span>
              <span v-if="track.message" class="track-message">{{ track.message }}</span>
              <span v-if="track.address" class="track-address">{{ track.address }}</span>
            </div>
          </div>
        </div>
        <el-empty v-else description="暂无追踪记录" />
      </el-card>

      <el-card v-if="order.proof_images && order.proof_images.length > 0" class="proof-card">
        <template #header>
          <div class="card-header">
            <el-icon><Picture /></el-icon>
            <span>送达凭证</span>
          </div>
        </template>
        <div class="proof-images">
          <el-image
            v-for="(img, index) in order.proof_images"
            :key="index"
            :src="img.image_url"
            fit="cover"
            class="proof-image"
            :preview-src-list="order.proof_images.map(i => i.image_url)"
          />
        </div>
      </el-card>

      <div class="action-bar">
        <template v-if="order.status === 'accepted' && isCourier">
          <el-button type="primary" size="large" @click="handleStart">
            开始任务
          </el-button>
        </template>
        <template v-else-if="order.status === 'in_progress' && isCourier">
          <el-button type="primary" size="large" @click="showLocationDialog = true">
            上报位置
          </el-button>
          <el-upload
            :auto-upload="false"
            :show-file-list="false"
            @change="handleUploadProof"
          >
            <el-button type="success" size="large">
              完成任务
            </el-button>
          </el-upload>
        </template>
        <template v-else-if="order.status === 'completed' && isPublisher">
          <el-button type="primary" size="large" @click="goToReview">
            评价跑腿员
          </el-button>
        </template>
      </div>
    </div>

    <el-dialog v-model="showLocationDialog" title="上报位置" width="500px">
      <el-form label-width="80px">
        <el-form-item label="经度">
          <el-input v-model="locationForm.lng" />
        </el-form-item>
        <el-form-item label="纬度">
          <el-input v-model="locationForm.lat" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="locationForm.address" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="locationForm.message" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showLocationDialog = false">取消</el-button>
        <el-button type="primary" @click="submitLocation">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type UploadFile } from 'element-plus'
import { Location, ArrowDown, LocationInformation, Picture } from '@element-plus/icons-vue'
import { orderApi } from '@/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import type { Order, OrderStatus, TaskType } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const order = ref<Order | null>(null)
const showLocationDialog = ref(false)
const proofImages = ref<string[]>([])

const locationForm = reactive({
  lat: 0,
  lng: 0,
  address: '',
  message: ''
})

const orderId = computed(() => Number(route.params.id))
const isPublisher = computed(() => order.value?.publisher_id === userStore.userInfo?.id)
const isCourier = computed(() => order.value?.courier_id === userStore.userInfo?.id)

const statusLabels: Record<OrderStatus, string> = {
  pending: '待处理',
  accepted: '已接单',
  in_progress: '进行中',
  delivered: '已送达',
  completed: '已完成',
  cancelled: '已取消'
}

const taskTypeLabels: Record<TaskType, string> = {
  buy: '代购',
  pickup: '代取',
  deliver: '代送',
  queue: '排队代办',
  errand: '其他代办'
}

const getStatusLabel = (status: OrderStatus) => statusLabels[status] || status
const getTaskTypeLabel = (type?: TaskType) => type ? taskTypeLabels[type] : '代办'

const formatTime = (time?: string) => {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

const getStepActive = () => {
  const status = order.value?.status
  switch (status) {
    case 'pending': return 0
    case 'accepted': return 1
    case 'in_progress': return 2
    case 'completed': return 3
    default: return 0
  }
}

const getEventLabel = (eventType: string) => {
  const labels: Record<string, string> = {
    start: '开始任务',
    location: '位置更新',
    complete: '完成任务',
    pickup: '已取件',
    delivery: '已送达'
  }
  return labels[eventType] || eventType
}

const fetchOrder = async () => {
  try {
    const res = await orderApi.get(orderId.value)
    if (res.code === 200) {
      order.value = res.data as Order
    }
  } catch (error) {
    console.error('Failed to fetch order:', error)
  }
}

const handleStart = async () => {
  try {
    const res = await orderApi.track(orderId.value, {
      latitude: 0,
      longitude: 0,
      event_type: 'start',
      message: '开始执行任务'
    })
    if (res.code === 200) {
      ElMessage.success('任务已开始')
      fetchOrder()
    }
  } catch (error) {
    console.error('Failed to start task:', error)
  }
}

const submitLocation = async () => {
  try {
    const res = await orderApi.track(orderId.value, {
      latitude: locationForm.lat,
      longitude: locationForm.lng,
      address: locationForm.address,
      message: locationForm.message,
      event_type: 'location'
    })
    if (res.code === 200) {
      ElMessage.success('位置已上报')
      showLocationDialog.value = false
      fetchOrder()
    }
  } catch (error) {
    console.error('Failed to submit location:', error)
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
    const res = await orderApi.track(orderId.value, {
      latitude: 0,
      longitude: 0,
      event_type: 'complete',
      message: '任务已完成'
    })
    if (res.code === 200) {
      ElMessage.success('任务已完成')
      fetchOrder()
    }
  } catch (error) {
    console.error('Failed to complete task:', error)
  }
}

const goToReview = () => {
  router.push(`/reviews?order_id=${order.value?.id}`)
}

onMounted(() => {
  fetchOrder()
})
</script>

<style lang="scss" scoped>
.order-detail-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.order-card,
.tracking-card,
.proof-card {
  margin-bottom: 20px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  .order-no {
    color: #909399;
    font-size: 14px;
  }
}

.order-steps {
  margin-bottom: 30px;
}

.task-info {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;

  h3 {
    font-size: 18px;
    margin-bottom: 10px;
  }

  .task-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;

    .reward {
      color: #ff5722;
      font-size: 20px;
      font-weight: bold;
    }
  }

  .task-desc {
    color: #606266;
    line-height: 1.6;
  }
}

.location-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;

  .location-item {
    flex: 1;
    display: flex;
    align-items: flex-start;
    gap: 8px;
    padding: 12px;
    background: #f5f7fa;
    border-radius: 8px;

    .label {
      color: #909399;
      font-size: 12px;
      display: block;
    }

    .address {
      color: #303133;
      font-size: 14px;
    }
  }

  .location-arrow {
    color: #667eea;
    font-size: 20px;
  }
}

.time-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;

  .time-item {
    display: flex;
    justify-content: space-between;

    .reward {
      color: #ff5722;
      font-weight: bold;
    }
  }
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.track-list {
  .track-item {
    display: flex;
    gap: 16px;
    padding: 12px 0;
    border-bottom: 1px solid #ebeef5;

    &:last-child {
      border-bottom: none;
    }

    .track-time {
      color: #909399;
      font-size: 12px;
      min-width: 120px;
    }

    .track-content {
      flex: 1;

      .track-event {
        color: #667eea;
        font-weight: 500;
        margin-right: 8px;
      }

      .track-message,
      .track-address {
        color: #606266;
      }
    }
  }
}

.proof-images {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;

  .proof-image {
    width: 120px;
    height: 120px;
    border-radius: 8px;
  }
}

.action-bar {
  display: flex;
  gap: 12px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

@media (max-width: 768px) {
  .location-info {
    flex-direction: column;
  }

  .time-info {
    grid-template-columns: 1fr;
  }
}
</style>
