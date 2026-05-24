<template>
  <div class="page-container" v-loading="loading">
    <div v-if="equipment" class="equipment-detail">
      <el-row :gutter="30">
        <el-col :span="12">
          <el-carousel
            v-if="equipment.images && equipment.images.length > 0"
            height="400px"
            indicator-position="outside"
          >
            <el-carousel-item v-for="img in equipment.images" :key="img.id">
              <img
                :src="`/uploads/${img.imageUrl}`"
                :alt="equipment.name"
                class="detail-image"
              />
            </el-carousel-item>
          </el-carousel>
          <div v-else class="detail-image-placeholder">
            <el-icon :size="96"><Camera /></el-icon>
          </div>
        </el-col>

        <el-col :span="12">
          <div class="equipment-info">
            <h1 class="equipment-title">{{ equipment.name }}</h1>
            <div class="equipment-brand">
              {{ equipment.brand }} / {{ equipment.model }}
            </div>

            <div class="equipment-rating-section">
              <el-rate :model-value="equipment.rating" disabled />
              <span class="rating-text">
                {{ equipment.rating.toFixed(1) }} 分 ({{ equipment.reviewCount }} 条评价)
              </span>
            </div>

            <div class="equipment-price-section">
              <span class="price-label">日租金</span>
              <span class="price-value">¥{{ equipment.dailyRent }}</span>
            </div>

            <div class="equipment-deposit-section">
              <span class="deposit-label">押金</span>
              <span class="deposit-value">¥{{ equipment.deposit }}</span>
            </div>

            <div class="equipment-meta">
              <div class="meta-item">
                <span class="meta-label">分类：</span>
                <el-tag size="small">{{ equipment.category }}</el-tag>
              </div>
              <div class="meta-item">
                <span class="meta-label">状态：</span>
                <el-tag
                  :type="equipment.status === 'available' ? 'success' : 'info'"
                  size="small"
                >
                  {{ getStatusText(equipment.status) }}
                </el-tag>
              </div>
              <div v-if="equipment.purchaseDate" class="meta-item">
                <span class="meta-label">购买时间：</span>
                <span>{{ equipment.purchaseDate }}</span>
              </div>
            </div>

            <div v-if="equipment.description" class="equipment-description">
              <h3>设备描述</h3>
              <p>{{ equipment.description }}</p>
            </div>

            <div class="owner-info" v-if="equipment.owner">
              <h3>出租方信息</h3>
              <div class="owner-detail">
                <el-avatar :size="48" :src="equipment.owner.avatar">
                  {{ equipment.owner.username?.charAt(0)?.toUpperCase() }}
                </el-avatar>
                <div class="owner-text">
                  <div class="owner-name">{{ equipment.owner.realName || equipment.owner.username }}</div>
                  <div class="owner-role">{{ equipment.owner.role === 'owner' ? '认证出租方' : '管理员' }}</div>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>

      <div class="card calendar-section">
        <h3>选择租赁日期</h3>
        <p class="calendar-hint">红色标记的日期表示已被预订，不可选择</p>
        <el-date-picker
          v-model="rentalDates"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :disabled-date="disabledDate"
          value-format="YYYY-MM-DD"
          @change="handleDateChange"
        />
      </div>

      <div v-if="userStore.isLoggedIn && canRent" class="card rental-form">
        <h3>提交租赁申请</h3>
        <el-form :model="rentalForm" :rules="rentalRules" ref="rentalFormRef">
          <el-form-item label="配送方式" prop="deliveryMethod">
            <el-radio-group v-model="rentalForm.deliveryMethod">
              <el-radio value="pickup">自取</el-radio>
              <el-radio value="delivery">配送</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item
            v-if="rentalForm.deliveryMethod === 'delivery'"
            label="配送地址"
            prop="deliveryAddress"
          >
            <el-input
              v-model="rentalForm.deliveryAddress"
              placeholder="请输入配送地址"
              type="textarea"
              :rows="2"
            />
          </el-form-item>
          <el-form-item v-if="totalRent > 0">
            <div class="rental-summary">
              <div class="summary-item">
                <span>租赁天数：</span>
                <span>{{ rentalDays }} 天</span>
              </div>
              <div class="summary-item">
                <span>租金总额：</span>
                <span class="highlight">¥{{ totalRent }}</span>
              </div>
              <div class="summary-item">
                <span>押金：</span>
                <span>¥{{ equipment.deposit }}</span>
              </div>
              <div class="summary-item total">
                <span>合计支付：</span>
                <span class="highlight">¥{{ totalRent + equipment.deposit }}</span>
              </div>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="submitting"
              :disabled="!canSubmit"
              @click="handleSubmitRental"
            >
              提交租赁申请
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <div v-if="!userStore.isLoggedIn" class="card login-prompt">
        <p>请先登录后再提交租赁申请</p>
        <el-button type="primary" @click="goLogin">立即登录</el-button>
      </div>

      <div class="card reviews-section">
        <div class="reviews-header">
          <h3>用户评价</h3>
          <el-button type="primary" link @click="goToReviews">查看全部评价</el-button>
        </div>
        <div v-if="reviews.length > 0" class="reviews-list">
          <div v-for="review in reviews" :key="review.id" class="review-item">
            <div class="review-header">
              <el-avatar :size="36">
                {{ review.fromUserId }}
              </el-avatar>
              <div class="review-user">
                <el-rate :model-value="review.rating" disabled size="small" />
                <span class="review-date">{{ formatDate(review.createdAt) }}</span>
              </div>
            </div>
            <div class="review-content">{{ review.content }}</div>
          </div>
        </div>
        <div v-else class="empty-state">
          <el-icon><ChatDotRound /></el-icon>
          <p>暂无评价</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import dayjs from 'dayjs'
