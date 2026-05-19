<template>
  <div class="schedule-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span style="font-weight: 600">课程排期</span>
          <div class="header-actions">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              style="width: 300px; margin-right: 12px"
              @change="loadSchedules"
            />
          </div>
        </div>
      </template>

      <el-table :data="schedules" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="course.name" label="课程名称" />
        <el-table-column prop="course.coach.name" label="教练" />
        <el-table-column prop="start_time" label="开始时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.start_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="end_time" label="结束时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.end_time) }}
          </template>
        </el-table-column>
        <el-table-column label="预约情况" width="140">
          <template #default="{ row }">
            <el-progress
              :percentage="Math.round((row.booked_count / row.capacity) * 100)"
              :stroke-width="12"
            />
            <span style="font-size: 12px; color: #909399">
              {{ row.booked_count }}/{{ row.capacity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1 && row.booked_count >= row.capacity" type="warning">已满</el-tag>
            <el-tag v-else-if="row.status === 1" type="success">可预约</el-tag>
            <el-tag v-else-if="row.status === 3" type="danger">已取消</el-tag>
            <el-tag v-else type="info">已结束</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="bookCourse(row)"
              :disabled="row.status !== 1 || row.booked_count >= row.capacity"
            >
              预约
            </el-button>
            <el-button
              v-if="row.status === 1 && row.booked_count >= row.capacity"
              type="warning"
              link
              size="small"
              @click="addToWaitlist(row)"
            >
              候补
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
        @size-change="loadSchedules"
        @current-change="loadSchedules"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { scheduleApi } from '@/api/course'
import { bookingApi } from '@/api/booking'
import type { CourseSchedule } from '@/types'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const userStore = useUserStore()
const loading = ref(false)
const schedules = ref<CourseSchedule[]>([])
const dateRange = ref<string[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')

const loadSchedules = async () => {
  try {
    loading.value = true
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const res = await scheduleApi.getList(params)
    schedules.value = res.data
    pagination.total = res.pagination.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const bookCourse = async (row: CourseSchedule) => {
  if (!userStore.userInfo) {
    ElMessage.warning('请先登录')
    return
  }
  
  try {
    await bookingApi.book({
      member_id: userStore.userInfo.id,
      schedule_id: row.id
    })
    ElMessage.success('预约成功')
    loadSchedules()
  } catch (error) {
    console.error(error)
  }
}

const addToWaitlist = async (row: CourseSchedule) => {
  if (!userStore.userInfo) {
    ElMessage.warning('请先登录')
    return
  }
  
  try {
    await bookingApi.addToWaitlist({
      member_id: userStore.userInfo.id,
      schedule_id: row.id
    })
    ElMessage.success('已加入等待列表，有位置时会自动通知您')
  } catch (error) {
    console.error(error)
  }
}

onMounted(loadSchedules)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}
</style>
