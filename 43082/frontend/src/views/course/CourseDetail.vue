<template>
  <div class="course-detail">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card v-if="course">
          <template #header>
            <span style="font-weight: 600">{{ course.name }}</span>
          </template>
          <div class="course-info">
            <div class="info-item">
              <span class="label">教练：</span>
              <span class="value">{{ course.coach?.name }}</span>
            </div>
            <div class="info-item">
              <span class="label">课程类型：</span>
              <span class="value">
                <el-tag :type="getCourseTagType(course.type)">
                  {{ getCourseTypeName(course.type) }}
                </el-tag>
              </span>
            </div>
            <div class="info-item">
              <span class="label">容量：</span>
              <span class="value">{{ course.capacity }} 人</span>
            </div>
            <div class="info-item">
              <span class="label">时长：</span>
              <span class="value">{{ course.duration }} 分钟</span>
            </div>
            <div class="info-item">
              <span class="label">开始时间：</span>
              <span class="value">{{ course.start_time }}</span>
            </div>
            <div class="info-item" v-if="course.weekdays">
              <span class="label">上课周几：</span>
              <span class="value">{{ formatWeekdays(course.weekdays) }}</span>
            </div>
            <div class="info-item">
              <span class="label">地点：</span>
              <span class="value">{{ course.location || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">状态：</span>
              <span class="value">
                <el-tag :type="course.status === 1 ? 'success' : course.status === 2 ? 'danger' : 'info'">
                  {{ getStatusName(course.status) }}
                </el-tag>
              </span>
            </div>
            <div class="info-item" v-if="course.description">
              <span class="label">描述：</span>
              <span class="value">{{ course.description }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="16">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span style="font-weight: 600">课程排期</span>
              <el-button type="primary" size="small" @click="generateSchedules">
                生成排期
              </el-button>
            </div>
          </template>
          
          <el-table :data="schedules" style="width: 100%" v-loading="loadingSchedules">
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
            <el-table-column prop="booked_count" label="已预约" width="100" />
            <el-table-column prop="capacity" label="容量" width="80" />
            <el-table-column label="预约状态" width="120">
              <template #default="{ row }">
                <el-tag v-if="row.booked_count >= row.capacity" type="warning">已满</el-tag>
                <el-tag v-else type="success">可预约</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : row.status === 3 ? 'danger' : 'info'">
                  {{ getScheduleStatusName(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="bookCourse(row.id)"
                  :disabled="row.booked_count >= row.capacity || row.status !== 1"
                >
                  预约
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { courseApi, scheduleApi } from '@/api/course'
import { bookingApi } from '@/api/booking'
import type { Course, CourseSchedule } from '@/types'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const route = useRoute()
const userStore = useUserStore()
const courseId = computed(() => Number(route.params.id))

const course = ref<Course | null>(null)
const schedules = ref<CourseSchedule[]>([])
const loadingSchedules = ref(false)

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getCourseTypeName = (type: string) => {
  const map: Record<string, string> = { single: '单次课', weekly: '周课', monthly: '月课' }
  return map[type] || type
}

const getCourseTagType = (type: string) => {
  const map: Record<string, string> = { single: '', weekly: 'warning', monthly: 'success' }
  return map[type] || ''
}

const getStatusName = (status: number) => {
  const map: Record<number, string> = { 1: '正常', 2: '取消', 3: '结束' }
  return map[status] || '未知'
}

const getScheduleStatusName = (status: number) => {
  const map: Record<number, string> = { 1: '可预约', 2: '已满', 3: '已取消', 4: '已结束' }
  return map[status] || '未知'
}

const formatWeekdays = (weekdays: string) => {
  const dayNames: Record<string, string> = {
    '1': '周一', '2': '周二', '3': '周三', '4': '周四', '5': '周五', '6': '周六', '7': '周日'
  }
  return weekdays.split(',').map(d => dayNames[d] || d).join('、')
}

const loadCourseDetail = async () => {
  try {
    const res = await courseApi.getById(courseId.value)
    course.value = res.data
  } catch (error) {
    console.error(error)
  }
}

const loadSchedules = async () => {
  try {
    loadingSchedules.value = true
    const res = await scheduleApi.getList({ course_id: courseId.value, page_size: 100 })
    schedules.value = res.data
  } catch (error) {
    console.error(error)
  } finally {
    loadingSchedules.value = false
  }
}

const generateSchedules = async () => {
  try {
    await courseApi.generateSchedules(courseId.value)
    ElMessage.success('排期生成成功')
    loadSchedules()
  } catch (error) {
    console.error(error)
  }
}

const bookCourse = async (scheduleId: number) => {
  if (!userStore.userInfo) {
    ElMessage.warning('请先登录')
    return
  }
  
  try {
    await bookingApi.book({
      member_id: userStore.userInfo.id,
      schedule_id: scheduleId
    })
    ElMessage.success('预约成功')
    loadSchedules()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  loadCourseDetail()
  loadSchedules()
})
</script>

<style scoped>
.info-item {
  display: flex;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item .label {
  width: 100px;
  color: #909399;
  flex-shrink: 0;
}

.info-item .value {
  flex: 1;
  color: #303133;
}
</style>
