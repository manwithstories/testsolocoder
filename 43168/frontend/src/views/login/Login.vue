<template>
  <div class="login-container">
    <el-card class="login-card" shadow="always">
      <template #header>
        <div class="login-header">
          <h2 class="title">家具定制平台</h2>
          <p class="subtitle">欢迎登录</p>
        </div>
      </template>
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="rules"
        label-position="top"
        size="large"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-radio-group v-model="loginForm.role">
            <el-radio-button label="manufacturer">厂商</el-radio-button>
            <el-radio-button label="designer">设计师</el-radio-button>
            <el-radio-button label="owner">业主</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            class="submit-btn"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
        <div class="footer-link">
          还没有账号？
          <router-link to="/register" class="link">立即注册</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { login, type LoginResult, type UserRole } from '@/api/user'
import { setToken, setUser } from '@/utils/request'

interface LoginForm {
  username: string
  password: string
  role: UserRole
}

const router = useRouter()
const loginFormRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive<LoginForm>({
  username: '',
  password: '',
  role: 'owner'
})

const rules: FormRules<LoginForm> = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 64, message: '用户名长度为 3-64 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 64, message: '密码长度为 6-64 个字符', trigger: 'blur' }
  ],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

const getHomePath = (role: UserRole): string => {
  switch (role) {
    case 'manufacturer':
      return '/manufacturer/home'
    case 'designer':
      return '/designer/home'
    case 'owner':
      return '/owner/home'
    default:
      return '/'
  }
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      const res = (await login({
        username: loginForm.username,
        password: loginForm.password
      })) as unknown as LoginResult
      if (res?.token) {
        setToken(res.token)
        setUser(res.user)
        ElMessage.success('登录成功')
        const path = getHomePath(res.user?.role || loginForm.role)
        router.replace(path)
      } else {
        ElMessage.error('登录失败：未返回 token')
      }
    } catch (e) {
      // error handled by interceptor
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  border-radius: 12px;
}

.login-header {
  text-align: center;
}

.title {
  margin: 0;
  font-size: 22px;
  color: #303133;
}

.subtitle {
  margin: 6px 0 0;
  color: #909399;
  font-size: 14px;
}

.submit-btn {
  width: 100%;
}

.footer-link {
  text-align: center;
  font-size: 14px;
  color: #606266;
}

.link {
  color: #409eff;
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .login-card {
    max-width: 100%;
  }
}
</style>
