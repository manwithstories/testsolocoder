<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">收到的简历</div>
        <div class="page-subtitle">管理所有投递到您职位的简历</div>
      </div>

      <el-card class="search-card mb-20">
        <el-form :inline="true" :model="searchForm" @submit.prevent="handleSearch">
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部状态" clearable style="width: 150px">
              <el-option label="待处理" value="pending" />
              <el-option label="已查看" value="viewed" />
              <el-option label="感兴趣" value="interested" />
              <el-option label="面试中" value="interview" />
              <el-option label="已录用" value="accepted" />
              <el-option label="未通过" value="rejected" />
            </el-select>
          </el-form-item>
          <el-form-item label="关键词">
            <el-input v-model="searchForm.keyword" placeholder="搜索求职者" clearable style="width: 200px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card>
        <el-table
          v-loading="loading"
          :data="applications"
          @selection-change="handleSelectionChange"
          style="width: 100%"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="求职者" min-width="150">
            <template #default="{ row }">
              <div class="applicant-info">
                <el-avatar :size="40">{{ getApplicantName(row).charAt(0) }}</el-avatar>
                <div>
                  <div class="name">{{ getApplicantName(row) }}</div>
                  <div class="email">{{ row.applicant?.email }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="应聘职位" prop="job.title" min-width="150" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="投递时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.applied_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="goToDetail(row.id)">查看</el-button>
              <el-dropdown trigger="click" @command="(cmd: string) => handleStatusChange(row.id, cmd)">
                <el-button type="primary" link>
                  更新状态
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="viewed">标记已查看</el-dropdown-item>
                    <el-dropdown-item command="interested">标记感兴趣</el-dropdown-item>
                    <el-dropdown-item command="rejected">标记不合适</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <el-button type="success" link @click="scheduleInterview(row)">安排面试</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="action-bar">
          <el-button
            type="primary"
            :disabled="selectedApplications.length === 0"
            @click="handleBulkStatusChange('interested')"
          >
            批量标记感兴趣
          </el-button>
          <el-button
            type="danger"
            :disabled="selectedApplications.length === 0"
            @click="handleBulkStatusChange('rejected')"
          >
            批量标记不合适
          </el-button>
        </div>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchApplications"
            @size-change="fetchApplications"
          />
        </div>
      </el-card>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { listCompanyApplications, updateApplicationStatus, bulkUpdateStatus } from '@/api/application'
import type { Application } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const applications = ref<Application[]>([])
const total = ref(0)
const selectedApplications = ref<Application[]>([])

const searchForm = reactive({
  status: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10
})

async function fetchApplications() {
  loading.value = true
  try {
    const res = await listCompanyApplications({
      page: pagination.page,
      page_size: pagination.pageSize,
      status: searchForm.status || undefined,
      keyword: searchForm.keyword || undefined
    })
    if (res.data) {
      applications.value = res.data.items
      total.value = res.data.pagination.total
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleSelectionChange(selection: Application[]) {
  selectedApplications.value = selection
}

function handleSearch() {
  pagination.page = 1
  fetchApplications()
}

function goToDetail(id: number) {
  router.push(`/applications/${id}`)
}

async function handleStatusChange(id: number, status: string) {
  try {
    await updateApplicationStatus(id, status)
    ElMessage.success('状态更新成功')
    fetchApplications()
  } catch (e) {
    // error handled
  }
}

async function handleBulkStatusChange(status: string) {
  if (selectedApplications.value.length === 0) return

  ElMessageBox.confirm(`确定要将选中的 ${selectedApplications.value.length} 条投递标记为"${getStatusText(status)}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await bulkUpdateStatus(selectedApplications.value.map(a => a.id), status)
      ElMessage.success('批量更新成功')
      fetchApplications()
    } catch (e) {
      // error handled
    }
  }).catch(() => {})
}

function scheduleInterview(application: Application) {
  router.push({
    path: '/company/interviews/schedule',
    query: { applicationId: application.id }
  })
}

function getApplicantName(row: Application) {
  return row.applicant?.profile?.full_name || row.applicant?.email || '未知'
}

function getStatusType(status: string) {
  const types: Record<string, string> = {
    'pending': 'warning',
    'viewed': 'info',
    'interested': 'primary',
    'interview': 'success',
    'accepted': 'success',
    'rejected': 'danger',
    'withdrawn': 'info'
  }
  return types[status] || ''
}

function getStatusText(status: string) {
  const texts: Record<string, string> = {
    'pending': '待处理',
    'viewed': '已查看',
    'interested': '感兴趣',
    'interview': '面试中',
    'accepted': '已录用',
    'rejected': '未通过',
    'withdrawn': '已撤回'
  }
  return texts[status] || status
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchApplications()
})
</script>

<style scoped>
.applicant-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.applicant-info .name {
  font-weight: 500;
}

.applicant-info .email {
  color: #909399;
  font-size: 12px;
}

.action-bar {
  margin-top: 20px;
  display: flex;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
