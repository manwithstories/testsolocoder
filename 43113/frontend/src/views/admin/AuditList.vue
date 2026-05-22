<template>
  <div class="audit-list-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>内容审核</span>
          <div class="filters">
            <el-select v-model="targetType" placeholder="选择类型" clearable style="width: 150px" @change="fetchAudits">
              <el-option label="问题" value="question" />
              <el-option label="回答" value="answer" />
              <el-option label="评论" value="comment" />
            </el-select>
            <el-select v-model="status" placeholder="选择状态" clearable style="width: 150px" @change="fetchAudits">
              <el-option label="待审核" value="pending" />
              <el-option label="已通过" value="approved" />
              <el-option label="已拒绝" value="rejected" />
            </el-select>
          </div>
        </div>
      </template>

      <el-alert
        v-if="pendingCount"
        :title="`待审核: 问题 ${pendingCount.questions} | 回答 ${pendingCount.answers} | 评论 ${pendingCount.comments}`"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
      />

      <el-table :data="audits" style="width: 100%">
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getTargetTypeName(row.targetType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="目标ID" prop="targetId" width="100" />
        <el-table-column label="操作" prop="action" width="100" />
        <el-table-column label="原因" prop="reason" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作人" width="150">
          <template #default="{ row }">
            {{ row.admin?.nickname || row.admin?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchAudits"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { auditApi } from '@/api'
import type { AuditRecord } from '@/types'
import dayjs from 'dayjs'

const audits = ref<AuditRecord[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const targetType = ref('')
const status = ref('')
const pendingCount = ref<{ questions: number; answers: number; comments: number } | null>(null)

const fetchAudits = async () => {
  try {
    const res = await auditApi.getAuditList({
      page: page.value,
      pageSize: pageSize.value,
      targetType: targetType.value || undefined,
      status: status.value || undefined
    })
    audits.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const fetchPendingCount = async () => {
  try {
    const res = await auditApi.getPendingAuditCount()
    pendingCount.value = res.data || null
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
    approved: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

const getStatusName = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return map[status] || status
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchAudits()
  fetchPendingCount()
})
</script>

<style scoped lang="scss">
.audit-list-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .filters {
      display: flex;
      gap: 12px;
    }
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
