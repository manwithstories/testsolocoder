<template>
  <el-card>
    <template #header>数据统计</template>
    <el-row :gutter="16">
      <el-col :span="6"><el-statistic title="总交易数" :value="overview.total_trades" /></el-col>
      <el-col :span="6"><el-statistic title="已完成" :value="overview.completed" /></el-col>
      <el-col :span="6"><el-statistic title="总成交金额" :value="overview.total_amount" :precision="2" /></el-col>
      <el-col :span="6"><el-statistic title="用户数" :value="overview.users" /></el-col>
    </el-row>
    <el-divider />
    <h3>热门品牌</h3>
    <el-table :data="brands">
      <el-table-column prop="brand" label="品牌" />
      <el-table-column prop="count" label="成交数" />
      <el-table-column prop="total" label="成交总额" />
    </el-table>
    <el-button type="primary" style="margin-top: 16px" @click="exportExcel">导出Excel</el-button>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const overview = ref<any>({})
const brands = ref<any[]>([])

onMounted(async () => {
  overview.value = await request.get('/stats/overview') || {}
  brands.value = await request.get('/stats/brands') || []
})
function exportExcel() {
  window.location.href = '/api/stats/export?token=' + localStorage.getItem('token')
}
</script>
