<template>
  <div class="page-container">
    <div class="card">
      <div class="page-header">
        <h2>我的订单</h2>
        <div class="filter-bar">
          <el-radio-group v-model="status" size="default" @change="loadOrders">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="pending">待支付</el-radio-button>
            <el-radio-button label="paid">已支付</el-radio-button>
            <el-radio-button label="shipped">已发货</el-radio-button>
            <el-radio-button label="completed">已完成</el-radio-button>
            <el-radio-button label="cancelled">已取消</el-radio-button>
          </el-radio-group>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <el-table
      v-else-if="orders.length > 0"
      :data="orders"
      style="width: 100%"
      stripe
    >
      <el-table-column label="订单号" prop="order_number" width="200" />
      <el-table-column label="商品" min-width="200">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.product?.images?.[0]" :src="row.product.images[0].image_url" />
            <div class="product-info">
              <span class="product-title text-ellipsis">{{ row.product?.title }}</span>
              <span class="product-price">¥{{ row.price.toFixed(2) }}</span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getOrderStatusType(row.status)">
            {{ getOrderStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="支付状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.payment_status === 'success' ? 'success' : 'warning'">
            {{ row.payment_status === 'success' ? '已支付' : '待支付' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="鉴定" width="80">
        <template #default="{ row }">
          <span v-if="row.need_auth">需要</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="下单时间" prop="created_at" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'pending'"
            type="primary"
            size="small"
            @click="handlePay(row)"
          >
            去支付
          </el-button>
          <el-button
            v-if="row.status === 'shipped' && userStore.userRole === 'buyer'"
            type="success"
            size="small"
            @click="handleConfirm(row)"
          >
            确认收货
          </el-button>
          <el-button
            v-if="row.status === 'paid' && userStore.userRole === 'seller'"
            type="primary"
            size="small"
            @click="handleShip(row)"
          >
            去发货
          </el-button>
          <el-button
            v-if="row.status === 'pending' || row.status === 'paid'"
            type="danger"
            size="small"
            @click="handleCancel(row)"
          >
            取消
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-else class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>暂无订单</p>
    </div>

    <div v-if="!loading && total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadOrders"
      />
    </div>

    <el-dialog v-model="payDialogVisible" title="支付订单" width="400px">
      <el-form :model="payForm">
        <el-form-item label="支付方式">
          <el-radio-group v-model="payForm.payment_method">
            <el-radio-button label="alipay">支付宝</el-radio-button>
            <el-radio-button label="wechat">微信</el-radio-button>
            <el-radio-button label="card">银行卡</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="payDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="paying" @click="confirmPay">确认支付</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="shipDialogVisible" title="发货" width="400px">
      <el-form :model="shipForm">
        <el-form-item label="快递单号">
          <el-input v-model="shipForm.tracking_number" placeholder="请输入快递单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="shipping" @click="confirmShip">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { orderApi } from '@/api/order'
import { authServiceApi } from '@/api/authentication'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Order } from '@/types'
import { ORDER_STATUS_OPTIONS } from '@/types'
import dayjs from 'dayjs'
import { Loading, Box } from '@element-plus/icons-vue'

const userStore = useUserStore()

const orders = ref<Order[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const status = ref('')

const payDialogVisible = ref(false)
const shipDialogVisible = ref(false)
const currentOrder = ref<Order | null>(null)
const paying = ref(false)
const shipping = ref(false)

const payForm = ref({ payment_method: 'alipay' })
const shipForm = ref({ tracking_number: '' })

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await orderApi.listOrders({
      page: page.value,
      page_size: pageSize.value,
      status: status.value || undefined
    })
    if (res.code === 200) {
      const data = res.data as any
      orders.value = data?.list || []
      total.value = data?.total || 0
    }
  } catch (error) {
    console.error('Load orders error:', error)
  } finally {
    loading.value = false
  }
}

const getOrderStatusType = (status: string) => {
  const opt = ORDER_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.type || 'info'
}

const getOrderStatusLabel = (status: string) => {
  const opt = ORDER_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.label || status
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const handlePay = (order: Order) => {
  currentOrder.value = order
  payDialogVisible.value = true
}

const confirmPay = async () => {
  if (!currentOrder.value) return

  paying.value = true
  try {
    const res = await orderApi.payOrder(currentOrder.value.id, payForm.value)
    if (res.code === 200) {
      ElMessage.success('支付成功')
      payDialogVisible.value = false
      loadOrders()

      if (currentOrder.value.need_auth) {
        await authServiceApi.createAuthentication({ order_id: currentOrder.value.id })
      }
    }
  } catch (error) {
    console.error('Pay error:', error)
  } finally {
    paying.value = false
  }
}

const handleShip = (order: Order) => {
  currentOrder.value = order
  shipDialogVisible.value = true
}

const confirmShip = async () => {
  if (!currentOrder.value) return

  shipping.value = true
  try {
    const res = await orderApi.shipOrder(currentOrder.value.id, shipForm.value)
    if (res.code === 200) {
      ElMessage.success('发货成功')
      shipDialogVisible.value = false
      shipForm.value.tracking_number = ''
      loadOrders()
    }
  } catch (error) {
    console.error('Ship error:', error)
  } finally {
    shipping.value = false
  }
}

const handleConfirm = async (order: Order) => {
  try {
    await ElMessageBox.confirm('确认收到商品吗？', '确认收货', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await orderApi.confirmDelivery(order.id)
    if (res.code === 200) {
      ElMessage.success('确认收货成功')
      loadOrders()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Confirm error:', error)
    }
  }
}

const handleCancel = async (order: Order) => {
  try {
    await ElMessageBox.confirm('确定取消此订单吗？', '取消订单', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await orderApi.cancelOrder(order.id)
    if (res.code === 200) {
      ElMessage.success('订单已取消')
      loadOrders()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Cancel error:', error)
    }
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  h2 {
    font-size: 20px;
    font-weight: 600;
  }
}

.filter-bar {
  display: flex;
  gap: 12px;
}

.product-cell {
  display: flex;
  gap: 12px;
  align-items: center;
  
  img {
    width: 60px;
    height: 60px;
    object-fit: cover;
    border-radius: 4px;
  }
  
  .product-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    
    .product-title {
      font-size: 14px;
      max-width: 200px;
    }
    
    .product-price {
      font-size: 14px;
      font-weight: 600;
      color: var(--danger-color);
    }
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style>
