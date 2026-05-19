<template>
  <div class="my-bookings">
    <el-card>
      <template #header>
        <div class="card-header">
          <span style="font-weight: 600">我的预约</span>
          <el-radio-group v-model="filterStatus" @change="loadBookings">
            <el-radio-button value="">全部</el-radio-button>
            <el-radio-button value="1">已预约</el-radio-button>
            <el-radio-button value="3">已签到</el-radio-button>
            <el-radio-button value="2">已取消</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="filteredBookings" style="width: 100%" v-loading="loading">
        <el-table-column prop="schedule.course.name" label="课程名称" />
        <el-table-column prop="schedule.course.coach.name" label="教练" />
        <el-table-column prop="schedule.start_time" label="开始时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.schedule?.start_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="schedule.end_time" label="结束时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.schedule?.end_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="schedule.location" label="地点" />
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
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 1 && canCancel(row)"
              type="danger"
              link
              size="small"
              @click="cancelBooking(row.id)"
            >
              取消预约
            </el-button>
            <el-button
              v-if="row.status === 1 && !canCheckIn(row)"
              type="info"
              link
              size="small"
              disabled
            >
              未到签到时间
            </el-button>
            <el-button
              v-if="row.status === 1 && canCheckIn(row)"
              type="success"
              link
              size="small"
              @click="checkIn(row)"
            >
              签到
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="filteredBookings.length === 0 && !loading" description="暂无预约记录" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { bookingApi, checkInApi } from '@/api/booking'
import type { Booking } from '@/types'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const userStore = useUserStore()
const loading = ref(false)
const bookings = ref<Booking[]>([])
const filterStatus = ref('')

const filteredBookings = computed(() => {
  if (!filterStatus.value) return bookings.value
  return bookings.value.filter(b => b.status === Number(filterStatus.value))
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

const canCancel = (row: Booking) => {
  if (!row.schedule) return false
  const hoursUntil = dayjs(row.schedule.start_time).diff(dayjs(), 'hour')
  return hoursUntil > 2
}

const canCheckIn = (row: Booking) => {
  if (!row.schedule) return false
  const hoursUntil = dayjs(row.schedule.start_time).diff(dayjs(), 'hour')
  return hoursUntil <= 2 && hoursUntil >= -1
}

const loadBookings = async () => {
  if (!userStore.userInfo) return
  
  try {
    loading.value = true
    const res = await bookingApi.getByMember(userStore.userInfo.id, { page_size: 100 })
    bookings.value = res.data.sort((a, b) => 
      new Date(b.schedule?.start_time || '').getTime() - new Date(a.schedule?.start_time || '').getTime()
    )
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

const checkIn = async (row: Booking) => {
  try {
    await checkInApi.checkIn({
      member_id: userStore.userInfo!.id,
      schedule_id: row.schedule_id
    })
    ElMessage.success('签到成功')
    loadBookings()
  } catch (error) {
    console.error(error)
  }
}

onMounted(loadBookings)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
