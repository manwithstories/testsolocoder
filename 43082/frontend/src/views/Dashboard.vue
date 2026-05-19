<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-cards">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon blue">
              <el-icon size="28"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.total_members || 0 }}</div>
              <div class="stat-label">会员总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon green">
              <el-icon size="28"><Medal /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.active_memberships || 0 }}</div>
              <div class="stat-label">有效会员卡</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon orange">
              <el-icon size="28"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.today_check_ins || 0 }}</div>
              <div class="stat-label">今日签到</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon purple">
              <el-icon size="28"><Tickets /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.today_bookings || 0 }}</div>
              <div class="stat-label">今日预约</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span style="font-weight: 600">即将开始的课程</span>
          </template>
          <el-table :data="upcomingCourses" style="width: 100%" v-loading="loading">
            <el-table-column prop="course.name" label="课程名称" />
            <el-table-column prop="course.coach.name" label="教练" />
            <el-table-column prop="start_time" label="开始时间">
              <template #default="{ row }">
                {{ formatDateTime(row.start_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="booked_count" label="已预约" />
            <el-table-column prop="capacity" label="容量" />
            <el-table-column label="状态">
              <template #default="{ row }">
                <el-tag v-if="row.booked_count >= row.capacity" type="warning">已满</el-tag>
                <el-tag v-else type="success">可预约</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <span style="font-weight: 600">快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/schedules')">
              <el-icon><Calendar /></el-icon>
              查看课程排期
            </el-button>
            <el-button type="success" @click="$router.push('/my-bookings')">
              <el-icon><Document /></el-icon>
              我的预约
            </el-button>
            <el-button type="warning" @click="$router.push('/check-ins')">
              <el-icon><CircleCheck /></el-icon>
              签到记录
            </el-button>
            <el-button type="info" @click="$router.push('/stats')">
              <el-icon><DataAnalysis /></el-icon>
              数据统计
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { statsApi, scheduleApi } from '@/api/booking'
import dayjs from 'dayjs'

const loading = ref(false)
const dashboardData = ref<any>({})
const upcomingCourses = ref<any[]>([])

const loadDashboardData = async () => {
  try {
    const res = await statsApi.getDashboard()
    dashboardData.value = res.data
  } catch (error) {
    console.error('Failed to load dashboard:', error)
  }
}

const loadUpcomingCourses = async () => {
  try {
    loading.value = true
    const res = await scheduleApi.getAvailable()
    upcomingCourses.value = res.data.slice(0, 5)
  } catch (error) {
    console.error('Failed to load courses:', error)
  } finally {
    loading.value = false
  }
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  loadDashboardData()
  loadUpcomingCourses()
})
</script>

<style scoped>
.stats-cards {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-icon.blue {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.green {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.stat-icon.orange {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.purple {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quick-actions .el-button {
  width: 100%;
  justify-content: center;
  height: 44px;
}
</style>
