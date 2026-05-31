<template>
  <div class="orders-page">
    <div class="page-header">
      <h2 class="page-title">我的订单</h2>
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部" name="all" />
        <el-tab-pane label="待支付" name="pending" />
        <el-tab-pane label="已支付" name="paid" />
        <el-tab-pane label="已发货" name="shipped" />
        <el-tab-pane label="已完成" name="completed" />
      </el-tabs>
    </div>

    <div class="order-list" v-loading="loading">
      <div
        v-for="order in orders"
        :key="order.id"
        class="order-card card"
      >
        <div class="order-header">
          <span class="order-no">订单号：{{ order.orderNo }}</span>
          <span class="order-time">{{ formatTime(order.createdAt) }}</span>
          <el-tag :type="getStatusType(order.status)">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>
        <div class="order-content">
          <div class="product-info" @click="goToProduct(order.productId)">
            <img :src="getFirstImage(order.productImage)" class="product-image" />
            <div class="product-detail">
              <h4 class="product-title text-ellipsis">{{ order.productTitle }}</h4>
              <p class="product-price">¥{{ order.finalPrice.toFixed(2) }}</p>
            </div>
          </div>
          <div class="order-actions">
            <el-button
              v-if="order.status === 1"
              type="primary"
              @click="handlePay(order)"
            >
              去支付
            </el-button>
            <el-button
              v-if="order.status === 3"
              type="success"
              @click="handleConfirm(order)"
            >
              确认收货
            </el-button>
            <el-button
              v-if="order.status === 1 || order.status === 2"
              @click="handleCancel(order)"
            >
              取消订单
            </el-button>
            <el-button
              v-if="order.status === 3 || order.status === 4"
              @click="handleRefund(order)"
            >
              申请退款
            </el-button>
            <el-button
              v-if="order.status === 5 && !hasReviewed(order)"
              type="warning"
              @click="showReviewDialog(order)"
            >
              去评价
            </el-button>
          </div>
        </div>
      </div>

      <div class="empty-state" v-if="!loading && orders.length === 0">
        <el-empty description="暂无订单" />
      </div>
    </div>

    <div class="pagination-wrapper" v-if="total > 0">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchOrders"
        @current-change="fetchOrders"
      />
    </div>

    <el-dialog v-model="showPayDialog" title="支付订单" width="400px">
      <div class="pay-info">
        <p>订单金额：<span class="price-text">¥{{ currentOrder?.finalPrice.toFixed(2) }}</span></p>
      </div>
      <el-form label-width="80px">
        <el-form-item label="支付方式">
          <el-radio-group v-model="payMethod">
            <el-radio value="alipay">支付宝</el-radio>
            <el-radio value="wechat">微信支付</el-radio>
            <el-radio value="wallet">钱包支付</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPayDialog = false">取消</el-button>
        <el-button type="primary" :loading="paying" @click="submitPay">确认支付</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showReviewDialog" title="评价订单" width="500px">
      <el-form :model="reviewForm" label-width="80px">
        <el-form-item label="评分">
          <el-rate v-model="reviewForm.rating" />
        </el-form-item>
        <el-form-item label="评价内容">
          <el-input
            v-model="reviewForm.content"
            type="textarea"
            :rows="3"
            placeholder="请输入评价内容"
          />
        </el-form-item>
        <el-form-item label="商品质量">
          <el-rate v-model="reviewForm.qualityScore" />
        </el-form-item>
        <el-form-item label="服务态度">
          <el-rate v-model="reviewForm.serviceScore" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReviewDialog = false">取消</el-button>
        <el-button type="primary" :loading="submittingReview" @click="submitReview">
          提交评价
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orderApi, reviewApi } from '@/api'
import { OrderStatus, OrderStatusText } from '@/types'
import type { Order } from '@/types'

const router = useRouter()

const loading = ref(false)
const orders = ref<Order[]>([])
const total = ref(0)
const activeTab = ref('all')
const showPayDialog = ref(false)
const showReviewDialog = ref(false)
const currentOrder = ref<Order | null>(null)
const payMethod = ref('alipay')
const paying = ref(false)
const submittingReview = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10
})

const reviewForm = reactive({
  rating: 5,
  content: '',
  qualityScore: 5,
  serviceScore: 5
})

function getStatusText(status: number): string {
  return OrderStatusText[status] || '未知'
}

function getStatusType(status: number): string {
  const typeMap: Record<number, string> = {
    1: 'warning',
    2: 'primary',
    3: 'primary',
    4: 'success',
    5: 'success',
    6: 'info',
    7: 'warning',
    8: 'danger',
    9: 'warning'
  }
  return typeMap[status] || 'info'
}

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

