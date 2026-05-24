<template>
  <div class="page-container">
    <h1 class="page-title">个人中心</h1>

    <div v-loading="loading" class="profile-content">
      <div class="card user-info-card">
        <div class="avatar-section">
          <el-avatar :size="100" :src="user?.avatar">
            {{ user?.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <el-button type="primary" link @click="showAvatarDialog = true">更换头像</el-button>
        </div>
        <div class="user-info">
          <h2>{{ user?.realName || user?.username }}</h2>
          <div class="user-role">
            <el-tag :type="getRoleTagType(user?.role)">
              {{ getRoleText(user?.role) }}
            </el-tag>
            <el-tag v-if="user?.verified" type="success" style="margin-left: 8px">
              已认证
            </el-tag>
            <el-tag v-else type="info" style="margin-left: 8px">
              未认证
            </el-tag>
          </div>
          <div class="user-meta">
            <p><span>邮箱：</span>{{ user?.email }}</p>
            <p><span>手机：</span>{{ user?.phone || '未绑定' }}</p>
            <p><span>注册时间：</span>{{ formatDate(user?.createdAt) }}</p>
          </div>
        </div>
        <div class="balance-section">
          <div class="balance-card">
            <div class="balance-label">押金余额</div>
            <div class="balance-value">¥{{ user?.depositBalance?.toFixed(2) || '0.00' }}</div>
          </div>
        </div>
      </div>

      <div class="card edit-profile-card">
        <h3>编辑个人信息</h3>
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
          style="max-width: 500px"
        >
          <el-form-item label="真实姓名" prop="realName">
            <el-input v-model="form.realName" placeholder="请输入真实姓名" />
          </el-form-item>
          <el-form-item label="身份证号" prop="idCard">
            <el-input v-model="form.idCard" placeholder="请输入身份证号" />
          </el-form-item>
          <el-form-item label="手机号码" prop="phone">
            <el-input v-model="form.phone" placeholder="请输入手机号码" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="saving" @click="handleSave">
              保存
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <el-dialog v-model="showAvatarDialog" title="更换头像" width="400px">
      <el-upload
        class="avatar-uploader"
        :show-file-list="false"
        :auto-upload="false"
        :on-change="handleAvatarChange"
        accept="image/jpeg,image/png"
      >
        <el-avatar :size="100" :src="avatarPreview">
          {{ user?.username?.charAt(0)?.toUpperCase() }}
        </el-avatar>
      </el-upload>
      <template #footer>
        <el-button @click="showAvatarDialog = false">取消</el-button>
        <el-button type="primary" @click="saveAvatar">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import dayjs from 'dayjs'
import { userApi } from '@/api/equipment'
import { useUserStore } from '@/stores/user'
import type { User } from '@/types'

const userStore = useUserStore()
const formRef = ref<FormInstance>()

const loading = ref(false)
const saving = ref(false)
const user = ref<User | null>(null)
const showAvatarDialog = ref(false)
const avatarPreview = ref('')
const avatarFile = ref<File | null>(null)

const form = reactive({
  realName: '',
  idCard: '',
  phone: ''
})

const rules: FormRules = {
  realName: [
    { max: 50, message: '姓名长度不能超过50个字符', trigger: 'blur' }
  ],
  idCard: [
    { len: 18, message: '身份证号必须为18位', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

onMounted(() => {
  loadProfile()
})

async function loadProfile() {
  loading.value = true
  try {
    const response = await userApi.getProfile()
    user.value = response.data
    form.realName = response.data.realName || ''
    form.idCard = response.data.idCard || ''
    form.phone = response.data.phone || ''
    avatarPreview.value = response.data.avatar || ''
    userStore.setUser(response.data)
  } catch (error) {
    console.error('Failed to load profile:', error)
    ElMessage.error('加载用户信息失败')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        const response = await userApi.updateProfile(form)
        user.value = response.data
        userStore.setUser(response.data)
        ElMessage.success('保存成功')
      } catch (error) {
        console.error('Failed to update profile:', error)
      } finally {
        saving.value = false
      }
    }
  })
}

function handleAvatarChange(file: any) {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      avatarPreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file.raw)
    avatarFile.value = file.raw
  }
}

async function saveAvatar() {
  if (!avatarFile.value || !user.value) return

  try {
    ElMessage.info('头像上传功能需要后端支持')
    showAvatarDialog.value = false
  } catch (error) {
    console.error('Failed to save avatar:', error)
  }
}

function getRoleText(role?: string) {
  const textMap: Record<string, string> = {
    renter: '租借方',
    owner: '出租方',
    admin: '管理员'
  }
  return textMap[role || ''] || role
}

function getRoleTagType(role?: string) {
  const typeMap: Record<string, string> = {
    renter: '',
    owner: 'success',
    admin: 'danger'
  }
  return typeMap[role || ''] || ''
}

function formatDate(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD') : '-'
}
</script>

<style scoped>
.profile-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-info-card {
  display: flex;
  align-items: center;
  gap: 30px;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.user-info {
  flex: 1;
}

.user-info h2 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.user-role {
  margin-bottom: 16px;
}

.user-meta p {
  margin-bottom: 8px;
  color: #606266;
}

.user-meta span {
  color: #909399;
}

.balance-section {
  min-width: 200px;
}

.balance-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
}

.balance-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 8px;
}

.balance-value {
  font-size: 28px;
  font-weight: 600;
}

.edit-profile-card h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
  color: #303133;
}

.avatar-uploader :deep(.el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.avatar-uploader :deep(.el-upload:hover) {
  border-color: #409eff;
}
</style>
