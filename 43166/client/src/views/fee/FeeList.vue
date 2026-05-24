<template>
  <div class="fee-list">
    <div class="page-header">
      <h2 class="page-title">费用管理</h2>
    </div>

    <div class="filter-bar">
      <el-form :inline="true" :model="filterForm" @submit.prevent>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 160px">
            <el-option label="全部" value="" />
            <el-option label="待支付" value="pending" />
            <el-option label="已支付" value="paid" />
            <el-option label="已退款" value="refunded" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-wrapper">
      <el-table :data="fees" style="width: 100%" v-loading="loading">
        <el-table-column prop="application.applicationNo" label="申请编号" width="160">
          <template #default="{ row }">
            {{ row.application?.applicationNo || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="application.companyName" label="公司名称">
          <template #default="{ row }">
            {{ row.application?.companyName || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="totalAmount" label="总金额" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.totalAmount) }}
          </template>
        </el-table-column>
        <el-table-column prop="discountAmount" label="优惠金额" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.discountAmount) }}
          </template>
        </el-table-column>
        <el-table-column prop="paidAmount" label="实付金额" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.paidAmount) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-tag', row.status === 'paid' ? 'completed' : 'pending']">
              {{ row.status === 'paid' ? '已支付' : row.status === 'pending' ? '待支付' : '已退款' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="paymentMethod" label="支付方式" width="120">
          <template #default="{ row }">
            {{ getPaymentMethodText(row.paymentMethod) }}
          </template>
        </el-table-column>
        <el-table-column prop="paymentTime" label="支付时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.paymentTime) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper mt-16">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { feeApi } from '@/api/fee'
import { ApplicationFee } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const fees = ref<ApplicationFee[]>([])

const filterForm = reactive({
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await feeApi.getList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: filterForm.status || undefined
    })
    fees.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error('获取费用列表失败:', error)
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.status = ''
  pagination.page = 1
  fetchData()
}

const formatMoney = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

const formatTime = (time?: string | null) => {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

const getPaymentMethodText = (method?: string) => {
  if (!method) return '-'
  const map: Record<string, string> = {
    alipay: '支付宝',
    wechat: '微信支付',
    bank: '银行卡'
  }
  return map[method] || method
}

onMounted(fetchData)
</script>

<style scoped>
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
}
</style>
