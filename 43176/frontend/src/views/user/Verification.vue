<template>
  <div class="verification-container">
    <el-card v-if="verificationStatus === 'pending'" class="status-card">
      <el-result icon="info" title="认证审核中" sub-title="您的实名认证申请正在审核中，请耐心等待">
        <template #extra>
          <el-button type="primary" @click="fetchStatus">刷新状态</el-button>
        </template>
      </el-result>
    </el-card>

    <el-card v-else-if="verificationStatus === 'approved'" class="status-card">
      <el-result icon="success" title="认证已通过" sub-title="您的实名认证已通过审核">
        <template #extra>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="姓名">{{ verificationInfo?.id_card_name }}</el-descriptions-item>
            <el-descriptions-item label="身份证号">{{ maskIDCard(verificationInfo?.id_card_no) }}</el-descriptions-item>
            <el-descriptions-item label="审核时间">{{ formatTime(verificationInfo?.reviewed_at) }}</el-descriptions-item>
          </el-descriptions>
        </template>
      </el-result>
    </el-card>

    <el-card v-else-if="verificationStatus === 'rejected'" class="status-card">
      <el-result icon="error" title="认证被拒绝" :sub-title="`拒绝原因：${verificationInfo?.reason || '未知'}`">
        <template #extra>
          <el-button type="primary" @click="resetForm">重新提交</el-button>
        </template>
      </el-result>
    </el-card>

    <el-card v-else class="form-card">
      <template #header>
        <div class="card-header">
          <el-icon><Postcard /></el-icon>
          <span>实名认证</span>
        </div>
      </template>

      <el-alert
        title="请上传真实有效的身份证信息，信息将被严格保密"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      />

      <el-form :model="verifyForm" :rules="rules" ref="verifyFormRef" label-width="100px">
        <el-form-item label="真实姓名" prop="id_card_name">
          <el-input v-model="verifyForm.id_card_name" placeholder="请输入身份证上的姓名" />
        </el-form-item>
        <el-form-item label="身份证号" prop="id_card_no">
          <el-input v-model="verifyForm.id_card_no" placeholder="请输入18位身份证号" />
        </el-form-item>
        <el-form-item label="身份证正面" prop="id_card_front">
          <el-upload
            :auto-upload="false"
            :show-file-list="false"
            accept="image/*"
            @change="(file: UploadFile) => handleImageUpload(file, 'front')"
          >
            <div class="upload-area">
              <img v-if="verifyForm.id_card_front" :src="verifyForm.id_card_front" class="preview-img" />
              <div v-else class="upload-placeholder">
                <el-icon :size="48"><Plus /></el-icon>
                <p>点击上传身份证正面（人像面）</p>
              </div>
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item label="身份证背面" prop="id_card_back">
          <el-upload
            :auto-upload="false"
            :show-file-list="false"
            accept="image/*"
            @change="(file: UploadFile) => handleImageUpload(file, 'back')"
          >
            <div class="upload-area">
              <img v-if="verifyForm.id_card_back" :src="verifyForm.id_card_back" class="preview-img" />
              <div v-else class="upload-placeholder">
                <el-icon :size="48"><Plus /></el-icon>
                <p>点击上传身份证背面（国徽面）</p>
              </div>
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitVerification" :loading="submitting">
            提交认证
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules, type UploadFile } from 'element-plus'
import { Postcard, Plus } from '@element-plus/icons-vue'
import { userApi } from '@/api'
import dayjs from 'dayjs'

const verifyFormRef = ref<FormInstance>()
const submitting = ref(false)
const verificationStatus = ref<'none' | 'pending' | 'approved' | 'rejected'>('none')
const verificationInfo = ref<any>(null)

const verifyForm = reactive({
  id_card_name: '',
  id_card_no: '',
  id_card_front: '',
  id_card_back: ''
})

const rules: FormRules = {
  id_card_name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  id_card_no: [
    { required: true, message: '请输入身份证号', trigger: 'blur' },
    { pattern: /^\d{17}[\dXx]$/, message: '请输入正确的18位身份证号', trigger: 'blur' }
  ],
  id_card_front: [
    { required: true, message: '请上传身份证正面', trigger: 'change' }
  ],
  id_card_back: [
    { required: true, message: '请上传身份证背面', trigger: 'change' }
  ]
}

const maskIDCard = (idCard?: string) => {
  if (!idCard) return ''
  return idCard.replace(/^(\d{4})\d+(\d{4})$/, '$1**********$2')
}

const formatTime = (time?: string) => {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

const handleImageUpload = (file: UploadFile, type: 'front' | 'back') => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      if (type === 'front') {
        verifyForm.id_card_front = e.target?.result as string
      } else {
        verifyForm.id_card_back = e.target?.result as string
      }
    }
    reader.readAsDataURL(file.raw)
  }
}

const fetchStatus = async () => {
  try {
    const res = await userApi.getProfile()
    if (res.code === 200) {
      const data: any = res.data
      if (data.courier_profile?.status === 'approved') {
        verificationStatus.value = 'approved'
      } else if (data.courier_profile?.status === 'rejected') {
        verificationStatus.value = 'rejected'
      } else if (data.courier_profile?.status === 'pending') {
        verificationStatus.value = 'pending'
      }
    }
  } catch (error) {
    console.error('Failed to fetch status:', error)
  }
}

const submitVerification = async () => {
  if (!verifyFormRef.value) return

  await verifyFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const res = await userApi.submitVerification(verifyForm)
        if (res.code === 200) {
          ElMessage.success('认证申请已提交，请等待审核')
          verificationStatus.value = 'pending'
        }
      } catch (error) {
        console.error('Verification failed:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const resetForm = () => {
  verifyForm.id_card_name = ''
  verifyForm.id_card_no = ''
  verifyForm.id_card_front = ''
  verifyForm.id_card_back = ''
  verificationStatus.value = 'none'
}

onMounted(() => {
  fetchStatus()
})
</script>

<style lang="scss" scoped>
.verification-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.status-card {
  margin-bottom: 20px;
}

.form-card {
  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
  }

  .upload-area {
    width: 200px;
    height: 120px;
    border: 2px dashed #dcdfe6;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.3s;

    &:hover {
      border-color: #667eea;
    }

    .upload-placeholder {
      text-align: center;
      color: #909399;

      p {
        margin-top: 8px;
        font-size: 12px;
      }
    }

    .preview-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      border-radius: 6px;
    }
  }
}
</style>
