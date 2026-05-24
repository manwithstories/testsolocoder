<template>
  <div class="page-container">
    <div class="page-header">
      <h2>个人中心</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="8">
        <div class="chart-container">
          <div class="profile-avatar">
            <el-avatar :size="100" :src="user?.avatar">
              {{ user?.username?.charAt(0) || 'U' }}
            </el-avatar>
          </div>
          <div class="profile-info">
            <h3>{{ user?.username }}</h3>
            <p>{{ user?.email }}</p>
            <p v-if="user?.phone">📱 {{ user?.phone }}</p>
          </div>
        </div>
      </el-col>
      <el-col :span="16">
        <div class="form-container">
          <h3 class="mb-20">基本信息</h3>
          <el-form :model="profileForm" ref="profileFormRef" label-width="80px">
            <el-form-item label="用户名">
              <el-input v-model="profileForm.username" />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="profileForm.phone" />
            </el-form-item>
            <el-form-item label="头像URL">
              <el-input v-model="profileForm.avatar" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveProfile">保存修改</el-button>
            </el-form-item>
          </el-form>
        </div>

        <div class="form-container mt-20">
          <h3 class="mb-20">修改密码</h3>
          <el-form :model="passwordForm" ref="passwordFormRef" label-width="80px">
            <el-form-item label="当前密码">
              <el-input v-model="passwordForm.oldPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="passwordForm.newPassword" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="changing" @click="changePassword">修改密码</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getProfile, updateProfile, changePassword as apiChangePassword } from '@/api/auth'

const userStore = useUserStore()
const user = computed(() => userStore.user)
const saving = ref(false)
const changing = ref(false)
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

const profileForm = reactive({
  username: '',
  phone: '',
  avatar: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: ''
})

onMounted(async () => {
  if (userStore.token) {
    await loadProfile()
  }
})

async function loadProfile() {
  try {
    const res = await getProfile()
    Object.assign(profileForm, {
      username: res.username,
      phone: res.phone || '',
      avatar: res.avatar || ''
    })
  } catch (e) {
    console.error(e)
  }
}

async function saveProfile() {
  saving.value = true
  try {
    await updateProfile(profileForm)
    ElMessage.success('保存成功')
    await userStore.fetchProfile()
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

async function changePassword() {
  if (!passwordForm.oldPassword || !passwordForm.newPassword) {
    ElMessage.warning('请填写完整的密码信息')
    return
  }
  if (passwordForm.newPassword.length < 6) {
    ElMessage.warning('新密码至少6位')
    return
  }
  changing.value = true
  try {
    await apiChangePassword(passwordForm)
    ElMessage.success('密码修改成功')
    Object.assign(passwordForm, { oldPassword: '', newPassword: '' })
  } catch (e) {
    console.error(e)
  } finally {
    changing.value = false
  }
}
</script>

<style lang="scss" scoped>
.profile-avatar {
  text-align: center;
  margin-bottom: 20px;
}

.profile-info {
  text-align: center;

  h3 {
    margin-bottom: 8px;
    color: #303133;
  }

  p {
    color: #909399;
    margin-bottom: 4px;
  }
}
</style>
