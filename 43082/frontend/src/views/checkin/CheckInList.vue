<template>
  <div class="checkin-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span style="font-weight: 600">签到记录</span>
          <el-date-picker
            v-model="filterDate"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 200px"
            @change="loadByDate"
          />
        </div>
      </template>

      <el-table :data="checkIns" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="member.name" label="会员姓名" />
        <el-table-column prop="member.phone" label="会员手机号" />
        <el-table-column prop="schedule.course.name" label="课程名称">
          <template #default="{ row }">
            {{ row.schedule?.course?.name || '自由训练' }}
          </template>
        </el-table-column>
        <el-table-column prop="check_in_time" label="签到时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.check_in_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="check_type" label="签到类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.check_type === 1 ? 'success' : 'warning'">
              {{ row.check_type === 1 ? '正常签到' : '补签' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
        @size-change="loadCheckIns"
        @current-change="loadCheckIns"
      />
    </el-card>

    <el-card style="margin-top: 20px" v-if="userStore.userInfo">
      <template #header>
        <span style="font-weight: 600">快速签到</span>
      </template>
      <div class="quick-checkin">
        <el-select v-model="selectedSchedule" placeholder="选择课程（可选）" style="width: 300px; margin-right: 12px">
          <el-option
            v-for="schedule in todaySchedules"
            :key="schedule.id"
            :label="`${schedule.course?.name} - ${formatTime(schedule.start_time)}`"
            :value="schedule.id"
          />
        </el-select>
        <el-button type="primary" @click="quickCheckIn" :loading="checkingIn">
          <el-icon><CircleCheck /></el-icon>
          立即签到
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { checkInApi, bookingApi } from '@/api/booking'
import { scheduleApi } from '@/api/course'
import type { CheckIn, CourseSchedule } from '@/types'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const userStore = useUserStore()
const loading = ref(false)
const checkingIn = ref(false)
const checkIns = ref<CheckIn[]>([])
const filterDate = ref('')
const selectedSchedule = ref<number | null>(null)
const todaySchedules = ref<CourseSchedule[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')
const formatTime = (date: string) => dayjs(date).format('HH:mm')

const loadCheckIns = async () => {
  if (!userStore.userInfo) return
  
  try {
    loading.value = true
    const res = await checkInApi.getByMember(userStore.userInfo.id, {
      page: pagination.page,
      page_size: pagination.pageSize
    })
    checkIns.value = res.data
    pagination.total = res.pagination.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadByDate = async () => {
  if (!filterDate.value) {
    loadCheckIns()
    return
  }
  
  try {
    loading.value = true
    const res = await checkInApi.getByDate(filterDate.value)
    checkIns.value = res.data
    pagination.total = res.data.length
    pagination.page = 1
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadTodayBookings = async () => {
  if (!userStore.userInfo) return
  
  try {
    const res = await bookingApi.getByMember(userStore.userInfo.id, { page_size: 100 })
    const today = dayjs().format('YYYY-MM-DD')
    todaySchedules.value = res.data
      .filter(b => b.status === 1 && dayjs(b.schedule?.start_time).format('YYYY-MM-DD') === today)
      .map(b => b.schedule!)
      .filter(Boolean)
  } catch (error) {
    console.error(error)
  }
}

const quickCheckIn = async () => {
  if (!userStore.userInfo) return
  
  try {
    checkingIn.value = true
    await checkInApi.checkIn({
      member_id: userStore.userInfo.id,
      schedule_id: selectedSchedule.value || undefined
    })
    ElMessage.success('签到成功')
    selectedSchedule.value = null
    loadCheckIns()
    loadTodayBookings()
  } catch (error) {
    console.error(error)
  } finally {
    checkingIn.value = false
  }
}

onMounted(() => {
  loadCheckIns()
  loadTodayBookings()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quick-checkin {
  display: flex;
  align-items: center;
}
</style>
