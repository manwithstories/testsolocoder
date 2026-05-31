<template>
  <AppLayout>
    <div class="page">
      <h2>工单</h2>
      <div class="row" style="margin-bottom:16px">
        <el-select v-model="status" placeholder="全部状态" clearable style="width:160px" @change="load">
          <el-option value="open" label="待处理" />
          <el-option value="assigned" label="已分配" />
          <el-option value="pending" label="处理中" />
          <el-option value="resolved" label="已解决" />
          <el-option value="closed" label="已关闭" />
          <el-option value="escalated" label="已升级" />
        </el-select>
        <el-button v-if="userStore.role==='customer'" type="primary" @click="$router.push('/tickets/new')">新建工单</el-button>
      </div>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type==='refund' ? 'danger' : 'warning'">{{ row.type==='refund' ? '退款' : '投诉' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status==='escalated' ? 'danger' : ''">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
      </el-table>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '../components/AppLayout.vue'
import { listTickets } from '../api/ticket'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const list = ref<any[]>([])
const status = ref('')

async function load() {
  const params: Record<string, string> = {}
  if (status.value) params.status = status.value
  const res = await listTickets(params)
  list.value = (res.data as any).data || []
}

function statusLabel(s: string) {
  return { open: '待处理', assigned: '已分配', pending: '处理中', resolved: '已解决', closed: '已关闭', escalated: '已升级' }[s] || s
}

onMounted(load)
</script>
