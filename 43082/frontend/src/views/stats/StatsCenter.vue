<template>
  <div class="stats-center">
    <el-card>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="会员统计" name="members">
          <div class="stats-header">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="loadMemberStats"
            />
          </div>
          <el-row :gutter="20" style="margin-bottom: 20px">
            <el-col :span="12">
              <el-card>
                <div style="text-align: center">
              <div style="font-size: 32px; font-weight: 700; color: #409eff">
                {{ memberStats?.new_members || 0 }}
              </div>
              <div style="color: #909399; margin-top: 8px">新增会员</div>
            </div>
            </el-card>
            </el-col>
            <el-col :span="12">
              <el-card>
                <div style="text-align: center">
              <div style="font-size: 32px; font-weight: 700; color: #67c23a">
                {{ memberStats?.active_members || 0 }}
              </div>
              <div style="color: #909399; margin-top: 8px">活跃会员</div>
            </div>
            </el-card>
            </el-col>
          </el-row>
          <el-card>
            <div ref="memberChart" style="height: 300px" />
          </el-card>
        </el-tab-pane>
        
        <el-tab-pane label="课程统计" name="courses">
          <div class="stats-header">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="loadCourseStats"
            />
          </div>
          <el-row :gutter="20" style="margin-bottom: 20px">
            <el-col :span="12">
              <el-card>
                <div style="text-align: center">
              <div style="font-size: 32px; font-weight: 700; color: #e6a23c">
                {{ courseStats?.total_bookings || 0 }}
              </div>
              <div style="color: #909399; margin-top: 8px">总预约数</div>
            </div>
            </el-card>
            </el-col>
            <el-col :span="12">
              <el-card>
                <div style="text-align: center">
              <div style="font-size: 32px; font-weight: 700; color: #f56c6c">
                {{ courseStats?.total_check_ins || 0 }}
              </div>
              <div style="color: #909399; margin-top: 8px">总签到数</div>
            </div>
            </el-card>
            </el-col>
          </el-row>
          <el-card>
            <div ref="courseChart" style="height: 300px" />
          </el-card>
        </el-tab-pane>
        
        <el-tab-pane label="教练统计" name="coaches">
          <div class="stats-header">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="loadCoachStats"
            />
          </div>
          <el-card>
            <el-table :data="coachStats?.coach_stats || []" style="width: 100%">
              <el-table-column prop="coach_name" label="教练姓名" />
              <el-table-column prop="total_courses" label="课程数" />
              <el-table-column prop="total_students" label="学员数" />
              <el-table-column prop="total_hours" label="课时(h)" />
              <el-table-column label="操作">
                <template #default="{ row }">
                  <el-progress :percentage="Math.min(100, row.total_hours * 10)" />
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { statsApi } from '@/api/booking'
import * as echarts from 'echarts'
import dayjs from 'dayjs'

const activeTab = ref('members')
const dateRange = ref<string[]>([])
const memberStats = ref<any>(null)
const courseStats = ref<any>(null)
const coachStats = ref<any>(null)

const memberChart = ref<HTMLElement>()
const courseChart = ref<HTMLElement>()

let memberChartInstance: echarts.ECharts | null = null
let courseChartInstance: echarts.ECharts | null = null

const initMemberChart = () => {
  if (!memberChart.value) return
  
  memberChartInstance = echarts.init(memberChart.value)
  memberChartInstance.setOption({
    title: { text: '每日新增会员', left: 'center' },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: [] },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: [], itemStyle: { color: '#409eff' } }]
  })
}

const initCourseChart = () => {
  if (!courseChart.value) return
  
  courseChartInstance = echarts.init(courseChart.value)
  courseChartInstance.setOption({
    title: { text: '每日预约趋势', left: 'center' },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: [] },
    yAxis: { type: 'value' },
    series: [
      { name: '预约数', type: 'line', data: [], smooth: true, itemStyle: { color: '#67c23a' } },
      { name: '签到数', type: 'line', data: [], smooth: true, itemStyle: { color: '#f56c6c' } }
    ],
    legend: { data: ['预约数', '签到数'], bottom: 0 }
  })
}

const loadMemberStats = async () => {
  try {
    const params: any = {}
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const res = await statsApi.getMemberStats(params)
    memberStats.value = res.data
    
    await nextTick()
    if (memberChartInstance && res.data.daily_new) {
      memberChartInstance.setOption({
        xAxis: { data: res.data.daily_new.map((d: any) => d.date) },
        series: [{ data: res.data.daily_new.map((d: any) => d.count) }]
      })
    }
  } catch (error) {
    console.error(error)
  }
}

const loadCourseStats = async () => {
  try {
    const params: any = {}
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const res = await statsApi.getCourseStats(params)
    courseStats.value = res.data
    
    await nextTick()
    if (courseChartInstance && res.data.daily_bookings) {
      const bookingData: Record<string, number> = {}
      res.data.daily_bookings.forEach((d: any) => {
        bookingData[d.date] = d.count
      })
      
      const dates = Array.from(new Set([
        ...res.data.daily_bookings.map((d: any) => d.date)
      ])).sort()
      
      courseChartInstance.setOption({
        xAxis: { data: dates },
        series: [
          { data: dates.map(d => bookingData[d] || 0) }
        ]
      })
    }
  } catch (error) {
    console.error(error)
  }
}

const loadCoachStats = async () => {
  try {
    const params: any = {}
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const res = await statsApi.getCoachStats(params)
    coachStats.value = res.data
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  nextTick(() => {
    initMemberChart()
    initCourseChart()
    loadMemberStats()
    loadCourseStats()
    loadCoachStats()
  })
})
</script>

<style scoped>
.stats-header {
  margin-bottom: 20px;
}
</style>
