<template>
  <div class="admin-orders">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="订单号">
          <el-input v-model="filterForm.keyword" placeholder="请输入订单号" clearable />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 150px;">
            <el-option label="待支付" :value="0" />
            <el-option label="已支付" :value="1" />
            <el-option label="已取消" :value="2" />
            <el-option label="已退款" :value="3" />
            <el-option label="退款中" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="下单时间">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchOrders">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
          <el-button type="success" @click="handleExport">导出Excel</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="stats-card">
      <el-row :gutter="24">
        <el-col :span="6">
          <div class="stat-item">
            <div class="stat-value">{{ stats.totalOrders }}</div>
            <div class="stat-label">总订单数</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item paid">
            <div class="stat-value">¥{{ stats.totalAmount.toFixed(2) }}</div>
            <div class="stat-label">总销售额</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item pending">
            <div class="stat-value">{{ stats.pendingOrders }}</div>
            <div class="stat-label">待支付订单</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item refund">
            <div class="stat-value">{{ stats.refundOrders }}</div>
            <div class="stat-label">退款中订单</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <el-card class="table-card">
      <el-table :data="orders" style="width: 100%" v-loading="loading">
        <el-table-column prop="order_no" label="订单号" width="200" />
        <el-table-column prop="real_name" label="实名人" width="120" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column label="购票数量" width="100">
          <template #default="{ row }">
            {{ row.tickets?.length || 0 }} 张
          </template>
        </el-table-column>
        <el-table-column prop="total_amount" label="订单金额" width="120">
          <template #default="{ row }">
            ¥{{ row.total_amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="pay_amount" label="实付金额" width="120">
          <template #default="{ row }">
            ¥{{ row.pay_amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="支付方式" width="100">
          <template #default="{ row }">
            {{ getPayTypeText(row.pay_type) }}
          </template>
        </el-table-column>
        <el-table-column label="订单状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewDetail(row)">详情</el-button>
            <el-button
              size="small"
              type="warning"
              v-if="row.status === 4 && row.refund?.status === 0"
              @click="handleAuditRefund(row)"
            >
              退款审核
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        class="pagination"
        @current-change="fetchOrders"
      />
    </el-card>

    <el-dialog v-model="showDetailDialog" title="订单详情" width="800px">
      <template v-if="currentOrder">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ currentOrder.order_no }}</el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ formatTime(currentOrder.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ getPayTypeText(currentOrder.pay_type) }}</el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ currentOrder.pay_time ? formatTime(currentOrder.pay_time) : '-' }}</el-descriptions-item>
          <el-descriptions-item label="订单总额">¥{{ currentOrder.total_amount.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="优惠金额">¥{{ currentOrder.discount.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="实付金额" class="pay-amount">¥{{ currentOrder.pay_amount.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-tag :type="getStatusTagType(currentOrder.status)">
              {{ getStatusText(currentOrder.status) }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <h4 style="margin: 20px 0 10px;">票券信息</h4>
        <el-table :data="currentOrder.tickets" style="width: 100%">
          <el-table-column prop="ticket_no" label="票号" width="200" />
          <el-table-column prop="seat_info" label="座位" />
          <el-table-column prop="price" label="票价" width="100">
            <template #default="{ row }">¥{{ row.price.toFixed(2) }}</template>
          </el-table-column>
          <el-table-column prop="real_name" label="观演人" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.status === 1" type="success">已使用</el-tag>
              <el-tag v-else-if="row.status === 2" type="info">已退款</el-tag>
              <el-tag v-else type="primary">有效</el-tag>
            </template>
          </el-table-column>
        </el-table>

        <h4 style="margin: 20px 0 10px;">联系信息</h4>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="实名人">{{ currentOrder.real_name }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ currentOrder.id_card }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ currentOrder.phone }}</el-descriptions-item>
          <el-descriptions-item label="电子邮箱">{{ currentOrder.email || '-' }}</el-descriptions-item>
        </el-descriptions>
      </template>
    </el-dialog>

    <el-dialog v-model="showAuditDialog" title="退款审核" width="500px">
      <template v-if="currentOrder && currentOrder.refund">
        <el-alert
          title="退款申请信息"
          type="info"
          :closable="false"
          style="margin-bottom: 20px;"
        >
          <p>退款单号：{{ currentOrder.refund.refund_no }}</p>
          <p>退款金额：¥{{ currentOrder.refund.refund_amount.toFixed(2) }}</p>
          <p>退款原因：{{ currentOrder.refund.reason }}</p>
        </el-alert>

        <el-form :model="auditForm" label-width="100px">
          <el-form-item label="审核结果">
            <el-radio-group v-model="auditForm.status">
              <el-radio :value="1">通过</el-radio>
              <el-radio :value="2">拒绝</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="审核备注">
            <el-input v-model="auditForm.audit_remark" type="textarea" :rows="3" placeholder="请输入审核备注" />
          </el-form-item>
        </el-form>
      </template>
      <template #footer>
        <el-button @click="showAuditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitAudit">确认审核</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { orderApi } from '@/api'
import type { Order } from '@/types'
import { OrderStatusText, PayTypeText } from '@/types'

const loading = ref(false)
const orders = ref<Order[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filterForm = reactive({
  keyword: '',
  status: undefined as number | undefined,
  dateRange: [] as string[]
})

const showDetailDialog = ref(false)
const showAuditDialog = ref(false)
const currentOrder = ref<Order | null>(null)

const auditForm = reactive({
  status: 1,
  audit_remark: ''
})

const stats = computed(() => {
  let totalOrders = orders.value.length
  let totalAmount = 0
  let pendingOrders = 0
  let refundOrders = 0

  orders.value.forEach(order => {
    if (order.status === 1) {
      totalAmount += order.pay_amount
    }
    if (order.status === 0) {
      pendingOrders++
    }
    if (order.status === 4) {
      refundOrders++
    }
  })

  return {
    totalOrders,
    totalAmount,
    pendingOrders,
    refundOrders
  }
})

async function fetchOrders() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filterForm.keyword) {
      params.keyword = filterForm.keyword
    }
    if (filterForm.status !== undefined) {
      params.status = filterForm.status
    }
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }

    const res = await orderApi.list(params)
    orders.value = res.list
    total.value = res.pagination?.total || 0
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function resetFilter() {
  filterForm.keyword = ''
  filterForm.status = undefined
  filterForm.dateRange = []
  page.value = 1
  fetchOrders()
}

async function handleExport() {
  try {
    const params: any = {}
    if (filterForm.status !== undefined) {
      params.status = filterForm.status
    }
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }

    const blob = await orderApi.export(params)
    const url = window.URL.createObjectURL(new Blob([blob as any]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `订单报表_${dayjs().format('YYYYMMDDHHmmss')}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    ElMessage.success('导出成功')
  } catch (err) {
    ElMessage.error('导出失败')
  }
}

function handleViewDetail(row: Order) {
  currentOrder.value = row
  showDetailDialog.value = true
}

function handleAuditRefund(row: Order) {
  currentOrder.value = row
  auditForm.status = 1
  auditForm.audit_remark = ''
  showAuditDialog.value = true
}

async function handleSubmitAudit() {
  try {
    await ElMessageBox.confirm('确认提交审核结果吗？', '提示', {
      type: 'warning'
    })
    await orderApi.auditRefund({
      refund_no: currentOrder.value!.refund!.refund_no,
      status: auditForm.status,
      audit_remark: auditForm.audit_remark
    })
    ElMessage.success('审核成功')
    showAuditDialog.value = false
    fetchOrders()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('审核失败')
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

fetchOrders()
</script>

<style lang="scss" scoped>
.admin-orders {
  .filter-card {
    margin-bottom: 20px;
  }

  .stats-card {
    margin-bottom: 20px;

    .stat-item {
      text-align: center;
      padding: 20px;
      border-radius: 8px;
      background: #f5f7fa;

      .stat-value {
        font-size: 28px;
        font-weight: bold;
        color: #409eff;
        margin-bottom: 8px;
      }

      .stat-label {
        color: #999;
        font-size: 14px;
      }

      &.paid .stat-value {
        color: #67c23a;
      }

      &.pending .stat-value {
        color: #e6a23c;
      }

      &.refund .stat-value {
        color: #f56c6c;
      }
    }
  }

  .table-card {
    .pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }
  }

  :deep(.pay-amount .el-descriptions-item__content) {
    color: #f56c6c;
    font-weight: bold;
    font-size: 18px;
  }
}
</style>
