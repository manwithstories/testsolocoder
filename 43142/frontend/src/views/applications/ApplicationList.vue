<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">我的投递</div>
        <div class="page-subtitle">查看您的所有投递记录</div>
      </div>

      <el-card>
        <el-table v-loading="loading" :data="applications" style="width: 100%">
          <el-table-column label="职位名称" prop="job.title" min-width="200">
            <template #default="{ row }">
              <el-link type="primary" @click="goToJob(row.job_id)">{{ row.job?.title }}</el-link>
            </template>
          </el-table-column>
          <el-table-column label="公司" prop="job.company.name" width="150">
            <template #default="{ row }">
              {{ row.job?.company?.name }}
            </template>
          </el-table-column>
          <el-table-column label="工作地点" prop="job.location" width="120">
            <template #default="{ row }">
              {{ row.job?.location }}
            </template>
          </el-table-column>
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
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="goToDetail(row.id)">查看详情</el-button>
              <el-button
                v-if="row.status === 'pending' || row.status === 'viewed'"
                type="danger"
                link
                @click="handleWithdraw(row.id)"
              >
                撤回
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
import { listMyApplications, withdrawApplication } from '@/api/application'
import type { Application } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const applications = ref<Application[]>([])
const total = ref(0)

const pagination = reactive({
  page: 1,
  pageSize: 10
})

async function fetchApplications() {
  loading.value = true
  try {
    const res = await listMyApplications({
      page: pagination.page,
      page_size: pagination.pageSize
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

function goToJob(id: number) {
  router.push(`/jobs/${id}`)
}

function goToDetail(id: number) {
  router.push(`/applications/${id}`)
}

async function handleWithdraw(id: number) {
  ElMessageBox.confirm('确定要撤回该投递吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await withdrawApplication(id)
      ElMessage.success('撤回成功')
      fetchApplications()
    } catch (e) {
      // error handled
    }
  }).catch(() => {})
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
    'interested': 'HR感兴趣',
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
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
