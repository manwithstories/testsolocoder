<template>
  <AppLayout>
    <div class="page">
      <h2>公司数据统计</h2>
      <el-row :gutter="16">
        <el-col :span="6"><el-card>
          <div class="muted">订单总数</div><h1>{{ d.orders || 0 }}</h1>
        </el-card></el-col>
        <el-col :span="6"><el-card>
          <div class="muted">累计收入</div><h1>¥{{ d.revenue?.toFixed(2) || '0.00' }}</h1>
        </el-card></el-col>
        <el-col :span="6"><el-card>
          <div class="muted">待处理预约</div><h1>{{ d.pending || 0 }}</h1>
        </el-card></el-col>
        <el-col :span="6"><el-card>
          <div class="muted">时间</div><h3 style="margin:0">{{ d.date }}</h3>
        </el-card></el-col>
      </el-row>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '../../components/AppLayout.vue'
import { companyDashboard } from '../../api/stats'

const d = ref<any>({})
onMounted(async () => {
  const res = await companyDashboard()
  d.value = (res.data as any).data || {}
})
</script>
