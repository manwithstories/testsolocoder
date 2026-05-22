<template>
  <div class="payment-list">
    <div class="header">
      <h2>支付记录</h2>
      <el-button type="primary" @click="exportVoucher">导出凭证</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="支付方式">
          <el-select v-model="filterForm.method" placeholder="全部" clearable @change="fetchList">
            <el-option label="现金" :value="PaymentMethod.CASH" />
            <el-option label="微信支付" :value="PaymentMethod.WECHAT" />
            <el-option label="支付宝" :value="PaymentMethod.ALIPAY" />
            <el-option label="信用卡" :value="PaymentMethod.CREDIT_CARD" />
            <el-option label="借记卡" :value="PaymentMethod.DEBIT_CARD" />
            <el-option label="转账" :value="PaymentMethod.TRANSFER" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable @change="fetchList">
            <el-option label="待支付" :value="PaymentStatus.PENDING" />
            <el-option label="已支付" :value="PaymentStatus.PAID" />
            <el-option label="已退款" :value="PaymentStatus.REFUNDED" />
            <el-option label="支付失败" :value="PaymentStatus.FAILED" />
          </el-select>
        </el-form-item>
        <el-form-item label="交易时间">
          <el-date-picker
            v-model="filterForm.startDate"
            type="date"
            placeholder="开始日期"
            value-format="YYYY-MM-DD"
            @change="fetchList"
          />
          <span class="date-separator">至</span>
          <el-date-picker
            v-model="filterForm.endDate"
            type="date"
            placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="fetchList"
          />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filterForm.keyword" placeholder="支付号/交易号" @keyup.enter="fetchList" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="paymentNo" label="支付号" width="180" />
        <el-table-column label="订单类型" width="100">
          <template #default="{ row }">
            {{ row.bookingId ? '预订' : '入住' }}
          </template>
        </el-table-column>
        <el-table-column label="订单号" width="180">
          <template #default="{ row }">
            {{ row.booking?.bookingNo || row.checkIn?.checkInNo || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ row.amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="支付方式" width="120">
          <template #default="{ row }">
            {{ getMethodText(row.method) }}
          </template>
        </el-table-column>
        <el-table-column label="支付类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === PaymentStatus.REFUNDED ? 'danger' : 'success'">
              {{ row.status === PaymentStatus.REFUNDED ? '退款' : '支付' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="transactionId" label="交易号" width="180" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === PaymentStatus.PAID"
              size="small"
              type="danger"
              link
              @click="openRefundDialog(row)"
            >退款</el-button>
            <el-button size="small" type="success" link @click="exportSingleVoucher(row)">
              导出凭证
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="refundDialogVisible" title="退款" width="500px">
      <el-form :model="refundForm" label-width="100px">
        <el-form-item label="支付号">
          <el-input :value="currentPayment?.paymentNo" disabled />
        </el-form-item>
        <el-form-item label="原支付金额">
          <el-input :value="'¥' + currentPayment?.amount" disabled />
        </el-form-item>
        <el-form-item label="退款金额">
          <el-input-number v-model="refundForm.amount" :min="0" :max="currentPayment?.amount || 0" />
        </el-form-item>
        <el-form-item label="退款原因">
          <el-input v-model="refundForm.reason" type="textarea" :rows="3" placeholder="请输入退款原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="refundDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleRefund">确认退款</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="支付详情" width="600px">
      <el-descriptions v-if="currentPayment" :column="2" border>
        <el-descriptions-item label="支付号">
          {{ currentPayment.paymentNo }}
        </el-descriptions-item>
        <el-descriptions-item label="金额">
          ¥{{ currentPayment.amount }}
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">
          {{ getMethodText(currentPayment.method) }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTagType(currentPayment.status)">
            {{ getStatusText(currentPayment.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="交易号">
          {{ currentPayment.transactionId || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDateTime(currentPayment.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="关联订单" :span="2">
          {{ currentPayment.booking?.bookingNo || currentPayment.checkIn?.checkInNo || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ currentPayment.remark || '-' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { PaymentMethod, PaymentStatus, type Payment } from '@/types'
import { getPaymentList, refundPayment } from '@/api/payment'

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Payment[]>([])

const filterForm = reactive({
  method: '',
  status: '',
  startDate: '',
  endDate: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const refundDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentPayment = ref<Payment | null>(null)

const refundForm = reactive({
  amount: 0,
  reason: ''
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getPaymentList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: filterForm.keyword || undefined,
      method: filterForm.method as PaymentMethod || undefined,
      status: filterForm.status as PaymentStatus || undefined,
      startDate: filterForm.startDate || undefined,
      endDate: filterForm.endDate || undefined
    })
    tableData.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.method = ''
  filterForm.status = ''
  filterForm.startDate = ''
  filterForm.endDate = ''
  filterForm.keyword = ''
  pagination.page = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchList()
}

const formatDateTime = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

const getMethodText = (method: string) => {
  switch (method) {
    case PaymentMethod.CASH:
      return '现金'
    case PaymentMethod.WECHAT:
      return '微信支付'
    case PaymentMethod.ALIPAY:
      return '支付宝'
    case PaymentMethod.CREDIT_CARD:
      return '信用卡'
    case PaymentMethod.DEBIT_CARD:
      return '借记卡'
    case PaymentMethod.TRANSFER:
      return '转账'
    default:
      return method
  }
}

const getStatusTagType = (status: string) => {
  switch (status) {
    case PaymentStatus.PENDING:
      return 'warning'
    case PaymentStatus.PAID:
      return 'success'
    case PaymentStatus.REFUNDED:
      return 'info'
    case PaymentStatus.FAILED:
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case PaymentStatus.PENDING:
      return '待支付'
    case PaymentStatus.PAID:
      return '已支付'
    case PaymentStatus.REFUNDED:
      return '已退款'
    case PaymentStatus.FAILED:
      return '支付失败'
    default:
      return status
  }
}

const viewDetail = (row: Payment) => {
  currentPayment.value = row
  detailDialogVisible.value = true
}

const openRefundDialog = (row: Payment) => {
  currentPayment.value = row
  refundForm.amount = row.amount
  refundForm.reason = ''
  refundDialogVisible.value = true
}

const handleRefund = async () => {
  if (!currentPayment.value) return
  if (!refundForm.amount) {
    ElMessage.warning('请输入退款金额')
    return
  }
  if (!refundForm.reason) {
    ElMessage.warning('请输入退款原因')
    return
  }
  submitting.value = true
  try {
    await refundPayment(currentPayment.value.id, {
      amount: refundForm.amount,
      reason: refundForm.reason
    })
    ElMessage.success('退款成功')
    refundDialogVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '退款失败')
  } finally {
    submitting.value = false
  }
}

const exportVoucher = () => {
  ElMessage.info('导出功能开发中...')
}

const exportSingleVoucher = (row: Payment) => {
  ElMessage.info(`正在导出支付凭证：${row.paymentNo}`)
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.payment-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-form {
  margin: 0;
}

.date-separator {
  margin: 0 10px;
}

.table-card {
  margin-bottom: 20px;
}

.amount {
  color: #f56c6c;
  font-weight: 600;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
