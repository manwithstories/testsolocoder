<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">安排面试</div>
      </div>

      <el-card>
        <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" class="form-container">
          <el-form-item label="选择投递" prop="application_id">
            <el-select
              v-model="form.application_id"
              placeholder="请选择投递记录"
              filterable
              style="width: 100%"
              @change="handleApplicationChange"
            >
              <el-option
                v-for="app in applications"
                :key="app.id"
                :label="`${app.applicant?.profile?.full_name || app.applicant?.email} - ${app.job?.title}`"
                :value="app.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="面试时间" prop="scheduled_at">
            <el-date-picker
              v-model="form.scheduled_at"
              type="datetime"
              placeholder="选择面试时间"
              style="width: 100%"
              :disabled-date="disabledDate"
            />
          </el-form-item>
          <el-form-item label="面试时长">
            <el-input-number v-model="form.duration" :min="15" :max="240" :step="15" />
            <span style="margin-left: 10px; color: #909399;">分钟</span>
          </el-form-item>
          <el-form-item label="面试方式">
            <el-select v-model="form.interview_type" placeholder="请选择面试方式" style="width: 100%">
              <el-option label="现场面试" value="onsite" />
              <el-option label="远程面试" value="online" />
              <el-option label="电话面试" value="phone" />
            </el-select>
          </el-form-item>
          <el-form-item label="面试地点">
            <el-input v-model="form.location" placeholder="请输入面试地点" />
          </el-form-item>
          <el-form-item label="会议链接" v-if="form.interview_type === 'online'">
            <el-input v-model="form.meeting_link" placeholder="请输入会议链接" />
          </el-form-item>
          <el-form-item label="面试官">
            <el-input v-model="form.interviewer" placeholder="请输入面试官姓名" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input
              v-model="form.notes"
              type="textarea"
              :rows="3"
              placeholder="请输入备注信息"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="handleSubmit">
              安排面试
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
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { scheduleInterview } from '@/api/interview'
import { listCompanyApplications } from '@/api/application'
import type { Application } from '@/types'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const applications = ref<Application[]>([])

const form = reactive({
  application_id: null as number | null,
  scheduled_at: null as string | null,
  duration: 60,
  interview_type: 'onsite',
  location: '',
  meeting_link: '',
  interviewer: '',
  notes: ''
})

const rules: FormRules = {
  application_id: [{ required: true, message: '请选择投递记录', trigger: 'change' }],
  scheduled_at: [{ required: true, message: '请选择面试时间', trigger: 'change' }]
}

async function fetchApplications() {
  try {
    const res = await listCompanyApplications({
      page: 1,
      page_size: 100,
      status: 'interested'
    })
    if (res.data) {
      applications.value = res.data.items
    }
  } catch (e) {
    console.error(e)
  }
}

function handleApplicationChange() {
  // Optional: do something when application changes
}

function disabledDate(date: Date) {
  return date.getTime() < Date.now() - 86400000
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await scheduleInterview({
          application_id: form.application_id!,
          scheduled_at: form.scheduled_at!,
          duration: form.duration,
          interview_type: form.interview_type,
          location: form.location,
          meeting_link: form.meeting_link,
          interviewer: form.interviewer,
          notes: form.notes
        })
        ElMessage.success('面试安排成功，已发送通知邮件')
        router.push('/company/interviews')
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

onMounted(() => {
  fetchApplications()
  const applicationId = route.query.applicationId
  if (applicationId) {
    form.application_id = Number(applicationId)
  }
})
</script>

<style scoped>
.form-container {
  max-width: 700px;
}
</style>
