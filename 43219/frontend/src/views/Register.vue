<template>
  <div class="login-wrap">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">注册账号</div>
      </template>
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="onSubmit">
        <el-form-item label="角色" prop="role">
          <el-radio-group v-model="form.role">
            <el-radio value="customer">服务客户</el-radio>
            <el-radio value="staff">家政人员</el-radio>
            <el-radio value="company">家政公司</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="3-64位" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="至少6位" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model="form.real_name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item v-if="form.role === 'staff'" label="所属公司ID" prop="company_id">
          <el-input-number v-model="form.company_id" :min="0" />
        </el-form-item>
        <el-form-item v-if="form.role === 'staff'" label="资质证书链接" prop="cert_files">
          <el-input v-model="form.cert_files" placeholder="多个以逗号分隔" />
        </el-form-item>
        <el-form-item v-if="form.role === 'staff'" label="健康证明链接" prop="health_files">
          <el-input v-model="form.health_files" />
        </el-form-item>
        <el-form-item v-if="form.role === 'staff'" label="技能标签" prop="skills">
          <el-input v-model="form.skills" placeholder="如:保洁,月嫂" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" style="width:100%">注 册</el-button>
        </el-form-item>
        <div class="row">
          <span class="muted">已有账号？</span>
          <router-link to="/login">去登录</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { register } from '../api/auth'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({
  role: 'customer' as 'customer' | 'staff' | 'company',
  username: '',
  password: '',
  confirmPassword: '',
  real_name: '',
  phone: '',
  company_id: 0,
  cert_files: '',
  health_files: '',
  skills: '',
})
const rules: FormRules = {
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  username: [{ required: true, min: 3, max: 64, message: '3-64位', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '至少6位', trigger: 'blur' }],
  confirmPassword: [{
    validator: (_r, v, cb) => v === form.password ? cb() : cb(new Error('两次密码不一致')),
    trigger: 'blur',
  }],
}

async function onSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    const payload: any = {
      role: form.role,
      username: form.username,
      password: form.password,
      real_name: form.real_name,
      phone: form.phone,
    }
    if (form.role === 'staff') {
      if (form.company_id > 0) payload.company_id = form.company_id
      if (form.cert_files) payload.cert_files = form.cert_files
      if (form.health_files) payload.health_files = form.health_files
      if (form.skills) payload.skills = form.skills
    }
    await register(payload)
    ElMessage.success('注册成功，请登录')
    router.replace('/login')
  } catch (e: any) {
    ElMessage.error(e.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-wrap {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}
.login-card { width: 100%; max-width: 480px; }
.login-header { text-align: center; font-weight: 600; font-size: 18px; }
</style>
