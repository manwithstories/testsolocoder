<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElForm, ElFormItem, ElInput, ElButton, ElCard, ElMessage, ElSelect, ElOption, ElSteps, ElStep } from 'element-plus'
import { User, Lock, Iphone, ArrowLeft } from '@element-plus/icons-vue'
import { register } from '@/api/auth'
import type { UserRole } from '@/types'

const router = useRouter()
const loading = ref(false)
const activeStep = ref(0)

const registerType = ref<UserRole>('hr')

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  real_name: '',
  phone: '',
  email: '',
  companyName: '',
  unifiedCode: '',
  agencyName: ''
})

const registerRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }, { min: 3, max: 50, message: '用户名长度在3-50个字符', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, max: 50, message: '密码长度在6-50个字符', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请确认密码', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }]
}

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const handleRegister = async () => {
  if (registerForm.password !== registerForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  loading.value = true
  try {
    await register({
      username: registerForm.username,
      password: registerForm.password,
      real_name: registerForm.real_name,
      phone: registerForm.phone,
      email: registerForm.email,
      role: registerType.value
    })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (error) {
    console.error('Registration failed:', error)
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<template>
  <div class="register-container">
    <div class="register-bg"></div>
    <ElCard class="register-card">
      <div class="register-header">
        <h1>用户注册</h1>
        <p>请填写以下信息完成注册</p>
      </div>

      <ElSteps :active="activeStep" finish-status="success" align-center class="register-steps">
        <ElStep title="选择角色" />
        <ElStep title="填写信息" />
        <ElStep title="完成注册" />
      </ElSteps>

      <div v-if="activeStep === 0" class="step-content">
        <h3>请选择您的身份</h3>
        <div class="role-cards">
          <div 
            class="role-card" 
            :class="{ active: registerType === 'hr' }"
            @click="registerType = 'hr'"
          >
            <div class="role-icon">🏢</div>
            <div class="role-title">企业HR</div>
            <div class="role-desc">管理企业员工、预算和预约分配</div>
          </div>
          <div 
            class="role-card"
            :class="{ active: registerType === 'agency' }"
            @click="registerType = 'agency'"
          >
            <div class="role-icon">🏥</div>
            <div class="role-title">体检机构</div>
            <div class="role-desc">管理套餐、预约和报告上传</div>
          </div>
          <div 
            class="role-card"
            :class="{ active: registerType === 'employee' }"
            @click="registerType = 'employee'"
          >
            <div class="role-icon">👤</div>
            <div class="role-title">企业员工</div>
            <div class="role-desc">预约体检、查看报告和健康档案</div>
          </div>
        </div>
        <ElButton type="primary" size="large" class="next-button" @click="activeStep = 1">
          下一步
        </ElButton>
      </div>

      <div v-else-if="activeStep === 1" class="step-content">
        <ElForm :model="registerForm" :rules="registerRules" label-width="100px" class="register-form">
          <ElFormItem label="用户名" prop="username">
            <ElInput v-model="registerForm.username" :prefix-icon="User" placeholder="请输入用户名" />
          </ElFormItem>
          
          <ElFormItem label="密码" prop="password">
            <ElInput v-model="registerForm.password" type="password" :prefix-icon="Lock" placeholder="请输入密码" show-password />
          </ElFormItem>

          <ElFormItem label="确认密码" prop="confirmPassword" :rules="[{ validator: validateConfirmPassword, trigger: 'blur' }]">
            <ElInput v-model="registerForm.confirmPassword" type="password" :prefix-icon="Lock" placeholder="请再次输入密码" show-password />
          </ElFormItem>

          <ElFormItem label="真实姓名" prop="real_name">
            <ElInput v-model="registerForm.real_name" placeholder="请输入真实姓名" />
          </ElFormItem>

          <ElFormItem label="手机号" prop="phone">
            <ElInput v-model="registerForm.phone" :prefix-icon="Iphone" placeholder="请输入手机号" />
          </ElFormItem>

          <ElFormItem label="邮箱">
            <ElInput v-model="registerForm.email" placeholder="请输入邮箱（选填）" />
          </ElFormItem>
        </ElForm>

        <div class="form-actions">
          <ElButton @click="activeStep = 0">
            <ElIcon><ArrowLeft /></ElIcon>
            上一步
          </ElButton>
          <ElButton type="primary" :loading="loading" @click="handleRegister">
            完成注册
          </ElButton>
        </div>
      </div>

      <div v-else class="step-content success-content">
        <div class="success-icon">✓</div>
        <h3>注册成功！</h3>
        <p>您的账号已创建成功，请使用用户名和密码登录</p>
        <ElButton type="primary" size="large" @click="goToLogin">
          立即登录
        </ElButton>
      </div>

      <div class="register-footer">
        <span>已有账号？</span>
        <a @click="goToLogin">返回登录</a>
      </div>
    </ElCard>
  </div>
</template>

<style scoped lang="scss">
.register-container {
  position: relative;
  width: 100vw;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.register-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  z-index: 0;
}

.register-card {
  position: relative;
  z-index: 1;
  width: 500px;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}

.register-header {
  text-align: center;
  margin-bottom: 20px;

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

.register-steps {
  margin-bottom: 30px;
}

.step-content {
  .role-cards {
    display: flex;
    gap: 15px;
    margin: 20px 0;

    .role-card {
      flex: 1;
      padding: 20px;
      border: 2px solid #e4e7ed;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;
      text-align: center;

      &:hover {
        border-color: #409eff;
        transform: translateY(-2px);
      }

      &.active {
        border-color: #409eff;
        background-color: #ecf5ff;
      }

      .role-icon {
        font-size: 40px;
        margin-bottom: 10px;
      }

      .role-title {
        font-size: 16px;
        font-weight: bold;
        margin-bottom: 8px;
        color: #303133;
      }

      .role-desc {
        font-size: 12px;
        color: #909399;
      }
    }
  }

  .next-button {
    width: 100%;
    margin-top: 20px;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    border: none;
  }

  .success-content {
    text-align: center;
    padding: 40px 0;

    .success-icon {
      width: 80px;
      height: 80px;
      line-height: 80px;
      border-radius: 50%;
      background-color: #67c23a;
      color: #fff;
      font-size: 40px;
      margin: 0 auto 20px;
    }

    h3 {
      margin: 0 0 10px 0;
      color: #303133;
    }

    p {
      color: #909399;
      margin-bottom: 30px;
    }
  }
}

.form-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.register-footer {
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
