<template>
  <div class="profile-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <div class="user-card">
            <el-avatar :size="80" :src="userStore.userInfo?.avatar">
              {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <h2>{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</h2>
            <el-tag v-if="userStore.isExpert" type="primary" size="large" effect="dark">
              <el-icon><Medal /></el-icon>
              认证专家
            </el-tag>
            <div class="user-stats">
              <div class="stat">
                <span class="num">Lv.{{ userStore.userInfo?.level || 1 }}</span>
                <span class="label">等级</span>
              </div>
              <div class="stat">
                <span class="num">{{ userStore.userInfo?.points || 0 }}</span>
                <span class="label">积分</span>
              </div>
            </div>
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>快捷操作</template>
          <div class="quick-actions">
            <router-link to="/user/points">
              <el-button style="width: 100%; margin-bottom: 8px">
                <el-icon><Coin /></el-icon>
                积分中心
              </el-button>
            </router-link>
            <router-link to="/user/favorites">
              <el-button style="width: 100%; margin-bottom: 8px">
                <el-icon><Star /></el-icon>
                我的收藏
              </el-button>
            </router-link>
            <router-link to="/user/notifications">
              <el-button style="width: 100%">
                <el-icon><Bell /></el-icon>
                消息通知
              </el-button>
            </router-link>
          </div>
        </el-card>

        <el-card v-if="!userStore.isExpert && (userStore.userInfo?.level || 0) >= 5" style="margin-top: 20px">
          <template #header>专家认证</template>
          <p>您的等级已达到申请专家认证的要求</p>
          <el-button type="primary" @click="showExpertDialog = true">
            申请成为专家
          </el-button>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>个人信息</template>
          <el-form :model="form" label-width="100px">
            <el-form-item label="昵称">
              <el-input v-model="form.nickname" placeholder="请输入昵称" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="form.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="form.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item label="个人简介">
              <el-input
                v-model="form.bio"
                type="textarea"
                :rows="4"
                placeholder="介绍一下自己..."
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveProfile">
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="showExpertDialog" title="申请专家认证" width="500px">
      <el-form :model="expertForm" label-width="100px">
        <el-form-item label="擅长领域">
          <el-input v-model="expertForm.field" placeholder="如：前端开发、人工智能等" />
        </el-form-item>
        <el-form-item label="个人描述">
          <el-input
            v-model="expertForm.description"
            type="textarea"
            :rows="4"
            placeholder="请描述您的专业背景和擅长领域..."
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showExpertDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitExpertApplication">
          提交申请
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user'

const userStore = useUserStore()

const saving = ref(false)
const submitting = ref(false)
const showExpertDialog = ref(false)

const form = reactive({
  nickname: '',
  email: '',
  phone: '',
  bio: ''
})

const expertForm = reactive({
  field: '',
  description: ''
})

const loadProfile = () => {
  if (userStore.userInfo) {
    form.nickname = userStore.userInfo.nickname || ''
    form.email = userStore.userInfo.email || ''
    form.phone = userStore.userInfo.phone || ''
    form.bio = userStore.userInfo.bio || ''
  }
}

const saveProfile = async () => {
  saving.value = true
  try {
    const res = await userApi.updateProfile(form)
    if (res.data) {
      userStore.setUserInfo(res.data)
    }
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

const submitExpertApplication = async () => {
  if (!expertForm.field || !expertForm.description) return
  submitting.value = true
  try {
    await userApi.applyExpert(expertForm)
    showExpertDialog.value = false
    expertForm.field = ''
    expertForm.description = ''
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped lang="scss">
.profile-page {
  .user-card {
    text-align: center;
    padding: 20px 0;

    h2 {
      margin: 12px 0;
    }

    .user-stats {
      display: flex;
      justify-content: center;
      gap: 40px;
      margin-top: 20px;

      .stat {
        text-align: center;

        .num {
          display: block;
          font-size: 24px;
          font-weight: bold;
          color: #409eff;
        }

        .label {
          display: block;
          font-size: 14px;
          color: #909399;
        }
      }
    }
  }

  .quick-actions {
    .el-button {
      justify-content: flex-start;
    }
  }
}
</style>