import { equipmentApi } from '@/api/equipment'
import { orderApi, reviewApi } from '@/api/order'
import { useUserStore } from '@/stores/user'
import type { Equipment, Review, CreateOrderRequest } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const rentalFormRef = ref<FormInstance>()

const loading = ref(false)
const submitting = ref(false)
const equipment = ref<Equipment | null>(null)
const reviews = ref<Review[]>([])
const reservedDates = ref<string[]>([])
const rentalDates = ref<string[]>([])

const rentalForm = reactive<{
  deliveryMethod: string
  deliveryAddress: string
}>({
  deliveryMethod: 'pickup',
  deliveryAddress: ''
})

const rentalRules: FormRules = {
  deliveryMethod: [
    { required: true, message: '请选择配送方式', trigger: 'change' }
  ],
  deliveryAddress: [
    { required: true, message: '请输入配送地址', trigger: 'blur' }
  ]
}

const rentalDays = computed(() => {
  if (rentalDates.value.length === 2) {
    const start = dayjs(rentalDates.value[0])
    const end = dayjs(rentalDates.value[1])
    return end.diff(start, 'day') + 1
  }
  return 0
})

const totalRent = computed(() => {
  if (equipment.value && rentalDays.value > 0) {
    return equipment.value.dailyRent * rentalDays.value
  }
  return 0
})

const canRent = computed(() => {
  return userStore.isLoggedIn &&
         equipment.value &&
         equipment.value.status === 'available' &&
         userStore.user?.id !== equipment.value.ownerId
})

const canSubmit = computed(() => {
  return rentalDates.value.length === 2 &&
         rentalForm.deliveryMethod &&
         (rentalForm.deliveryMethod === 'pickup' || rentalForm.deliveryAddress)
})

onMounted(async () => {
  const id = parseInt(route.params.id as string)
  await loadEquipment(id)
  await loadReservedDates(id)
  await loadReviews(id)
})

async function loadEquipment(id: number) {
  loading.value = true
  try {
    const response = await equipmentApi.getEquipment(id)
    equipment.value = response.data
  } catch (error) {
    console.error('Failed to load equipment:', error)
    ElMessage.error('加载设备信息失败')
  } finally {
    loading.value = false
  }
}

