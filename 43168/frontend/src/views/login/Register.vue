<template>
  <div class="register-container">
    <el-card class="register-card" shadow="always">
      <template #header>
        <div class="register-header">
          <h2 class="title">家具定制平台</h2>
          <p class="subtitle">创建新账号</p>
        </div>
      </template>
      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="rules"
        label-position="top"
        size="large"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="registerForm.username"
            placeholder="3-64 个字符"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="6-64 个字符"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="再次输入密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input
            v-model="registerForm.nickname"
            placeholder="请输入昵称（选填）"
            clearable
          />
        </el-form-item>
        <el-form-item label="手机" prop="phone">
          <el-input
            v-model="registerForm.phone"
            placeholder="请输入手机号"
            clearable
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="registerForm.email"
            placeholder="请输入邮箱（选填）"
            clearable
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-radio-group v-model="registerForm.role">
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
            @click="handleRegister"
          >
            注册
          </el-button>
        </el-form-item>
        <div class="footer-link">
          已有账号？
          <router-link to="/login" class="link">前往登录</router-link>
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
import { register, type LoginResult, type UserRole } from '@/api/user'
import { setToken, setUser } from '@/utils/request'

interface RegisterForm {
  username: string
  password: string
  confirmPassword: string
  nickname: string
  phone: string
  email: string
  role: UserRole
}

const router = useRouter()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)

const registerForm = reactive<RegisterForm>({
  username: '',
  password: '',
  confirmPassword: '',
  nickname: '',
  phone: '',
  email: '',
  role: 'owner'
})

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const validatePhone = (_rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请输入手机号'))
  } else if (!/^1[3-9]\d{9}$/.test(value)) {
    callback(new Error('请输入正确的手机号'))
  } else {
    callback()
  }
}

const rules: FormRules<RegisterForm> = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 64, message: '用户名长度为 3-64 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 64, message: '密码长度为 6-64 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
  ],
  phone: [{ required: true, validator: validatePhone, trigger: 'blur' }],
  email: [
    {
      type: 'email',
      message: '请输入正确的邮箱格式',
      trigger: 'blur'
    }
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

const handleRegister = async () => {
  if (!registerFormRef.value) return
  await registerFormRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      const res = (await register({
        username: registerForm.username,
        password: registerForm.password,
        role: registerForm.role,
        nickname: registerForm.nickname || undefined,
        phone: registerForm.phone || undefined,
        email: registerForm.email || undefined
      })) as unknown as LoginResult
      if (res?.token) {
        setToken(res.token)
        setUser(res.user)
        ElMessage.success('注册成功')
        router.replace(getHomePath(res.user?.role || registerForm.role))
      } else {
        ElMessage.success('注册成功，请登录')
        router.replace('/login')
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
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 460px;
  border-radius: 12px;
}

.register-header {
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
  .register-card {
    max-width: 100%;
  }
}
</style>
