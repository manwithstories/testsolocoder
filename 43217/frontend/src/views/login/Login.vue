<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElForm, ElFormItem, ElInput, ElButton, ElCard, ElMessage, ElSelect, ElOption } from 'element-plus'
import { User, Lock, Iphone, ArrowRight } from '@element-plus/icons-vue'
import type { UserRole } from '@/types'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const loginType = ref<UserRole>('hr')

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  loading.value = true
  try {
    await userStore.login(loginForm.username, loginForm.password)
    ElMessage.success('登录成功')
    
    const redirect = route.query.redirect as string
    if (redirect) {
      router.push(redirect)
    } else {
      const role = userStore.userRole
      if (role === 'hr') {
        router.push('/hr/employees')
      } else if (role === 'agency') {
        router.push('/agency/packages')
      } else if (role === 'employee') {
        router.push('/employee/appointments')
      } else {
        router.push('/dashboard')
      }
    }
  } catch (error) {
    console.error('Login failed:', error)
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<template>
  <div class="login-container">
    <div class="login-bg"></div>
    <ElCard class="login-card">
      <div class="login-header">
        <h1>健康管理平台</h1>
        <p>企业员工体检预约系统</p>
      </div>
      
      <div class="login-tabs">
        <ElSelect v-model="loginType" placeholder="选择登录类型" class="role-select">
          <ElOption label="企业HR登录" value="hr" />
          <ElOption label="体检机构登录" value="agency" />
          <ElOption label="员工登录" value="employee" />
        </ElSelect>
      </div>

      <ElForm
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
      >
        <ElFormItem prop="username">
          <ElInput
            v-model="loginForm.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
            size="large"
          />
        </ElFormItem>
        
        <ElFormItem prop="password">
          <ElInput
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </ElFormItem>

        <ElFormItem>
          <ElButton
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </ElButton>
        </ElFormItem>
      </ElForm>

      <div class="login-footer">
        <span>还没有账号？</span>
        <a @click="goToRegister">立即注册</a>
      </div>
    </ElCard>
  </div>
</template>

<style scoped lang="scss">
.login-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  z-index: 0;
}

.login-card {
  position: relative;
  z-index: 1;
  width: 400px;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;

  h1 {
    margin: 0 0 8px 0;
    font-size: 24px;
    color: #303133;
  }

  p {
    margin: 0;
    color: #909399;
    font-size: 14px;
  }
}

.login-tabs {
  margin-bottom: 20px;

  .role-select {
    width: 100%;
  }
}

.login-form {
  .login-button {
    width: 100%;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    border: none;
  }
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #909399;

  a {
    color: #409eff;
    cursor: pointer;
    margin-left: 5px;

    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
