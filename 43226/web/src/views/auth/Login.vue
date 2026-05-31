<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-left">
        <div class="auth-brand">
          <el-icon size="48" color="#fff"><Museum /></el-icon>
          <h1>在线博物馆</h1>
          <p>探索艺术与历史的数字世界</p>
        </div>
        <div class="auth-features">
          <div class="feature-item">
            <el-icon size="24"><Tickets /></el-icon>
            <span>在线预约展览</span>
          </div>
          <div class="feature-item">
            <el-icon size="24"><Picture /></el-icon>
            <span>海量藏品展示</span>
          </div>
          <div class="feature-item">
            <el-icon size="24"><Microphone /></el-icon>
            <span>智能语音导览</span>
          </div>
        </div>
      </div>
      <div class="auth-right">
        <div class="auth-form">
          <h2>欢迎回来</h2>
          <p class="auth-subtitle">请登录您的账号</p>
          <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent>
            <el-form-item prop="username">
              <el-input
                v-model="form.username"
                placeholder="用户名"
                size="large"
                :prefix-icon="User"
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="form.password"
                type="password"
                placeholder="密码"
                size="large"
                :prefix-icon="Lock"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                size="large"
                style="width: 100%"
                :loading="loading"
                @click="handleLogin"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>
          <div class="auth-footer">
            <span>还没有账号？</span>
            <router-link to="/register">立即注册</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  try {
    loading.value = true
    await userStore.login(form.username, form.password)
    ElMessage.success('登录成功')
    const redirect = route.query.redirect as string || '/dashboard'
    router.push(redirect)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.auth-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.auth-container {
  display: flex;
  max-width: 1000px;
  width: 100%;
  height: 600px;
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.auth-left {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60px 50px;
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: space-between;

  .auth-brand {
    h1 {
      font-size: 36px;
      margin: 16px 0 8px;
    }
    p {
      opacity: 0.9;
      font-size: 16px;
    }
  }

  .auth-features {
    display: flex;
    flex-direction: column;
    gap: 20px;

    .feature-item {
      display: flex;
      align-items: center;
      gap: 12px;
      font-size: 16px;
    }
  }
}

.auth-right {
  width: 420px;
  padding: 60px 50px;
  display: flex;
  align-items: center;

  .auth-form {
    width: 100%;

    h2 {
      font-size: 28px;
      margin-bottom: 8px;
    }

    .auth-subtitle {
      color: #909399;
      margin-bottom: 32px;
    }
  }

  .auth-footer {
    text-align: center;
    margin-top: 24px;
    color: #909399;

    a {
      color: #409eff;
      margin-left: 4px;
    }
  }
}
</style>
