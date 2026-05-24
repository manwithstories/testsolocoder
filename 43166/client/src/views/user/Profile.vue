<template>
  <div class="profile">
    <div class="page-header">
      <h2 class="page-title">个人中心</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="8">
        <div class="card">
          <div class="profile-header text-center">
            <el-avatar :size="100" :src="userStore.userInfo?.avatar">
              {{ userStore.userInfo?.realName?.charAt(0) || 'U' }}
            </el-avatar>
            <h3 class="mt-16">{{ userStore.userInfo?.realName }}</h3>
            <p class="text-muted">{{ getRoleText(userStore.userInfo?.role) }}</p>
          </div>
          <el-divider />
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">{{ userStore.userInfo?.username }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ userStore.userInfo?.email || '-' }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ userStore.userInfo?.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="最后登录">{{ formatTime(userStore.userInfo?.lastLoginAt) }}</el-descriptions-item>
            <el-descriptions-item label="注册时间">{{ formatTime(userStore.userInfo?.createdAt) }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-col>
      <el-col :span="16">
        <div class="card">
          <el-tabs v-model="activeTab">
            <el-tab-pane label="编辑资料" name="edit">
              <el-form :model="editForm" label-width="100px" @submit.prevent>
                <el-form-item label="真实姓名">
                  <el-input v-model="editForm.realName" />
                </el-form-item>
                <el-form-item label="邮箱">
                  <el-input v-model="editForm.email" />
                </el-form-item>
                <el-form-item label="手机号">
                  <el-input v-model="editForm.phone" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="saving" @click="handleSaveProfile">
                    保存修改
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            <el-tab-pane label="修改密码" name="password">
              <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px">
                <el-form-item label="原密码" prop="oldPassword">
                  <el-input v-model="passwordForm.oldPassword" type="password" show-password />
                </el-form-item>
                <el-form-item label="新密码" prop="newPassword">
                  <el-input v-model="passwordForm.newPassword" type="password" show-password />
                </el-form-item>
                <el-form-item label="确认密码" prop="confirmPassword">
                  <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="changingPassword" @click="handleChangePassword">
                    确认修改
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/store/user'
import { UserRole } from '@/types'
import dayjs from 'dayjs'

const userStore = useUserStore()
const activeTab = ref('edit')
const saving = ref(false)
const changingPassword = ref(false)
const passwordFormRef = ref<FormInstance>()

const editForm = reactive({
  realName: '',
  email: '',
  phone: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const getRoleText = (role?: UserRole) => {
  const map: Record<UserRole, string> = {
    admin: '平台管理员',
    entrepreneur: '创业者',
    agent: '代办专员'
  }
  return map[role || 'entrepreneur'] || ''
}

const formatTime = (time?: string | null) => {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

const handleSaveProfile = async () => {
  saving.value = true
  try {
    await userStore.updateProfile({
      realName: editForm.realName,
      email: editForm.email,
      phone: editForm.phone
    })
    ElMessage.success('保存成功')
    await userStore.fetchUserInfo()
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      changingPassword.value = true
      try {
        await userStore.changePassword(passwordForm.oldPassword, passwordForm.newPassword)
        ElMessage.success('密码修改成功')
        passwordForm.oldPassword = ''
        passwordForm.newPassword = ''
        passwordForm.confirmPassword = ''
      } catch (error: any) {
        ElMessage.error(error.message || '密码修改失败')
      } finally {
        changingPassword.value = false
      }
    }
  })
}

onMounted(() => {
  if (userStore.userInfo) {
    editForm.realName = userStore.userInfo.realName || ''
    editForm.email = userStore.userInfo.email || ''
    editForm.phone = userStore.userInfo.phone || ''
  }
})
</script>

<style scoped>
.profile-header {
  padding: 20px 0;
}
</style>
