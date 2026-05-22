<template>
  <div class="profile">
    <div class="page-header">
      <h2 class="page-title">个人中心</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="8">
        <div class="card avatar-card">
          <el-avatar :size="120" :src="userStore.userInfo?.avatar">
            {{ userStore.userInfo?.realName?.charAt(0) || 'U' }}
          </el-avatar>
          <h2>{{ userStore.userInfo?.realName }}</h2>
          <p>{{ getRoleText(userStore.userInfo?.role) }}</p>
          <el-upload
            :show-file-list="false"
            :before-upload="beforeUpload"
            :http-request="handleAvatarUpload"
            accept="image/*"
          >
            <el-button type="primary" plain>更换头像</el-button>
          </el-upload>
        </div>
      </el-col>

      <el-col :span="16">
        <div class="card">
          <el-form
            ref="formRef"
            :model="form"
            :rules="rules"
            label-width="100px"
          >
            <el-form-item label="用户名">
              <el-input v-model="form.username" disabled />
            </el-form-item>
            <el-form-item label="真实姓名" prop="realName">
              <el-input v-model="form.realName" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="form.phone" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="handleSubmit">
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules, type UploadProps } from 'element-plus'
import { useUserStore } from '@/store/user'
import { updateProfile } from '@/api/auth'
import { uploadImage } from '@/api/property'

const userStore = useUserStore()
const loading = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  username: '',
  realName: '',
  phone: '',
  email: '',
  avatar: ''
})

const rules: FormRules = {
  realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }]
}

onMounted(() => {
  if (userStore.userInfo) {
    form.username = userStore.userInfo.username
    form.realName = userStore.userInfo.realName
    form.phone = userStore.userInfo.phone
    form.email = userStore.userInfo.email
    form.avatar = userStore.userInfo.avatar
  }
})

function getRoleText(role?: string) {
  switch (role) {
    case 'admin': return '管理员'
    case 'landlord': return '房东'
    case 'maintainer': return '维修人员'
    default: return '租户'
  }
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  return true
}

async function handleAvatarUpload(options: any) {
  try {
    const res = await uploadImage(options.file)
    form.avatar = res.data.url
    await updateProfile({ avatar: form.avatar })
    ElMessage.success('头像更新成功')
    await userStore.fetchUserInfo()
  } catch (error) {
    console.error('Upload failed:', error)
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await updateProfile({
          realName: form.realName,
          phone: form.phone,
          email: form.email
        })
        ElMessage.success('保存成功')
        await userStore.fetchUserInfo()
      } catch (error) {
        console.error('Submit failed:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.profile {
  padding: 0;
}

.avatar-card {
  text-align: center;
  padding: 30px;
}

.avatar-card .el-avatar {
  margin-bottom: 15px;
}

.avatar-card h2 {
  margin: 10px 0;
  font-size: 20px;
}

.avatar-card p {
  color: #909399;
  margin-bottom: 20px;
}
</style>
