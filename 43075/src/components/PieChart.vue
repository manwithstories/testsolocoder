<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import * as echarts from 'echarts'
import { useRecordsStore } from '@/stores/records'
import { useSettingsStore } from '@/stores/settings'
import { getStartOfDay, getEndOfDay } from '@/utils/date'

interface Props {
  range: 'week' | 'month' | 'all'
}

const props = defineProps<Props>()

const chartRef = ref<HTMLDivElement | null>(null)
let chartInstance: echarts.ECharts | null = null

const recordsStore = useRecordsStore()
const settingsStore = useSettingsStore()

const dateRange = computed(() => {
  const now = Date.now()
  if (props.range === 'week') {
    return {
      start: getStartOfDay(now - 6 * 24 * 60 * 60 * 1000),
      end: getEndOfDay(now)
    }
  } else if (props.range === 'month') {
    return {
      start: getStartOfDay(now - 29 * 24 * 60 * 60 * 1000),
      end: getEndOfDay(now)
    }
  }
  return {
    start: 0,
    end: now
  }
})

const chartData = computed(() => {
  return recordsStore.getCategoryStats(dateRange.value.start, dateRange.value.end)
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

  const data = chartData.value

  if (data.length === 0) {
    chartInstance.setOption({
      graphic: {
        type: 'text',
        left: 'center',
        top: 'middle',
        style: {
          text: '暂无数据',
          fill: isDark ? '#6b7280' : '#9ca3af',
          fontSize: 14
        }
      }
    })
    return
  }

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'item',
      backgroundColor: isDark ? '#1f2937' : '#ffffff',
      borderColor: isDark ? '#374151' : '#e5e7eb',
      textStyle: {
        color: textColor
      },
      formatter: (params: any) => {
        const minutes = Math.floor(params.value / 60)
        const seconds = params.value % 60
        return `${params.name}<br/>时长: ${minutes}分${seconds}秒<br/>占比: ${params.percent}%`
      }
    },
    legend: {
      orient: 'vertical',
      right: '5%',
      top: 'center',
      textStyle: {
        color: textColor,
        fontSize: 12
      },
      itemWidth: 12,
      itemHeight: 12
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: isDark ? '#1f2937' : '#ffffff',
          borderWidth: 2
        },
        label: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold',
            color: textColor
          }
        },
        data: data.map(d => ({
          value: d.value,
          name: d.name,
          itemStyle: { color: d.color }
        }))
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
  <div ref="chartRef" class="w-full h-64 lg:h-80"></div>
</template>
