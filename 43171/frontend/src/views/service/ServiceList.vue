<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>航拍服务</span>
        <el-button type="primary" @click="$router.push('/service/create')" v-if="role === 'client'">
          <el-icon><Plus /></el-icon>
          发布需求
        </el-button>
      </div>
    </template>

    <el-table :data="services" v-loading="loading">
      <el-table-column prop="title" label="需求标题" min-width="150" />
      <el-table-column prop="region" label="区域" width="100" />
      <el-table-column label="预算" width="150">
        <template #default="{ row }">¥{{ row.budget_min }} ~ ¥{{ row.budget_max }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="发布时间" width="180" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="$router.push(`/service/${row.id}`)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 20px; justify-content: flex-end"
      layout="prev, pager, next, total"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="currentPage"
      @current-change="fetchServices"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const userStore = useUserStore()
const role = computed(() => userStore.role)

const loading = ref(false)
const services = ref<AerialService[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

onMounted(() => {
  fetchServices()
})

async function fetchServices() {
  loading.value = true
  try {
    const res: any = await request.get('/services', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    services.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function statusText(status: string) {
  const map: Record<string, string> = {
    open: '招募中', assigned: '已分配', progress: '进行中',
    completed: '已完成', cancelled: '已取消'
  }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = {
    open: 'success', assigned: 'warning', progress: 'primary',
    completed: 'info', cancelled: 'info'
  }
  return map[status] || ''
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
