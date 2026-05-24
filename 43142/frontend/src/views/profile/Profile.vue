<template>
  <MainLayout>
    <div class="page-container">
      <div class="page-header">
        <div class="page-title">个人中心</div>
        <div class="page-subtitle">管理您的个人信息</div>
      </div>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-card class="user-card">
            <div class="user-avatar">
              <el-avatar :size="100">{{ userStore.userName.charAt(0).toUpperCase() }}</el-avatar>
            </div>
            <h2 class="user-name">{{ userStore.userName }}</h2>
            <p class="user-email">{{ userStore.user?.email }}</p>
            <el-tag :type="getRoleType(userStore.userRole)" size="large">
              {{ getRoleText(userStore.userRole) }}
            </el-tag>
          </el-card>
        </el-col>

        <el-col :span="16">
          <el-card v-if="userStore.hasRole('applicant')">
            <template #header>
              <span class="card-title">个人信息</span>
            </template>
            <el-form
              v-if="profileForm"
              :model="profileForm"
              ref="formRef"
              label-width="100px"
              @submit.prevent="handleSaveProfile"
            >
              <el-form-item label="姓名">
                <el-input v-model="profileForm.full_name" />
              </el-form-item>
              <el-form-item label="手机号">
                <el-input v-model="profileForm.phone" />
              </el-form-item>
              <el-form-item label="性别">
                <el-select v-model="profileForm.gender" placeholder="请选择" style="width: 100%">
                  <el-option label="男" value="male" />
                  <el-option label="女" value="female" />
                </el-select>
              </el-form-item>
              <el-form-item label="出生日期">
                <el-date-picker
                  v-model="profileForm.birth_date"
                  type="date"
                  placeholder="选择日期"
                  style="width: 100%"
                />
              </el-form-item>
              <el-form-item label="所在地">
                <el-input v-model="profileForm.location" />
              </el-form-item>
              <el-form-item label="学历">
                <el-input v-model="profileForm.education" />
              </el-form-item>
              <el-form-item label="工作经验">
                <el-input v-model="profileForm.experience" type="textarea" :rows="3" />
              </el-form-item>
              <el-form-item label="技能">
                <el-input v-model="profileForm.skills" placeholder="多个技能用逗号分隔" />
              </el-form-item>
              <el-form-item label="个人简介">
                <el-input v-model="profileForm.summary" type="textarea" :rows="3" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="saving" @click="handleSaveProfile">
                  保存
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>

          <el-card v-if="userStore.hasRole('company')" class="mt-20">
            <template #header>
              <span class="card-title">公司信息</span>
            </template>
            <el-form
              v-if="companyForm"
              :model="companyForm"
              ref="companyFormRef"
              label-width="100px"
              @submit.prevent="handleSaveCompany"
            >
              <el-form-item label="公司名称">
                <el-input v-model="companyForm.name" />
              </el-form-item>
              <el-form-item label="所属行业">
                <el-input v-model="companyForm.industry" />
              </el-form-item>
              <el-form-item label="公司规模">
                <el-select v-model="companyForm.size" placeholder="请选择" style="width: 100%">
                  <el-option label="少于20人" value="0-20" />
                  <el-option label="20-99人" value="20-99" />
                  <el-option label="100-499人" value="100-499" />
                  <el-option label="500-999人" value="500-999" />
                  <el-option label="1000人以上" value="1000+" />
                </el-select>
              </el-form-item>
              <el-form-item label="公司地址">
                <el-input v-model="companyForm.address" />
              </el-form-item>
              <el-form-item label="公司网站">
                <el-input v-model="companyForm.website" />
              </el-form-item>
              <el-form-item label="公司简介">
                <el-input v-model="companyForm.description" type="textarea" :rows="4" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="saving" @click="handleSaveCompany">
                  保存
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/stores/user'
import { getProfile, updateProfile, getCompany, updateCompany } from '@/api/user'
import type { ApplicantProfile, Company } from '@/types'

const userStore = useUserStore()
const formRef = ref<FormInstance>()
const companyFormRef = ref<FormInstance>()
const saving = ref(false)

const profileForm = ref<ApplicantProfile | null>(null)
const companyForm = ref<Company | null>(null)

async function fetchProfile() {
  try {
    const res = await getProfile()
    if (res.data?.profile) {
      profileForm.value = { ...res.data.profile }
    } else {
      profileForm.value = {
        id: 0,
        user_id: userStore.userId,
        full_name: '',
        phone: '',
        gender: '',
        birth_date: '',
        location: '',
        education: '',
        experience: '',
        skills: '',
        summary: '',
        created_at: '',
        updated_at: ''
      }
    }
  } catch (e) {
    console.error(e)
  }
}

async function fetchCompany() {
  try {
    const res = await getCompany()
    if (res.data) {
      companyForm.value = { ...res.data }
    }
  } catch (e) {
    console.error(e)
  }
}

async function handleSaveProfile() {
  if (!profileForm.value) return

  saving.value = true
  try {
    await updateProfile(profileForm.value)
    ElMessage.success('保存成功')
    await userStore.fetchProfile()
  } catch (e) {
    // error handled
  } finally {
    saving.value = false
  }
}

async function handleSaveCompany() {
  if (!companyForm.value) return

  saving.value = true
  try {
    await updateCompany(companyForm.value)
    ElMessage.success('保存成功')
    await userStore.fetchProfile()
  } catch (e) {
    // error handled
  } finally {
    saving.value = false
  }
}

function getRoleType(role: string) {
  const types: Record<string, string> = {
    'admin': 'danger',
    'company': 'success',
    'applicant': 'primary'
  }
  return types[role] || ''
}

function getRoleText(role: string) {
  const texts: Record<string, string> = {
    'admin': '管理员',
    'company': '企业用户',
    'applicant': '求职者'
  }
  return texts[role] || role
}

void [formRef, companyFormRef]

onMounted(() => {
  if (userStore.hasRole('applicant')) {
    fetchProfile()
  }
  if (userStore.hasRole('company')) {
    fetchCompany()
  }
})
</script>

<style scoped>
.user-card {
  text-align: center;
}

.user-avatar {
  margin-bottom: 20px;
}

.user-name {
  margin: 10px 0;
}

.user-email {
  color: #909399;
  margin: 10px 0;
}

.card-title {
  font-weight: 600;
}

.mt-20 {
  margin-top: 20px;
}
</style>
