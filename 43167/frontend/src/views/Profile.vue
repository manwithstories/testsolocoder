<template>
  <el-card>
    <template #header>个人中心</template>
    <div class="avatar-section" v-if="user">
      <el-avatar :size="80" :src="user.avatar">
        {{ user.username?.[0]?.toUpperCase() }}
      </el-avatar>
      <div class="avatar-actions">
        <el-upload
          :show-file-list="false"
          :before-upload="beforeAvatarUpload"
          :http-request="uploadAvatar"
          accept="image/*"
        >
          <el-button size="small">上传新头像</el-button>
        </el-upload>
      </div>
    </div>
    <el-descriptions :column="2" border v-if="user" style="margin-top: 16px">
      <el-descriptions-item label="用户名">{{ user.username }}</el-descriptions-item>
      <el-descriptions-item label="角色">{{ user.role }}</el-descriptions-item>
      <el-descriptions-item label="邮箱">{{ user.email }}</el-descriptions-item>
      <el-descriptions-item label="手机">{{ user.phone }}</el-descriptions-item>
      <el-descriptions-item label="信誉分">{{ user.credit_score }}</el-descriptions-item>
      <el-descriptions-item label="评价数">{{ user.review_count }}</el-descriptions-item>
    </el-descriptions>
    <el-divider />
    <h3>修改信息</h3>
    <el-form :model="form" label-width="80px">
      <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
      <el-form-item label="手机"><el-input v-model="form.phone" /></el-form-item>
      <el-form-item label="真实姓名"><el-input v-model="form.real_name" /></el-form-item>
      <el-form-item>
        <el-button type="primary" @click="save">保存</el-button>
      </el-form-item>
    </el-form>
    <el-divider />
    <h3>修改密码</h3>
    <el-form :model="pwdForm" label-width="80px">
      <el-form-item label="旧密码"><el-input v-model="pwdForm.old_password" type="password" /></el-form-item>
      <el-form-item label="新密码"><el-input v-model="pwdForm.new_password" type="password" /></el-form-item>
      <el-form-item>
        <el-button type="primary" @click="changePwd">修改</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type UploadRequestOptions } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import request from '@/utils/request'
import type { User } from '@/types'
import axios from 'axios'

const auth = useAuthStore()
const user = ref<User | null>(null)
const form = reactive({ email: '', phone: '', real_name: '' })
const pwdForm = reactive({ old_password: '', new_password: '' })

onMounted(async () => {
  user.value = await auth.me()
  form.email = user.value?.email || ''
  form.phone = user.value?.phone || ''
  form.real_name = user.value?.real_name || ''
})

function beforeAvatarUpload(file: File) {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB')
  }
  return isImage && isLt5M
}

async function uploadAvatar(options: UploadRequestOptions) {
  const formData = new FormData()
  formData.append('file', options.file)
  try {
    const token = localStorage.getItem('token')
    const res = await axios.post('/api/me/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': token ? `Bearer ${token}` : ''
      }
    })
    if (res.data?.code === 0) {
      ElMessage.success('头像上传成功')
      if (user.value) {
        user.value.avatar = res.data.data.url
      }
      auth.me()
    } else {
      ElMessage.error(res.data?.message || '上传失败')
    }
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '上传失败')
  }
}

async function save() {
  await request.put('/me', form)
  ElMessage.success('已保存')
  auth.me()
}
async function changePwd() {
  await request.put('/me/password', pwdForm)
  ElMessage.success('密码已修改')
  pwdForm.old_password = ''
  pwdForm.new_password = ''
}
</script>

<style scoped>
.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
}
.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
