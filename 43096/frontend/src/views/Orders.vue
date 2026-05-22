<template>
  <Layout>
    <div class="orders-page">
      <div class="page-header">
        <h2>我的订单</h2>
        <el-tabs v-model="activeTab" class="order-tabs" @tab-change="handleTabChange">
          <el-tab-pane label="全部" name="-1" />
          <el-tab-pane label="待支付" name="0" />
          <el-tab-pane label="已支付" name="1" />
          <el-tab-pane label="已取消" name="2" />
          <el-tab-pane label="退款中" name="4" />
          <el-tab-pane label="已退款" name="3" />
        </el-tabs>
      </div>

      <div class="orders-list" v-if="orders.length > 0">
        <div
          v-for="order in orders"
          :key="order.id"
          class="order-card"
          @click="goToDetail(order.order_no)"
        >
          <div class="order-header">
            <span class="order-no">订单号：{{ order.order_no }}</span>
            <el-tag :type="getStatusTagType(order.status)">{{ getStatusText(order.status) }}</el-tag>
          </div>

          <div class="order-content">
            <div class="tickets" v-if="order.tickets">
              <div v-for="ticket in order.tickets" :key="ticket.id" class="ticket-item">
                <span class="seat-info">{{ ticket.seat_info }}</span>
                <span class="ticket-price">¥{{ ticket.price }}</span>
              </div>
            </div>

            <div class="order-info">
              <p>实名人：{{ order.real_name }}</p>
              <p>下单时间：{{ formatTime(order.created_at) }}</p>
            </div>
          </div>

          <div class="order-footer">
            <span class="total">共 {{ order.tickets?.length || 0 }} 张，合计：<span class="price">¥{{ order.pay_amount.toFixed(2) }}</span></span>
            <div class="actions">
              <el-button size="small" v-if="order.status === 0" type="primary" @click.stop="goToPay(order.order_no)">
                去支付
              </el-button>
              <el-button size="small" v-if="order.status === 0" @click.stop="cancelOrder(order.order_no)">
                取消订单
              </el-button>
              <el-button size="small" v-if="order.status === 1" type="danger" @click.stop="requestRefund(order.order_no)">
                申请退款
              </el-button>
              <el-button size="small" @click.stop="goToDetail(order.order_no)">
                查看详情
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <el-empty v-else description="暂无订单数据" />

      <el-pagination
        v-if="total > 0"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        class="pagination"
        @current-change="fetchOrders"
      />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { orderApi } from '@/api'
import Layout from '@/components/Layout.vue'
import type { Order } from '@/types'
import { OrderStatusText } from '@/types'

const router = useRouter()

const orders = ref<Order[]>([])
const activeTab = ref('-1')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

async function fetchOrders() {
  try {
    const status = Number(activeTab.value)
    const res = await orderApi.list({
      page: page.value,
      page_size: pageSize.value,
      status: status >= 0 ? status : undefined
    })
    orders.value = res.list
    total.value = res.pagination.total
  } catch (err) {
    console.error(err)
  }
}

function handleTabChange() {
  page.value = 1
  fetchOrders()
}

function goToDetail(orderNo: string) {
  router.push(`/order/${orderNo}`)
}

function goToPay(orderNo: string) {
  router.push(`/order/${orderNo}`)
}

async function cancelOrder(orderNo: string) {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      type: 'warning'
    })
    await orderApi.cancel(orderNo)
    ElMessage.success('取消成功')
    fetchOrders()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

async function requestRefund(orderNo: string) {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入退款原因', '申请退款', {
      inputPlaceholder: '请输入退款原因',
      inputValidator: (value: string) => {
        if (!value || value.length < 5) {
          return '请输入至少5个字符的退款原因'
        }
        return true
      }
    })
    await orderApi.refund({
      order_no: orderNo,
      reason: reason
    })
    ElMessage.success('退款申请已提交，请等待审核')
    fetchOrders()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('申请失败')
    }
  }
}

function getStatusText(status: number) {
  return OrderStatusText[status] || '未知'
}

function getStatusTagType(status: number) {
  const types: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'info',
    4: 'warning'
  }
  return types[status] || 'info'
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchOrders()
})
</script>

<style lang="scss" scoped>
.orders-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  background: white;
  padding: 20px 24px;
  border-radius: 12px;
  margin-bottom: 20px;

  h2 {
    margin: 0 0 16px 0;
  }

  :deep(.el-tabs__header) {
    margin: 0;
  }
}

.order-card {
  background: white;
  padding: 20px 24px;
  border-radius: 12px;
  margin-bottom: 16px;
  cursor: pointer;
  transition: box-shadow 0.2s;

  &:hover {
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  }
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;

  .order-no {
    color: #666;
    font-size: 14px;
  }
}

.order-content {
  padding: 16px 0;

  .tickets {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 12px;
  }

  .ticket-item {
    background: #f5f7fa;
    padding: 8px 12px;
    border-radius: 6px;
    display: flex;
    gap: 16px;
    font-size: 14px;

    .ticket-price {
      color: #f56c6c;
      font-weight: 600;
    }
  }

  .order-info p {
    margin: 4px 0;
    color: #666;
    font-size: 14px;
  }
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;

  .total {
    font-size: 14px;

    .price {
      color: #f56c6c;
      font-size: 18px;
      font-weight: bold;
    }
  }

  .actions {
    display: flex;
    gap: 8px;
  }
}

.pagination {
  margin-top: 20px;
  justify-content: center;
}
</style>
