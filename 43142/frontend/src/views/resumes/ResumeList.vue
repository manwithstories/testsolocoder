<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">我的简历</div>
        <div class="page-subtitle">管理您的所有简历</div>
      </div>

      <div class="action-bar mb-20">
        <el-button type="primary" @click="goCreate">
          <el-icon><Plus /></el-icon>
          创建简历
        </el-button>
      </div>

      <div v-loading="loading" class="resume-list">
        <div v-for="resume in resumes" :key="resume.id" class="resume-card">
          <div class="resume-header">
            <h3 class="resume-title">{{ resume.title }}</h3>
            <el-tag v-if="resume.is_default" type="success" size="small">默认</el-tag>
          </div>
          <div class="resume-info">
            <span><el-icon><User /></el-icon> {{ resume.full_name }}</span>
            <span v-if="resume.email"><el-icon><Message /></el-icon> {{ resume.email }}</span>
            <span v-if="resume.phone"><el-icon><Phone /></el-icon> {{ resume.phone }}</span>
          </div>
          <div class="resume-skills" v-if="resume.skills">
            <el-tag v-for="skill in parseSkills(resume.skills)" :key="skill" size="small" type="info">
              {{ skill }}
            </el-tag>
          </div>
          <div class="resume-footer">
            <span class="resume-time">更新于 {{ formatDate(resume.updated_at) }}</span>
            <div class="resume-actions">
              <el-button type="primary" link @click="goDetail(resume.id)">查看</el-button>
              <el-button type="primary" link @click="goEdit(resume.id)">编辑</el-button>
              <el-button
                v-if="!resume.is_default"
                type="success"
                link
                @click="setDefault(resume.id)"
              >
                设为默认
              </el-button>
              <el-button
                v-if="resume.file_url"
                type="primary"
                link
                @click="downloadFile(resume.file_url)"
              >
                下载附件
              </el-button>
              <el-button type="danger" link @click="handleDelete(resume.id)">删除</el-button>
            </div>
          </div>
        </div>

        <el-empty v-if="!loading && resumes.length === 0" description="暂无简历，请点击创建简历" />
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { listResumes, deleteResume, setDefaultResume } from '@/api/resume'
import type { Resume } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const resumes = ref<Resume[]>([])

async function fetchResumes() {
  loading.value = true
  try {
    const res = await listResumes()
    if (res.data) {
      resumes.value = res.data
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function goCreate() {
  router.push('/resumes/create')
}

function goDetail(id: number) {
  router.push(`/resumes/${id}`)
}

function goEdit(id: number) {
  router.push(`/resumes/${id}/edit`)
}

async function setDefault(id: number) {
  try {
    await setDefaultResume(id)
    ElMessage.success('设置成功')
    fetchResumes()
  } catch (e) {
    // error handled
  }
}

async function handleDelete(id: number) {
  ElMessageBox.confirm('确定要删除该简历吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteResume(id)
      ElMessage.success('删除成功')
      fetchResumes()
    } catch (e) {
      // error handled
    }
  }).catch(() => {})
}

function downloadFile(url: string) {
  window.open(url, '_blank')
}

function parseSkills(skills: string) {
  return skills.split(',').map(s => s.trim()).filter(Boolean)
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD')
}

onMounted(() => {
  fetchResumes()
})
</script>

<style scoped>
.action-bar {
  display: flex;
  gap: 10px;
}

.resume-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 16px;
}

.resume-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  border: 1px solid #ebeef5;
}

.resume-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.resume-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.resume-info {
  display: flex;
  gap: 16px;
  color: #606266;
  font-size: 14px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.resume-info span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.resume-skills {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.resume-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}

.resume-time {
  color: #c0c4cc;
  font-size: 12px;
}

.resume-actions {
  display: flex;
  gap: 8px;
}
</style>
