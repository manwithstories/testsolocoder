<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">创建简历</div>
      </div>

      <el-card>
        <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" class="form-container">
          <el-form-item label="简历标题" prop="title">
            <el-input v-model="form.title" placeholder="请输入简历标题，如：前端工程师简历" />
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
              placeholder="请输入技能，多个技能用逗号分隔"
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
          <el-form-item label="设为默认">
            <el-switch v-model="form.is_default" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="handleSubmit">
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { createResume } from '@/api/resume'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

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

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await createResume(form)
        if (res.data) {
          ElMessage.success('创建成功')
          router.push(`/resumes/${res.data.id}`)
        }
      } catch (e) {
        // error handled
      } finally {
        loading.value = false
      }
    }
  })
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.form-container {
  max-width: 800px;
}
</style>
