<template>
  <div class="dashboard-index">
    <el-row :gutter="20">
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-blue">
            <el-icon size="28"><Tickets /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.todayVisitors }}</div>
            <div class="stat-label">今日参观</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-green">
            <el-icon size="28"><Calendar /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.todayReservations }}</div>
            <div class="stat-label">今日预约</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-orange">
            <el-icon size="28"><Picture /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalCollections }}</div>
            <div class="stat-label">藏品总数</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-red">
            <el-icon size="28"><Money /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">¥{{ stats.monthRevenue }}</div>
            <div class="stat-label">本月营收</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="16">
        <div class="chart-card card-shadow p-20">
          <h3 class="card-title">参观趋势</h3>
          <div ref="visitChartRef" class="chart-container"></div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="chart-card card-shadow p-20">
          <h3 class="card-title">展览热度排名</h3>
          <div ref="rankChartRef" class="chart-container rank-chart"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="list-card card-shadow p-20">
          <div class="flex-between mb-20">
            <h3 class="card-title">最新预约</h3>
            <el-button type="primary" text @click="$router.push('/dashboard/reservations')">查看全部</el-button>
          </div>
          <el-table :data="recentReservations" v-loading="loading">
            <el-table-column prop="user.nickname" label="用户" width="100" />
            <el-table-column prop="exhibition.title" label="展览" show-overflow-tooltip />
            <el-table-column prop="visitor_count" label="人数" width="80" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)" size="small">
                  {{ statusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="160">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="list-card card-shadow p-20">
          <div class="flex-between mb-20">
            <h3 class="card-title">热门藏品</h3>
            <el-button type="primary" text @click="$router.push('/dashboard/collections')">查看全部</el-button>
          </div>
          <el-table :data="hotCollections" v-loading="loading">
            <el-table-column prop="name" label="藏品名称" show-overflow-tooltip />
            <el-table-column prop="category.name" label="分类" width="100" />
            <el-table-column prop="view_count" label="浏览量" width="100">
              <template #default="{ row }">
                <el-tag type="warning" size="small">{{ row.view_count }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                  {{ row.status === 'active' ? '展出' : '其他' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import type { Reservation, Collection } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const visitChartRef = ref<HTMLElement>()
const rankChartRef = ref<HTMLElement>()
let visitChart: echarts.ECharts | null = null
let rankChart: echarts.ECharts | null = null

const stats = ref({
  todayVisitors: 128,
  todayReservations: 45,
  totalCollections: 12580,
  monthRevenue: '12.5万'
})

const recentReservations = ref<Reservation[]>([])
const hotCollections = ref<Collection[]>([])

const statusType = (status: string) => {
  if (status === 'confirmed') return 'success'
  if (status === 'cancelled') return 'info'
  if (status === 'completed') return 'success'
  return 'warning'
}

const statusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    cancelled: '已取消',
    completed: '已完成'
  }
  return map[status] || status
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const initCharts = () => {
  if (visitChartRef.value) {
    visitChart = echarts.init(visitChartRef.value)
    visitChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['参观人数', '预约数'] },
      xAxis: {
        type: 'category',
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '参观人数',
          type: 'line',
          smooth: true,
          data: [120, 132, 101, 134, 90, 230, 210],
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
            ])
          }
        },
        {
          name: '预约数',
          type: 'line',
          smooth: true,
          data: [80, 92, 71, 104, 60, 180, 160]
        }
      ]
    })
  }

  if (rankChartRef.value) {
    rankChart = echarts.init(rankChartRef.value)
    rankChart.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'value' },
      yAxis: {
        type: 'category',
        data: ['古代书画', '瓷器珍品', '青铜器', '玉器', '近现代艺术']
      },
      series: [
        {
          type: 'bar',
          data: [820, 732, 601, 534, 490],
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: '#667eea' },
              { offset: 1, color: '#764ba2' }
            ]),
            borderRadius: [0, 4, 4, 0]
          }
        }
      ]
    })
  }
}

const handleResize = () => {
  visitChart?.resize()
  rankChart?.resize()
}

onMounted(() => {
  nextTick(() => {
    initCharts()
  })
  window.addEventListener('resize', handleResize)
})
</script>

<style scoped lang="scss">
.dashboard-index {
  .stat-card {
    display: flex;
    align-items: center;
    gap: 20px;
    border-radius: 8px;

    .stat-icon {
      width: 60px;
      height: 60px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;

      &.icon-blue { background: linear-gradient(135deg, #667eea, #764ba2); }
      &.icon-green { background: linear-gradient(135deg, #43e97b, #38f9d7); }
      &.icon-orange { background: linear-gradient(135deg, #fa709a, #fee140); }
      &.icon-red { background: linear-gradient(135deg, #f093fb, #f5576c); }
    }

    .stat-info {
      .stat-value {
        font-size: 28px;
        font-weight: 700;
        margin-bottom: 4px;
      }

      .stat-label {
        color: #909399;
        font-size: 14px;
      }
    }
  }

  .chart-card {
    border-radius: 8px;

    .card-title {
      font-size: 18px;
      margin-bottom: 20px;
    }

    .chart-container {
      height: 300px;

      &.rank-chart {
        height: 300px;
      }
    }
  }

  .list-card {
    border-radius: 8px;

    .card-title {
      font-size: 18px;
      margin: 0;
    }
  }
}
</style>
