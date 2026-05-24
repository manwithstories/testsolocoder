<template>
  <div class="auth-container">
    <el-card class="auth-card">
      <template #header>
        <div class="auth-header">
          <el-icon :size="40"><Briefcase /></el-icon>
          <h2>注册招聘平台</h2>
        </div>
      </template>
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="用户类型" prop="role">
          <el-radio-group v-model="form.role">
            <el-radio value="applicant">求职者</el-radio>
            <el-radio value="company">企业</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" :prefix-icon="Message" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" :prefix-icon="Lock" show-password />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" placeholder="请再次输入密码" :prefix-icon="Lock" show-password />
        </el-form-item>

        <el-form-item v-if="form.role === 'applicant'" label="姓名" prop="full_name">
          <el-input v-model="form.full_name" placeholder="请输入姓名" />
        </el-form-item>

        <el-form-item v-if="form.role === 'applicant'" label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>

        <template v-if="form.role === 'company'">
          <el-form-item label="公司名称" prop="company.name">
            <el-input v-model="form.company!.name" placeholder="请输入公司名称" />
          </el-form-item>
          <el-form-item label="所属行业" prop="company.industry">
            <el-input v-model="form.company!.industry" placeholder="请输入所属行业" />
          </el-form-item>
          <el-form-item label="公司规模" prop="company.size">
            <el-select v-model="form.company!.size" placeholder="请选择公司规模" style="width: 100%">
              <el-option label="少于20人" value="0-20" />
              <el-option label="20-99人" value="20-99" />
              <el-option label="100-499人" value="100-499" />
              <el-option label="500-999人" value="500-999" />
              <el-option label="1000人以上" value="1000+" />
            </el-select>
          </el-form-item>
          <el-form-item label="公司地址" prop="company.address">
            <el-input v-model="form.company!.address" placeholder="请输入公司地址" />
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" class="submit-btn" :loading="loading" @click="handleRegister">
            注册
          </el-button>
        </el-form-item>
        <div class="auth-links">
          <span>已有账号？</span>
          <router-link to="/login">立即登录</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Message, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: '',
  password: '',
  confirmPassword: '',
  role: 'applicant' as 'applicant' | 'company',
  full_name: '',
  phone: '',
  company: {
    name: '',
    industry: '',
    size: '',
    address: ''
  }
})

const rules: FormRules = {
  role: [{ required: true, message: '请选择用户类型', trigger: 'change' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== form.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  full_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

async function handleRegister() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const registerData: any = {
          email: form.email,
          password: form.password,
          role: form.role
        }

        if (form.role === 'applicant') {
          registerData.full_name = form.full_name
          registerData.phone = form.phone
        } else {
          registerData.full_name = form.company.name
          registerData.company = { ...form.company }
        }

        await userStore.register(registerData)
        ElMessage.success('注册成功')
        router.push('/jobs')
      } catch (e) {
        // error handled in interceptor
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.auth-card {
  width: 500px;
  border-radius: 12px;
}

.auth-header {
  text-align: center;
  color: #409eff;
}

.auth-header h2 {
  margin-top: 10px;
  color: #303133;
}

.submit-btn {
  width: 100%;
}

.auth-links {
  text-align: center;
  margin-top: 10px;
}

.auth-links a {
  color: #409eff;
  text-decoration: none;
}
</style>
