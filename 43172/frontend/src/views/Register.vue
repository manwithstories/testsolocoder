<template>
  <div class="register-container">
    <el-card class="register-card" shadow="hover">
      <div class="register-header">
        <div class="logo">👜</div>
        <h2>注册账号</h2>
        <p>加入奢侈品交易平台</p>
      </div>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
            :prefix-icon="Message"
            size="large"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请确认密码"
            :prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号（选填）"
            :prefix-icon="Phone"
            size="large"
          />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input
            v-model="form.real_name"
            placeholder="请输入真实姓名（选填）"
            :prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item label="注册角色" prop="role">
          <el-radio-group v-model="form.role" size="large">
            <el-radio-button label="buyer">买家</el-radio-button>
            <el-radio-button label="seller">卖家</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            style="width: 100%"
            @click="handleRegister"
          >
            注册
          </el-button>
        </el-form-item>
        <div class="register-footer">
          <span>已有账号？</span>
          <el-link type="primary" @click="$router.push('/login')">立即登录</el-link>
        </div>
        <div class="register-footer">
          <el-link type="primary" @click="$router.push('/register/authenticator')">
            鉴定师注册入口
          </el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { User, Lock, Message, Phone } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  phone: '',
  real_name: '',
  role: 'buyer'
})

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 100, message: '用户名长度在 3 到 100 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择注册角色', trigger: 'change' }
  ]
}

const handleRegister = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const success = await userStore.register({
          username: form.username,
          email: form.email,
          password: form.password,
          phone: form.phone || undefined,
          real_name: form.real_name || undefined,
          role: form.role
        })
        if (success) {
          router.push('/login')
        }
      } catch (error) {
        console.error('Register error:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  padding: 40px 20px;
}

.register-card {
  width: 480px;
  border-radius: 12px;
  overflow: hidden;
  
  :deep(.el-card__body) {
    padding: 40px;
  }
}

.register-header {
  text-align: center;
  margin-bottom: 24px;
  
  .logo {
    font-size: 48px;
    margin-bottom: 12px;
  }
  
  h2 {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 8px;
    color: var(--primary-color);
  }
  
  p {
    color: var(--text-light);
    font-size: 14px;
  }
}

.register-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
  color: var(--text-secondary);
}
</style>
