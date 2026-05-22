<template>
  <Layout>
    <div class="order-detail-page" v-if="order">
      <div class="page-header">
        <el-button type="text" @click="$router.back()">
          <el-icon><ArrowLeft /></el-icon> 返回
        </el-button>
        <h2>订单详情</h2>
      </div>

      <el-card class="order-info-card">
        <template #header>
          <div class="card-header">
            <span>订单信息</span>
            <el-tag :type="getStatusTagType(order.status)" size="large">{{ getStatusText(order.status) }}</el-tag>
          </div>
        </template>

        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ formatTime(order.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ getPayTypeText(order.pay_type) }}</el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ order.pay_time ? formatTime(order.pay_time) : '-' }}</el-descriptions-item>
          <el-descriptions-item label="订单总额">¥{{ order.total_amount.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="优惠金额">¥{{ order.discount.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="实付金额">¥{{ order.pay_amount.toFixed(2) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card class="ticket-info-card">
        <template #header>
          <span>票券信息 ({{ order.tickets?.length || 0 }}张)</span>
        </template>

        <div class="tickets-list">
          <div v-for="ticket in order.tickets" :key="ticket.id" class="ticket-card">
            <div class="ticket-left">
              <div class="seat-info">{{ ticket.seat_info }}</div>
              <div class="price">¥{{ ticket.price.toFixed(2) }}</div>
            </div>
            <div class="ticket-right">
              <div class="ticket-no">{{ ticket.ticket_no }}</div>
              <div class="ticket-status">
                <el-tag v-if="ticket.status === 1" type="success">已入场</el-tag>
                <el-tag v-else-if="ticket.status === 2" type="info">已退款</el-tag>
                <el-tag v-else type="primary">有效</el-tag>
              </div>
            </div>
            <div class="ticket-seat-info">
              <p>观演人：{{ ticket.real_name }}</p>
              <p>身份证：{{ maskIDCard(ticket.id_card) }}</p>
            </div>
          </div>
        </div>
      </el-card>

      <el-card class="refund-info-card" v-if="order.refund">
        <template #header>
          <span>退款信息</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="退款单号">{{ order.refund.refund_no }}</el-descriptions-item>
          <el-descriptions-item label="退款金额">¥{{ order.refund.refund_amount.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="退款原因">{{ order.refund.reason }}</el-descriptions-item>
          <el-descriptions-item label="退款状态">
            <el-tag :type="order.refund.status === 1 ? 'success' : order.refund.status === 2 ? 'danger' : 'warning'">
              {{ order.refund.status === 1 ? '已通过' : order.refund.status === 2 ? '已拒绝' : '审核中' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="审核备注" :span="2">{{ order.refund.audit_remark || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card class="contact-info-card">
        <template #header>
          <span>联系信息</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="实名人">{{ order.real_name }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ maskIDCard(order.id_card) }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ order.phone }}</el-descriptions-item>
          <el-descriptions-item label="电子邮箱">{{ order.email || '-' }}</el-descriptions-item>
          <el-descriptions-item label="订单备注" :span="2">{{ order.remark || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <div class="action-bar" v-if="order.status === 0 || order.status === 1">
        <el-button size="large" v-if="order.status === 0" @click="cancelOrder">取消订单</el-button>
        <el-button size="large" type="primary" v-if="order.status === 0" @click="payOrder">立即支付</el-button>
        <el-button size="large" type="danger" v-if="order.status === 1" @click="requestRefund">申请退款</el-button>
      </div>
    </div>

    <el-empty v-else description="加载中..." />
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { orderApi } from '@/api'
import Layout from '@/components/Layout.vue'
import type { Order } from '@/types'
import { OrderStatusText, PayTypeText } from '@/types'

const route = useRoute()
const router = useRouter()

const order = ref<Order | null>(null)

async function fetchOrderDetail() {
  try {
    const orderNo = route.params.orderNo as string
    const res = await orderApi.get(orderNo)
    order.value = res
  } catch (err) {
    console.error(err)
  }
}

async function cancelOrder() {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      type: 'warning'
    })
    await orderApi.cancel(order.value!.order_no)
    ElMessage.success('取消成功')
    fetchOrderDetail()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

async function payOrder() {
  try {
    const { value: payType } = await ElMessageBox.prompt('请选择支付方式：1-支付宝，2-微信支付', '支付', {
      inputPattern: /^[12]$/,
      inputErrorMessage: '请输入1或2'
    })
    await orderApi.pay({
      order_no: order.value!.order_no,
      pay_type: Number(payType)
    })
    ElMessage.success('支付成功')
    fetchOrderDetail()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('支付失败')
    }
  }
}

async function requestRefund() {
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
      order_no: order.value!.order_no,
      reason: reason
    })
    ElMessage.success('退款申请已提交，请等待审核')
    fetchOrderDetail()
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

function getPayTypeText(payType: number) {
  return PayTypeText[payType] || '未支付'
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

function maskIDCard(idCard: string) {
  if (!idCard || idCard.length < 8) return idCard
  return idCard.substring(0, 6) + '********' + idCard.substring(idCard.length - 4)
}

onMounted(() => {
  fetchOrderDetail()
})
</script>

<style scoped>
.order-detail-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
}

.order-info-card,
.ticket-info-card,
.refund-info-card,
.contact-info-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tickets-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ticket-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: #fafafa;
  border-radius: 8px;
}

.ticket-left {
  flex: 1;
}

.ticket-left .seat-info {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.ticket-left .price {
  color: #f56c6c;
  font-size: 18px;
  font-weight: bold;
}

.ticket-right {
  text-align: center;
  margin-left: 20px;
}

.ticket-right .ticket-no {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
}

.ticket-seat-info {
  margin-left: 20px;
  border-left: 1px dashed #ddd;
  padding-left: 20px;
}

.ticket-seat-info p {
  margin: 4px 0;
  color: #666;
  font-size: 14px;
}

.action-bar {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}
</style>
