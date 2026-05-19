<template>
  <div class="booking-list">
    <el-card>
      <template #header>
        <span style="font-weight: 600">预约管理</span>
      </template>

      <el-table :data="bookings" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="member.name" label="会员姓名" />
        <el-table-column prop="member.phone" label="会员手机号" />
        <el-table-column prop="schedule.course.name" label="课程名称" />
        <el-table-column prop="schedule.course.coach.name" label="教练" />
        <el-table-column prop="schedule.start_time" label="课程时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.schedule?.start_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="booking_time" label="预约时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.booking_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 1"
              type="danger"
              link
              size="small"
              @click="cancelBooking(row.id)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
        @size-change="loadBookings"
        @current-change="loadBookings"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { bookingApi } from '@/api/booking'
import type { Booking } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const bookings = ref<Booking[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const formatDateTime = (date: string) => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'

const getStatusName = (status: number) => {
  const map: Record<number, string> = { 1: '已预约', 2: '已取消', 3: '已签到', 4: '未到场' }
  return map[status] || '未知'
}

const getStatusType = (status: number) => {
  const map: Record<number, string> = { 1: 'primary', 2: 'info', 3: 'success', 4: 'danger' }
  return map[status] || ''
}

const loadBookings = async () => {
  try {
    loading.value = true
    const res = await bookingApi.getList({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    bookings.value = res.data
    pagination.total = res.pagination.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const cancelBooking = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要取消该预约吗？开课前2小时内无法取消', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await bookingApi.cancel(id)
    ElMessage.success('取消成功')
    loadBookings()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

onMounted(loadBookings)
</script>
