<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import { useRecordsStore } from '@/stores/records'
import { useSettingsStore } from '@/stores/settings'

const chartRef = ref<HTMLDivElement | null>(null)
let chartInstance: echarts.ECharts | null = null

const recordsStore = useRecordsStore()
const settingsStore = useSettingsStore()

const hours = Array.from({ length: 24 }, (_, i) => `${i.toString().padStart(2, '0')}:00`)
const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

const initChart = () => {
  if (!chartRef.value) return

  chartInstance = echarts.init(chartRef.value)
  updateChart()
}

const updateChart = () => {
  if (!chartInstance) return

  const isDark = settingsStore.settings.darkMode
  const textColor = isDark ? '#e5e7eb' : '#374151'
  const gridColor = isDark ? '#374151' : '#e5e7eb'

  const heatmapData = recordsStore.getHeatmapData()
  const data = heatmapData.map(item => [item.hour, item.day, item.value])
  const maxValue = Math.max(...heatmapData.map(d => d.value), 1)

  const option: echarts.EChartsOption = {
    tooltip: {
      position: 'top',
      backgroundColor: isDark ? '#1f2937' : '#ffffff',
      borderColor: isDark ? '#374151' : '#e5e7eb',
      textStyle: {
        color: textColor
      },
      formatter: (params: any) => {
        const [hour, day, value] = params.value
        const minutes = Math.floor(value / 60)
        const seconds = value % 60
        return `${day} ${hour}:00-${hour + 1}:00<br/>专注时长: ${minutes}分${seconds}秒`
      }
    },
    grid: {
      left: '10%',
      right: '5%',
      bottom: '15%',
      top: '5%'
    },
    xAxis: {
      type: 'category',
      data: hours,
      splitArea: {
        show: true
      },
      axisLine: {
        lineStyle: { color: gridColor }
      },
      axisLabel: {
        color: textColor,
        fontSize: 10,
        interval: 2
      }
    },
    yAxis: {
      type: 'category',
      data: days,
      splitArea: {
        show: true
      },
      axisLine: {
        lineStyle: { color: gridColor }
      },
      axisLabel: {
        color: textColor,
        fontSize: 11
      }
    },
    visualMap: {
      min: 0,
      max: maxValue,
      calculable: false,
      orient: 'horizontal',
      left: 'center',
      bottom: '0%',
      inRange: {
        color: ['#e0e7ff', '#a5b4fc', '#818cf8', '#6366f1', '#4f46e5', '#4338ca']
      },
      textStyle: {
        color: textColor
      },
      show: false
    },
    series: [
      {
        name: '专注时长',
        type: 'heatmap',
        data: data,
        label: {
          show: false
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }

  chartInstance.setOption(option)
}

const handleResize = () => {
  chartInstance?.resize()
}

watch(
  () => recordsStore.records,
  () => updateChart(),
  { deep: true }
)

watch(
  () => settingsStore.settings.darkMode,
  () => updateChart()
)

onMounted(() => {
  initChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
})
</script>

<template>
  <div ref="chartRef" class="w-full h-64 lg:h-72"></div>
</template>
