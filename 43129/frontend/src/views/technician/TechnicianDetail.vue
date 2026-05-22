<template>
  <div class="page-container" v-loading="loading">
    <el-card shadow="never" v-if="technician">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button :icon="ArrowLeft" @click="$router.back()">返回</el-button>
            <span style="margin-left: 10px; font-weight: 600">技师详情</span>
          </div>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="姓名">{{ technician.name }}</el-descriptions-item>
        <el-descriptions-item label="职称">{{ technician.title || '-' }}</el-descriptions-item>
        <el-descriptions-item label="擅长项目">{{ technician.specialties || '-' }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ technician.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="评分">
          <el-rate :model-value="technician.rating" disabled />
          <span style="margin-left: 8px">{{ technician.rating }} ({{ technician.review_count }}条评价)</span>
        </el-descriptions-item>
        <el-descriptions-item label="工作时间">{{ technician.work_start_time }} - {{ technician.work_end_time }}</el-descriptions-item>
      </el-descriptions>

      <el-divider />

      <h3 style="margin-bottom: 16px">日历视图</h3>
      <el-date-picker
        v-model="selectedDate"
        type="date"
        placeholder="选择日期"
        :clearable="false"
        @change="fetchSchedule"
      />

      <el-table :data="schedule" stripe style="margin-top: 16px" v-loading="scheduleLoading">
        <el-table-column prop="customer.name" label="顾客" />
        <el-table-column prop="service.name" label="服务" />
        <el-table-column prop="start_time" label="开始时间" width="120" />
        <el-table-column prop="end_time" label="结束时间" width="120" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getTechnician, getTechnicianSchedule } from '@/api/technician'
import { ArrowLeft } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Technician } from '@/types'

const route = useRoute()
const loading = ref(false)
const scheduleLoading = ref(false)
const technician = ref<Technician | null>(null)
const selectedDate = ref(dayjs().format('YYYY-MM-DD'))
const schedule = ref<any[]>([])

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'info',
    confirmed: 'primary',
    paid: 'success',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    paid: '已支付',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const fetchTechnician = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const res = await getTechnician(id)
    technician.value = res.data
    fetchSchedule()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchSchedule = async () => {
  if (!technician.value) return
  scheduleLoading.value = true
  try {
    const res = await getTechnicianSchedule(technician.value.id, { date: selectedDate.value })
    schedule.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    scheduleLoading.value = false
  }
}

onMounted(fetchTechnician)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-left {
    display: flex;
    align-items: center;
  }
}
</style>
