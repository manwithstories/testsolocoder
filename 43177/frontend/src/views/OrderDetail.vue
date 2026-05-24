<template>
  <div class="order-detail-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">工单详情</h2>
        <el-button @click="router.back()">返回</el-button>
      </div>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="工单号">{{ order?.order_no }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :class="`status-${order?.status}`">
            {{ getStatusText(order?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="服务项目">{{ order?.service_item?.name }}</el-descriptions-item>
        <el-descriptions-item label="服务分类">{{ order?.service_item?.category?.name }}</el-descriptions-item>
        <el-descriptions-item label="问题标题">{{ order?.title }}</el-descriptions-item>
        <el-descriptions-item label="紧急程度">
          {{ getUrgentLevelText(order?.urgent_level) }}
        </el-descriptions-item>
        <el-descriptions-item label="详细描述" :span="2">{{ order?.description }}</el-descriptions-item>
        <el-descriptions-item label="联系人">{{ order?.contact_name }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ order?.contact_phone }}</el-descriptions-item>
        <el-descriptions-item label="服务地址" :span="2">{{ order?.address }}</el-descriptions-item>
        <el-descriptions-item label="预约时间">{{ formatTime(order?.appointment_time) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(order?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="报价">{{ order?.quoted_price ? '¥' + order.quoted_price : '-' }}</el-descriptions-item>
        <el-descriptions-item label="最终费用">{{ order?.final_price ? '¥' + order.final_price : '-' }}</el-descriptions-item>
      </el-descriptions>

      <el-card v-if="order?.technician" class="mt-20">
        <template #header>
          <div class="card-header">技师信息</div>
        </template>
        <div class="tech-info">
          <el-avatar :size="50" :src="order.technician.avatar">
            {{ order.technician.username?.charAt(0) }}
          </el-avatar>
          <div class="tech-detail">
            <div class="tech-name">{{ order.technician.real_name || order.technician.username }}</div>
            <div class="tech-phone">{{ order.technician.phone }}</div>
          </div>
        </div>
      </el-card>

      <el-card v-if="order && orderLogs.length > 0" class="mt-20">
        <template #header>
          <div class="card-header">工单日志</div>
        </template>
        <el-timeline>
          <el-timeline-item
            v-for="log in orderLogs"
            :key="log.id"
            :timestamp="formatTime(log.created_at)"
            placement="top"
          >
            <div class="log-item">
              <div class="log-action">{{ getActionText(log.action) }}</div>
              <div class="log-content">{{ log.content }}</div>
            </div>
          </el-timeline-item>
        </el-timeline>
      </el-card>

      <el-card v-if="order && reviews.length > 0" class="mt-20">
        <template #header>
          <div class="card-header">评价信息</div>
        </template>
        <div v-for="review in reviews" :key="review.id" class="review-item">
          <div class="review-header">
            <el-rate :model-value="review.rating" disabled size="small" />
            <span class="review-time">{{ formatTime(review.created_at) }}</span>
          </div>
          <div class="review-content">{{ review.content }}</div>
          <div v-if="review.reply" class="review-reply">
            <strong>技师回复：</strong>{{ review.reply }}
          </div>
          <div v-if="review.is_intervened" class="review-intervene">
            <strong>平台介入：</strong>{{ review.intervene_note }}
          </div>
        </div>
      </el-card>

      <div v-if="order" class="action-bar">
        <el-button
          v-if="order.status === 'assigned' && userStore.isTechnician"
          type="primary"
          @click="handleAccept"
        >
          接单
        </el-button>
        <el-button
          v-if="order.status === 'accepted' && userStore.isTechnician"
          type="success"
          @click="handleArrive"
        >
          已到达
        </el-button>
        <el-button
          v-if="order.status === 'on_site' && userStore.isTechnician"
          type="warning"
          @click="handleStart"
        >
          开始维修
        </el-button>
        <el-button
          v-if="order.status === 'repairing' && userStore.isTechnician"
          type="success"
          @click="handleComplete"
        >
          完工
        </el-button>
        <el-button
          v-if="(order.status === 'pending' || order.status === 'assigned') && userStore.isCustomer"
          type="danger"
          @click="handleCancel"
        >
          取消工单
        </el-button>
        <el-button
          v-if="order.status === 'completed' && userStore.isCustomer && !hasReview"
          type="primary"
          @click="showReviewDialog = true"
        >
          评价
        </el-button>
        <el-button
          v-if="order.status === 'completed' && userStore.isCustomer"
          type="warning"
          @click="showRefundDialog = true"
        >
          申请退款
        </el-button>
      </div>
    </div>

    <el-dialog v-model="showReviewDialog" title="评价工单" width="500px">
      <el-form :model="reviewForm" label-width="80px">
        <el-form-item label="评分">
          <el-rate v-model="reviewForm.rating" />
        </el-form-item>
        <el-form-item label="评价内容">
          <el-input v-model="reviewForm.content" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReviewDialog = false">取消</el-button>
        <el-button type="primary" @click="submitReview" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showRefundDialog" title="申请退款" width="500px">
      <el-form :model="refundForm" label-width="80px">
        <el-form-item label="退款原因">
          <el-input v-model="refundForm.reason" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="退款金额">
          <el-input-number v-model="refundForm.amount" :min="0" :max="order?.final_price || 0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRefundDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRefund" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'
import { orderApi } from '@/api/order'
import { userApi } from '@/api/user'
import type { Order, OrderLog, Review } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const order = ref<Order | null>(null)
const orderLogs = ref<OrderLog[]>([])
const reviews = ref<Review[]>([])
const showReviewDialog = ref(false)
const showRefundDialog = ref(false)
const submitting = ref(false)

const reviewForm = reactive({
  rating: 5,
  content: ''
})

const refundForm = reactive({
  reason: '',
  amount: 0
})

const hasReview = computed(() => reviews.value.length > 0)

onMounted(() => {
  loadOrderDetail()
})

async function loadOrderDetail() {
  const id = route.params.id
  try {
    const res = await orderApi.getOrderDetail(Number(id))
    if (res.data) {
      order.value = res.data.order
      orderLogs.value = res.data.logs || []
      reviews.value = res.data.reviews || []
    }
  } catch (error) {
    console.error('Failed to load order detail:', error)
  }
}

function getStatusText(status?: string): string {
  if (!status) return ''
  const statusMap: Record<string, string> = {
    pending: '待分配',
    assigned: '待接单',
    accepted: '已接单',
    on_site: '已到达',
    repairing: '维修中',
    completed: '已完成',
    cancelled: '已取消',
    refunding: '退款中',
    refunded: '已退款'
  }
  return statusMap[status] || status
}

function getUrgentLevelText(level?: number): string {
  const levels = ['普通', '加急', '特急']
  return levels[level || 0] || '普通'
}

function getActionText(action: string): string {
  const actionMap: Record<string, string> = {
    create: '创建工单',
    assign: '分配技师',
    accept: '技师接单',
    arrive: '到达现场',
    start_repair: '开始维修',
    complete: '维修完成',
    cancel: '取消工单',
    refund_request: '申请退款',
    refund_approved: '退款通过',
    refund_rejected: '退款拒绝'
  }
  return actionMap[action] || action
}

function formatTime(time?: string): string {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

async function handleAccept() {
  try {
    await ElMessageBox.confirm('确定要接这个工单吗？', '确认')
    await orderApi.acceptOrder(order.value!.id)
    ElMessage.success('接单成功')
    loadOrderDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to accept:', error)
    }
  }
}

async function handleArrive() {
  try {
    await orderApi.arriveAtSite(order.value!.id)
    ElMessage.success('已标记到达')
    loadOrderDetail()
  } catch (error) {
    console.error('Failed to arrive:', error)
  }
}

async function handleStart() {
  try {
    await orderApi.startRepair(order.value!.id)
    ElMessage.success('已开始维修')
    loadOrderDetail()
  } catch (error) {
    console.error('Failed to start:', error)
  }
}

async function handleComplete() {
  try {
    const { value: price } = await ElMessageBox.prompt('请输入最终费用', '完工确认', {
      inputValidator: (value) => {
        if (!value || isNaN(Number(value)) || Number(value) <= 0) return '请输入有效金额'
        return true
      }
    })
    await orderApi.completeOrder(order.value!.id, { final_price: Number(price) })
    ElMessage.success('工单已完工')
    loadOrderDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to complete:', error)
    }
  }
}

async function handleCancel() {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入取消原因', '取消工单', {
      inputValidator: (value) => {
        if (!value) return '请输入取消原因'
        return true
      }
    })
    await orderApi.cancelOrder(order.value!.id, { cancel_reason: reason })
    ElMessage.success('工单已取消')
    loadOrderDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to cancel:', error)
    }
  }
}

async function submitReview() {
  submitting.value = true
  try {
    await userApi.createReview({
      order_id: order.value!.id,
      rating: reviewForm.rating,
      content: reviewForm.content
    })
    ElMessage.success('评价成功')
    showReviewDialog.value = false
    loadOrderDetail()
  } catch (error) {
    console.error('Failed to submit review:', error)
  } finally {
    submitting.value = false
  }
}

async function submitRefund() {
  submitting.value = true
  try {
    await orderApi.requestRefund(order.value!.id, {
      refund_reason: refundForm.reason,
      refund_amount: refundForm.amount
    })
    ElMessage.success('退款申请已提交')
    showRefundDialog.value = false
    loadOrderDetail()
  } catch (error) {
    console.error('Failed to submit refund:', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.order-detail-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  font-weight: 600;
}

.tech-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.tech-detail {
  flex: 1;
}

.tech-name {
  font-weight: 600;
  color: #303133;
  margin-bottom: 5px;
}

.tech-phone {
  color: #909399;
  font-size: 14px;
}

.log-item {
  padding: 5px 0;
}

.log-action {
  font-weight: 600;
  color: #303133;
  margin-bottom: 5px;
}

.log-content {
  color: #606266;
  font-size: 14px;
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
  align-items: center;
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

.review-reply, .review-intervene {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin-top: 10px;
  color: #606266;
}

.action-bar {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-top: 30px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
}
</style>
