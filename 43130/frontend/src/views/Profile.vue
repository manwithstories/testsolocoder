<template>
  <div class="profile-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="user-info">
            <el-avatar :size="80" :src="user?.avatar">
              {{ user?.full_name?.charAt(0) || 'U' }}
            </el-avatar>
            <h3>{{ user?.full_name }}</h3>
            <p class="user-role">{{ roleText(user?.role) }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
              <el-button type="primary" :icon="Edit" @click="editing = true" v-if="!editing">
                编辑
              </el-button>
            </div>
          </template>
          
          <el-form v-if="!editing" label-width="100px">
            <el-form-item label="用户名">{{ user?.username }}</el-form-item>
            <el-form-item label="邮箱">{{ user?.email }}</el-form-item>
            <el-form-item label="姓名">{{ user?.full_name }}</el-form-item>
            <el-form-item label="电话">{{ user?.phone || '-' }}</el-form-item>
            <el-form-item label="角色">{{ roleText(user?.role) }}</el-form-item>
            <el-form-item label="状态">
              <el-tag :type="user?.status === 'active' ? 'success' : 'danger'">
                {{ user?.status === 'active' ? '正常' : '禁用' }}
              </el-tag>
            </el-form-item>
            <el-form-item label="注册时间">{{ formatDate(user?.created_at) }}</el-form-item>
          </el-form>

          <el-form v-else ref="profileForm" :model="profileForm" :rules="profileRules" label-width="100px">
            <el-form-item label="姓名" prop="full_name">
              <el-input v-model="profileForm.full_name" />
            </el-form-item>
            <el-form-item label="电话" prop="phone">
              <el-input v-model="profileForm.phone" />
            </el-form-item>
            <el-form-item>
              <el-button @click="editing = false">取消</el-button>
              <el-button type="primary" :loading="saving" @click="saveProfile">保存</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card shadow="hover" style="margin-top: 20px">
          <template #header>
            <span>修改密码</span>
          </template>
          <el-form ref="passwordForm" :model="passwordForm" :rules="passwordRules" label-width="100px">
            <el-form-item label="当前密码" prop="old_password">
              <el-input v-model="passwordForm.old_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
              <el-input v-model="passwordForm.new_password" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="changing" @click="changePassword">修改密码</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { userApi } from '@/api/auth'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Edit } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { User } from '@/types'

const userStore = useUserStore()

const user = ref<User | null>(null)
const editing = ref(false)
const saving = ref(false)
const changing = ref(false)

const profileForm = reactive({
  full_name: '',
  phone: ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: ''
})

const profileRules: FormRules = {
  full_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const validateNewPassword = (_rule: any, value: string, callback: any) => {
  if (value.length < 6) {
    callback(new Error('密码长度至少6位'))
  } else {
    callback()
  }
}

const passwordRules: FormRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { validator: validateNewPassword, trigger: 'blur' }
  ]
}

const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

function roleText(role?: string) {
  const texts: Record<string, string> = {
    couple: '准新人',
    planner: '婚礼策划师',
    vendor: '婚庆服务商',
    admin: '管理员'
  }
  return texts[role || ''] || role || '-'
}

function formatDate(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
}

async function fetchUserInfo() {
  try {
    const res = await userApi.getProfile()
    user.value = res.data
    profileForm.full_name = res.data.full_name
    profileForm.phone = res.data.phone || ''
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  }
}

async function saveProfile() {
  if (!profileFormRef.value) return
  
  await profileFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        await userApi.updateProfile(profileForm)
        ElMessage.success('更新成功')
        editing.value = false
        await userStore.fetchUserInfo()
        fetchUserInfo()
      } catch (error: any) {
        ElMessage.error(error.message || '更新失败')
      } finally {
        saving.value = false
      }
    }
  })
}

async function changePassword() {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      changing.value = true
      try {
        await userApi.changePassword(passwordForm)
        ElMessage.success('密码修改成功')
        passwordForm.old_password = ''
        passwordForm.new_password = ''
      } catch (error: any) {
        ElMessage.error(error.message || '修改失败')
      } finally {
        changing.value = false
      }
    }
  })
}

onMounted(fetchUserInfo)
</script>

<style scoped>
.profile-page {
  padding: 0;
}

.user-info {
  text-align: center;
  padding: 20px;
}

.user-info .el-avatar {
  margin-bottom: 16px;
}

.user-info h3 {
  margin: 0 0 8px;
}

.user-role {
  color: #909399;
  margin: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
