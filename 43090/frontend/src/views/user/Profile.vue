<template>
  <div class="profile">
    <el-card>
      <template #header>
        <h3>个人信息</h3>
      </template>
      <el-form :model="form" label-width="100px">
        <el-form-item label="头像">
          <el-avatar :size="80" :src="form.avatar">
            {{ form.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="form.username" disabled />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" disabled />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="账户余额">
          <span class="balance">¥{{ userStore.userInfo?.balance?.toFixed(2) || '0.00' }}</span>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveProfile">保存修改</el-button>
          <el-button @click="showPasswordDialog = true">修改密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-dialog v-model="showPasswordDialog" title="修改密码" width="400px">
      <el-form :model="passwordForm" label-width="100px">
        <el-form-item label="原密码">
          <el-input v-model="passwordForm.old_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.new_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认新密码">
          <el-input v-model="passwordForm.confirm_password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="changePassword">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { authApi } from '@/api'
import type { User } from '@/types'

const userStore = useUserStore()
const showPasswordDialog = ref(false)

const form = reactive<Partial<User>>({
  username: '',
  email: '',
  phone: '',
  nickname: '',
  avatar: '',
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const loadUserInfo = async () => {
  const res = await authApi.getUserInfo()
  Object.assign(form, res)
}

const saveProfile = async () => {
  try {
    await authApi.updateUserInfo(form)
    userStore.fetchUserInfo()
    ElMessage.success('保存成功')
  } catch (e) {}
}

const changePassword = async () => {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  try {
    await authApi.changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password,
    })
    ElMessage.success('密码修改成功')
    showPasswordDialog.value = false
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (e) {}
}

onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.balance {
  font-size: 20px;
  font-weight: bold;
  color: #f56c6c;
}
</style>