async function loadReservedDates(id: number) {
  try {
    const response = await equipmentApi.getReservedDates(id)
    reservedDates.value = response.data.map(date => {
      const parts = date.split(':')
      return parts[parts.length - 1]
    })
  } catch (error) {
    console.error('Failed to load reserved dates:', error)
  }
}

async function loadReviews(id: number) {
  try {
    const response = await reviewApi.getEquipmentReviews(id, 1, 5)
    reviews.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load reviews:', error)
  }
}

function disabledDate(date: Date) {
  const dateStr = dayjs(date).format('YYYY-MM-DD')
  return reservedDates.value.includes(dateStr)
}

function handleDateChange() {
  // Date change handler
}

async function handleSubmitRental() {
  if (!rentalFormRef.value || !equipment.value) return

  const equipmentId = equipment.value.id

  await rentalFormRef.value.validate(async (valid) => {
    if (valid && canSubmit.value) {
      submitting.value = true
      try {
        const orderData: CreateOrderRequest = {
          equipmentId,
          startDate: rentalDates.value[0],
          endDate: rentalDates.value[1],
          deliveryMethod: rentalForm.deliveryMethod,
          deliveryAddress: rentalForm.deliveryAddress
        }

        const response = await orderApi.createOrder(orderData)
        ElMessage.success('租赁申请提交成功')
        router.push(`/orders/${response.data.id}`)
      } catch (error) {
        console.error('Failed to create order:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

function getStatusText(status: string) {
  const statusMap: Record<string, string> = {
    available: '可出租',
    rented: '已出租',
    maintenance: '维护中'
  }
  return statusMap[status] || status
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function goToReviews() {
  router.push(`/reviews/${equipment.value?.id}`)
}

function goLogin() {
  router.push('/login')
}
</script>

<style scoped>
.equipment-detail {
  margin-top: 20px;
}

.detail-image {
  width: 100%;
  height: 400px;
  object-fit: cover;
  border-radius: 8px;
}

.detail-image-placeholder {
  width: 100%;
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  color: #c0c4cc;
  border-radius: 8px;
}

.equipment-title {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.equipment-brand {
  font-size: 16px;
  color: #909399;
  margin-bottom: 16px;
}

.equipment-rating-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.rating-text {
  font-size: 14px;
  color: #606266;
}

.equipment-price-section {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 12px;
}

.price-label {
  font-size: 14px;
  color: #909399;
}

.price-value {
  font-size: 32px;
  font-weight: 600;
  color: #f56c6c;
}

.equipment-deposit-section {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 20px;
}

.deposit-label {
  font-size: 14px;
  color: #909399;
}

.deposit-value {
  font-size: 18px;
  color: #606266;
}

.equipment-meta {
  margin-bottom: 20px;
}

.meta-item {
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
}

.meta-label {
  color: #909399;
}

.equipment-description {
  margin-bottom: 20px;
}

.equipment-description h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #303133;
}

.equipment-description p {
  color: #606266;
  line-height: 1.6;
}

.owner-info h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #303133;
}

.owner-detail {
  display: flex;
  align-items: center;
  gap: 12px;
}

.owner-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.owner-role {
  font-size: 12px;
  color: #909399;
}

.calendar-section {
  margin-top: 30px;
}

.calendar-section h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #303133;
}

.calendar-hint {
  font-size: 12px;
  color: #f56c6c;
  margin-bottom: 16px;
}

.rental-form {
  margin-top: 30px;
}

.rental-form h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
  color: #303133;
}

.rental-summary {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
  width: 100%;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 14px;
  color: #606266;
}

.summary-item.total {
  border-top: 1px solid #e4e7ed;
  padding-top: 12px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.highlight {
  color: #f56c6c;
  font-weight: 600;
}

.login-prompt {
  margin-top: 30px;
  text-align: center;
}

.login-prompt p {
  margin-bottom: 16px;
  color: #606266;
}

.reviews-section {
  margin-top: 30px;
}

.reviews-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.reviews-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.review-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.review-user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.review-date {
  font-size: 12px;
  color: #909399;
}

.review-content {
  color: #606266;
  line-height: 1.6;
  padding-left: 48px;
}
</style>
