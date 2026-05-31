<template>
  <AppLayout>
    <div class="page">
      <h2>运营数据统计</h2>
      <el-row :gutter="16" style="margin-bottom:16px">
        <el-col :span="6"><el-card shadow="hover">
          <div class="muted">总订单数</div><h2>{{ overview.orders || 0 }}</h2>
        </el-card></el-col>
        <el-col :span="6"><el-card shadow="hover">
          <div class="muted">总营收</div><h2>¥{{ overview.revenue?.toFixed(2) || '0.00' }}</h2>
        </el-card></el-col>
        <el-col :span="6"><el-card shadow="hover">
          <div class="muted">客户数</div><h2>{{ overview.customers || 0 }}</h2>
        </el-card></el-col>
        <el-col :span="6"><el-card shadow="hover">
          <div class="muted">家政人员</div><h2>{{ overview.staff || 0 }}</h2>
        </el-card></el-col>
      </el-row>

      <el-card style="margin-bottom:16px">
        <template #header>收入趋势(近30天)</template>
        <v-chart :option="revenueOpt" style="height:340px" autoresize />
      </el-card>
      <el-row :gutter="16">
        <el-col :span="12"><el-card>
          <template #header>订单按分类</template>
          <v-chart :option="categoryOpt" style="height:300px" autoresize />
        </el-card></el-col>
        <el-col :span="12"><el-card>
          <template #header>家政人员绩效 TOP20</template>
          <el-table :data="staffPerf" size="small">
            <el-table-column prop="real_name" label="人员" />
            <el-table-column prop="order_count" label="订单" width="80" />
            <el-table-column prop="revenue" label="营收" width="120" />
            <el-table-column prop="rating" label="评分" width="100" />
          </el-table>
        </el-card></el-col>
      </el-row>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, markRaw } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import AppLayout from '../../components/AppLayout.vue'
import { statsOverview, statsRevenue, statsCategory, statsStaff } from '../../api/stats'

use([
  CanvasRenderer,
  markRaw(BarChart),
  markRaw(LineChart),
  markRaw(TitleComponent),
  markRaw(TooltipComponent),
  markRaw(GridComponent),
  markRaw(LegendComponent),
])

const overview = ref<any>({})
const revenue = ref<any[]>([])
const category = ref<any[]>([])
const staffPerf = ref<any[]>([])

const revenueOpt = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: revenue.value.map(r => r.date) },
  yAxis: { type: 'value' },
  series: [{ type: 'line', data: revenue.value.map(r => r.revenue), smooth: true, areaStyle: {} }],
}))
const categoryOpt = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: category.value.map(r => r.category) },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: category.value.map(r => r.count) }],
}))

onMounted(async () => {
  const [o, r, c, s] = await Promise.all([statsOverview(), statsRevenue(), statsCategory(), statsStaff()])
  overview.value = (o.data as any).data || {}
  revenue.value = (r.data as any).data || []
  category.value = (c.data as any).data || []
  staffPerf.value = (s.data as any).data || []
})
</script>
