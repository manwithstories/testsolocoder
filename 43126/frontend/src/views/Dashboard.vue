<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409EFF">
              <el-icon :size="24" color="#fff"><OfficeBuilding /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ roomCount }}</div>
              <div class="stat-label">会议室总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67C23A">
              <el-icon :size="24" color="#fff"><Calendar /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ todayBookingCount }}</div>
              <div class="stat-label">今日预订</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #E6A23C">
              <el-icon :size="24" color="#fff"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ weeklyBookingCount }}</div>
              <div class="stat-label">本周预订</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #F56C6C">
              <el-icon :size="24" color="#fff"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ upcomingCount }}</div>
              <div class="stat-label">即将开始</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card class="panel-card">
          <template #header>
            <div class="card-header">
              <span>最近预订</span>
              <el-button type="primary" text @click="$router.push('/bookings')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentBookings" style="width: 100%" empty-text="暂无预订记录">
            <el-table-column prop="title" label="会议名称" min-width="150" />
            <el-table-column prop="room_name" label="会议室" width="120" />
            <el-table-column label="时间" width="200">
              <template #default="{ row }">
                {{ formatDateTime(row.start_time) }} - {{ formatTime(row.end_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="panel-card">
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/bookings')">
              <el-icon><Plus /></el-icon>
              新建预订
            </el-button>
            <el-button size="large" @click="$router.push('/calendar')">
              <el-icon><Calendar /></el-icon>
              查看日历
            </el-button>
            <el-button v-if="isAdmin" size="large" @click="$router.push('/rooms')">
              <el-icon><OfficeBuilding /></el-icon>
              管理会议室
            </el-button>
            <el-button v-if="isAdmin" size="large" @click="$router.push('/stats')">
              <el-icon><DataAnalysis /></el-icon>
              查看统计
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const userStore = useUserStore()
const isAdmin = ref(userStore.isAdmin)

const roomCount = ref(0)
const todayBookingCount = ref(0)
const weeklyBookingCount = ref(0)
const upcomingCount = ref(0)
const recentBookings = ref<any[]>([])

onMounted(() => {
  loadData()
})

async function loadData() {
  try {
    const [roomsRes, bookingsRes]: any[] = await Promise.all([
      api.listAllRooms(),
      api.getBookings({ page: 1, page_size: 5 })
    ])
    roomCount.value = roomsRes.data?.length || 0
    recentBookings.value = bookingsRes.data?.bookings || []
    todayBookingCount.value = recentBookings.value.filter(b =>
      dayjs(b.start_time).isSame(dayjs(), 'day')
    ).length
    weeklyBookingCount.value = recentBookings.value.filter(b =>
      dayjs(b.start_time).isSame(dayjs(), 'week')
    ).length
    upcomingCount.value = recentBookings.value.filter(b =>
      dayjs(b.start_time).isAfter(dayjs()) && dayjs(b.start_time).diff(dayjs(), 'hour') < 2
    ).length
  } catch (e) {
    console.error(e)
  }
}

function formatDateTime(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function formatTime(date: string) {
  return dayjs(date).format('HH:mm')
}

function getStatusType(status: number) {
  const map: Record<number, string> = { 0: 'warning', 1: 'success', 2: 'info', 3: 'primary' }
  return map[status] || 'info'
}

function getStatusText(status: number) {
  const map: Record<number, string> = { 0: '待确认', 1: '已确认', 2: '已取消', 3: '已完成' }
  return map[status] || '未知'
}
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-card {
  border-radius: 8px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.panel-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quick-actions .el-button {
  justify-content: flex-start;
}
</style>
