<template>
  <div class="page-container">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h1 class="page-title">订单详情</h1>
    </div>

    <div v-loading="loading" v-if="order" class="order-detail">
      <div class="card order-info">
        <div class="order-header">
          <div class="order-id">订单编号：{{ order.id }}</div>
          <el-tag :class="getStatusClass(order.status)" size="large">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>

        <el-descriptions :column="2" border>
          <el-descriptions-item label="设备名称">
            {{ order.equipment?.name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="品牌型号">
            {{ order.equipment?.brand }} / {{ order.equipment?.model }}
          </el-descriptions-item>
          <el-descriptions-item label="租期">
            {{ order.startDate }} 至 {{ order.endDate }}
          </el-descriptions-item>
          <el-descriptions-item label="配送方式">
            {{ order.deliveryMethod === 'pickup' ? '自取' : '配送' }}
          </el-descriptions-item>
          <el-descriptions-item v-if="order.deliveryAddress" label="配送地址" :span="2">
            {{ order.deliveryAddress }}
          </el-descriptions-item>
          <el-descriptions-item label="租金总额">
            <span class="highlight">¥{{ order.totalRent }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="押金">
            ¥{{ order.deposit }}
          </el-descriptions-item>
          <el-descriptions-item v-if="order.rejectReason" label="拒绝原因" :span="2">
            {{ order.rejectReason }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="card user-info">
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="user-card">
              <div class="user-card-title">出租方</div>
              <div class="user-card-content" v-if="order.owner">
                <el-avatar :size="48" :src="order.owner.avatar">
                  {{ order.owner.username?.charAt(0)?.toUpperCase() }}
                </el-avatar>
                <div class="user-info-text">
                  <div class="user-name">{{ order.owner.realName || order.owner.username }}</div>
                  <div class="user-phone">{{ order.owner.phone || '-' }}</div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="user-card">
              <div class="user-card-title">租借方</div>
              <div class="user-card-content" v-if="order.renter">
                <el-avatar :size="48" :src="order.renter.avatar">
                  {{ order.renter.username?.charAt(0)?.toUpperCase() }}
                </el-avatar>
                <div class="user-info-text">
                  <div class="user-name">{{ order.renter.realName || order.renter.username }}</div>
                  <div class="user-phone">{{ order.renter.phone || '-' }}</div>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <div class="card actions" v-if="canShowActions">
        <el-button
          v-if="canConfirm"
          type="success"
          :loading="actionLoading"
          @click="handleConfirm"
        >确认订单</el-button>
        <el-button
          v-if="canReject"
          type="danger"
          :loading="actionLoading"
          @click="handleReject"
        >拒绝订单</el-button>
        <el-button
          v-if="canStart"
          type="primary"
          :loading="actionLoading"
          @click="handleStart"
        >开始租赁</el-button>
        <el-button
          v-if="canComplete"
          type="success"
          :loading="actionLoading"
          @click="handleComplete"
        >完成订单</el-button>
        <el-button
          v-if="canCancel"
          type="danger"
          :loading="actionLoading"
          @click="handleCancel"
        >取消订单</el-button>
      </div>

      <div class="card review-section" v-if="order.status === 'completed'">
        <h3>评价</h3>
        <div v-if="!hasReviewed">
          <el-form :model="reviewForm" ref="reviewFormRef" :rules="reviewRules">
            <el-form-item label="评分" prop="rating">
              <el-rate v-model="reviewForm.rating" />
            </el-form-item>
            <el-form-item label="评价内容" prop="content">
              <el-input
                v-model="reviewForm.content"
                type="textarea"
                :rows="4"
                placeholder="请输入您的评价"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="submittingReview" @click="handleSubmitReview">
                提交评价
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        <div v-else class="already-reviewed">
          <p>您已完成对该订单的评价</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { orderApi, reviewApi } from '@/api/order'
import { useUserStore } from '@/stores/user'
import type { Order } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const reviewFormRef = ref<FormInstance>()

const loading = ref(false)
const actionLoading = ref(false)
const submittingReview = ref(false)
const order = ref<Order | null>(null)
const hasReviewed = ref(false)

const reviewForm = reactive({
  rating: 5,
  content: ''
})

const reviewRules: FormRules = {
  rating: [
    { required: true, message: '请选择评分', trigger: 'change' }
  ]
}

const canConfirm = computed(() => {
  return order.value?.status === 'pending' &&
         userStore.isOwner() &&
         order.value.ownerId === userStore.user?.id
})

const canReject = computed(() => {
  return order.value?.status === 'pending' &&
         userStore.isOwner() &&
         order.value.ownerId === userStore.user?.id
})

const canStart = computed(() => {
  return order.value?.status === 'confirmed' &&
         order.value.renterId === userStore.user?.id
})

const canComplete = computed(() => {
  return order.value?.status === 'rented' &&
         userStore.isOwner() &&
         order.value.ownerId === userStore.user?.id
})

const canCancel = computed(() => {
  return (order.value?.status === 'pending' || order.value?.status === 'confirmed') &&
         (order.value?.renterId === userStore.user?.id ||
          order.value?.ownerId === userStore.user?.id)
})

const canShowActions = computed(() => {
  return canConfirm.value || canReject.value || canStart.value || canComplete.value || canCancel.value
})

onMounted(async () => {
  const id = parseInt(route.params.id as string)
  await loadOrder(id)
  await checkReviewStatus(id)
})

async function loadOrder(id: number) {
  loading.value = true
  try {
    const response = await orderApi.getOrder(id)
    order.value = response.data
  } catch (error) {
    console.error('Failed to load order:', error)
    ElMessage.error('加载订单详情失败')
  } finally {
    loading.value = false
  }
}

async function checkReviewStatus(orderId: number) {
  try {
    const response = await reviewApi.getMyReviews()
    hasReviewed.value = response.data.some(r => r.orderId === orderId)
  } catch (error) {
    console.error('Failed to check review status:', error)
  }
}

async function handleConfirm() {
  if (!order.value) return

  try {
    await ElMessageBox.confirm('确认要接受此租赁申请吗？', '确认订单', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    actionLoading.value = true
    await orderApi.confirmOrder(order.value.id)
    ElMessage.success('订单已确认')
    await loadOrder(order.value.id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to confirm order:', error)
    }
  } finally {
    actionLoading.value = false
  }
}

async function handleReject() {
  if (!order.value) return

  try {
    const { value: reason } = await ElMessageBox.prompt(
      '请输入拒绝原因',
      '拒绝订单',
      {
        confirmButtonText: '确认拒绝',
        cancelButtonText: '取消',
        inputPlaceholder: '请输入拒绝原因'
      }
    )

    actionLoading.value = true
    await orderApi.rejectOrder(order.value.id, reason)
    ElMessage.success('订单已拒绝')
    await loadOrder(order.value.id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to reject order:', error)
    }
  } finally {
    actionLoading.value = false
  }
}

async function handleStart() {
  if (!order.value) return

  try {
    await ElMessageBox.confirm('确认要开始租赁吗？', '开始租赁', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    actionLoading.value = true
    await orderApi.startRental(order.value.id)
    ElMessage.success('租赁已开始')
    await loadOrder(order.value.id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to start rental:', error)
    }
  } finally {
    actionLoading.value = false
  }
}

async function handleComplete() {
  if (!order.value) return

  try {
    await ElMessageBox.confirm(
      '确认要完成此订单吗？设备已归还且无损坏？',
      '完成订单',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    actionLoading.value = true
    await orderApi.completeOrder(order.value.id)
    ElMessage.success('订单已完成')
    await loadOrder(order.value.id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to complete order:', error)
    }
  } finally {
    actionLoading.value = false
  }
}

async function handleCancel() {
  if (!order.value) return

  try {
    await ElMessageBox.confirm('确认要取消此订单吗？', '取消订单', {
      confirmButtonText: '确认取消',
      cancelButtonText: '继续操作',
      type: 'warning'
    })

    actionLoading.value = true
    await orderApi.cancelOrder(order.value.id)
    ElMessage.success('订单已取消')
    await loadOrder(order.value.id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to cancel order:', error)
    }
  } finally {
    actionLoading.value = false
  }
}

async function handleSubmitReview() {
  if (!reviewFormRef.value || !order.value) return

  const orderData = order.value

  await reviewFormRef.value.validate(async (valid) => {
    if (valid) {
      submittingReview.value = true
      try {
        const targetUserId = orderData.renterId === userStore.user?.id
          ? orderData.ownerId
          : orderData.renterId

        await reviewApi.createReview({
          orderId: orderData.id,
          toUserId: targetUserId,
          equipmentId: orderData.equipmentId,
          rating: reviewForm.rating,
          content: reviewForm.content
        })

        ElMessage.success('评价提交成功')
        hasReviewed.value = true
      } catch (error) {
        console.error('Failed to submit review:', error)
      } finally {
        submittingReview.value = false
      }
    }
  })
}

function goBack() {
  router.back()
}

function getStatusText(status: string) {
  const textMap: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    rented: '租赁中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return textMap[status] || status
}

function getStatusClass(status: string) {
  return `status-${status}`
}
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.order-id {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.highlight {
  color: #f56c6c;
  font-weight: 600;
}

.user-card {
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
}

.user-card-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 12px;
}

.user-card-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.user-phone {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.actions {
  display: flex;
  gap: 12px;
}

.status-pending {
  color: #e6a23c;
}

.status-confirmed {
  color: #409eff;
}

.status-rented {
  color: #67c23a;
}

.status-completed {
  color: #909399;
}

.status-cancelled {
  color: #f56c6c;
}

.already-reviewed {
  text-align: center;
  color: #909399;
  padding: 20px;
}

.review-section h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #303133;
}
</style>
