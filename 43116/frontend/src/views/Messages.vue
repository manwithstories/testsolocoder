<template>
  <div class="container" style="padding-top: 20px;">
    <div class="page-header">
      <h1 class="page-title">消息中心</h1>
      <el-button type="primary" :icon="Check" @click="markAllRead">全部已读</el-button>
    </div>

    <el-table :data="messages" v-loading="loading" style="width: 100%">
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)">{{ getTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="content" label="内容" min-width="300" show-overflow-tooltip />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag v-if="!row.is_read" type="danger" size="small">未读</el-tag>
          <el-tag v-else type="info" size="small">已读</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="!row.is_read"
            type="primary"
            link
            @click="handleMarkRead(row.id)"
          >
            标为已读
          </el-button>
          <el-button type="danger" link @click="handleDelete(row.id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadMessages"
      />
    </div>

    <el-empty v-if="messages.length === 0 && !loading" description="暂无消息" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { messageApi } from '@/api'
import type { Message } from '@/types'

const messages = ref<Message[]>([])
const loading = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadMessages()
})

const loadMessages = async () => {
  loading.value = true
  try {
    const res = await messageApi.getMessages({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    messages.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const handleMarkRead = async (id: number) => {
  try {
    await messageApi.markAsRead(id)
    ElMessage.success('已标为已读')
    loadMessages()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const markAllRead = async () => {
  try {
    await messageApi.markAllAsRead()
    ElMessage.success('全部标为已读')
    loadMessages()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该消息吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await messageApi.deleteMessage(id)
    ElMessage.success('删除成功')
    loadMessages()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getTypeTagType = (type: string) => {
  const map: Record<string, string> = {
    booking: 'primary',
    pickup: 'success',
    return: 'warning',
    review: 'info',
    maintenance: 'danger',
    system: ''
  }
  return map[type] || ''
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    booking: '预订',
    pickup: '取车',
    return: '还车',
    review: '评价',
    maintenance: '维护',
    system: '系统'
  }
  return map[type] || type
}
</script>
