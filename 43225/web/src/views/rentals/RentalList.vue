<template>
  <div class="rental-list">
    <div class="page-header">
      <h2 class="page-title">租赁订单</h2>
    </div>

    <div class="filter-bar">
      <el-select v-model="filters.status" placeholder="订单状态" clearable style="width: 140px">
        <el-option label="待确认" value="pending" />
        <el-option label="已确认" value="confirmed" />
        <el-option label="进行中" value="active" />
        <el-option label="已完成" value="completed" />
        <el-option label="已取消" value="cancelled" />
      </el-select>
      <el-button type="primary" @click="fetchRentals">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
    </div>

    <div class="card-container">
      <el-table :data="rentals" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="订单号" width="200" show-overflow-tooltip />
        <el-table-column label="船只">
          <template #default="{ row }">
            {{ row.ship?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="租客">
          <template #default="{ row }">
            {{ row.tenant?.full_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="租赁时间">
          <template #default="{ row }">
            {{ formatDate(row.start_date) }} ~ {{ formatDate(row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column label="金额">
          <template #default="{ row }">
            {{ row.currency }} {{ row.total_amount }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/rentals/${row.id}`)">
              详情
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              link
              @click="handleUpdateStatus(row.id, 'confirmed')"
            >
              确认
            </el-button>
            <el-button
              v-if="row.status === 'confirmed'"
              type="primary"
              link
              @click="handleUpdateStatus(row.id, 'active')"
            >
              开始
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              type="success"
              link
              @click="handleUpdateStatus(row.id, 'completed')"
            >
              完成
            </el-button>
            <el-button
              v-if="['pending', 'confirmed'].includes(row.status)"
              type="danger"
              link
              @click="handleCancel(row.id)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="pagination.total"
          layout="total, prev, pager, next"
          @current-change="fetchRentals"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRentalsApi, updateRentalStatusApi, cancelRentalApi } from '@/api/rental'
import type { Rental } from '@/types/rental'
import dayjs from 'dayjs'

const loading = ref(false)
const rentals = ref<Rental[]>([])

const filters = reactive({
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    confirmed: 'primary',
    active: 'success',
    completed: 'info',
    cancelled: 'danger',
    refunded: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || status
}

const fetchRentals = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (filters.status) params.status = filters.status

    const res: any = await getRentalsApi(params)
    rentals.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch rentals:', error)
  } finally {
    loading.value = false
  }
}

const handleUpdateStatus = async (id: string, status: string) => {
  try {
    await ElMessageBox.confirm('确定要更新订单状态吗？', '提示', {
      type: 'warning'
    })
    await updateRentalStatusApi(id, { status: status as any })
    ElMessage.success('状态更新成功')
    fetchRentals()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('状态更新失败')
    }
  }
}

const handleCancel = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      type: 'warning'
    })
    await cancelRentalApi(id)
    ElMessage.success('订单已取消')
    fetchRentals()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

onMounted(fetchRentals)
</script>

<style lang="scss" scoped>
.rental-list {
  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 16px;
  }
}
</style>
