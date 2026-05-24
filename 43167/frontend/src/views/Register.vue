<template>
  <div class="login-page">
    <el-card class="login-card">
      <template #header><h2>注册</h2></template>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role">
            <el-option label="买家" value="buyer" />
            <el-option label="卖家" value="seller" />
            <el-option label="鉴定师" value="appraiser" />
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-button type="primary" class="full" :loading="loading" @click="submit">注册</el-button>
        <div class="tip">已有账号？<router-link to="/login">去登录</router-link></div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({ username: '', password: '', role: 'buyer', email: '' })
const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名' }],
  password: [{ required: true, message: '请输入密码' }],
  role: [{ required: true, message: '请选择角色' }]
}
async function submit() {
  await formRef.value?.validate()
  try {
    loading.value = true
    await auth.register({ ...form })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page { display: flex; justify-content: center; align-items: center; min-height: 100vh; background: #f5f7fa; }
.login-card { width: 380px; }
.full { width: 100%; }
.tip { text-align: center; margin-top: 12px; }
</style>
