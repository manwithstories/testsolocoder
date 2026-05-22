<template>
  <div class="my-orders">
    <h3>我的订单</h3>
    <el-tabs v-model="activeTab" @tab-change="fetchOrders">
      <el-tab-pane label="我买到的" name="buyer" />
      <el-tab-pane label="我卖出的" name="seller" />
    </el-tabs>
    <el-table :data="orders" v-loading="loading">
      <el-table-column prop="order_no" label="订单号" width="180" />
      <el-table-column label="拍卖品" min-width="200">
        <template #default="{ row }">
          <div class="item-title">
            <img :src="getImage(row)" />
            <span>{{ row.auction_item?.title }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="price" label="价格" width="120">
        <template #default="{ row }">¥{{ row.price.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="viewOrder(row)">查看</el-button>
          <el-button
            v-if="activeTab === 'buyer' && row.status === 0"
            link
            type="success"
            size="small"
            @click="payOrder(row)"
          >支付</el-button>
          <el-button
            v-if="activeTab === 'seller' && row.status === 1"
            link
            type="primary"
            size="small"
            @click="shipOrder(row)"
          >发货</el-button>
          <el-button
            v-if="activeTab === 'buyer' && row.status === 2"
            link
            type="success"
            size="small"
            @click="confirmDelivery(row)"
          >确认收货</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchOrders"
    />

    <el-dialog v-model="showShipDialog" title="发货" width="400px">
      <el-form>
        <el-form-item label="物流单号">
          <el-input v-model="trackingNo" placeholder="请输入物流单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showShipDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmShip">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { Order } from '@/types'
import { orderApi } from '@/api'

const activeTab = ref('buyer')
const orders = ref<Order[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const showShipDialog = ref(false)
const trackingNo = ref('')
const currentOrder = ref<Order | null>(null)

const fetchOrders = async () => {
  loading.value = true
  try {
    const api = activeTab.value === 'buyer' ? orderApi.getBuyerOrders : orderApi.getSellerOrders
    const res = await api({ page: page.value, page_size: pageSize.value })
    orders.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const getImage = (order: Order) => {
  if (order.auction_item?.images?.[0]?.url) {
    return order.auction_item.images[0].url
  }
  return 'https://via.placeholder.com/50x50?text=No'
}

const statusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'warning', 1: 'primary', 2: 'info', 3: 'success', 4: 'success', 5: 'danger',
  }
  return map[status] || 'info'
}

const statusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待支付', 1: '已支付', 2: '已发货', 3: '已送达', 4: '已完成', 5: '已取消',
  }
  return map[status] || ''
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

const viewOrder = (order: Order) => {
  ElMessage.info(`查看订单: ${order.order_no}`)
}

const payOrder = async (order: Order) => {
  try {
    await orderApi.pay(order.id, { method: 'balance' })
    ElMessage.success('支付成功')
    fetchOrders()
  } catch (e) {}
}

const shipOrder = (order: Order) => {
  currentOrder.value = order
  trackingNo.value = ''
  showShipDialog.value = true
}

const confirmShip = async () => {
  if (!currentOrder.value) return
  try {
    await orderApi.ship(currentOrder.value.id, trackingNo.value)
    ElMessage.success('发货成功')
    showShipDialog.value = false
    fetchOrders()
  } catch (e) {}
}

const confirmDelivery = async (order: Order) => {
  try {
    await orderApi.confirmDelivery(order.id)
    ElMessage.success('已确认收货')
    fetchOrders()
  } catch (e) {}
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
h3 {
  margin: 0 0 20px;
}

.item-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.item-title img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 4px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
