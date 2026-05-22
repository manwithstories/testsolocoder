<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>报表分析</span>
          <div class="header-actions">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="fetchData"
            />
            <el-button type="success" @click="exportExcel">导出Excel</el-button>
            <el-button type="primary" @click="exportPDF">导出PDF</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20" v-if="reportData">
        <el-col :span="24">
          <el-card shadow="never" class="stat-card">
            <div class="stat-item">
              <span class="label">总收入</span>
              <span class="value">¥{{ reportData.revenue_report?.total_revenue?.toFixed(2) || '0.00' }}</span>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <span>收入趋势</span>
            </template>
            <div ref="revenueChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <span>技师业绩</span>
            </template>
            <div ref="techChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <span>热门服务排行</span>
            </template>
            <el-table :data="reportData?.service_ranks || []" stripe>
              <el-table-column type="index" label="#" width="60" />
              <el-table-column prop="service_name" label="服务名称" />
              <el-table-column prop="revenue" label="收入">
                <template #default="{ row }">¥{{ row.revenue?.toFixed(2) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <span>技师业绩明细</span>
            </template>
            <el-table :data="reportData?.technician_performances || []" stripe>
              <el-table-column type="index" label="#" width="60" />
              <el-table-column prop="technician_name" label="技师" />
              <el-table-column prop="revenue" label="业绩">
                <template #default="{ row }">¥{{ row.revenue?.toFixed(2) }}</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import * as echarts from 'echarts'
import { getFullReport, exportExcel as apiExportExcel, exportPDF as apiExportPDF } from '@/api/report'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { ReportData } from '@/types'

const dateRange = ref<string[]>([dayjs().subtract(7, 'day').format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')])
const reportData = ref<ReportData | null>(null)
const revenueChartRef = ref<HTMLElement>()
const techChartRef = ref<HTMLElement>()
let revenueChart: echarts.ECharts | null = null
let techChart: echarts.ECharts | null = null

const fetchData = async () => {
  if (!dateRange.value || dateRange.value.length !== 2) return

  try {
    const res = await getFullReport({
      start_date: dateRange.value[0],
      end_date: dateRange.value[1]
    })
    reportData.value = res.data
    await nextTick()
    initCharts()
  } catch (e) {
    console.error(e)
  }
}

const initCharts = () => {
  if (revenueChartRef.value && reportData.value) {
    if (revenueChart) revenueChart.dispose()
    revenueChart = echarts.init(revenueChartRef.value)
    revenueChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: reportData.value.revenue_report?.daily_data?.map(d => d.date) || []
      },
      yAxis: {
        type: 'value',
        axisLabel: { formatter: '¥{value}' }
      },
      series: [{
        type: 'line',
        smooth: true,
        areaStyle: {},
        data: reportData.value.revenue_report?.daily_data?.map(d => d.revenue) || []
      }]
    })
  }

  if (techChartRef.value && reportData.value) {
    if (techChart) techChart.dispose()
    techChart = echarts.init(techChartRef.value)
    techChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: reportData.value.technician_performances?.map(t => t.technician_name) || []
      },
      yAxis: {
        type: 'value',
        axisLabel: { formatter: '¥{value}' }
      },
      series: [{
        type: 'bar',
        data: reportData.value.technician_performances?.map(t => t.revenue) || []
      }]
    })
  }
}

const exportExcel = async () => {
  try {
    const res: any = await apiExportExcel({
      start_date: dateRange.value[0],
      end_date: dateRange.value[1]
    })
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `report_${dayjs().format('YYYYMMDD')}.xlsx`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

const exportPDF = async () => {
  try {
    const res: any = await apiExportPDF({
      start_date: dateRange.value[0],
      end_date: dateRange.value[1]
    })
    const blob = new Blob([res], { type: 'application/pdf' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `report_${dayjs().format('YYYYMMDD')}.pdf`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

const handleResize = () => {
  revenueChart?.resize()
  techChart?.resize()
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  revenueChart?.dispose()
  techChart?.dispose()
})
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }

  .stat-card {
    text-align: center;
    padding: 20px 0;

    .stat-item {
      .label {
        font-size: 16px;
        color: #909399;
        margin-right: 20px;
      }

      .value {
        font-size: 32px;
        font-weight: bold;
        color: #409EFF;
      }
    }
  }
}
</style>
