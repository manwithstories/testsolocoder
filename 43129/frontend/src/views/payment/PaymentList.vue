<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>支付记录</span>
        </div>
      </template>

      <el-table :data="payments" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="顾客">
          <template #default="{ row }">
            {{ row.appointment?.customer?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额">
          <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="pay_method" label="支付方式">
          <template #default="{ row }">{{ getPayMethod(row.pay_method) }}</template>
        </el-table-column>
        <el-table-column prop="points_used" label="积分抵扣" />
        <el-table-column prop="transaction_no" label="交易号" show-overflow-tooltip />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="支付时间" width="180">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPayments } from '@/api/payment'
import dayjs from 'dayjs'
import type { Payment } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const payments = ref<Payment[]>([])

const getPayMethod = (method: string) => {
  const methods: Record<string, string> = {
    cash: '现金',
    card: '会员卡',
    points: '积分抵扣',
    wechat: '微信支付',
    alipay: '支付宝'
  }
  return methods[method] || method
}

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getPayments({ page: page.value, page_size: pageSize.value })
    payments.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
