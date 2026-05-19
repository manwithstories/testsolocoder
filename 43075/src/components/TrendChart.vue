<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import * as echarts from 'echarts'
import { useRecordsStore } from '@/stores/records'
import { useSettingsStore } from '@/stores/settings'

interface Props {
  type: 'week' | 'month'
}

const props = defineProps<Props>()

const chartRef = ref<HTMLDivElement | null>(null)
let chartInstance: echarts.ECharts | null = null

const recordsStore = useRecordsStore()
const settingsStore = useSettingsStore()

const chartData = computed(() => {
  if (props.type === 'week') {
    return recordsStore.getWeeklyTrendData()
  }
  return recordsStore.getMonthlyTrendData()
})

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

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: isDark ? '#1f2937' : '#ffffff',
      borderColor: isDark ? '#374151' : '#e5e7eb',
      textStyle: {
        color: textColor
      },
      formatter: (params: any) => {
        const data = params[0]
        const minutes = Math.floor(data.value / 60)
        const seconds = data.value % 60
        return `${data.name}<br/>专注时长: ${minutes}分${seconds}秒`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: chartData.value.map(d => d.date),
      axisLine: {
        lineStyle: { color: gridColor }
      },
      axisLabel: {
        color: textColor,
        fontSize: 12
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        lineStyle: { color: gridColor }
      },
      axisLabel: {
        color: textColor,
        fontSize: 12,
        formatter: (value: number) => {
          return Math.floor(value / 60) + 'm'
        }
      },
      splitLine: {
        lineStyle: { color: gridColor, type: 'dashed' }
      }
    },
    series: [
      {
        data: chartData.value.map(d => d.duration),
        type: 'bar',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#6366f1' },
            { offset: 1, color: '#8b5cf6' }
          ]),
          borderRadius: [6, 6, 0, 0]
        },
        barWidth: props.type === 'week' ? '40%' : '60%'
      }
    ]
  }

  chartInstance.setOption(option)
}

const handleResize = () => {
  chartInstance?.resize()
}

watch(
  () => chartData.value,
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
  <div ref="chartRef" class="w-full h-64 lg:h-80"></div>
</template>
