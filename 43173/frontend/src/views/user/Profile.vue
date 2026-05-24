<template>
  <div class="profile">
    <div class="page-header">
      <h2>个人中心</h2>
    </div>
    
    <el-form :model="form" label-width="100px" style="max-width: 500px;">
      <el-form-item label="头像">
        <el-avatar :size="80" :src="form.avatar">
          {{ form.nickname?.charAt(0) || 'U' }}
        </el-avatar>
      </el-form-item>
      <el-form-item label="用户名">
        <el-input v-model="form.username" disabled />
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="form.nickname" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="form.phone" />
      </el-form-item>
      <el-form-item label="简介">
        <el-input 
          v-model="form.bio" 
          type="textarea" 
          :rows="4"
          placeholder="介绍一下自己"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">
          保存
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { authApi } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const saving = ref(false)

const form = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: '',
  bio: '',
  avatar: ''
})

onMounted(() => {
  loadProfile()
})

async function loadProfile() {
  try {
    const profile = await authApi.getProfile()
    Object.assign(form, {
      username: profile.username,
      nickname: profile.nickname,
      email: profile.email,
      phone: profile.phone,
      bio: profile.bio,
      avatar: profile.avatar
    })
  } catch (e) {
    console.error(e)
  }
}

async function save() {
  saving.value = true
  try {
    await authApi.updateProfile({
      nickname: form.nickname,
      email: form.email,
      phone: form.phone,
      bio: form.bio
    })
    ElMessage.success('保存成功')
    userStore.updateUser({
      nickname: form.nickname,
      avatar: form.avatar
    })
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped lang="scss">
.profile {
  .el-avatar {
    margin-right: 12px;
  }
}
</style>
