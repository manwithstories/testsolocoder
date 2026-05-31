<template>
  <div class="dashboard-profile">
    <div class="card-shadow p-20">
      <h2 class="page-title mb-20">个人中心</h2>

      <el-row :gutter="20">
        <el-col :span="6">
          <div class="avatar-section">
            <div class="avatar-wrapper">
              <el-avatar
                :size="120"
                :src="userStore.user?.avatar"
                class="profile-avatar"
              >
                {{ userStore.user?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <el-upload
                class="avatar-uploader"
                :show-file-list="false"
                :before-upload="handleAvatarUpload"
                action="#"
              >
                <div class="avatar-edit">
                  <el-icon><Edit /></el-icon>
                </div>
              </el-upload>
            </div>
            <div class="user-info mt-20">
              <h3 class="user-name">{{ userStore.user?.nickname || userStore.user?.username }}</h3>
              <p class="user-role">
                <el-tag :type="roleType(userStore.user?.role || 'normal')" size="small">
                  {{ roleText(userStore.user?.role || 'normal') }}
                </el-tag>
              </p>
              <p class="user-join">
                注册时间：{{ formatDate(userStore.user?.created_at || '') }}
              </p>
            </div>
          </div>
        </el-col>

        <el-col :span="18">
          <div class="form-section">
            <el-tabs v-model="activeTab">
              <el-tab-pane label="基本信息" name="basic">
                <el-form
                  ref="basicFormRef"
                  :model="basicForm"
                  :rules="basicRules"
                  label-width="80px"
                  class="mt-20"
                >
                  <el-form-item label="昵称" prop="nickname">
                    <el-input v-model="basicForm.nickname" placeholder="请输入昵称" style="max-width: 400px" />
                  </el-form-item>
                  <el-form-item label="邮箱" prop="email">
                    <el-input v-model="basicForm.email" placeholder="请输入邮箱" style="max-width: 400px" />
                  </el-form-item>
                  <el-form-item label="电话" prop="phone">
                    <el-input v-model="basicForm.phone" placeholder="请输入手机号" style="max-width: 400px" />
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" :loading="saving" @click="handleSaveBasic">保存修改</el-button>
                  </el-form-item>
                </el-form>
              </el-tab-pane>

              <el-tab-pane label="修改密码" name="password">
                <el-form
                  ref="passwordFormRef"
                  :model="passwordForm"
                  :rules="passwordRules"
                  label-width="100px"
                  class="mt-20"
                >
                  <el-form-item label="当前密码" prop="oldPassword">
                    <el-input
                      v-model="passwordForm.oldPassword"
                      type="password"
                      placeholder="请输入当前密码"
                      show-password
                      style="max-width: 400px"
                    />
                  </el-form-item>
                  <el-form-item label="新密码" prop="newPassword">
                    <el-input
                      v-model="passwordForm.newPassword"
                      type="password"
                      placeholder="请输入新密码"
                      show-password
                      style="max-width: 400px"
                    />
                  </el-form-item>
                  <el-form-item label="确认密码" prop="confirmPassword">
                    <el-input
                      v-model="passwordForm.confirmPassword"
                      type="password"
                      placeholder="请再次输入新密码"
                      show-password
                      style="max-width: 400px"
                    />
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" :loading="changingPassword" @click="handleChangePassword">修改密码</el-button>
                  </el-form-item>
                </el-form>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Edit } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import * as authApi from '@/api/auth'
import dayjs from 'dayjs'
import type { User } from '@/types'

const userStore = useUserStore()
const activeTab = ref('basic')
const saving = ref(false)
const changingPassword = ref(false)

const basicFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

const basicForm = reactive({
  nickname: '',
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

const basicRules: FormRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const roleType = (role: string) => {
  if (role === 'admin') return 'danger'
  if (role === 'guide') return 'primary'
  if (role === 'researcher') return 'success'
  return 'info'
}

const roleText = (role: string) => {
  const map: Record<string, string> = {
    normal: '普通用户',
    admin: '管理员',
    guide: '讲解员',
    researcher: '研究员'
  }
  return map[role] || role
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const handleAvatarUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB')
    return false
  }

  const reader = new FileReader()
  reader.onload = async (e) => {
    const avatar = e.target?.result as string
    try {
      await authApi.updateProfile({ avatar } as Partial<User>)
      userStore.user!.avatar = avatar
      ElMessage.success('头像更新成功')
    } catch (e) {
      console.error(e)
      ElMessage.error('头像更新失败')
    }
  }
  reader.readAsDataURL(file)
  return false
}

const handleSaveBasic = async () => {
  if (!basicFormRef.value) return
  await basicFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        saving.value = true
        await authApi.updateProfile({
          nickname: basicForm.nickname,
          email: basicForm.email,
          phone: basicForm.phone
        } as Partial<User>)
        userStore.user!.nickname = basicForm.nickname
        userStore.user!.email = basicForm.email
        userStore.user!.phone = basicForm.phone
        localStorage.setItem('user', JSON.stringify(userStore.user))
        ElMessage.success('保存成功')
      } catch (e) {
        console.error(e)
        ElMessage.error('保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        changingPassword.value = true
        await authApi.updateProfile({
          password: passwordForm.oldPassword,
          newPassword: passwordForm.newPassword
        } as Partial<User>)
        ElMessage.success('密码修改成功')
        passwordForm.oldPassword = ''
        passwordForm.newPassword = ''
        passwordForm.confirmPassword = ''
      } catch (e) {
        console.error(e)
        ElMessage.error('密码修改失败')
      } finally {
        changingPassword.value = false
      }
    }
  })
}

onMounted(() => {
  if (userStore.user) {
    basicForm.nickname = userStore.user.nickname || ''
    basicForm.email = userStore.user.email || ''
    basicForm.phone = userStore.user.phone || ''
  }
})
</script>

<style scoped lang="scss">
.dashboard-profile {
  .avatar-section {
    text-align: center;
    padding: 20px 0;
    border-right: 1px solid #ebeef5;

    .avatar-wrapper {
      position: relative;
      display: inline-block;

      .profile-avatar {
        border: 4px solid #fff;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
      }

      .avatar-uploader {
        .avatar-edit {
          position: absolute;
          right: 0;
          bottom: 0;
          width: 32px;
          height: 32px;
          background: #409eff;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
          cursor: pointer;
          border: 2px solid #fff;
          transition: all 0.3s;

          &:hover {
            background: #66b1ff;
          }
        }
      }
    }

    .user-info {
      .user-name {
        font-size: 20px;
        font-weight: 600;
        margin-bottom: 8px;
        color: #303133;
      }

      .user-role {
        margin-bottom: 8px;
      }

      .user-join {
        color: #909399;
        font-size: 14px;
      }
    }
  }

  .form-section {
    padding-left: 20px;

    :deep(.el-tabs__header) {
      margin-bottom: 0;
    }
  }
}
</style>
