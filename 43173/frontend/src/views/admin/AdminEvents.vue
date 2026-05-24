<template>
  <div class="admin-events">
    <div class="page-header">
      <h2>演出管理</h2>
    </div>
    
    <el-table :data="events" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="演出名称" />
      <el-table-column prop="venue" label="场地" width="150" />
      <el-table-column prop="city" label="城市" width="100" />
      <el-table-column prop="start_time" label="开始时间" width="180" />
      <el-table-column label="票务" width="120">
        <template #default="{ row }">
          {{ row.sold_count || 0 }} / {{ row.total_tickets || 0 }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button 
            type="primary" 
            text 
            size="small"
            @click="viewDetail(row)"
          >
            详情
          </el-button>
          <el-button 
            type="danger" 
            text 
            size="small"
            @click="deleteEvent(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadEvents"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { eventApi } from '@/api/event'

const loading = ref(false)
const events = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

onMounted(() => {
  loadEvents()
})

async function loadEvents() {
  loading.value = true
  try {
    const res = await eventApi.getAll({
      page: page.value,
      page_size: pageSize.value
    })
    events.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function viewDetail(event: any) {
  window.open(`/events/${event.id}`, '_blank')
}

async function deleteEvent(event: any) {
  try {
    await ElMessageBox.confirm('确定要删除该演出吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await eventApi.delete(event.id)
    ElMessage.success('删除成功')
    loadEvents()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function getStatusTag(status: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const tags: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    0: 'info',
    1: 'success',
    2: 'danger',
    3: 'info'
  }
  return tags[status] || 'info'
}

function getStatusText(status: number) {
  const texts: Record<number, string> = {
    0: '未开始',
    1: '售票中',
    2: '已结束',
    3: '已取消'
  }
  return texts[status] || '未知'
}
</script>

<style scoped lang="scss">
.admin-events {
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
