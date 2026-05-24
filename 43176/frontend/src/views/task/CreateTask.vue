<template>
  <div class="create-task-container">
    <el-card class="form-card">
      <template #header>
        <div class="card-header">
          <el-icon><Edit /></el-icon>
          <span>发布任务</span>
        </div>
      </template>

      <el-form
        ref="taskForm"
        :model="taskForm"
        :rules="rules"
        label-width="100px"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="任务类型" prop="type">
          <el-radio-group v-model="taskForm.type">
            <el-radio value="buy">代购</el-radio>
            <el-radio value="pickup">代取</el-radio>
            <el-radio value="deliver">代送</el-radio>
            <el-radio value="queue">排队代办</el-radio>
            <el-radio value="errand">其他代办</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="任务标题" prop="title">
          <el-input v-model="taskForm.title" placeholder="请输入任务标题" maxlength="100" show-word-limit />
        </el-form-item>

        <el-form-item label="详细说明" prop="description">
          <el-input
            v-model="taskForm.description"
            type="textarea"
            :rows="4"
            placeholder="请详细描述任务要求"
            maxlength="1000"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="起点地址" prop="start_addr">
          <el-input v-model="taskForm.start_addr" placeholder="请输入起点地址" />
        </el-form-item>

        <el-form-item label="终点地址" prop="end_addr">
          <el-input v-model="taskForm.end_addr" placeholder="请输入终点地址" />
        </el-form-item>

        <el-form-item label="截止时间" prop="deadline">
          <el-date-picker
            v-model="taskForm.deadline"
            type="datetime"
            placeholder="选择截止时间"
            :disabled-date="disabledDate"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="报酬金额" prop="reward">
          <el-input-number
            v-model="taskForm.reward"
            :min="1"
            :max="1000"
            :step="5"
            controls-position="right"
            style="width: 100%"
          />
          <span class="reward-hint">预估服务费 ¥{{ serviceFee.toFixed(2) }}</span>
        </el-form-item>

        <el-form-item label="参考图片">
          <el-upload
            v-model:file-list="fileList"
            :auto-upload="false"
            list-type="picture-card"
            :limit="5"
            accept="image/*"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            发布任务
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Edit, Plus } from '@element-plus/icons-vue'
import { taskApi } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const taskFormRef = ref<FormInstance>()
const loading = ref(false)
const fileList = ref<any[]>([])

const taskForm = reactive({
  type: 'deliver',
  title: '',
  description: '',
  start_addr: '',
  start_lat: 39.9042,
  start_lng: 116.4074,
  end_addr: '',
  end_lat: 39.9042,
  end_lng: 116.4074,
  deadline: new Date(Date.now() + 3600000),
  reward: 20
})

const serviceFee = computed(() => {
  const feeRate = taskForm.reward >= 100 ? 0.08 : 0.10
  return taskForm.reward * feeRate
})

const disabledDate = (time: Date) => {
  return time.getTime() < Date.now()
}

const rules: FormRules = {
  type: [
    { required: true, message: '请选择任务类型', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' },
    { min: 5, max: 100, message: '标题长度为5-100个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入详细说明', trigger: 'blur' },
    { min: 10, max: 1000, message: '说明长度为10-1000个字符', trigger: 'blur' }
  ],
  start_addr: [
    { required: true, message: '请输入起点地址', trigger: 'blur' }
  ],
  end_addr: [
    { required: true, message: '请输入终点地址', trigger: 'blur' }
  ],
  deadline: [
    { required: true, message: '请选择截止时间', trigger: 'change' }
  ],
  reward: [
    { required: true, message: '请输入报酬金额', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!taskFormRef.value) return

  await taskFormRef.value.validate(async (valid) => {
    if (valid) {
      const balance = userStore.userInfo?.balance || 0
      if (balance < taskForm.reward) {
        ElMessage.warning('余额不足，请先充值')
        return
      }

      loading.value = true
      try {
        const images = fileList.value.map(f => f.url || f.response?.url || '')
        const res = await taskApi.create({
          ...taskForm,
          images
        })
        if (res.code === 200) {
          ElMessage.success('任务发布成功')
          userStore.updateUserInfo({ balance: balance - taskForm.reward })
          router.push('/my-tasks')
        }
      } catch (error) {
        console.error('Failed to create task:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

const handleCancel = () => {
  router.back()
}
</script>

<style lang="scss" scoped>
.create-task-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.form-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.reward-hint {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

@media (max-width: 768px) {
  .create-task-container {
    padding: 10px;
  }

  :deep(.el-form-item__label) {
    width: 80px !important;
  }
}
</style>
