<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6" v-for="card in statCards" :key="card.label">
        <el-card shadow="hover" class="stat-card">
          <div class="card-content">
            <div class="card-info">
              <p class="card-label">{{ card.label }}</p>
              <p class="card-value">{{ card.value }}</p>
            </div>
            <el-icon :size="48" :color="card.color">
              <component :is="card.icon" />
            </el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>近7日收入趋势</span>
            </div>
          </template>
          <div ref="revenueChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>热门服务排行</span>
            </div>
          </template>
          <div ref="serviceChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最新预约</span>
              <el-link type="primary" @click="$router.push('/appointments')">查看全部</el-link>
            </div>
          </template>
          <el-table :data="latestAppointments" stripe>
            <el-table-column prop="customer.name" label="顾客" />
            <el-table-column prop="service.name" label="服务" />
            <el-table-column prop="technician.name" label="技师" />
            <el-table-column label="时间">
              <template #default="{ row }">
                {{ formatDateTime(row.appointment_date, row.start_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>低库存预警</span>
              <el-link type="primary" @click="$router.push('/products')">查看全部</el-link>
            </div>
          </template>
          <el-table :data="lowStockProducts" stripe>
            <el-table-column prop="name" label="产品名称" />
            <el-table-column prop="stock" label="当前库存" />
            <el-table-column prop="threshold" label="预警阈值" />
            <el-table-column label="状态">
              <template #default="{ row }">
                <el-tag type="danger" size="small">库存不足</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { getAppointments } from '@/api/appointment'
import { getLowStockProducts } from '@/api/product'
import { User, Service, Calendar, Warning } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const revenueChartRef = ref<HTMLElement>()
const serviceChartRef = ref<HTMLElement>()
let revenueChart: echarts.ECharts | null = null
let serviceChart: echarts.ECharts | null = null

const latestAppointments = ref<any[]>([])
const lowStockProducts = ref<any[]>([])

const statCards = computed(() => [
  { label: '今日预约', value: '0', icon: Calendar, color: '#409EFF' },
  { label: '今日收入', value: '¥0', icon: Service, color: '#67C23A' },
  { label: '顾客总数', value: '0', icon: User, color: '#E6A23C' },
  { label: '低库存预警', value: lowStockProducts.value.length, icon: Warning, color: '#F56C6C' }
])

const formatDateTime = (date: string, time: string) => {
  return dayjs(date).format('MM-DD') + ' ' + time
}

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

const fetchLatestAppointments = async () => {
  try {
    const res = await getAppointments({ page: 1, page_size: 5 })
    latestAppointments.value = res.data.list
  } catch (e) {
    console.error(e)
  }
}

const fetchLowStockProducts = async () => {
  try {
    const res = await getLowStockProducts()
    lowStockProducts.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const initCharts = () => {
  if (revenueChartRef.value) {
    revenueChart = echarts.init(revenueChartRef.value)
    revenueChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '¥{value}'
        }
      },
      series: [{
        name: '收入',
        type: 'line',
        smooth: true,
        areaStyle: {},
        data: [1200, 1900, 1500, 2200, 1800, 2800, 2400]
      }]
    })
  }

  if (serviceChartRef.value) {
    serviceChart = echarts.init(serviceChartRef.value)
    serviceChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [{
        name: '服务占比',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: { show: false, position: 'center' },
        emphasis: {
          label: { show: true, fontSize: 20, fontWeight: 'bold' }
        },
        labelLine: { show: false },
        data: [
          { value: 1048, name: '剪发' },
          { value: 735, name: '染发' },
          { value: 580, name: '烫发' },
          { value: 484, name: '美容' },
          { value: 300, name: 'SPA' }
        ]
      }]
    })
  }
}

const handleResize = () => {
  revenueChart?.resize()
  serviceChart?.resize()
}

onMounted(async () => {
  await Promise.all([fetchLatestAppointments(), fetchLowStockProducts()])
  await nextTick()
  initCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  revenueChart?.dispose()
  serviceChart?.dispose()
})
</script>

<style scoped lang="scss">
.dashboard {
  .stat-card {
    .card-content {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .card-info {
      .card-label {
        font-size: 14px;
        color: #909399;
        margin-bottom: 8px;
      }

      .card-value {
        font-size: 28px;
        font-weight: bold;
        color: #303133;
      }
    }
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