function getFirstImage(image: string | undefined): string {
  if (!image) return 'https://picsum.photos/80/80'
  try {
    const arr = JSON.parse(image)
    return Array.isArray(arr) && arr.length > 0 ? arr[0] : 'https://picsum.photos/80/80'
  } catch {
    return 'https://picsum.photos/80/80'
  }
}

function goToProduct(id: number | undefined) {
  if (id) router.push(`/products/${id}`)
}

function hasReviewed(order: Order): boolean {
  return order.reviews && order.reviews.length > 0
}

function handleTabChange() {
  pagination.page = 1
  fetchOrders()
}

function getStatusFilter(): number | undefined {
  const statusMap: Record<string, number> = {
    pending: 1,
    paid: 2,
    shipped: 3,
    completed: 5
  }
  return activeTab.value === 'all' ? undefined : statusMap[activeTab.value]
}

async function fetchOrders() {
  loading.value = true
  try {
    const res = await orderApi.getList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: getStatusFilter()
    })
    orders.value = res.data
    total.value = res.pagination.total
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loading.value = false
  }
}

function handlePay(order: Order) {
  currentOrder.value = order
  showPayDialog.value = true
}

async function submitPay() {
  if (!currentOrder.value) return
  paying.value = true
  try {
    await orderApi.pay({
      orderNo: currentOrder.value.orderNo,
      paymentMethod: payMethod.value
    })
    ElMessage.success('支付成功')
    showPayDialog.value = false
    fetchOrders()
  } catch (error: any) {
    ElMessage.error(error.message || '支付失败')
  } finally {
    paying.value = false
  }
}

async function handleConfirm(order: Order) {
  try {
    await ElMessageBox.confirm('确认收到商品？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    await orderApi.confirm(order.orderNo)
    ElMessage.success('已确认收货')
    fetchOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to confirm:', error)
    }
  }
}

async function handleCancel(order: Order) {
  try {
    await ElMessageBox.confirm('确定取消订单？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    await orderApi.cancel(order.orderNo)
    ElMessage.success('订单已取消')
    fetchOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to cancel:', error)
    }
  }
}

async function handleRefund(order: Order) {
  try {
    const { value } = await ElMessageBox.prompt('请输入退款原因', '申请退款', {
      confirmButtonText: '提交',
      cancelButtonText: '取消',
      inputPlaceholder: '请输入退款原因'
    })
    await orderApi.refund({ orderNo: order.orderNo, reason: value })
    ElMessage.success('退款申请已提交')
    fetchOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to refund:', error)
    }
  }
}

function showReviewDialog(order: Order) {
  currentOrder.value = order
  reviewForm.rating = 5
  reviewForm.content = ''
  reviewForm.qualityScore = 5
  reviewForm.serviceScore = 5
  showReviewDialog.value = true
}

async function submitReview() {
  if (!currentOrder.value) return
  submittingReview.value = true
  try {
    await reviewApi.create({
      orderId: currentOrder.value.id,
      revieweeId: currentOrder.value.sellerId,
      reviewType: 'product',
      rating: reviewForm.rating,
      content: reviewForm.content,
      qualityScore: reviewForm.qualityScore,
      serviceScore: reviewForm.serviceScore
    })
    ElMessage.success('评价成功')
    showReviewDialog.value = false
    fetchOrders()
  } catch (error: any) {
    ElMessage.error(error.message || '评价失败')
  } finally {
    submittingReview.value = false
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style lang="scss" scoped>
.orders-page {
  .order-list {
    .order-card {
      margin-bottom: 16px;

      .order-header {
        display: flex;
        align-items: center;
        gap: 16px;
        padding-bottom: 12px;
        border-bottom: 1px solid #f0f0f0;
        margin-bottom: 12px;

        .order-no {
          font-weight: 500;
        }

        .order-time {
          color: var(--text-lighter-color);
          font-size: 13px;
        }
      }

      .order-content {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .product-info {
          display: flex;
          gap: 16px;
          cursor: pointer;

          .product-image {
            width: 80px;
            height: 80px;
            border-radius: 4px;
            object-fit: cover;
          }

          .product-detail {
            .product-title {
              font-size: 14px;
              max-width: 300px;
              margin-bottom: 8px;
            }

            .product-price {
              color: var(--danger-color);
              font-weight: 500;
            }
          }
        }

        .order-actions {
          display: flex;
          gap: 8px;
        }
      }
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 20px 0;
  }

  .pay-info {
    margin-bottom: 20px;
    text-align: center;
    font-size: 18px;
  }
}
</style>
