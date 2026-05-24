<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">发布职位</div>
      </div>

      <el-card>
        <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" class="form-container">
          <el-form-item label="职位名称" prop="title">
            <el-input v-model="form.title" placeholder="请输入职位名称" />
          </el-form-item>
          <el-form-item label="职位描述" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="6"
              placeholder="请输入职位描述"
            />
          </el-form-item>
          <el-form-item label="薪资范围">
            <el-input-group>
              <el-input v-model.number="form.salary_min" placeholder="最低薪资" style="width: 150px" />
              <template #append>-</template>
              <el-input v-model.number="form.salary_max" placeholder="最高薪资" style="width: 150px" />
              <el-select v-model="form.salary_type" placeholder="薪资类型" style="width: 120px">
                <el-option label="月薪" value="monthly" />
                <el-option label="年薪" value="yearly" />
                <el-option label="日薪" value="daily" />
              </el-select>
            </el-input-group>
          </el-form-item>
          <el-form-item label="工作地点" prop="location">
            <el-input v-model="form.location" placeholder="请输入工作地点" />
          </el-form-item>
          <el-form-item label="工作类型">
            <el-select v-model="form.job_type" placeholder="请选择工作类型" style="width: 100%">
              <el-option label="全职" value="full-time" />
              <el-option label="兼职" value="part-time" />
              <el-option label="合同" value="contract" />
              <el-option label="实习" value="intern" />
              <el-option label="远程" value="remote" />
            </el-select>
          </el-form-item>
          <el-form-item label="经验要求">
            <el-select v-model="form.experience" placeholder="请选择经验要求" style="width: 100%">
              <el-option label="不限" value="不限" />
              <el-option label="应届生" value="应届生" />
              <el-option label="1年以内" value="1年以内" />
              <el-option label="1-3年" value="1-3年" />
              <el-option label="3-5年" value="3-5年" />
              <el-option label="5-10年" value="5-10年" />
              <el-option label="10年以上" value="10年以上" />
            </el-select>
          </el-form-item>
          <el-form-item label="学历要求">
            <el-select v-model="form.education" placeholder="请选择学历要求" style="width: 100%">
              <el-option label="不限" value="不限" />
              <el-option label="大专" value="大专" />
              <el-option label="本科" value="本科" />
              <el-option label="硕士" value="硕士" />
              <el-option label="博士" value="博士" />
            </el-select>
          </el-form-item>
          <el-form-item label="技能要求">
            <el-input
              v-model="form.skills"
              placeholder="请输入技能要求，多个技能用逗号分隔"
            />
          </el-form-item>
          <el-form-item label="任职要求">
            <el-input
              v-model="form.requirements"
              type="textarea"
              :rows="4"
              placeholder="请输入任职要求"
            />
          </el-form-item>
          <el-form-item label="福利待遇">
            <el-input
              v-model="form.benefits"
              type="textarea"
              :rows="3"
              placeholder="请输入福利待遇"
            />
          </el-form-item>
          <el-form-item label="截止日期">
            <el-date-picker
              v-model="form.deadline"
              type="date"
              placeholder="选择截止日期"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="handleSaveDraft">
              保存草稿
            </el-button>
            <el-button type="success" :loading="loading" @click="handlePublish">
              保存并发布
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
import { createJob } from '@/api/job'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  title: '',
  description: '',
  salary_min: null as number | null,
  salary_max: null as number | null,
  salary_type: 'monthly',
  location: '',
  job_type: 'full-time',
  experience: '',
  education: '',
  skills: '',
  requirements: '',
  benefits: '',
  deadline: null as string | null
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入职位名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入职位描述', trigger: 'blur' }],
  location: [{ required: true, message: '请输入工作地点', trigger: 'blur' }]
}

async function handleSubmit(publish = false) {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const submitData: any = { ...form }
        if (publish) {
          submitData.status = 'open'
        }
        await createJob(submitData)
        ElMessage.success(publish ? '发布成功' : '保存成功')
        router.push('/my/jobs')
      } catch (e) {
        // error handled
      } finally {
        loading.value = false
      }
    }
  })
}

function handleSaveDraft() {
  handleSubmit(false)
}

function handlePublish() {
  handleSubmit(true)
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
