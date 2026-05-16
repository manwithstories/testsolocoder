<script setup lang="ts">
import { computed } from 'vue'
import { useStore } from '../stores/useStore'
import { formatMinutes } from '../utils/storage'

const store = useStore()

const colors = ['#e74c3c', '#3498db', '#f39c12', '#9b59b6', '#1abc9c', '#e67e22', '#34495e']

const chartData = computed(() => {
  if (store.taskStats.value.length === 0) {
    return []
  }
  
  const total = store.taskStats.value.reduce((sum, item) => sum + item.minutes, 0)
  
  return store.taskStats.value.map((item, index) => ({
    ...item,
    percentage: total > 0 ? (item.minutes / total) * 100 : 0,
    color: colors[index % colors.length]
  }))
})

const totalMinutes = computed(() => {
  return store.taskStats.value.reduce((sum, item) => sum + item.minutes, 0)
})

const maxMinutes = computed(() => {
  if (store.taskStats.value.length === 0) return 1
  return Math.max(...store.taskStats.value.map(s => s.minutes), 1)
})

const pieSegments = computed(() => {
  const total = chartData.value.reduce((sum, item) => sum + item.minutes, 0)
  if (total === 0) return []
  
  const circumference = 2 * Math.PI * 40
  let offset = 0
  
  return chartData.value.map((item) => {
    const dasharray = (item.minutes / total) * circumference
    const segment = {
      color: item.color,
      dasharray: `${dasharray} ${circumference}`,
      offset: -offset
    }
    offset += dasharray
    return segment
  })
})
</script>

<template>
  <div class="statistics">
    <h3>今日统计</h3>
    
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon">🍅</div>
        <div class="stat-info">
          <span class="stat-value">{{ store.todayStats.value.completedPomodoros }}</span>
          <span class="stat-label">完成番茄</span>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">⏱️</div>
        <div class="stat-info">
          <span class="stat-value">{{ formatMinutes(store.todayStats.value.totalWorkMinutes) }}</span>
          <span class="stat-label">专注时长</span>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">📋</div>
        <div class="stat-info">
          <span class="stat-value">{{ store.taskStats.value.length }}</span>
          <span class="stat-label">进行中任务</span>
        </div>
      </div>
    </div>

    <div class="chart-section">
      <h4>任务时间分布</h4>
      
      <div v-if="chartData.length > 0" class="chart-container">
        <div class="bar-chart">
          <div 
            v-for="item in chartData" 
            :key="item.task.id"
            class="bar-item"
          >
            <div class="bar-label">
              <span class="bar-title">{{ item.task.title }}</span>
              <span class="bar-time">{{ formatMinutes(item.minutes) }}</span>
            </div>
            <div class="bar-track">
              <div 
                class="bar-fill"
                :style="{ 
                  width: (item.minutes / maxMinutes) * 100 + '%',
                  backgroundColor: item.color 
                }"
              ></div>
            </div>
          </div>
        </div>
        
        <div class="pie-chart-container">
          <svg class="pie-chart" viewBox="0 0 100 100">
            <circle 
              v-for="(item, index) in pieSegments" 
              :key="index"
              cx="50"
              cy="50"
              r="40"
              fill="transparent"
              :stroke="item.color"
              stroke-width="20"
              :stroke-dasharray="item.dasharray"
              :stroke-dashoffset="item.offset"
              transform="rotate(-90 50 50)"
            />
          </svg>
          <div class="pie-center">
            <span class="pie-total">{{ formatMinutes(totalMinutes) }}</span>
            <span class="pie-label">总计</span>
          </div>
        </div>
      </div>
      
      <div v-else class="empty-chart">
        <p>今天还没有专注记录</p>
        <p class="subtext">开始一个番茄钟，数据会显示在这里</p>
      </div>
    </div>

    <div v-if="chartData.length > 0" class="legend">
      <div 
        v-for="item in chartData" 
        :key="item.task.id"
        class="legend-item"
      >
        <span class="legend-color" :style="{ backgroundColor: item.color }"></span>
        <span class="legend-text">{{ item.task.title }}</span>
        <span class="legend-percent">{{ item.percentage.toFixed(1) }}%</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.statistics {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.statistics h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  font-size: 28px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.chart-section {
  flex: 1;
}

.chart-section h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 600;
  color: #666;
}

.chart-container {
  display: grid;
  grid-template-columns: 1fr 180px;
  gap: 24px;
}

.bar-chart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.bar-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.bar-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.bar-title {
  font-size: 13px;
  color: #333;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bar-time {
  font-size: 12px;
  color: #999;
}

.bar-track {
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s ease;
}

.pie-chart-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pie-chart {
  width: 160px;
  height: 160px;
}

.pie-center {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.pie-total {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.pie-label {
  font-size: 11px;
  color: #999;
}

.empty-chart {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.empty-chart p {
  margin: 0;
}

.empty-chart .subtext {
  font-size: 12px;
  margin-top: 4px;
}

.legend {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.legend-text {
  font-size: 12px;
  color: #666;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.legend-percent {
  font-size: 12px;
  color: #999;
}
</style>
