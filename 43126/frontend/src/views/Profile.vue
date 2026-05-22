<template>
  <div class="profile">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="profile-card">
          <div class="avatar-section">
            <el-avatar :size="80" class="profile-avatar">
              {{ userStore.userInfo?.real_name?.charAt(0)?.toUpperCase() || userStore.userInfo?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <h3>{{ userStore.userInfo?.real_name || userStore.userInfo?.username }}</h3>
            <el-tag :type="getRoleType(userStore.role)">{{ getRoleText(userStore.role) }}</el-tag>
          </div>
          <div class="info-section">
            <el-descriptions :column="1" border>
              <el-descriptions-item label="用户名">{{ userStore.userInfo?.username }}</el-descriptions-item>
              <el-descriptions-item label="邮箱">{{ userStore.userInfo?.email }}</el-descriptions-item>
              <el-descriptions-item label="部门">{{ userStore.userInfo?.department || '-' }}</el-descriptions-item>
              <el-descriptions-item label="电话">{{ userStore.userInfo?.phone || '-' }}</el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>编辑个人信息</span>
          </template>
          <el-form ref="formRef" :model="profileForm" :rules="rules" label-width="100px" style="max-width: 500px">
            <el-form-item label="真实姓名" prop="real_name">
              <el-input v-model="profileForm.real_name" placeholder="请输入真实姓名" />
            </el-form-item>
            <el-form-item label="电话" prop="phone">
              <el-input v-model="profileForm.phone" placeholder="请输入电话" />
            </el-form-item>
            <el-form-item label="部门" prop="department">
              <el-input v-model="profileForm.department" placeholder="请输入部门" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveProfile">保存</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { api } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const formRef = ref<FormInstance>()
const saving = ref(false)

const profileForm = reactive({
  real_name: '',
  phone: '',
  department: ''
})

const rules: FormRules = {
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }]
}

onMounted(() => {
  if (userStore.userInfo) {
    profileForm.real_name = userStore.userInfo.real_name || ''
    profileForm.phone = userStore.userInfo.phone || ''
    profileForm.department = userStore.userInfo.department || ''
  }
})

async function saveProfile() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        const res: any = await api.updateProfile(profileForm)
        ElMessage.success('保存成功')
        userStore.fetchProfile()
      } catch (e: any) {
        console.error(e)
      } finally {
        saving.value = false
      }
    }
  })
}

function getRoleType(role: string) {
  const map: Record<string, string> = { admin: 'danger', space_admin: 'warning', user: '' }
  return map[role] || ''
}

function getRoleText(role: string) {
  const map: Record<string, string> = { admin: '管理员', space_admin: '空间管理员', user: '普通用户' }
  return map[role] || role
}
</script>

<style scoped>
.profile-card {
  border-radius: 8px;
}

.avatar-section {
  text-align: center;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 20px;
}

.profile-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 32px;
  margin-bottom: 12px;
}

.avatar-section h3 {
  margin: 0 0 8px;
  color: #303133;
}

.info-section {
  margin-top: 20px;
}
</style>
