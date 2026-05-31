<template>
  <AppLayout>
    <div class="page">
      <h2>个人中心</h2>
      <el-form :model="form" label-width="100px" @submit.prevent="save">
        <el-form-item label="姓名"><el-input v-model="form.real_name" /></el-form-item>
        <el-form-item label="手机"><el-input v-model="form.phone" /></el-form-item>
        <el-form-item label="技能"><el-input v-model="form.skills" /></el-form-item>
        <el-form-item label="资质证书"><el-input v-model="form.cert_files" /></el-form-item>
        <el-form-item label="健康证明"><el-input v-model="form.health_files" /></el-form-item>
        <el-form-item label="身份证号"><el-input v-model="form.id_card" /></el-form-item>
        <el-form-item label="简介"><el-input v-model="form.intro" type="textarea" :rows="3" /></el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit">保存</el-button>
        </el-form-item>
      </el-form>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { getMe, updateMe, updateStaffCert } from '../../api/user'

const form = reactive({
  real_name: '', phone: '', skills: '',
  cert_files: '', health_files: '', id_card: '', intro: '',
})

onMounted(async () => {
  const res = await getMe()
  const d = (res.data as any).data
  form.real_name = d.real_name || ''
  form.phone = d.phone || ''
  form.skills = d.skills || ''
  if (d.profile) {
    form.cert_files = d.profile.cert_files || ''
    form.health_files = d.profile.health_files || ''
    form.id_card = d.profile.id_card || ''
    form.intro = d.profile.intro || ''
  }
})

async function save() {
  await updateMe({ real_name: form.real_name, phone: form.phone, skills: form.skills })
  await updateStaffCert({
    cert_files: form.cert_files, health_files: form.health_files,
    id_card: form.id_card, intro: form.intro,
  })
  ElMessage.success('保存成功')
}
</script>
