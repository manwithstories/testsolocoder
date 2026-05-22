<template>
  <div class="home-page">
    <div class="banner">
      <h1>知识问答社区</h1>
      <p>分享知识，解决问题，共同成长</p>
      <el-button type="primary" size="large" @click="goToAsk" v-if="userStore.isLoggedIn">
        我要提问
      </el-button>
      <el-button size="large" @click="goToLogin" v-else>
        立即加入
      </el-button>
    </div>

    <div class="main-section">
      <div class="left-content">
        <div class="section-header">
          <h2>热门问题</h2>
          <el-radio-group v-model="sortType" size="small">
            <el-radio-button value="hot">热门</el-radio-button>
            <el-radio-button value="newest">最新</el-radio-button>
            <el-radio-button value="unanswered">待回答</el-radio-button>
          </el-radio-group>
        </div>

        <el-skeleton v-if="loading" :rows="5" animated />

        <div v-else class="question-list">
          <div
            v-for="question in questions"
            :key="question.id"
            class="question-item"
            @click="goToDetail(question.id)"
          >
            <div class="question-title">
              <el-tag v-if="question.rewardPoints > 0" type="warning" size="small">
                悬赏 {{ question.rewardPoints }}
              </el-tag>
              <el-tag v-if="question.isSolved" type="success" size="small">已解决</el-tag>
              {{ question.title }}
            </div>
            <div class="question-meta">
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ question.views }}
              </span>
              <span class="meta-item">
                <el-icon><ChatDotRound /></el-icon>
                {{ question.answerCount }}
              </span>
              <span class="meta-item">
                <el-icon><Star /></el-icon>
                {{ question.collectCount }}
              </span>
              <span class="meta-item">
                <el-icon><User /></el-icon>
                {{ question.user?.nickname || question.user?.username }}
              </span>
              <span class="meta-item">
                <el-icon><Clock /></el-icon>
                {{ formatTime(question.createdAt) }}
              </span>
            </div>
            <div class="question-tags">
              <el-tag
                v-for="tag in question.tags"
                :key="tag.id"
                size="small"
                type="info"
                effect="plain"
              >
                {{ tag.name }}
              </el-tag>
            </div>
          </div>
        </div>

        <el-pagination
          v-if="total > 0"
          class="pagination"
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="prev, pager, next, total"
          @current-change="fetchQuestions"
        />
      </div>

      <div class="right-sidebar">
        <div class="sidebar-card">
          <h3>热门分类</h3>
          <div class="category-list">
            <div
              v-for="category in categories"
              :key="category.id"
              class="category-item"
              @click="filterByCategory(category.id)"
            >
              <span class="category-icon">{{ category.icon }}</span>
              <span>{{ category.name }}</span>
            </div>
          </div>
        </div>

        <div class="sidebar-card">
          <h3>热门标签</h3>
          <div class="tag-cloud">
            <el-tag
              v-for="tag in tags"
              :key="tag.id"
              class="tag-item"
              effect="plain"
              @click="filterByTag(tag.id)"
            >
              {{ tag.name }}
            </el-tag>
          </div>
        </div>

        <div class="sidebar-card" v-if="userStore.isLoggedIn">
          <h3>推荐问题</h3>
          <div class="recommend-list">
            <div
              v-for="question in recommendations"
              :key="question.id"
              class="recommend-item"
              @click="goToDetail(question.id)"
            >
              {{ question.title }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { questionApi, searchApi } from '@/api'
import type { Question, Category, Tag } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const questions = ref<Question[]>([])
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const recommendations = ref<Question[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const sortType = ref('hot')
const categoryFilter = ref<number | null>(null)
const tagFilter = ref<number | null>(null)

const fetchQuestions = async () => {
  loading.value = true
  try {
    const res = await questionApi.getQuestionList({
      page: page.value,
      pageSize: pageSize.value,
      sort: sortType.value,
      categoryId: categoryFilter.value || undefined,
      tagId: tagFilter.value || undefined
    })
    questions.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const res = await questionApi.getCategories()
    categories.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const fetchTags = async () => {
  try {
    const res = await questionApi.getTags()
    tags.value = (res.data || []).slice(0, 20)
  } catch (e) {
    console.error(e)
  }
}

const fetchRecommendations = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const res = await searchApi.getRecommendations({ limit: 5 })
    recommendations.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const goToAsk = () => {
  router.push('/questions/ask')
}

const goToLogin = () => {
  router.push('/login')
}

const goToDetail = (id: number) => {
  router.push(`/questions/${id}`)
}

const filterByCategory = (id: number) => {
  categoryFilter.value = id
  page.value = 1
  fetchQuestions()
}

const filterByTag = (id: number) => {
  tagFilter.value = id
  page.value = 1
  fetchQuestions()
}

const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

watch(sortType, () => {
  page.value = 1
  fetchQuestions()
})

onMounted(() => {
  fetchQuestions()
  fetchCategories()
  fetchTags()
  fetchRecommendations()
})
</script>

<style scoped lang="scss">
.home-page {
  .banner {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 60px 20px;
    text-align: center;
    border-radius: 8px;
    margin-bottom: 20px;

    h1 {
      font-size: 36px;
      margin-bottom: 16px;
    }

    p {
      font-size: 18px;
      margin-bottom: 24px;
      opacity: 0.9;
    }
  }

  .main-section {
    display: flex;
    gap: 20px;

    .left-content {
      flex: 1;

      .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        h2 {
          margin: 0;
        }
      }

      .question-list {
        .question-item {
          background: white;
          padding: 20px;
          border-radius: 8px;
          margin-bottom: 12px;
          cursor: pointer;
          transition: all 0.2s;

          &:hover {
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            transform: translateY(-2px);
          }

          .question-title {
            font-size: 16px;
            font-weight: 500;
            margin-bottom: 12px;
            display: flex;
            align-items: center;
            gap: 8px;
          }

          .question-meta {
            display: flex;
            gap: 16px;
            color: #909399;
            font-size: 13px;
            margin-bottom: 12px;

            .meta-item {
              display: flex;
              align-items: center;
              gap: 4px;
            }
          }

          .question-tags {
            display: flex;
            gap: 8px;
            flex-wrap: wrap;
          }
        }
      }

      .pagination {
        justify-content: center;
        margin-top: 20px;
      }
    }

    .right-sidebar {
      width: 280px;
      flex-shrink: 0;

      .sidebar-card {
        background: white;
        padding: 20px;
        border-radius: 8px;
        margin-bottom: 16px;

        h3 {
          margin: 0 0 16px 0;
          font-size: 16px;
        }

        .category-list {
          .category-item {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 8px 0;
            cursor: pointer;
            transition: color 0.2s;

            &:hover {
              color: #409eff;
            }

            .category-icon {
              font-size: 18px;
            }
          }
        }

        .tag-cloud {
          display: flex;
          flex-wrap: wrap;
          gap: 8px;

          .tag-item {
            cursor: pointer;
          }
        }

        .recommend-list {
          .recommend-item {
            padding: 8px 0;
            cursor: pointer;
            color: #606266;
            transition: color 0.2s;

            &:hover {
              color: #409eff;
            }
          }
        }
      }
    }
  }
}
</style>
