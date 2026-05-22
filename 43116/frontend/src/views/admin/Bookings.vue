<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">预订管理</h1>
    </div>

    <div class="search-bar">
      <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px">
        <el-option label="待确认" value="pending" />
        <el-option label="已确认" value="confirmed" />
        <el-option label="已取消" value="cancelled" />
        <el-option label="已完成" value="completed" />
      </el-select>
      <el-button type="primary" @click="loadBookings">搜索</el-button>
    </div>

    <el-table :data="bookings" v-loading="loading" style="width: 100%">
      <el-table-column prop="booking_no" label="预订号" width="180" />
      <el-table-column label="用户" width="120">
        <template #default="{ row }">
          {{ row.user?.username }}
        </template>
      </el-table-column>
      <el-table-column label="车辆" min-width="150">
        <template #default="{ row }">
          {{ row.car?.brand }} {{ row.car?.model }}
        </template>
      </el-table-column>
      <el-table-column label="取车门店" min-width="120">
        <template #default="{ row }">
          {{ row.pickup_store?.name }}
        </template>
      </el-table-column>
      <el-table-column label="还车门店" min-width="120">
        <template #default="{ row }">
          {{ row.return_store?.name }}
        </template>
      </el-table-column>
      <el-table-column label="取车时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.pickup_time) }}
        </template>
      </el-table-column>
      <el-table-column label="还车时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.return_time) }}
        </template>
      </el-table-column>
      <el-table-column label="金额" width="100" align="right">
        <template #default="{ row }">
          ¥{{ row.final_price.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'pending'"
            type="success"
            link
            size="small"
            @click="handleConfirm(row)"
          >
            确认
          </el-button>
          <el-button
            v-if="row.status === 'confirmed'"
            type="primary"
            link
            size="small"
            @click="handleComplete(row)"
          >
            完成
          </el-button>
          <el-button
            v-if="row.status === 'pending' || row.status === 'confirmed'"
            type="danger"
            link
            size="small"
            @click="handleCancel(row)"
          >
            取消
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadBookings"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { bookingApi } from '@/api'
import type { Booking } from '@/types'

const bookings = ref<Booking[]>([])
const loading = ref(false)

const filters = reactive({
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadBookings()
})

const loadBookings = async () => {
  loading.value = true
  try {
    const res = await bookingApi.getBookings({
      page: pagination.page,
      page_size: pagination.pageSize,
      status: filters.status
    })
    bookings.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const handleConfirm = async (row: Booking) => {
  try {
    await ElMessageBox.confirm('确定要确认该预订吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await bookingApi.confirmBooking(row.id)
    ElMessage.success('确认成功')
    loadBookings()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '操作失败')
    }
  }
}

const handleComplete = async (row: Booking) => {
  try {
    await ElMessageBox.confirm('确定要完成该预订吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await bookingApi.completeBooking(row.id)
    ElMessage.success('完成成功')
    loadBookings()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '操作失败')
    }
  }
}

const handleCancel = async (row: Booking) => {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入取消原因', '取消预订', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '请输入取消原因'
    })

    await bookingApi.cancelBooking(row.id, reason)
    ElMessage.success('取消成功')
    loadBookings()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '操作失败')
    }
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    confirmed: 'primary',
    cancelled: 'info',
    completed: 'success',
    no_show: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    cancelled: '已取消',
    completed: '已完成',
    no_show: '未取车'
  }
  return map[status] || status
}
</script>
