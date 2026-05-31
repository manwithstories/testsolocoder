<template>
  <AppLayout>
    <div class="page">
      <h2>我的档期</h2>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column label="服务">
          <template #default="{ row }">{{ row.service?.name }}</template>
        </el-table-column>
        <el-table-column label="客户">
          <template #default="{ row }">{{ row.customer?.real_name || row.customer?.username }}</template>
        </el-table-column>
        <el-table-column prop="address" label="地址" />
        <el-table-column prop="start_at" label="开始时间" width="170" />
        <el-table-column prop="end_at" label="结束时间" width="170" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag>{{ { pending: '待确认', confirmed: '已确认', canceled: '已取消' }[row.status] || row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '../../components/AppLayout.vue'
import { listBookings } from '../../api/booking'

const list = ref<any[]>([])
onMounted(async () => {
  const res = await listBookings({})
  list.value = (res.data as any).data || []
})
</script>
