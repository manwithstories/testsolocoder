<template>
  <div class="my-stats">
    <div class="page-header">
      <h2>数据统计</h2>
      <el-select v-model="period" style="width: 120px;">
        <el-option label="近7天" value="7" />
        <el-option label="近30天" value="30" />
        <el-option label="近90天" value="90" />
      </el-select>
    </div>
    
    <el-row :gutter="20" class="stats-cards">
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <div class="stat-icon" style="background: rgba(64, 158, 255, 0.1);">
              <el-icon :size="24" color="#409eff"><Headset /></el-icon>
            </div>
            <div class="stat-info">
              <div class="value">{{ formatCount(stats.monthly_summary?.total_plays || 0) }}</div>
              <div class="label">播放量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <div class="stat-icon" style="background: rgba(103, 194, 58, 0.1);">
              <el-icon :size="24" color="#67c23a"><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="value">{{ formatCount(stats.monthly_summary?.new_followers || 0) }}</div>
              <div class="label">新增粉丝</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <div class="stat-icon" style="background: rgba(230, 162, 60, 0.1);">
              <el-icon :size="24" color="#e6a23c"><Star /></el-icon>
            </div>
            <div class="stat-info">
              <div class="value">{{ formatCount(stats.monthly_summary?.total_likes || 0) }}</div>
              <div class="label">收藏数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <div class="stat-icon" style="background: rgba(245, 108, 108, 0.1);">
              <el-icon :size="24" color="#f56c6c"><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="value">¥{{ formatAmount(stats.monthly_summary?.total_revenue || 0) }}</div>
              <div class="label">收益</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <div class="chart-section">
      <h3>播放趋势</h3>
      <div ref="chartRef" class="chart-container" v-loading="loading"></div>
    </div>
    
    <div class="chart-section">
      <h3>粉丝增长</h3>
      <div ref="chartRef2" class="chart-container" v-loading="loading"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import { revenueApi } from '@/api/revenue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)
const period = ref('30')
const stats = ref<any>({})
const chartRef = ref<HTMLElement>()
const chartRef2 = ref<HTMLElement>()
let chart1: echarts.ECharts | null = null
let chart2: echarts.ECharts | null = null

onMounted(() => {
  loadStats()
})

watch(period, () => {
  loadStats()
})

async function loadStats() {
  if (!userStore.user?.artist_info) return
  
  loading.value = true
  try {
    stats.value = await revenueApi.getArtistStats(userStore.user.artist_info.id)
    await nextTick()
    renderCharts()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function renderCharts() {
  const dailyData = stats.value.daily_data || []
  const dates = dailyData.map((d: any) => d.date)
  const plays = dailyData.map((d: any) => d.new_plays)
  const followers = dailyData.map((d: any) => d.new_followers)
  
  if (chartRef.value) {
    if (chart1) chart1.dispose()
    chart1 = echarts.init(chartRef.value)
    chart1.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: dates },
      yAxis: { type: 'value' },
      series: [{
        name: '播放量',
        type: 'line',
        smooth: true,
        areaStyle: {},
        data: plays,
        itemStyle: { color: '#409eff' }
      }]
    })
  }
  
  if (chartRef2.value) {
    if (chart2) chart2.dispose()
    chart2 = echarts.init(chartRef2.value)
    chart2.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: dates },
      yAxis: { type: 'value' },
      series: [{
        name: '新增粉丝',
        type: 'bar',
        data: followers,
        itemStyle: { color: '#67c23a' }
      }]
    })
  }
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}

function formatAmount(amount: number) {
  return amount.toFixed(2)
}
</script>

<style scoped lang="scss">
.my-stats {
  .stats-cards {
    margin-bottom: 24px;
    
    .stat-card {
      display: flex;
      align-items: center;
      gap: 16px;
      
      .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      
      .stat-info {
        .value {
          font-size: 20px;
          font-weight: 600;
        }
        
        .label {
          font-size: 13px;
          color: var(--text-light);
        }
      }
    }
  }
  
  .chart-section {
    margin-bottom: 24px;
    
    h3 {
      margin: 0 0 16px 0;
    }
    
    .chart-container {
      height: 300px;
      background: #fff;
      border-radius: 8px;
    }
  }
}
</style>
