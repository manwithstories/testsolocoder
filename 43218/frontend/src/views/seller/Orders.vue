<template>
  <div class="seller-orders-page">
    <div class="page-header">
      <h2 class="page-title">订单管理</h2>
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部订单" name="all" />
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
          <span class="buyer-info">买家：{{ order.buyer?.nickname || order.buyer?.username }}</span>
          <span class="order-time">{{ formatTime(order.createdAt) }}</span>
          <el-tag :type="getStatusType(order.status)">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>
        <div class="order-content">
          <div class="product-info">
            <img :src="getFirstImage(order.productImage)" class="product-image" />
            <div class="product-detail">
              <h4 class="product-title text-ellipsis">{{ order.productTitle }}</h4>
              <p class="product-price">¥{{ order.finalPrice.toFixed(2) }}</p>
              <p v-if="order.negotiated" class="negotiated-tag">
                <el-tag size="small" type="warning">议价订单</el-tag>
              </p>
            </div>
          </div>
          <div class="shipping-info" v-if="order.status >= 2">
            <p>收货人：{{ order.receiverName }}</p>
            <p>联系电话：{{ order.receiverPhone }}</p>
            <p>收货地址：{{ order.receiverAddress }}</p>
          </div>
          <div class="order-actions">
            <el-button
              v-if="order.status === 2"
              type="primary"
              @click="showShipDialog(order)"
            >
              发货
            </el-button>
            <el-button
              v-if="order.status === 1"
              @click="handleCancel(order)"
            >
              取消订单
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

    <el-dialog v-model="showShip" title="发货" width="500px">
      <el-form :model="shipForm" label-width="100px">
        <el-form-item label="物流公司" prop="trackingCompany">
          <el-select v-model="shipForm.trackingCompany" placeholder="请选择物流公司">
            <el-option label="顺丰速运" value="顺丰速运" />
            <el-option label="京东物流" value="京东物流" />
            <el-option label="中通快递" value="中通快递" />
            <el-option label="圆通速递" value="圆通速递" />
            <el-option label="申通快递" value="申通快递" />
            <el-option label="韵达快递" value="韵达快递" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号" prop="trackingNo">
          <el-input v-model="shipForm.trackingNo" placeholder="请输入物流单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showShip = false">取消</el-button>
        <el-button type="primary" :loading="shipping" @click="submitShip">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orderApi } from '@/api'
import { OrderStatus, OrderStatusText } from '@/types'
import type { Order } from '@/types'

const loading = ref(false)
const orders = ref<Order[]>([])
const total = ref(0)
const activeTab = ref('all')
const showShip = ref(false)
const currentOrder = ref<Order | null>(null)
const shipping = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10
})

const shipForm = reactive({
  trackingCompany: '',
  trackingNo: ''
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
    const res = await orderApi.getSellerOrders({
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

function showShipDialog(order: Order) {
  currentOrder.value = order
  shipForm.trackingCompany = ''
  shipForm.trackingNo = ''
  showShip.value = true
}

async function submitShip() {
  if (!currentOrder.value || !shipForm.trackingCompany || !shipForm.trackingNo) {
    ElMessage.warning('请填写完整的物流信息')
    return
  }

  shipping.value = true
  try {
    await orderApi.ship({
      orderNo: currentOrder.value.orderNo,
      trackingNo: shipForm.trackingNo,
      trackingCompany: shipForm.trackingCompany
    })
    ElMessage.success('发货成功')
    showShip.value = false
    fetchOrders()
  } catch (error: any) {
    ElMessage.error(error.message || '发货失败')
  } finally {
    shipping.value = false
  }
}

async function handleCancel(order: Order) {
  try {
    await ElMessageBox.confirm('确定取消该订单？', '提示', {
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

onMounted(() => {
  fetchOrders()
})
</script>

<style lang="scss" scoped>
.seller-orders-page {
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

        .buyer-info {
          color: var(--text-lighter-color);
        }

        .order-time {
          color: var(--text-lighter-color);
          font-size: 13px;
          margin-left: auto;
        }
      }

      .order-content {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;

        .product-info {
          display: flex;
          gap: 16px;

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
              margin-bottom: 8px;
            }
          }
        }

        .shipping-info {
          flex: 1;
          margin-left: 40px;

          p {
            margin-bottom: 4px;
            color: var(--text-light-color);
            font-size: 13px;
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
}
</style>
