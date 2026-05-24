<template>
  <div class="profile-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">个人中心</h2>
      </div>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-card>
            <el-form :model="form" label-width="100px">
              <el-form-item label="头像">
                <el-avatar :size="80" :src="form.avatar">
                  {{ form.username?.charAt(0) }}
                </el-avatar>
                <el-upload
                  class="avatar-uploader"
                  action="/api/upload"
                  :show-file-list="false"
                  :on-success="handleAvatarSuccess"
                >
                  <el-button size="small" type="primary" style="margin-left: 10px;">
                    上传
                  </el-button>
                </el-upload>
              </el-form-item>
              <el-form-item label="用户名">
                <el-input v-model="form.username" disabled />
              </el-form-item>
              <el-form-item label="真实姓名">
                <el-input v-model="form.real_name" />
              </el-form-item>
              <el-form-item label="手机号">
                <el-input v-model="form.phone" />
              </el-form-item>
              <el-form-item label="邮箱">
                <el-input v-model="form.email" />
              </el-form-item>
              <el-form-item label="地址">
                <el-input v-model="form.address" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveProfile" :loading="saving">
                  保存
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>

        <el-tab-pane v-if="userStore.isTechnician" label="技师认证" name="certificate">
          <el-card>
            <div v-if="techProfile" class="certificate-info">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="认证状态">
                  <el-tag :type="getVerifyStatusType(techProfile.verify_status)">
                    {{ getVerifyStatusText(techProfile.verify_status) }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="证书编号">
                  {{ techProfile.certificate_no || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="专业领域">
                  {{ techProfile.specialty || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="从业年限">
                  {{ techProfile.experience_years }}年
                </el-descriptions-item>
                <el-descriptions-item label="备注" :span="2">
                  {{ techProfile.verify_remark || '-' }}
                </el-descriptions-item>
              </el-descriptions>
            </div>

            <el-form :model="certForm" label-width="100px" class="mt-20">
              <el-form-item label="证书图片">
                <el-input v-model="certForm.certificate_image" placeholder="请输入证书图片URL" />
              </el-form-item>
              <el-form-item label="证书编号">
                <el-input v-model="certForm.certificate_no" placeholder="请输入证书编号" />
              </el-form-item>
              <el-form-item label="专业领域">
                <el-input v-model="certForm.specialty" placeholder="如：家电维修、数码维修等" />
              </el-form-item>
              <el-form-item label="从业年限">
                <el-input-number v-model="certForm.experience_years" :min="0" :max="50" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="submitCertificate" :loading="submitting">
                  提交认证
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="账户安全" name="security">
          <el-card>
            <el-form label-width="100px">
              <el-form-item label="修改密码">
                <el-button type="primary" @click="showPasswordDialog = true">
                  修改密码
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="showPasswordDialog" title="修改密码" width="400px">
      <el-form :model="passwordForm" label-width="100px">
        <el-form-item label="当前密码">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.newPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="changePassword" :loading="changingPassword">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { userApi } from '@/api/user'
import type { TechnicianProfile } from '@/types'

const userStore = useUserStore()

const activeTab = ref('basic')
const saving = ref(false)
const submitting = ref(false)
const showPasswordDialog = ref(false)
const changingPassword = ref(false)
const techProfile = ref<TechnicianProfile | null>(null)

const form = reactive({
  username: '',
  real_name: '',
  phone: '',
  email: '',
  address: '',
  avatar: ''
})

const certForm = reactive({
  certificate_image: '',
  certificate_no: '',
  specialty: '',
  experience_years: 0
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

onMounted(() => {
  loadUserInfo()
})

function loadUserInfo() {
  if (userStore.userInfo) {
    form.username = userStore.userInfo.username
    form.real_name = userStore.userInfo.real_name || ''
    form.phone = userStore.userInfo.phone || ''
    form.email = userStore.userInfo.email || ''
    form.address = userStore.userInfo.address || ''
    form.avatar = userStore.userInfo.avatar || ''

    const techData = (userStore.userInfo as any)?.technician_profile
    if (techData) {
      techProfile.value = techData
      certForm.certificate_image = techData.certificate_image || ''
      certForm.certificate_no = techData.certificate_no || ''
      certForm.specialty = techData.specialty || ''
      certForm.experience_years = techData.experience_years || 0
    }
  }
}

async function saveProfile() {
  saving.value = true
  try {
    await userApi.updateProfile({
      real_name: form.real_name,
      phone: form.phone,
      email: form.email,
      address: form.address,
      avatar: form.avatar
    })
    ElMessage.success('保存成功')
    await userStore.fetchUserInfo()
    loadUserInfo()
  } catch (error) {
    console.error('Failed to save profile:', error)
  } finally {
    saving.value = false
  }
}

function handleAvatarSuccess(response: any) {
  if (response.data?.url) {
    form.avatar = response.data.url
  }
}

async function submitCertificate() {
  submitting.value = true
  try {
    await userApi.submitCertificate({
      certificate_image: certForm.certificate_image,
      certificate_no: certForm.certificate_no,
      specialty: certForm.specialty,
      experience_years: certForm.experience_years
    })
    ElMessage.success('认证提交成功，请等待审核')
    await userStore.fetchUserInfo()
    loadUserInfo()
  } catch (error) {
    console.error('Failed to submit certificate:', error)
  } finally {
    submitting.value = false
  }
}

function getVerifyStatusType(status: string): string {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

function getVerifyStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return statusMap[status] || status
}

async function changePassword() {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  changingPassword.value = true
  try {
    ElMessage.success('密码修改成功')
    showPasswordDialog.value = false
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error) {
    console.error('Failed to change password:', error)
  } finally {
    changingPassword.value = false
  }
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.certificate-info {
  margin-bottom: 20px;
}

.mt-20 {
  margin-top: 20px;
}
</style>
