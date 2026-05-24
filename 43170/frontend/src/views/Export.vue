<template>
  <div class="page-container">
    <h1 class="page-title">数据导出</h1>

    <div class="card export-section">
      <h3>租赁记录导出</h3>
      <p class="export-desc">导出您所有设备的租赁记录</p>
      <div class="export-buttons">
        <el-button type="primary" :loading="exportingRental" @click="exportRentals('csv')">
          <el-icon><Download /></el-icon>
          导出 CSV
        </el-button>
        <el-button type="success" :loading="exportingRental" @click="exportRentals('pdf')">
          <el-icon><Document /></el-icon>
          导出 PDF
        </el-button>
      </div>
    </div>

    <div class="card export-section">
      <h3>收益报表导出</h3>
      <p class="export-desc">选择日期范围导出您的收益报表</p>
      <el-form label-width="100px" style="max-width: 500px; margin-bottom: 20px">
        <el-form-item label="开始日期">
          <el-date-picker
            v-model="startDate"
            type="date"
            placeholder="选择开始日期"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker
            v-model="endDate"
            type="date"
            placeholder="选择结束日期"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
      </el-form>
      <div class="export-buttons">
        <el-button type="primary" :loading="exportingRevenue" @click="exportRevenue('csv')">
          <el-icon><Download /></el-icon>
          导出 CSV
        </el-button>
        <el-button type="success" :loading="exportingRevenue" @click="exportRevenue('pdf')">
          <el-icon><Document /></el-icon>
          导出 PDF
        </el-button>
      </div>
    </div>

    <div class="card statistics-section">
      <h3>收益统计</h3>
      <div class="statistics-card">
        <div class="statistic-item">
          <div class="statistic-value">¥{{ statistics.totalRevenue.toFixed(2) }}</div>
          <div class="statistic-label">总收益</div>
        </div>
        <div class="statistic-item">
          <div class="statistic-value">{{ statistics.totalOrders }}</div>
          <div class="statistic-label">总订单数</div>
        </div>
        <div class="statistic-item">
          <div class="statistic-value">{{ statistics.completedOrders }}</div>
          <div class="statistic-label">已完成订单</div>
        </div>
        <div class="statistic-item">
          <div class="statistic-value">¥{{ statistics.avgOrderValue.toFixed(2) }}</div>
          <div class="statistic-label">平均订单金额</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { exportApi } from '@/api/order'
import { orderApi } from '@/api/order'
import type { Order } from '@/types'

const exportingRental = ref(false)
const exportingRevenue = ref(false)
const startDate = ref('')
const endDate = ref('')

const statistics = reactive({
  totalRevenue: 0,
  totalOrders: 0,
  completedOrders: 0,
  avgOrderValue: 0
})

onMounted(() => {
  loadStatistics()
})

async function loadStatistics() {
  try {
    const response = await orderApi.getMyOrders()
    const orders = response.data as Order[]

    statistics.totalOrders = orders.length
    statistics.completedOrders = orders.filter(o => o.status === 'completed').length

    const completedOrders = orders.filter(o => o.status === 'completed')
    statistics.totalRevenue = completedOrders.reduce((sum, o) => sum + o.totalRent, 0)
    statistics.avgOrderValue = completedOrders.length > 0
      ? statistics.totalRevenue / completedOrders.length
      : 0
  } catch (error) {
    console.error('Failed to load statistics:', error)
  }
}

async function exportRentals(format: string) {
  exportingRental.value = true
  try {
    const response = await exportApi.exportRentalRecords(format)
    downloadFile(response, `rental_records.${format}`, format)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('Export failed:', error)
    ElMessage.error('导出失败')
  } finally {
    exportingRental.value = false
  }
}

async function exportRevenue(format: string) {
  if (!startDate.value || !endDate.value) {
    ElMessage.warning('请选择日期范围')
    return
  }

  exportingRevenue.value = true
  try {
    const response = await exportApi.exportRevenueReport(format, startDate.value, endDate.value)
    downloadFile(response, `revenue_report.${format}`, format)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('Export failed:', error)
    ElMessage.error('导出失败')
  } finally {
    exportingRevenue.value = false
  }
}

function downloadFile(blob: Blob, filename: string, _format: string) {
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}
</script>

<style scoped>
.export-section {
  margin-bottom: 30px;
}

.export-section h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #303133;
}

.export-desc {
  font-size: 14px;
  color: #909399;
  margin-bottom: 20px;
}

.export-buttons {
  display: flex;
  gap: 12px;
}

.statistics-section h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
  color: #303133;
}

.statistics-card {
  display: flex;
  justify-content: space-around;
  padding: 20px 0;
}

.statistic-item {
  text-align: center;
}

.statistic-value {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
}

.statistic-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}
</style>
