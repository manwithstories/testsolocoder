<template>
  <MainLayout>
    <div class="page-container" v-loading="loading">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>

      <div v-if="application" class="application-detail">
        <el-card class="status-card">
          <div class="status-info">
            <div>
              <h2>{{ application.job?.title }}</h2>
              <p class="company-name">{{ application.job?.company?.name }}</p>
            </div>
            <el-tag :type="getStatusType(application.status)" size="large">
              {{ getStatusText(application.status) }}
            </el-tag>
          </div>
        </el-card>

        <el-row :gutter="20">
          <el-col :span="16">
            <el-card class="mb-20">
              <template #header>
                <span class="card-title">求职信</span>
              </template>
              <div class="cover-letter" v-if="application.cover_letter">
                {{ application.cover_letter }}
              </div>
              <el-empty v-else description="暂无求职信" :image-size="60" />
            </el-card>

            <el-card v-if="application.resume" class="mb-20">
              <template #header>
                <span class="card-title">简历信息</span>
              </template>
              <div class="resume-preview">
                <h3>{{ application.resume.full_name }}</h3>
                <p v-if="application.resume.email">邮箱: {{ application.resume.email }}</p>
                <p v-if="application.resume.phone">电话: {{ application.resume.phone }}</p>
                <p v-if="application.resume.summary">简介: {{ application.resume.summary }}</p>
                <el-button v-if="application.resume.file_url" type="primary" @click="downloadFile(application.resume.file_url!)">
                  下载简历附件
                </el-button>
              </div>
            </el-card>

            <el-card>
              <template #header>
                <span class="card-title">状态变更历史</span>
              </template>
              <el-timeline>
                <el-timeline-item
                  v-for="(item, index) in history"
                  :key="index"
                  :timestamp="formatDate(item.created_at)"
                  placement="top"
                >
                  <div class="history-item">
                    <el-tag :type="getStatusType(item.new_status)" size="small">
                      {{ getStatusText(item.new_status) }}
                    </el-tag>
                    <span v-if="item.change_reason" class="change-reason">
                      {{ item.change_reason }}
                    </span>
                  </div>
                </el-timeline-item>
              </el-timeline>
              <el-empty v-if="history.length === 0" description="暂无历史记录" :image-size="60" />
            </el-card>
          </el-col>

          <el-col :span="8">
            <el-card>
              <template #header>
                <span class="card-title">投递信息</span>
              </template>
              <div class="info-list">
                <div class="info-item">
                  <span class="label">投递时间</span>
                  <span class="value">{{ formatDate(application.applied_at) }}</span>
                </div>
                <div class="info-item">
                  <span class="label">最后更新</span>
                  <span class="value">{{ formatDate(application.last_update_at) }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { getApplication, getApplicationHistory } from '@/api/application'
import type { Application, ApplicationHistory } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const application = ref<Application | null>(null)
const history = ref<ApplicationHistory[]>([])

async function fetchApplication() {
  const id = Number(route.params.id)
  if (!id) return

  loading.value = true
  try {
    const [appRes, historyRes] = await Promise.all([
      getApplication(id),
      getApplicationHistory(id)
    ])
    if (appRes.data) {
      application.value = appRes.data
    }
    if (historyRes.data) {
      history.value = historyRes.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.back()
}

function downloadFile(url: string) {
  window.open(url, '_blank')
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
  fetchApplication()
})
</script>

<style scoped>
.application-detail {
  margin-top: 20px;
}

.status-card {
  margin-bottom: 20px;
}

.status-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-info h2 {
  margin: 0 0 5px 0;
}

.company-name {
  color: #909399;
  margin: 0;
}

.card-title {
  font-weight: 600;
}

.cover-letter {
  white-space: pre-wrap;
  line-height: 1.8;
  color: #606266;
}

.resume-preview {
  line-height: 1.8;
}

.resume-preview h3 {
  margin: 0 0 10px 0;
}

.resume-preview p {
  margin: 5px 0;
  color: #606266;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.change-reason {
  color: #909399;
  font-size: 12px;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
}

.label {
  color: #909399;
}

.value {
  color: #303133;
}
</style>
