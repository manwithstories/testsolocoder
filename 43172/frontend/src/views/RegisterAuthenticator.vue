<template>
  <div class="register-container">
    <el-card class="register-card" shadow="hover">
      <div class="register-header">
        <div class="logo">🔍</div>
        <h2>鉴定师注册</h2>
        <p>提交资质，开启专业鉴定之路</p>
      </div>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" size="large" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" size="large" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" size="large" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" placeholder="请确认密码" size="large" show-password />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="form.real_name" placeholder="请输入真实姓名" size="large" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" size="large" />
        </el-form-item>
        <el-divider>资质信息</el-divider>
        <el-form-item label="资质证书编号" prop="license_number">
          <el-input v-model="form.license_number" placeholder="请输入资质证书编号" size="large" />
        </el-form-item>
        <el-form-item label="资质证书文件" prop="license_file">
          <el-upload
            :auto-upload="false"
            :limit="1"
            :on-change="handleFileChange"
            accept=".pdf,.jpg,.jpeg,.png"
          >
            <el-button size="large">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">支持 PDF、JPG、PNG 格式</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="专业资质" prop="certifications">
          <el-input
            v-model="form.certifications"
            type="textarea"
            :rows="3"
            placeholder="请输入您的专业资质、证书名称等"
          />
        </el-form-item>
        <el-form-item label="擅长领域" prop="specialties">
          <el-input
            v-model="form.specialties"
            type="textarea"
            :rows="2"
            placeholder="请输入您擅长鉴定的品类，如：包包、首饰、手表等"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            style="width: 100%"
            @click="handleRegister"
          >
            提交注册
          </el-button>
        </el-form-item>
        <div class="register-footer">
          <span>已有账号？</span>
          <el-link type="primary" @click="$router.push('/login')">立即登录</el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules, UploadFile } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

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
  role: 'authenticator',
  license_number: '',
  license_file: '',
  certifications: '',
  specialties: ''
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
  real_name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  license_number: [
    { required: true, message: '请输入资质证书编号', trigger: 'blur' }
  ],
  certifications: [
    { required: true, message: '请输入专业资质', trigger: 'blur' }
  ],
  specialties: [
    { required: true, message: '请输入擅长领域', trigger: 'blur' }
  ]
}

const handleFileChange = (file: UploadFile) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      form.license_file = e.target?.result as string
    }
    reader.readAsDataURL(file.raw)
  }
}

const handleRegister = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const success = await userStore.registerAuthenticator({
          username: form.username,
          email: form.email,
          password: form.password,
          phone: form.phone || undefined,
          real_name: form.real_name,
          role: form.role,
          license_number: form.license_number,
          license_file: form.license_file || 'pending_upload',
          certifications: form.certifications,
          specialties: form.specialties
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
  width: 560px;
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
