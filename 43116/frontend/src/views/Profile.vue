<template>
  <div class="container" style="padding-top: 20px;">
    <div class="page-header">
      <h1 class="page-title">个人中心</h1>
    </div>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>头像</span>
          </template>
          <div style="text-align: center;">
            <el-avatar :size="120" :src="userStore.userInfo?.avatar">
              {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <el-upload
              :show-file-list="false"
              :before-upload="handleAvatarUpload"
              accept="image/*"
              style="margin-top: 16px;"
            >
              <el-button type="primary" :icon="Upload">上传头像</el-button>
            </el-upload>
          </div>
        </el-card>

        <el-card style="margin-top: 20px;">
          <template #header>
            <span>实名认证</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="认证状态">
              <el-tag :type="userStore.userInfo?.auth_status === 'active' ? 'success' : 'warning'">
                {{ getAuthStatusText(userStore.userInfo?.auth_status) }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
          <div v-if="userStore.userInfo?.auth_status !== 'active'" style="margin-top: 16px;">
            <el-upload
              :show-file-list="false"
              :before-upload="(file: File) => handleLicenseUpload(file, 'license_image')"
              accept="image/*"
            >
              <el-button type="primary" :icon="Upload">上传驾驶证</el-button>
            </el-upload>
            <el-upload
              :show-file-list="false"
              :before-upload="(file: File) => handleLicenseUpload(file, 'id_card_front')"
              accept="image/*"
              style="margin-left: 10px;"
            >
              <el-button :icon="Upload">上传身份证正面</el-button>
            </el-upload>
            <el-upload
              :show-file-list="false"
              :before-upload="(file: File) => handleLicenseUpload(file, 'id_card_back')"
              accept="image/*"
              style="margin-left: 10px;"
            >
              <el-button :icon="Upload">上传身份证反面</el-button>
            </el-upload>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>
            <span>基本信息</span>
          </template>
          <el-form :model="form" label-width="100px" style="max-width: 500px;">
            <el-form-item label="用户名">
              <el-input v-model="form.username" disabled />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="form.email" />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="form.phone" />
            </el-form-item>
            <el-form-item label="真实姓名">
              <el-input v-model="form.real_name" />
            </el-form-item>
            <el-form-item label="身份证号">
              <el-input v-model="form.id_card" />
            </el-form-item>
            <el-form-item label="驾驶证号">
              <el-input v-model="form.license_number" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveProfile">保存修改</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card style="margin-top: 20px;">
          <template #header>
            <span>修改密码</span>
          </template>
          <el-form :model="passwordForm" label-width="100px" style="max-width: 400px;">
            <el-form-item label="原密码">
              <el-input v-model="passwordForm.oldPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="passwordForm.newPassword" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="changePassword">修改密码</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api'

const userStore = useUserStore()

const form = reactive({
  username: '',
  email: '',
  phone: '',
  real_name: '',
  id_card: '',
  license_number: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: ''
})

onMounted(() => {
  loadProfile()
})

const loadProfile = async () => {
  try {
    const res = await userApi.getProfile()
    if (res.data) {
      form.username = res.data.username
      form.email = res.data.email
      form.phone = res.data.phone
      form.real_name = res.data.real_name
      form.id_card = res.data.id_card
      form.license_number = res.data.license_number
    }
  } catch {
    // ignore
  }
}

const saveProfile = async () => {
  try {
    await userApi.updateProfile(form)
    ElMessage.success('保存成功')
    await userStore.getProfile()
  } catch (err: any) {
    ElMessage.error(err.message || '保存失败')
  }
}

const changePassword = async () => {
  if (!passwordForm.oldPassword || !passwordForm.newPassword) {
    ElMessage.warning('请填写完整信息')
    return
  }
  try {
    await userApi.changePassword(passwordForm.oldPassword, passwordForm.newPassword)
    ElMessage.success('密码修改成功')
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
  } catch (err: any) {
    ElMessage.error(err.message || '修改失败')
  }
}

const handleAvatarUpload = async (file: File) => {
  try {
    const res = await userApi.uploadAvatar(file)
    ElMessage.success('头像上传成功')
    await userStore.getProfile()
  } catch (err: any) {
    ElMessage.error(err.message || '上传失败')
  }
  return false
}

const handleLicenseUpload = async (file: File, field: string) => {
  try {
    await userApi.uploadLicense(file, field)
    ElMessage.success('上传成功，等待审核')
  } catch (err: any) {
    ElMessage.error(err.message || '上传失败')
  }
  return false
}

const getAuthStatusText = (status?: string) => {
  const map: Record<string, string> = {
    pending: '待审核',
    active: '已认证',
    disabled: '已禁用',
    rejected: '已拒绝'
  }
  return map[status || 'pending'] || '待审核'
}
</script>
