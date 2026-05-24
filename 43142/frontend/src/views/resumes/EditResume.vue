<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">编辑简历</div>
      </div>

      <el-card v-loading="loading">
        <el-form v-if="resume" :model="form" :rules="rules" ref="formRef" label-width="100px" class="form-container">
          <el-form-item label="简历标题" prop="title">
            <el-input v-model="form.title" placeholder="请输入简历标题" />
          </el-form-item>
          <el-form-item label="姓名" prop="full_name">
            <el-input v-model="form.full_name" placeholder="请输入姓名" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="form.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="form.phone" placeholder="请输入手机号" />
          </el-form-item>
          <el-form-item label="所在地">
            <el-input v-model="form.location" placeholder="请输入所在地" />
          </el-form-item>
          <el-form-item label="教育经历">
            <el-input
              v-model="form.education"
              type="textarea"
              :rows="3"
              placeholder="请输入教育经历"
            />
          </el-form-item>
          <el-form-item label="工作经历">
            <el-input
              v-model="form.experience"
              type="textarea"
              :rows="5"
              placeholder="请输入工作经历"
            />
          </el-form-item>
          <el-form-item label="技能">
            <el-input
              v-model="form.skills"
              type="textarea"
              :rows="3"
              placeholder="请输入技能"
            />
          </el-form-item>
          <el-form-item label="个人简介">
            <el-input
              v-model="form.summary"
              type="textarea"
              :rows="3"
              placeholder="请输入个人简介"
            />
          </el-form-item>
          <el-form-item label="项目经历">
            <el-input
              v-model="form.projects"
              type="textarea"
              :rows="3"
              placeholder="请输入项目经历"
            />
          </el-form-item>
          <el-form-item label="上传附件">
            <el-upload
              :action="''"
              :auto-upload="false"
              :show-file-list="false"
              :on-change="handleFileChange"
              accept=".pdf,.doc,.docx"
            >
              <el-button>
                <el-icon><Upload /></el-icon>
                选择文件
              </el-button>
              <template #tip>
                <div class="upload-tip">支持 PDF、DOC、DOCX 格式，最大 10MB</div>
              </template>
            </el-upload>
            <span v-if="resume.file_url" class="file-name">
              <el-icon><Document /></el-icon>
              已上传附件
            </span>
          </el-form-item>
          <el-form-item label="设为默认">
            <el-switch v-model="form.is_default" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="saving" @click="handleSubmit">
              保存
            </el-button>
            <el-button @click="goBack">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadFile } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { getResume, updateResume, uploadResumeFile } from '@/api/resume'
import type { Resume } from '@/types'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)
const resume = ref<Resume | null>(null)

const form = reactive({
  title: '',
  full_name: '',
  email: '',
  phone: '',
  location: '',
  education: '',
  experience: '',
  skills: '',
  summary: '',
  projects: '',
  is_default: false
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入简历标题', trigger: 'blur' }],
  full_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

async function fetchResume() {
  const id = Number(route.params.id)
  if (!id) return

  loading.value = true
  try {
    const res = await getResume(id)
    if (res.data) {
      resume.value = res.data
      Object.assign(form, {
        title: res.data.title,
        full_name: res.data.full_name,
        email: res.data.email || '',
        phone: res.data.phone || '',
        location: res.data.location || '',
        education: res.data.education || '',
        experience: res.data.experience || '',
        skills: res.data.skills || '',
        summary: res.data.summary || '',
        projects: res.data.projects || '',
        is_default: res.data.is_default
      })
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleFileChange(file: UploadFile) {
  if (!file.raw || !resume.value) return

  try {
    await uploadResumeFile(resume.value.id, file.raw)
    ElMessage.success('文件上传成功')
    fetchResume()
  } catch (e) {
    // error handled
  }
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        const id = Number(route.params.id)
        await updateResume(id, form)
        ElMessage.success('保存成功')
        router.push('/resumes')
      } catch (e) {
        // error handled
      } finally {
        saving.value = false
      }
    }
  })
}

function goBack() {
  router.back()
}

onMounted(() => {
  fetchResume()
})
</script>

<style scoped>
.form-container {
  max-width: 800px;
}

.upload-tip {
  color: #909399;
  font-size: 12px;
  margin-top: 8px;
}

.file-name {
  margin-left: 10px;
  color: #67c23a;
}
</style>
