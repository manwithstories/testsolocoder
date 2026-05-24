<template>
  <div class="page-container">
    <h1 class="page-title">个人中心</h1>

    <el-row :gutter="20">
      <el-col :span="12">
        <div class="card">
          <div class="card-header" style="font-weight: 600;">基本信息</div>
          <div class="card-body">
            <el-form :model="profileForm" label-width="80px">
              <el-form-item label="邮箱">
                <el-input v-model="profileForm.email" disabled />
              </el-form-item>
              <el-form-item label="昵称">
                <el-input v-model="profileForm.nickname" />
              </el-form-item>
              <el-form-item label="头像">
                <el-avatar :size="64" :src="profileForm.avatar" style="margin-bottom: 8px;">
                  {{ profileForm.nickname?.charAt(0) || 'U' }}
                </el-avatar>
                <el-input v-model="profileForm.avatar" placeholder="头像URL" />
              </el-form-item>
              <el-form-item label="角色">
                <el-tag :type="getRoleTagType(profileForm.role)">
                  {{ getRoleText(profileForm.role) }}
                </el-tag>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="saving" @click="handleUpdateProfile">
                  保存
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </el-col>

      <el-col :span="12">
        <div class="card">
          <div class="card-header" style="font-weight: 600;">修改密码</div>
          <div class="card-body">
            <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="80px">
              <el-form-item label="原密码" prop="old_password">
                <el-input v-model="passwordForm.old_password" type="password" show-password />
              </el-form-item>
              <el-form-item label="新密码" prop="new_password">
                <el-input v-model="passwordForm.new_password" type="password" show-password />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="changing" @click="handleChangePassword">
                  修改密码
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { authApi } from '@/api/auth'
import { useUserStore } from '@/store'

const userStore = useUserStore()

const saving = ref(false)
const changing = ref(false)
const passwordFormRef = ref<FormInstance>()

const profileForm = reactive({
  email: '',
  nickname: '',
  avatar: '',
  role: ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: ''
})

const passwordRules: FormRules = {
  old_password: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const getRoleText = (role: string) => {
  const map: Record<string, string> = {
    admin: '管理员',
    editor: '编辑员',
    viewer: '查看员'
  }
  return map[role] || role
}

const getRoleTagType = (role: string) => {
  const map: Record<string, string> = {
    admin: 'danger',
    editor: 'warning',
    viewer: 'info'
  }
  return map[role] || ''
}

const loadProfile = async () => {
  try {
    const profile = await authApi.getProfile()
    profileForm.email = profile.email
    profileForm.nickname = profile.nickname
    profileForm.avatar = profile.avatar
    profileForm.role = profile.role
  } catch (error) {
    console.error('Failed to load profile')
  }
}

const handleUpdateProfile = async () => {
  saving.value = true
  try {
    await authApi.updateProfile({
      nickname: profileForm.nickname,
      avatar: profileForm.avatar
    })
    ElMessage.success('信息已更新')
    if (userStore.user) {
      userStore.user.nickname = profileForm.nickname
      userStore.user.avatar = profileForm.avatar
    }
  } catch (e: any) {
    ElMessage.error(e.message || '更新失败')
  } finally {
    saving.value = false
  }
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    changing.value = true

    await authApi.changePassword(passwordForm)
    ElMessage.success('密码已修改')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
  } catch (e: any) {
    if (e.message) {
      ElMessage.error(e.message)
    }
  } finally {
    changing.value = false
  }
}

onMounted(loadProfile)
</script>
