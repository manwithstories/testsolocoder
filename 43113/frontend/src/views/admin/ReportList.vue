<template>
  <div class="report-list-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>举报处理</span>
          <el-select v-model="statusFilter" placeholder="选择状态" clearable style="width: 150px" @change="fetchReports">
            <el-option label="待处理" value="pending" />
            <el-option label="已处理" value="resolved" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </div>
      </template>

      <el-table :data="reports" style="width: 100%">
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getTargetTypeName(row.targetType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="目标ID" prop="targetId" width="100" />
        <el-table-column label="举报人" width="150">
          <template #default="{ row }">
            {{ row.reporter?.nickname || row.reporter?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="举报原因" prop="reason" width="120" />
        <el-table-column label="详细描述" prop="description" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="举报时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              size="small"
              @click="handleReport(row.id, 'resolved', 'delete')"
            >
              删除内容
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="warning"
              size="small"
              @click="handleReport(row.id, 'resolved', 'keep')"
            >
              保留内容
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="info"
              size="small"
              @click="handleReport(row.id, 'rejected', '')"
            >
              驳回
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchReports"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { auditApi } from '@/api'
import type { Report } from '@/types'
import dayjs from 'dayjs'

const reports = ref<Report[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const statusFilter = ref('')

const fetchReports = async () => {
  try {
    const res = await auditApi.getReportList({
      page: page.value,
      pageSize: pageSize.value,
      status: statusFilter.value || undefined
    })
    reports.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const handleReport = async (id: number, status: string, result: string) => {
  try {
    await auditApi.handleReport(id, status, result)
    fetchReports()
  } catch (e) {
    console.error(e)
  }
}

const getTargetTypeName = (type: string) => {
  const map: Record<string, string> = {
    question: '问题',
    answer: '回答',
    comment: '评论'
  }
  return map[type] || type
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    resolved: 'success',
    rejected: 'info'
  }
  return map[status] || 'info'
}

const getStatusName = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    resolved: '已处理',
    rejected: '已驳回'
  }
  return map[status] || status
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchReports()
})
</script>

<style scoped lang="scss">
.report-list-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
