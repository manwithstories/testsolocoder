<template>
  <MainLayout>
    <div class="page-container" v-loading="loading">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>

      <div v-if="job" class="job-detail">
        <el-card class="job-header-card">
          <div class="job-header">
            <div>
              <h1 class="job-title">{{ job.title }}</h1>
              <div class="job-meta">
                <span class="salary">{{ formatSalary(job) }}</span>
                <span><el-icon><Location /></el-icon> {{ job.location }}</span>
                <span v-if="job.job_type"><el-icon><Clock /></el-icon> {{ formatJobType(job.job_type) }}</span>
                <span v-if="job.experience"><el-icon><User /></el-icon> {{ job.experience }}</span>
              </div>
            </div>
            <div>
              <el-button
                v-if="isApplicant"
                type="primary"
                size="large"
                :loading="applying"
                @click="handleApply"
              >
                立即投递
              </el-button>
              <el-button
                v-if="isApplicant && hasApplied"
                type="success"
                size="large"
                disabled
              >
                已投递
              </el-button>
            </div>
          </div>
        </el-card>

        <el-row :gutter="20">
          <el-col :span="16">
            <el-card class="mb-20">
              <template #header>
                <span class="card-title">职位描述</span>
              </template>
              <div class="job-description" v-html="formatDescription(job.description)"></div>
            </el-card>

            <el-card v-if="job.requirements" class="mb-20">
              <template #header>
                <span class="card-title">任职要求</span>
              </template>
              <div class="job-requirements" v-html="formatDescription(job.requirements)"></div>
            </el-card>

            <el-card v-if="job.benefits" class="mb-20">
              <template #header>
                <span class="card-title">福利待遇</span>
              </template>
              <div class="job-benefits">{{ job.benefits }}</div>
            </el-card>

            <el-card v-if="job.skills">
              <template #header>
                <span class="card-title">技能要求</span>
              </template>
              <div class="skill-tags">
                <el-tag v-for="skill in parseSkills(job.skills)" :key="skill" size="large">
                  {{ skill }}
                </el-tag>
              </div>
            </el-card>
          </el-col>

          <el-col :span="8">
            <el-card v-if="job.company" class="mb-20">
              <template #header>
                <span class="card-title">公司信息</span>
              </template>
              <div class="company-info">
                <h3 class="company-name">{{ job.company.name }}</h3>
                <div class="company-meta" v-if="job.company.industry">
                  <el-icon><OfficeBuilding /></el-icon>
                  {{ job.company.industry }}
                </div>
                <div class="company-meta" v-if="job.company.size">
                  <el-icon><UserFilled /></el-icon>
                  {{ job.company.size }}
                </div>
                <div class="company-meta" v-if="job.company.address">
                  <el-icon><Location /></el-icon>
                  {{ job.company.address }}
                </div>
                <p class="company-desc" v-if="job.company.description">{{ job.company.description }}</p>
              </div>
            </el-card>

            <el-card>
              <template #header>
                <span class="card-title">职位信息</span>
              </template>
              <div class="job-stats">
                <div class="stat-item">
                  <span class="stat-label">浏览量</span>
                  <span class="stat-value">{{ job.view_count }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">投递数</span>
                  <span class="stat-value">{{ job.apply_count }}</span>
                </div>
                <div class="stat-item" v-if="job.deadline">
                  <span class="stat-label">截止日期</span>
                  <span class="stat-value">{{ formatDate(job.deadline) }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>

    <el-dialog v-model="applyDialogVisible" title="投递简历" width="500px">
      <el-form :model="applyForm" label-width="100px">
        <el-form-item label="选择简历">
          <el-select v-model="applyForm.resume_id" placeholder="请选择简历" style="width: 100%">
            <el-option
              v-for="resume in resumes"
              :key="resume.id"
              :label="resume.title"
              :value="resume.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="求职信">
          <el-input
            v-model="applyForm.cover_letter"
            type="textarea"
            :rows="4"
            placeholder="请输入求职信（选填）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="applying" @click="confirmApply">确认投递</el-button>
      </template>
    </el-dialog>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { getJob } from '@/api/job'
import { listResumes } from '@/api/resume'
import { apply } from '@/api/application'
import { useUserStore } from '@/stores/user'
import type { Job, Resume } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const applying = ref(false)
const job = ref<Job | null>(null)
const resumes = ref<Resume[]>([])
const applyDialogVisible = ref(false)
const hasApplied = ref(false)

const isApplicant = computed(() => userStore.hasRole('applicant'))

const applyForm = reactive({
  resume_id: null as number | null,
  cover_letter: ''
})

async function fetchJob() {
  const id = Number(route.params.id)
  if (!id) return

  loading.value = true
  try {
    const res = await getJob(id)
    if (res.data) {
      job.value = res.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function fetchResumes() {
  try {
    const res = await listResumes()
    if (res.data) {
      resumes.value = res.data
      if (resumes.value.length > 0) {
        const defaultResume = resumes.value.find(r => r.is_default)
        applyForm.resume_id = defaultResume?.id || resumes.value[0].id
      }
    }
  } catch (e) {
    console.error(e)
  }
}

function handleApply() {
  if (!userStore.isLoggedIn) {
    ElMessageBox.confirm('请先登录后再投递简历', '提示', {
      confirmButtonText: '去登录',
      cancelButtonText: '取消',
      type: 'info'
    }).then(() => {
      router.push('/login')
    }).catch(() => {})
    return
  }

  if (resumes.value.length === 0) {
    ElMessageBox.confirm('您还没有创建简历，请先创建简历', '提示', {
      confirmButtonText: '去创建',
      cancelButtonText: '取消',
      type: 'info'
    }).then(() => {
      router.push('/resumes/create')
    }).catch(() => {})
    return
  }

  applyDialogVisible.value = true
}

async function confirmApply() {
  if (!applyForm.resume_id || !job.value) {
    ElMessage.warning('请选择简历')
    return
  }

  applying.value = true
  try {
    await apply({
      job_id: job.value.id,
      resume_id: applyForm.resume_id,
      cover_letter: applyForm.cover_letter
    })
    ElMessage.success('投递成功')
    applyDialogVisible.value = false
    hasApplied.value = true
  } catch (e) {
    // error handled in interceptor
  } finally {
    applying.value = false
  }
}

function goBack() {
  router.back()
}

function formatSalary(j: Job) {
  if (j.salary_min && j.salary_max) {
    return `${(j.salary_min / 1000).toFixed(0)}K-${(j.salary_max / 1000).toFixed(0)}K`
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

function formatDescription(text: string) {
  return text.replace(/\n/g, '<br>')
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
}

onMounted(() => {
  fetchJob()
  if (userStore.isLoggedIn && userStore.hasRole('applicant')) {
    fetchResumes()
  }
})
</script>

<style scoped>
.job-detail {
  margin-top: 20px;
}

.job-header-card {
  margin-bottom: 20px;
}

.job-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.job-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 10px 0;
}

.job-meta {
  display: flex;
  gap: 20px;
  color: #606266;
  font-size: 16px;
}

.job-meta .salary {
  color: #f56c6c;
  font-weight: 600;
  font-size: 20px;
}

.job-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-title {
  font-weight: 600;
}

.job-description,
.job-requirements,
.job-benefits {
  line-height: 1.8;
  color: #606266;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.company-info {
  text-align: center;
}

.company-name {
  font-size: 20px;
  margin: 0 0 10px 0;
}

.company-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: #909399;
  margin-bottom: 8px;
}

.company-desc {
  text-align: left;
  color: #606266;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}

.job-stats {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
}

.stat-label {
  color: #909399;
}

.stat-value {
  font-weight: 600;
  color: #303133;
}
</style>
