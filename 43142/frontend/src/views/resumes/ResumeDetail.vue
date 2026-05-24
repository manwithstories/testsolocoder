<template>
  <MainLayout>
    <div class="page-container" v-loading="loading">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>

      <div v-if="resume" class="resume-detail">
        <el-card class="resume-header-card">
          <div class="resume-header">
            <div>
              <h1 class="resume-title">{{ resume.title }}</h1>
              <div class="resume-meta">
                <span><el-icon><User /></el-icon> {{ resume.full_name }}</span>
                <span v-if="resume.email"><el-icon><Message /></el-icon> {{ resume.email }}</span>
                <span v-if="resume.phone"><el-icon><Phone /></el-icon> {{ resume.phone }}</span>
                <span v-if="resume.location"><el-icon><Location /></el-icon> {{ resume.location }}</span>
              </div>
            </div>
            <div>
              <el-tag v-if="resume.is_default" type="success" size="large">默认简历</el-tag>
            </div>
          </div>
        </el-card>

        <el-row :gutter="20">
          <el-col :span="16">
            <el-card v-if="resume.summary" class="mb-20">
              <template #header>
                <span class="card-title">个人简介</span>
              </template>
              <div class="resume-content">{{ resume.summary }}</div>
            </el-card>

            <el-card v-if="resume.experience" class="mb-20">
              <template #header>
                <span class="card-title">工作经历</span>
              </template>
              <div class="resume-content" v-html="formatContent(resume.experience)"></div>
            </el-card>

            <el-card v-if="resume.education" class="mb-20">
              <template #header>
                <span class="card-title">教育经历</span>
              </template>
              <div class="resume-content" v-html="formatContent(resume.education)"></div>
            </el-card>

            <el-card v-if="resume.projects" class="mb-20">
              <template #header>
                <span class="card-title">项目经历</span>
              </template>
              <div class="resume-content" v-html="formatContent(resume.projects)"></div>
            </el-card>

            <el-card v-if="resume.skills">
              <template #header>
                <span class="card-title">技能</span>
              </template>
              <div class="skill-tags">
                <el-tag v-for="skill in parseSkills(resume.skills)" :key="skill" size="large">
                  {{ skill }}
                </el-tag>
              </div>
            </el-card>
          </el-col>

          <el-col :span="8">
            <el-card>
              <template #header>
                <span class="card-title">操作</span>
              </template>
              <div class="action-buttons">
                <el-button type="primary" @click="goEdit">编辑简历</el-button>
                <el-button
                  v-if="resume.file_url"
                  @click="downloadFile(resume.file_url)"
                >
                  <el-icon><Download /></el-icon>
                  下载附件
                </el-button>
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
import { getResume } from '@/api/resume'
import type { Resume } from '@/types'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const resume = ref<Resume | null>(null)

async function fetchResume() {
  const id = Number(route.params.id)
  if (!id) return

  loading.value = true
  try {
    const res = await getResume(id)
    if (res.data) {
      resume.value = res.data
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

function goEdit() {
  router.push(`/resumes/${resume.value?.id}/edit`)
}

function downloadFile(url: string) {
  window.open(url, '_blank')
}

function parseSkills(skills: string) {
  return skills.split(',').map(s => s.trim()).filter(Boolean)
}

function formatContent(text: string) {
  return text.replace(/\n/g, '<br>')
}

onMounted(() => {
  fetchResume()
})
</script>

<style scoped>
.resume-header-card {
  margin-bottom: 20px;
}

.resume-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.resume-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 10px 0;
}

.resume-meta {
  display: flex;
  gap: 20px;
  color: #606266;
  font-size: 14px;
}

.resume-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-title {
  font-weight: 600;
}

.resume-content {
  line-height: 1.8;
  color: #606266;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>
