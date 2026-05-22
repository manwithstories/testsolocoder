<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon activity">
              <el-icon size="32"><Calendar /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.activities }}</p>
              <p class="stat-label">活动总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon order">
              <el-icon size="32"><List /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.orders }}</p>
              <p class="stat-label">订单总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon ticket">
              <el-icon size="32"><Ticket /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.tickets }}</p>
              <p class="stat-label">售票总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon revenue">
              <el-icon size="32"><Money /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">¥{{ stats.revenue.toFixed(2) }}</p>
              <p class="stat-label">总收入</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最近活动</span>
          </template>
          <el-table :data="recentActivities" style="width: 100%">
            <el-table-column prop="title" label="活动名称" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="capacity" label="容量" />
            <el-table-column prop="createdAt" label="创建时间" :formatter="formatDate" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最近订单</span>
          </template>
          <el-table :data="recentOrders" style="width: 100%">
            <el-table-column prop="orderNo" label="订单号" />
            <el-table-column prop="payAmount" label="金额">
              <template #default="{ row }">¥{{ row.payAmount.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)">{{ getOrderStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" :formatter="formatDate" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { getActivityList } from '@/api/activity'
import { getOrderList } from '@/api/order'
import dayjs from 'dayjs'

const stats = reactive({
  activities: 0,
  orders: 0,
  tickets: 0,
  revenue: 0
})

const recentActivities = ref<any[]>([])
const recentOrders = ref<any[]>([])

const loadData = async () => {
  try {
    const [activityRes, orderRes] = await Promise.all([
      getActivityList({ page: 1, pageSize: 5 }),
      getOrderList({ page: 1, pageSize: 5 })
    ])
    
    stats.activities = activityRes.total
    recentActivities.value = activityRes.list
    
    stats.orders = orderRes.total
    recentOrders.value = orderRes.list
    stats.revenue = orderRes.list.reduce((sum: number, o: any) => sum + o.payAmount, 0)
    stats.tickets = orderRes.list.reduce((sum: number, o: any) => {
      return sum + (o.orderItems?.reduce((s: number, item: any) => s + item.quantity, 0) || 0)
    }, 0)
  } catch (error) {
    console.error(error)
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = { draft: 'info', published: 'success', canceled: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { draft: '草稿', published: '已发布', canceled: '已取消' }
  return map[status] || status
}

const getOrderStatusType = (status: string) => {
  const map: Record<string, string> = { pending: 'warning', paid: 'success', canceled: 'info', refunded: 'danger' }
  return map[status] || 'info'
}

const getOrderStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待支付', paid: '已支付', canceled: '已取消', refunded: '已退款' }
  return map[status] || status
}

const formatDate = (_row: any, _column: any, value: string) => {
  return dayjs(value).format('YYYY-MM-DD HH:mm')
}

onMounted(loadData)
</script>

<style scoped lang="scss">
.stat-card {
  .stat-content {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;

    &.activity { background: linear-gradient(135deg, #667eea, #764ba2); }
    &.order { background: linear-gradient(135deg, #f093fb, #f5576c); }
    &.ticket { background: linear-gradient(135deg, #4facfe, #00f2fe); }
    &.revenue { background: linear-gradient(135deg, #43e97b, #38f9d7); }
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    margin: 0;
    color: #303133;
  }

  .stat-label {
    margin: 4px 0 0;
    color: #909399;
  }
}
</style>
