<template>
  <div class="order-list">
    <div class="page-header">
      <h2 class="page-title">订单管理</h2>
    </div>

    <el-card>
      <div class="search-bar">
        <el-input v-model="search.orderNo" placeholder="订单号" clearable style="width: 200px" />
        <el-select v-model="search.activityId" placeholder="选择活动" clearable style="width: 200px">
          <el-option v-for="act in activities" :key="act.id" :label="act.title" :value="act.id" />
        </el-select>
        <el-select v-model="search.status" placeholder="订单状态" clearable style="width: 140px">
          <el-option label="待支付" value="pending" />
          <el-option label="已支付" value="paid" />
          <el-option label="已取消" value="canceled" />
          <el-option label="已退款" value="refunded" />
        </el-select>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="orderNo" label="订单号" width="180" />
        <el-table-column prop="activity.title" label="活动名称" min-width="180" />
        <el-table-column prop="user.username" label="购买用户" width="120" />
        <el-table-column label="购票信息" min-width="200">
          <template #default="{ row }">
            <div v-for="item in row.orderItems" :key="item.id">
              {{ item.ticketType?.name }} x {{ item.quantity }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="totalAmount" label="原价(元)" width="100">
          <template #default="{ row }">¥{{ row.totalAmount.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="discount" label="优惠(元)" width="100">
          <template #default="{ row }">¥{{ row.discount.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="payAmount" label="实付(元)" width="100">
          <template #default="{ row }">
            <span class="pay-amount">¥{{ row.payAmount.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160" :formatter="formatDate" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <el-button size="small" @click="viewDetail(row)">详情</el-button>
              <el-button v-if="row.status === 'pending'" size="small" type="primary" @click="handlePay(row)">支付</el-button>
              <el-button v-if="row.status === 'pending'" size="small" type="danger" @click="handleCancel(row)">取消</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <el-dialog v-model="detailVisible" title="订单详情" width="600px">
      <el-descriptions v-if="orderDetail" :column="2" border>
        <el-descriptions-item label="订单号">{{ orderDetail.orderNo }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(orderDetail.status)">{{ getStatusText(orderDetail.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="活动">{{ orderDetail.activity?.title }}</el-descriptions-item>
        <el-descriptions-item label="购买用户">{{ orderDetail.user?.username }}</el-descriptions-item>
        <el-descriptions-item label="原价">¥{{ orderDetail.totalAmount.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="优惠">¥{{ orderDetail.discount.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="实付">¥{{ orderDetail.payAmount.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="优惠券">
          {{ orderDetail.coupon?.code || '无' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(null, null, orderDetail.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="支付时间" :span="2">
          {{ orderDetail.paidAt ? formatDate(null, null, orderDetail.paidAt) : '未支付' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider>购票明细</el-divider>
      <el-table :data="orderDetail?.orderItems || []" size="small">
        <el-table-column prop="ticketType.name" label="票型" />
        <el-table-column prop="unitPrice" label="单价" />
        <el-table-column prop="quantity" label="数量" />
        <el-table-column prop="subtotal" label="小计" />
        <el-table-column label="签到码">
          <template #default="{ row }">
            <div v-for="checkin in row.checkIns" :key="checkin.id" style="margin-bottom: 5px">
              <el-tag size="small" :type="checkin.checkedIn ? 'success' : 'info'">
                {{ checkin.qrCode }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrderList, getOrder, payOrder, cancelOrder } from '@/api/order'
import { getActivityList } from '@/api/activity'
import dayjs from 'dayjs'

const loading = ref(false)
const list = ref<any[]>([])
const activities = ref<any[]>([])
const detailVisible = ref(false)
const orderDetail = ref<any>(null)

const search = reactive({
  orderNo: '',
  activityId: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadActivities = async () => {
  try {
    const res = await getActivityList({ page: 1, pageSize: 100 })
    activities.value = res.list
  } catch (error) {
    console.error(error)
  }
}

const loadData = async () => {
  try {
    loading.value = true
    const res = await getOrderList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      orderNo: search.orderNo,
      activityId: search.activityId,
      status: search.status
    })
    list.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const viewDetail = async (row: any) => {
  try {
    orderDetail.value = await getOrder(row.id)
    detailVisible.value = true
  } catch (error) {
    console.error(error)
  }
}

const handlePay = async (row: any) => {
  try {
    await ElMessageBox.confirm('确认支付该订单吗？', '提示', { type: 'warning' })
    await payOrder(row.id)
    ElMessage.success('支付成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const handleCancel = async (row: any) => {
  try {
    await ElMessageBox.confirm('确认取消该订单吗？', '提示', { type: 'warning' })
    await cancelOrder(row.id)
    ElMessage.success('取消成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = { pending: 'warning', paid: 'success', canceled: 'info', refunded: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待支付', paid: '已支付', canceled: '已取消', refunded: '已退款' }
  return map[status] || status
}

const formatDate = (_row: any, _column: any, value: string) => value ? dayjs(value).format('YYYY-MM-DD HH:mm:ss') : '-'

onMounted(() => {
  loadActivities()
  loadData()
})
</script>

<style scoped lang="scss">
.pay-amount {
  color: #f56c6c;
  font-weight: 600;
}
</style>
