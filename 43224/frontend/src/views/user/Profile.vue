<template>
  <div class="page-container">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>基本信息</span>
          </template>
          <div class="avatar-section">
            <el-avatar :size="80" :src="userInfo?.avatar">
              {{ userInfo?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <h3>{{ userInfo?.real_name || userInfo?.username }}</h3>
            <el-tag :type="getRoleType(userInfo?.role)">{{ getRoleText(userInfo?.role) }}</el-tag>
          </div>
          <el-descriptions :column="1" border size="small" style="margin-top: 16px">
            <el-descriptions-item label="用户名">{{ userInfo?.username }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ userInfo?.email }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ userInfo?.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="评分">{{ userInfo?.rating?.toFixed(1) }}</el-descriptions-item>
            <el-descriptions-item label="已完成项目">{{ userInfo?.completed_count }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>
            <span>编辑资料</span>
          </template>

          <el-form :model="form" label-width="100px">
            <el-form-item label="真实姓名">
              <el-input v-model="form.real_name" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="form.email" />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="form.phone" />
            </el-form-item>
            <el-form-item v-if="userStore.role === 'translator'" label="语言对">
              <el-select v-model="form.language_pair_ids" multiple placeholder="选择语言对" style="width: 100%">
                <el-option
                  v-for="lp in languagePairs"
                  :key="lp.id"
                  :label="lp.display_name || lp.source_lang + ' - ' + lp.target_lang"
                  :value="lp.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item v-if="userStore.role === 'translator'" label="专业领域">
              <el-select v-model="form.expertise_tag_ids" multiple placeholder="选择专业领域" style="width: 100%">
                <el-option v-for="tag in expertiseTags" :key="tag.id" :label="tag.name" :value="tag.id" />
              </el-select>
            </el-form-item>
            <el-form-item v-if="userStore.role === 'translator'" label="日工作量">
              <el-input-number v-model="form.daily_capacity" :min="100" :max="10000" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSave" :loading="saving">保存修改</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>修改密码</span>
          </template>

          <el-form :model="passwordForm" label-width="100px" style="max-width: 500px">
            <el-form-item label="原密码">
              <el-input v-model="passwordForm.old_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="passwordForm.new_password" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleChangePassword" :loading="pwdLoading">修改密码</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { updateProfile, changePassword } from '@/api/auth'
import { listLanguagePairs, listExpertiseTags } from '@/api/project'

const userStore = useUserStore()
const userInfo = ref(userStore.userInfo)

const languagePairs = ref<any[]>([])
const expertiseTags = ref<any[]>([])

const saving = ref(false)
const pwdLoading = ref(false)

const form = reactive({
  real_name: '',
  email: '',
  phone: '',
  language_pair_ids: [] as number[],
  expertise_tag_ids: [] as number[],
  daily_capacity: 2000
})

const passwordForm = reactive({
  old_password: '',
  new_password: ''
})

async function loadOptions() {
  try {
    const [pairs, tags] = await Promise.all([
      listLanguagePairs(),
      listExpertiseTags()
    ])
    languagePairs.value = pairs || []
    expertiseTags.value = tags || []
  } catch (e) {
    console.error(e)
  }
}

function initForm() {
  if (userInfo.value) {
    form.real_name = userInfo.value.real_name || ''
    form.email = userInfo.value.email || ''
    form.phone = userInfo.value.phone || ''
    form.daily_capacity = userInfo.value.daily_capacity || 2000
    form.language_pair_ids = (userInfo.value.language_pairs || []).map((lp: any) => lp.id)
    form.expertise_tag_ids = (userInfo.value.expertise_tags || []).map((tag: any) => tag.id)
  }
}

async function handleSave() {
  saving.value = true
  try {
    await updateProfile(form)
    ElMessage.success('保存成功')
    userStore.fetchUserInfo()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function handleChangePassword() {
  if (!passwordForm.old_password || !passwordForm.new_password) {
    ElMessage.warning('请填写完整的密码信息')
    return
  }
  pwdLoading.value = true
  try {
    await changePassword(passwordForm)
    ElMessage.success('密码修改成功')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
  } catch (e: any) {
    ElMessage.error(e.message || '修改失败')
  } finally {
    pwdLoading.value = false
  }
}

function getRoleType(role: string) {
  const map: Record<string, string> = { client: '', translator: 'primary', pm: 'success', admin: 'danger' }
  return map[role] || ''
}

function getRoleText(role: string) {
  const map: Record<string, string> = { client: '客户', translator: '译者', pm: '项目经理', admin: '管理员' }
  return map[role] || role
}

onMounted(() => {
  loadOptions()
  initForm()
})
</script>

<style lang="scss" scoped>
.page-container {
  .avatar-section {
    text-align: center;

    h3 {
      margin: 16px 0 8px;
    }
  }
}
</style>
