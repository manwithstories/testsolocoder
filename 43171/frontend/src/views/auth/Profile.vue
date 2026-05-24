<template>
  <el-card>
    <template #header>
      <span>个人中心</span>
    </template>
    <el-descriptions :column="2" border v-if="userInfo">
      <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
      <el-descriptions-item label="角色">
        <el-tag :type="roleTagType">{{ roleText }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="昵称">{{ userInfo.nickname || '-' }}</el-descriptions-item>
      <el-descriptions-item label="认证状态">
        <el-tag :type="verifyTagType">{{ verifyText }}</el-tag>
      </el-descriptions-item>
    </el-descriptions>

    <el-divider />

    <h4>修改信息</h4>
    <el-form :model="form" label-width="80px" style="max-width: 500px">
      <el-form-item label="昵称">
        <el-input v-model="form.nickname" />
      </el-form-item>
      <el-form-item label="手机">
        <el-input v-model="form.phone" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleUpdate" :loading="loading">保存修改</el-button>
      </el-form-item>
    </el-form>

    <el-divider v-if="userInfo?.role === 'pilot'" />

    <div v-if="userInfo?.role === 'pilot'">
      <h4>飞手资质认证</h4>
      <el-form :model="verifyForm" label-width="100px" style="max-width: 500px">
        <el-form-item label="真实姓名">
          <el-input v-model="verifyForm.real_name" />
        </el-form-item>
        <el-form-item label="身份证号">
          <el-input v-model="verifyForm.id_card_no" />
        </el-form-item>
        <el-form-item label="执照编号">
          <el-input v-model="verifyForm.license_no" />
        </el-form-item>
        <el-form-item label="执照图片">
          <el-upload
            :action="'/api/upload/license'"
            :headers="{ Authorization: 'Bearer ' + token }"
            :show-file-list="false"
            :on-success="handleLicenseUpload"
          >
            <el-button>上传执照图片</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleVerify" :loading="verifyLoading">提交认证</el-button>
        </el-form-item>
      </el-form>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo)
const token = computed(() => userStore.token)

const loading = ref(false)
const verifyLoading = ref(false)

const form = reactive({
  nickname: '',
  phone: '',
  email: ''
})

const verifyForm = reactive({
  real_name: '',
  id_card_no: '',
  license_no: '',
  license_image: ''
})

const roleText = computed(() => {
  const map: Record<string, string> = { client: '客户', pilot: '飞手', owner: '机主' }
  return map[userInfo.value?.role || ''] || userInfo.value?.role
})

const roleTagType = computed(() => {
  const map: Record<string, string> = { client: '', pilot: 'success', owner: 'warning' }
  return map[userInfo.value?.role || ''] || ''
})

const verifyText = computed(() => {
  const map: Record<string, string> = { pending: '待审核', approved: '已通过', rejected: '已拒绝' }
  return map[userInfo.value?.verify_status || ''] || '未认证'
})

const verifyTagType = computed(() => {
  const map: Record<string, string> = { pending: 'warning', approved: 'success', rejected: 'danger' }
  return map[userInfo.value?.verify_status || ''] || 'info'
})

onMounted(() => {
  if (userInfo.value) {
    form.nickname = userInfo.value.nickname || ''
  }
})

async function handleUpdate() {
  loading.value = true
  try {
    await request.put('/user/profile', form)
    ElMessage.success('修改成功')
    userStore.fetchUserInfo()
  } catch (e: any) {
    ElMessage.error(e.message || '修改失败')
  } finally {
    loading.value = false
  }
}

function handleLicenseUpload(res: any) {
  verifyForm.license_image = res.data.url
  ElMessage.success('上传成功')
}

async function handleVerify() {
  verifyLoading.value = true
  try {
    await request.post('/user/verify-pilot', verifyForm)
    ElMessage.success('提交成功，请等待审核')
    userStore.fetchUserInfo()
  } catch (e: any) {
    ElMessage.error(e.message || '提交失败')
  } finally {
    verifyLoading.value = false
  }
}
</script>
