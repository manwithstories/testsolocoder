<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">我的面试</div>
        <div class="page-subtitle">查看您的所有面试安排</div>
      </div>

      <el-card>
        <el-table v-loading="loading" :data="interviews" style="width: 100%">
          <el-table-column label="职位名称" prop="application.job.title" min-width="200">
            <template #default="{ row }">
              {{ row.application?.job?.title }}
            </template>
          </el-table-column>
          <el-table-column label="公司" prop="application.job.company.name" width="150">
            <template #default="{ row }">
              {{ row.application?.job?.company?.name }}
            </template>
          </el-table-column>
          <el-table-column label="面试时间" width="180">
            <template #default="{ row }">
              {{ formatDateTime(row.scheduled_at) }}
            </template>
          </el-table-column>
          <el-table-column label="时长" width="100">
            <template #default="{ row }">
              {{ row.duration }}分钟
            </template>
          </el-table-column>
          <el-table-column label="面试方式" width="120">
            <template #default="{ row }">
              {{ getInterviewType(row.interview_type) }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button
                v-if="row.status === 'pending'"
                type="success"
                link
                @click="handleAccept(row.id)"
              >
                接受
              </el-button>
              <el-button
                v-if="row.status === 'pending'"
                type="danger"
                link
                @click="handleReject(row.id)"
              >
                拒绝
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchInterviews"
            @size-change="fetchInterviews"
          />
        </div>
      </el-card>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { listMyInterviews, acceptInterview, rejectInterview } from '@/api/interview'
import type { Interview } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const interviews = ref<Interview[]>([])
const total = ref(0)

const pagination = reactive({
  page: 1,
  pageSize: 10
})

async function fetchInterviews() {
  loading.value = true
  try {
    const res = await listMyInterviews({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.data) {
      interviews.value = res.data.items
      total.value = res.data.pagination.total
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleAccept(id: number) {
  ElMessageBox.confirm('确定要接受该面试吗？', '提示', {
    confirmButtonText: '接受',
    cancelButtonText: '取消',
    type: 'success'
  }).then(async () => {
    try {
      await acceptInterview(id)
      ElMessage.success('已接受面试')
      fetchInterviews()
    } catch (e) {
      // error handled
    }
  }).catch(() => {})
}

async function handleReject(id: number) {
  ElMessageBox.confirm('确定要拒绝该面试吗？', '提示', {
    confirmButtonText: '拒绝',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await rejectInterview(id)
      ElMessage.success('已拒绝面试')
      fetchInterviews()
    } catch (e) {
      // error handled
    }
  }).catch(() => {})
}

function getInterviewType(type: string) {
  const types: Record<string, string> = {
    'onsite': '现场面试',
    'online': '远程面试',
    'phone': '电话面试'
  }
  return types[type] || type || '待定'
}

function getStatusType(status: string) {
  const types: Record<string, string> = {
    'pending': 'warning',
    'accepted': 'success',
    'rejected': 'danger',
    'completed': 'success',
    'cancelled': 'info'
  }
  return types[status] || ''
}

function getStatusText(status: string) {
  const texts: Record<string, string> = {
    'pending': '待确认',
    'accepted': '已接受',
    'rejected': '已拒绝',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return texts[status] || status
}

function formatDateTime(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchInterviews()
})
</script>

<style scoped>
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
