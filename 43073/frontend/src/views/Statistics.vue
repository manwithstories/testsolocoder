<template>
  <div class="statistics">
    <div class="page-header">
      <h2 class="page-title">统计报表</h2>
      <el-button type="primary" @click="handleExport">
        <el-icon><Download /></el-icon>
        导出Excel
      </el-button>
    </div>

    <el-card>
      <div class="search-bar">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 300px"
        />
        <el-select v-model="search.activityId" placeholder="选择活动" clearable style="width: 200px">
          <el-option v-for="act in activities" :key="act.id" :label="act.title" :value="act.id" />
        </el-select>
        <el-button type="primary" @click="loadData">查询</el-button>
      </div>
    </el-card>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>按活动统计</template>
          <el-table :data="activityStats" style="width: 100%">
            <el-table-column prop="activityTitle" label="活动名称" min-width="180" />
            <el-table-column prop="totalOrders" label="订单数" width="100" />
            <el-table-column prop="totalTickets" label="售票数" width="100" />
            <el-table-column prop="totalAmount" label="收入(元)" width="120">
              <template #default="{ row }">¥{{ row.totalAmount.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>按票型统计</template>
          <el-table :data="ticketStats" style="width: 100%">
            <el-table-column prop="ticketTypeName" label="票型名称" />
            <el-table-column prop="totalSold" label="售出数量" width="120" />
            <el-table-column prop="totalAmount" label="收入(元)" width="120">
              <template #default="{ row }">¥{{ row.totalAmount.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 20px">
      <template #header>每日销售趋势</template>
      <div ref="chartRef" style="width: 100%; height: 400px"></div>
    </el-card>

    <el-card style="margin-top: 20px">
      <template #header>每日明细</template>
      <el-table :data="dailyStats" style="width: 100%">
        <el-table-column prop="date" label="日期" width="140" />
        <el-table-column prop="totalOrders" label="订单数" width="120" />
        <el-table-column prop="totalTickets" label="售票数" width="120" />
        <el-table-column prop="totalAmount" label="收入(元)">
          <template #default="{ row }">¥{{ row.totalAmount.toFixed(2) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getActivityStatistics, getTicketTypeStatistics, getDailyStatistics, exportStatistics } from '@/api/statistics'
import { getActivityList } from '@/api/activity'
import * as echarts from 'echarts'

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const dateRange = ref<string[]>([])
const activityStats = ref<any[]>([])
const ticketStats = ref<any[]>([])
const dailyStats = ref<any[]>([])
const activities = ref<any[]>([])

const search = reactive({
  activityId: ''
})

const loadActivities = async () => {
  try {
    const res = await getActivityList({ page: 1, pageSize: 100 })
    activities.value = res.list
  } catch (error) {
    console.error(error)
  }
}

const loadData = async () => {
  try {
    const params: any = {}
    if (dateRange.value?.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }
    if (search.activityId) {
      params.activityId = search.activityId
    }

    const [actRes, ticketRes, dailyRes] = await Promise.all([
      getActivityStatistics(params),
      getTicketTypeStatistics(params),
      getDailyStatistics(params)
    ])

    activityStats.value = actRes
    ticketStats.value = ticketRes
    dailyStats.value = dailyRes

    nextTick(() => {
      renderChart()
    })
  } catch (error) {
    console.error(error)
  }
}

const renderChart = () => {
  if (!chartRef.value) return

  if (!chart) {
    chart = echarts.init(chartRef.value)
  }

  const dates = dailyStats.value.map(d => d.date).reverse()
  const amounts = dailyStats.value.map(d => d.totalAmount).reverse()
  const orders = dailyStats.value.map(d => d.totalOrders).reverse()

  chart.setOption({
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['收入', '订单数']
    },
    xAxis: {
      type: 'category',
      data: dates
    },
    yAxis: [
      {
        type: 'value',
        name: '收入(元)'
      },
      {
        type: 'value',
        name: '订单数'
      }
    ],
    series: [
      {
        name: '收入',
        type: 'line',
        smooth: true,
        data: amounts,
        areaStyle: {}
      },
      {
        name: '订单数',
        type: 'bar',
        yAxisIndex: 1,
        data: orders
      }
    ]
  })
}

const handleExport = async () => {
  try {
    const params: any = {}
    if (dateRange.value?.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }
    if (search.activityId) {
      params.activityId = search.activityId
    }

    const blob = await exportStatistics(params) as unknown as Blob
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `statistics_${Date.now()}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  loadActivities()
  loadData()
})
</script>
