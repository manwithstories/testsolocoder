<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">职位列表</div>
        <div class="page-subtitle">浏览所有招聘职位</div>
      </div>

      <el-card class="search-card mb-20">
        <el-form :inline="true" :model="searchForm" @submit.prevent="handleSearch">
          <el-form-item label="关键词">
            <el-input v-model="searchForm.keyword" placeholder="搜索职位名称、描述、技能" clearable style="width: 300px" />
          </el-form-item>
          <el-form-item label="工作地点">
            <el-input v-model="searchForm.location" placeholder="工作地点" clearable style="width: 200px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <div v-loading="loading" class="job-list">
        <div v-for="job in jobs" :key="job.id" class="job-card" @click="goToDetail(job.id)">
          <div class="job-header">
            <h3 class="job-title">{{ job.title }}</h3>
            <span class="job-salary">
              {{ formatSalary(job) }}
            </span>
          </div>
          <div class="job-meta">
            <span><el-icon><Location /></el-icon> {{ job.location }}</span>
            <span v-if="job.job_type"><el-icon><Clock /></el-icon> {{ formatJobType(job.job_type) }}</span>
            <span v-if="job.experience"><el-icon><User /></el-icon> {{ job.experience }}</span>
          </div>
          <div class="job-skills" v-if="job.skills">
            <el-tag v-for="skill in parseSkills(job.skills)" :key="skill" size="small" type="info">
              {{ skill }}
            </el-tag>
          </div>
          <div class="job-footer">
            <span class="company-name" v-if="job.company">{{ job.company.name }}</span>
            <span class="job-time">{{ formatDate(job.created_at) }}</span>
          </div>
        </div>

        <el-empty v-if="!loading && jobs.length === 0" description="暂无职位" />
      </div>

      <div class="pagination-container" v-if="total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="fetchJobs"
          @size-change="fetchJobs"
        />
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { listJobs } from '@/api/job'
import type { Job } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const jobs = ref<Job[]>([])
const total = ref(0)

const searchForm = reactive({
  keyword: '',
  location: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10
})

async function fetchJobs() {
  loading.value = true
  try {
    const res = await listJobs({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      status: 'open'
    })
    if (res.data) {
      jobs.value = res.data.items
      total.value = res.data.pagination.total
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchJobs()
}

function resetSearch() {
  searchForm.keyword = ''
  searchForm.location = ''
  pagination.page = 1
  fetchJobs()
}

function goToDetail(id: number) {
  router.push(`/jobs/${id}`)
}

function formatSalary(job: Job) {
  if (job.salary_min && job.salary_max) {
    return `${(job.salary_min / 1000).toFixed(0)}K-${(job.salary_max / 1000).toFixed(0)}K`
  }
  return '面议'
}

function formatJobType(type: string) {
  const types: Record<string, string> = {
    'full-time': '全职',
    'part-time': '兼职',
    'contract': '合同',
    'intern': '实习',
    'remote': '远程'
  }
  return types[type] || type
}

function parseSkills(skills: string) {
  return skills.split(',').map(s => s.trim()).filter(Boolean)
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
}

onMounted(() => {
  fetchJobs()
})
</script>

<style scoped>
.search-card {
  margin-bottom: 20px;
}

.job-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 16px;
}

.job-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #ebeef5;
}

.job-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.job-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.job-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.job-salary {
  color: #f56c6c;
  font-weight: 600;
  font-size: 16px;
}

.job-meta {
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 14px;
  margin-bottom: 12px;
}

.job-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.job-skills {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.job-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}

.company-name {
  color: #606266;
  font-size: 14px;
}

.job-time {
  color: #c0c4cc;
  font-size: 12px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
