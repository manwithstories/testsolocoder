<template>
  <div class="user-detail-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <div class="user-card">
            <el-avatar :size="100" :src="user?.avatar">
              {{ user?.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <h2>{{ user?.nickname || user?.username }}</h2>
            <el-tag v-if="user?.isExpert" type="primary" size="large" effect="dark">
              <el-icon><Medal /></el-icon>
              认证专家
            </el-tag>
            <p v-if="user?.bio" class="bio">{{ user.bio }}</p>
            <div class="user-stats">
              <div class="stat">
                <span class="num">Lv.{{ user?.level || 1 }}</span>
                <span class="label">等级</span>
              </div>
              <div class="stat">
                <span class="num">{{ user?.points || 0 }}</span>
                <span class="label">积分</span>
              </div>
            </div>
            <el-button
              v-if="userStore.isLoggedIn && userStore.userInfo?.id !== user?.id"
              type="primary"
              @click="toggleFollow"
            >
              {{ isFollowing ? '已关注' : '关注' }}
            </el-button>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>
            <el-tabs v-model="activeTab">
              <el-tab-pane label="提问" name="questions" />
              <el-tab-pane label="回答" name="answers" />
            </el-tabs>
          </template>

          <el-table v-if="activeTab === 'questions'" :data="userQuestions" style="width: 100%">
            <el-table-column label="标题" prop="title">
              <template #default="{ row }">
                <span class="title-link" @click="goToQuestion(row.id)">{{ row.title }}</span>
              </template>
            </el-table-column>
            <el-table-column label="浏览" prop="views" width="80" />
            <el-table-column label="回答" prop="answerCount" width="80" />
            <el-table-column label="时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>

          <el-empty v-if="activeTab === 'questions' && userQuestions.length === 0" description="暂无提问" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi, questionApi, followApi } from '@/api'
import type { User, Question } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const user = ref<User | null>(null)
const userQuestions = ref<Question[]>([])
const activeTab = ref('questions')
const isFollowing = ref(false)

const fetchUser = async () => {
  try {
    const id = Number(route.params.id)
    const res = await userApi.getUserById(id)
    user.value = res.data || null
  } catch (e) {
    console.error(e)
  }
}

const fetchUserQuestions = async () => {
  try {
    const id = Number(route.params.id)
    const res = await questionApi.getQuestionList({
      userId: id,
      page: 1,
      pageSize: 20
    })
    userQuestions.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

const checkFollowStatus = async () => {
  if (!userStore.isLoggedIn || !user.value) return
  try {
    const res = await followApi.isFollowing({
      followingType: 'user',
      followingId: user.value.id
    })
    isFollowing.value = res.data?.isFollowing || false
  } catch (e) {
    console.error(e)
  }
}

const toggleFollow = async () => {
  if (!user.value) return
  try {
    if (isFollowing.value) {
      await followApi.unfollow({
        followingType: 'user',
        followingId: user.value.id
      })
    } else {
      await followApi.follow({
        followingType: 'user',
        followingId: user.value.id
      })
    }
    isFollowing.value = !isFollowing.value
  } catch (e) {
    console.error(e)
  }
}

const goToQuestion = (id: number) => {
  router.push(`/questions/${id}`)
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchUser()
  fetchUserQuestions()
  checkFollowStatus()
})
</script>

<style scoped lang="scss">
.user-detail-page {
  .user-card {
    text-align: center;
    padding: 20px 0;

    h2 {
      margin: 12px 0;
    }

    .bio {
      color: #606266;
      margin: 12px 0;
    }

    .user-stats {
      display: flex;
      justify-content: center;
      gap: 40px;
      margin: 20px 0;

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

  .title-link {
    color: #409eff;
    cursor: pointer;

    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
