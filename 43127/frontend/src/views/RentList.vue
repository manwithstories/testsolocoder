<template>
  <div class="rent-list">
    <div class="page-header">
      <h2 class="page-title">租金管理</h2>
      <div>
        <el-button type="success" @click="handleGenerate">
          <el-icon><Refresh /></el-icon>生成月度账单
        </el-button>
        <el-button type="warning" @click="handleCalculateLateFee">
          <el-icon><Money /></el-icon>计算滞纳金
        </el-button>
        <el-button type="primary" @click="handleExport">
          <el-icon><Download /></el-icon>导出Excel
        </el-button>
      </div>
    </div>

    <div class="search-bar card">
      <el-date-picker
        v-model="monthFilter"
        type="month"
        placeholder="选择月份"
        style="width: 150px"
        value-format="YYYY-MM"
      />
      <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 150px">
        <el-option label="未缴" :value="0" />
        <el-option label="已缴" :value="1" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <div class="summary">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="总金额">¥{{ summary.totalAmount }}</el-descriptions-item>
          <el-descriptions-item label="已缴金额">
            <span style="color: #67c23a;">¥{{ summary.paidAmount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="未缴金额">
            <span style="color: #f56c6c;">¥{{ summary.totalAmount - summary.paidAmount }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <el-table :data="bills" v-loading="loading">
        <el-table-column prop="month" label="月份" width="120" />
        <el-table-column label="房源" min-width="150">
          <template #default="{ row }">
            {{ row.contract?.property?.title }}
          </template>
        </el-table-column>
        <el-table-column label="租户" width="150">
          <template #default="{ row }">
            {{ row.contract?.tenant?.name }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.amount }}</template>
        </el-table-column>
        <el-table-column prop="lateFee" label="滞纳金" width="120">
          <template #default="{ row }">
            <span v-if="row.lateFee > 0" style="color: #f56c6c;">¥{{ row.lateFee }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="dueDate" label="应缴日期" width="120">
          <template #default="{ row }">{{ formatDate(row.dueDate) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '已缴' : '未缴' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handlePay(row)"
              v-if="row.status === 0"
            >缴费</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>

    <el-dialog v-model="payDialogVisible" title="缴费确认" width="400px">
      <el-form :model="payForm" label-width="80px">
        <el-form-item label="账单金额">¥{{ payForm.amount }}</el-form-item>
        <el-form-item label="实缴金额">
          <el-input-number v-model="payForm.paidAmount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="payForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="payDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmPay">确认缴费</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { RentRecord } from '@/types'
import { getRentBills, generateRentBills, payRentBill, calculateLateFee } from '@/api/business'
import { exportRentRecords } from '@/api/business'
import dayjs from 'dayjs'

const loading = ref(false)
const bills = ref<RentRecord[]>([])
const monthFilter = ref('')
const statusFilter = ref<number | ''>('')
const payDialogVisible = ref(false)
const payBillId = ref<number | null>(null)

const summary = reactive({
  totalAmount: 0,
  paidAmount: 0
})

const payForm = reactive({
  amount: 0,
  paidAmount: 0,
  remark: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getRentBills({
      page: pagination.page,
      pageSize: pagination.pageSize,
      month: monthFilter.value || undefined,
      status: statusFilter.value || undefined
    })
    bills.value = res.data.list
    pagination.total = res.data.total
    summary.totalAmount = res.data.totalAmount || 0
    summary.paidAmount = res.data.paidAmount || 0
  } catch (error) {
    console.error('Failed to load bills:', error)
  } finally {
    loading.value = false
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
}

async function handleGenerate() {
  try {
    const { value } = await ElMessageBox.prompt('请选择要生成账单的月份', '生成月度账单', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /^\d{4}-\d{2}$/,
      inputErrorMessage: '月份格式错误，请使用 YYYY-MM 格式',
      inputValue: monthFilter.value || dayjs().format('YYYY-MM')
    })
    const res = await generateRentBills(value)
    ElMessage.success(`已生成 ${res.data.generated} 条账单`)
    loadData()
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error(error)
    }
  }
}

async function handleCalculateLateFee() {
  try {
    await ElMessageBox.confirm('确定要计算滞纳金吗？', '提示', { type: 'warning' })
    const res = await calculateLateFee()
    ElMessage.success(`已处理 ${res.data.processed} 条账单`)
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

function handlePay(row: RentRecord) {
  payBillId.value = row.id
  payForm.amount = row.amount
  payForm.paidAmount = row.amount
  payForm.remark = ''
  payDialogVisible.value = true
}

async function confirmPay() {
  if (!payBillId.value) return
  try {
    await payRentBill(payBillId.value, {
      paidAmount: payForm.paidAmount,
      remark: payForm.remark
    })
    ElMessage.success('缴费成功')
    payDialogVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
  }
}

async function handleExport() {
  try {
    const res = await exportRentRecords()
    const blob = new Blob([res as any], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `rent_records_${dayjs().format('YYYY-MM')}.xlsx`
    link.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error(error)
  }
}
</script>

<style scoped>
.rent-list {
  padding: 0;
}

.summary {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
