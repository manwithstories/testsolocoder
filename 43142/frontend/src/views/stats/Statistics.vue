<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">统计分析</div>
        <div class="page-subtitle">查看招聘数据统计</div>
      </div>

      <el-card class="mb-20">
        <el-form :inline="true" :model="dateForm">
          <el-form-item label="开始日期">
            <el-date-picker
              v-model="dateForm.startDate"
              type="date"
              placeholder="选择开始日期"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item label="结束日期">
            <el-date-picker
              v-model="dateForm.endDate"
              type="date"
              placeholder="选择结束日期"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchStats">
              <el-icon><Search /></el-icon>
              查询
            </el-button>
            <el-button @click="exportStats">
              <el-icon><Download /></el-icon>
              导出数据
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-row :gutter="20" class="stats-cards mb-20">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">总投递数</div>
              <div class="stat-value">{{ applicationStats.total || 0 }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">面试数</div>
              <div class="stat-value">{{ interviewCount }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">录用数</div>
              <div class="stat-value">{{ hireCount }}</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-label">转化率</div>
              <div class="stat-value">{{ conversionRate }}%</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="mb-20">
            <template #header>
              <span class="card-title">投递状态分布</span>
            </template>
            <div v-loading="loading" style="height: 300px;">
              <v-chart class="chart" :option="statusChartOption" autoresize />
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="mb-20">
            <template #header>
              <span class="card-title">职位投递统计</span>
            </template>
            <div v-loading="loading" style="height: 300px;">
              <v-chart class="chart" :option="jobChartOption" autoresize />
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-card>
        <template #header>
          <span class="card-title">招聘周期分析</span>
        </template>
        <el-table v-loading="loading" :data="cycleStats" style="width: 100%">
          <el-table-column prop="title" label="职位名称" min-width="200" />
          <el-table-column prop="applied_count" label="投递数" width="100" />
          <el-table-column label="平均招聘周期(天)" width="150">
            <template #default="{ row }">
              {{ row.avg_days?.toFixed(1) || '-' }}
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, markRaw } from 'vue'
import { ElMessage } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart, BarChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import { getApplicationStats, getJobStats, getRecruitmentCycleStats, exportJobStats } from '@/api/stats'

use([
  CanvasRenderer,
  PieChart,
  BarChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const loading = ref(false)
const applicationStats = ref<any>({})
const jobStats = ref<any[]>([])
const cycleStats = ref<any[]>([])

const dateForm = reactive({
  startDate: '',
  endDate: ''
})

const interviewCount = computed(() => {
  return applicationStats.value.by_status?.interview || 0
})

const hireCount = computed(() => {
  return applicationStats.value.by_status?.accepted || 0
})

const conversionRate = computed(() => {
  const total = applicationStats.value.total || 0
  const hired = hireCount.value
  if (total === 0) return 0
  return ((hired / total) * 100).toFixed(1)
})

const statusChartOption = computed(() => {
  const byStatus = applicationStats.value.by_status || {}
  const data = Object.entries(byStatus).map(([key, value]) => ({
    name: getStatusText(key),
    value: value
  }))

  return markRaw({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: { show: false },
      data
    }]
  })
})

const jobChartOption = computed(() => {
  const jobs = jobStats.value.slice(0, 10)
  return markRaw({
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: jobs.map(j => j.title),
      axisLabel: { rotate: 30 }
    },
    yAxis: { type: 'value' },
    series: [{
      type: 'bar',
      data: jobs.map(j => j.apply_count),
      itemStyle: { color: '#409eff' }
    }]
  })
})

async function fetchStats() {
  loading.value = true
  try {
    const [appRes, jobRes, cycleRes] = await Promise.all([
      getApplicationStats(dateForm.startDate, dateForm.endDate),
      getJobStats(dateForm.startDate, dateForm.endDate),
      getRecruitmentCycleStats()
    ])

    if (appRes.data) applicationStats.value = appRes.data
    if (jobRes.data) jobStats.value = jobRes.data
    if (cycleRes.data) cycleStats.value = cycleRes.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function exportStats() {
  try {
    await exportJobStats(dateForm.startDate || undefined, dateForm.endDate || undefined)
    ElMessage.success('导出成功')
  } catch (e) {
    // error handled by interceptor
  }
}

function getStatusText(status: string) {
  const texts: Record<string, string> = {
    'pending': '待处理',
    'viewed': '已查看',
    'interested': '感兴趣',
    'interview': '面试中',
    'accepted': '已录用',
    'rejected': '未通过',
    'withdrawn': '已撤回'
  }
  return texts[status] || status
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.stat-card {
  text-align: center;
}

.stat-content {
  padding: 10px 0;
}

.stat-label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #409eff;
}

.card-title {
  font-weight: 600;
}

.chart {
  height: 100%;
}
</style>
