<template>
  <div class="profile">
    <div class="page-header">
      <h2 class="page-title">个人中心</h2>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :md="8">
        <div class="card-container text-center">
          <el-avatar :size="100" :src="userStore.user?.avatar_url">
            {{ userStore.user?.full_name?.[0] }}
          </el-avatar>
          <h3 style="margin-top: 16px">{{ userStore.user?.full_name }}</h3>
          <el-tag :type="getRoleType(userStore.userRole)">
            {{ getRoleText(userStore.userRole) }}
          </el-tag>
          <p style="margin-top: 8px; color: rgba(0, 0, 0, 0.45)">
            {{ userStore.user?.email }}
          </p>
        </div>
      </el-col>

      <el-col :xs="24" :md="16">
        <div class="card-container">
          <el-tabs v-model="activeTab">
            <el-tab-pane label="基本信息" name="basic">
              <el-form
                ref="basicFormRef"
                :model="basicForm"
                label-width="100px"
                @submit.prevent="handleUpdateBasic"
              >
                <el-form-item label="用户名">
                  <el-input v-model="basicForm.username" disabled />
                </el-form-item>
                <el-form-item label="邮箱">
                  <el-input v-model="basicForm.email" disabled />
                </el-form-item>
                <el-form-item label="姓名">
                  <el-input v-model="basicForm.full_name" />
                </el-form-item>
                <el-form-item label="手机号">
                  <el-input v-model="basicForm.phone" />
                </el-form-item>
                <el-form-item label="公司">
                  <el-input v-model="basicForm.company" />
                </el-form-item>
                <el-form-item label="时区">
                  <el-select v-model="basicForm.timezone" style="width: 100%">
                    <el-option label="UTC+8 北京" value="Asia/Shanghai" />
                    <el-option label="UTC+0 伦敦" value="Europe/London" />
                    <el-option label="UTC-5 纽约" value="America/New_York" />
                    <el-option label="UTC+9 东京" value="Asia/Tokyo" />
                  </el-select>
                </el-form-item>
                <el-form-item label="地址">
                  <el-input v-model="basicForm.address" />
                </el-form-item>
                <el-form-item label="城市">
                  <el-input v-model="basicForm.city" />
                </el-form-item>
                <el-form-item label="国家">
                  <el-input v-model="basicForm.country" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="loading" @click="handleUpdateBasic">
                    保存修改
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane label="修改密码" name="password">
              <el-form
                ref="passwordFormRef"
                :model="passwordForm"
                :rules="passwordRules"
                label-width="100px"
                @submit.prevent="handleChangePassword"
              >
                <el-form-item label="当前密码" prop="current_password">
                  <el-input v-model="passwordForm.current_password" type="password" show-password />
                </el-form-item>
                <el-form-item label="新密码" prop="new_password">
                  <el-input v-model="passwordForm.new_password" type="password" show-password />
                </el-form-item>
                <el-form-item label="确认密码" prop="confirm_password">
                  <el-input v-model="passwordForm.confirm_password" type="password" show-password />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="loading" @click="handleChangePassword">
                    修改密码
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
import { updateProfileApi, changePasswordApi } from '@/api/user'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const activeTab = ref('basic')
const loading = ref(false)

const basicFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

const basicForm = reactive({
  username: '',
  email: '',
  full_name: '',
  phone: '',
  company: '',
  timezone: 'UTC',
  address: '',
  city: '',
  country: ''
})

const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const passwordRules: FormRules = {
  current_password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.new_password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const getRoleType = (role: string) => {
  const map: Record<string, string> = {
    admin: 'danger',
    owner: 'primary',
    tenant: 'success'
  }
  return map[role] || 'info'
}

const getRoleText = (role: string) => {
  const map: Record<string, string> = {
    admin: '管理员',
    owner: '船东',
    tenant: '租客'
  }
  return map[role] || role
}

const handleUpdateBasic = async () => {
  loading.value = true
  try {
    await updateProfileApi(basicForm)
    ElMessage.success('信息更新成功')
    await userStore.fetchProfile()
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    loading.value = false
  }
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await changePasswordApi({
          current_password: passwordForm.current_password,
          new_password: passwordForm.new_password
        })
        ElMessage.success('密码修改成功')
        passwordForm.current_password = ''
        passwordForm.new_password = ''
        passwordForm.confirm_password = ''
      } catch (error: any) {
        ElMessage.error(error.message || '修改失败')
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  if (userStore.user) {
    basicForm.username = userStore.user.username
    basicForm.email = userStore.user.email
    basicForm.full_name = userStore.user.full_name
    basicForm.phone = userStore.user.phone || ''
    basicForm.company = userStore.user.company || ''
    basicForm.timezone = userStore.user.timezone
    basicForm.address = userStore.user.address || ''
    basicForm.city = userStore.user.city || ''
    basicForm.country = userStore.user.country || ''
  }
})
</script>

<style lang="scss" scoped>
.profile {
  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 24px;

    &.text-center {
      text-align: center;
    }
  }
}
</style>
