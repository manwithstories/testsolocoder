<template>
  <div class="question-list-page">
    <div class="page-header">
      <h1>问题列表</h1>
      <div class="filters">
        <el-select v-model="categoryFilter" placeholder="选择分类" clearable @change="fetchQuestions">
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>

        <el-input
          v-model="keyword"
          placeholder="搜索问题"
          clearable
          @keyup.enter="fetchQuestions"
          @clear="fetchQuestions"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <el-radio-group v-model="sortType" @change="fetchQuestions">
          <el-radio-button value="newest">最新</el-radio-button>
          <el-radio-button value="hot">热门</el-radio-button>
          <el-radio-button value="unanswered">待回答</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <el-skeleton v-if="loading" :rows="5" animated />

    <div v-else class="question-list">
      <div
        v-for="question in questions"
        :key="question.id"
        class="question-item"
        @click="goToDetail(question.id)"
      >
        <div class="question-stats">
          <div class="stat">
            <span class="num">{{ question.views }}</span>
            <span class="label">浏览</span>
          </div>
          <div class="stat" :class="{ solved: question.isSolved }">
            <span class="num">{{ question.answerCount }}</span>
            <span class="label">回答</span>
          </div>
        </div>
        <div class="question-content">
          <div class="question-title">
            <el-tag v-if="question.rewardPoints > 0" type="warning" size="small">
              悬赏 {{ question.rewardPoints }}
            </el-tag>
            <el-tag v-if="question.isSolved" type="success" size="small">已解决</el-tag>
            {{ question.title }}
          </div>
          <div class="question-excerpt">
            {{ question.content.slice(0, 150) }}...
          </div>
          <div class="question-meta">
            <el-tag
              v-for="tag in question.tags"
              :key="tag.id"
              size="small"
              type="info"
              effect="plain"
            >
              {{ tag.name }}
            </el-tag>
            <span class="user-info">
              <el-avatar :size="20" :src="question.user?.avatar">
                {{ question.user?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              {{ question.user?.nickname || question.user?.username }}
            </span>
            <span class="time">{{ formatTime(question.createdAt) }}</span>
          </div>
        </div>
      </div>
    </div>

    <el-empty v-if="!loading && questions.length === 0" description="暂无问题" />

    <el-pagination
      v-if="total > 0"
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50]"
      layout="prev, pager, next, total, jumper"
      @current-change="fetchQuestions"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { questionApi } from '@/api'
import type { Question, Category } from '@/types'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

const router = useRouter()

const questions = ref<Question[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const sortType = ref('newest')
const categoryFilter = ref<number | ''>('')
const keyword = ref('')

const fetchQuestions = async () => {
  loading.value = true
  try {
    const res = await questionApi.getQuestionList({
      page: page.value,
      pageSize: pageSize.value,
      sort: sortType.value,
      categoryId: categoryFilter.value ? Number(categoryFilter.value) : undefined,
      keyword: keywordFilter.value
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

const goToDetail = (id: number) => {
  router.push(`/questions/${id}`)
}

const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

onMounted(() => {
  fetchQuestions()
  fetchCategories()
})
</script>

<style scoped lang="scss">
.question-list-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 16px;

    h1 {
      margin: 0;
    }

    .filters {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;

      .el-input {
        width: 200px;
      }
    }
  }

  .question-list {
    .question-item {
      display: flex;
      background: white;
      padding: 20px;
      border-radius: 8px;
      margin-bottom: 12px;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }

      .question-stats {
        display: flex;
        flex-direction: column;
        gap: 8px;
        margin-right: 20px;
        min-width: 60px;

        .stat {
          text-align: center;
          padding: 8px;
          border-radius: 4px;
          background: #f5f7fa;

          &.solved {
            background: #f0f9eb;
            color: #67c23a;
          }

          .num {
            display: block;
            font-size: 18px;
            font-weight: bold;
          }

          .label {
            display: block;
            font-size: 12px;
            color: #909399;
          }
        }
      }

      .question-content {
        flex: 1;

        .question-title {
          font-size: 16px;
          font-weight: 500;
          margin-bottom: 8px;
          display: flex;
          align-items: center;
          gap: 8px;
        }

        .question-excerpt {
          color: #606266;
          font-size: 14px;
          margin-bottom: 12px;
        }

        .question-meta {
          display: flex;
          align-items: center;
          gap: 12px;
          flex-wrap: wrap;

          .user-info {
            display: flex;
            align-items: center;
            gap: 4px;
            font-size: 13px;
            color: #606266;
          }

          .time {
            font-size: 13px;
            color: #909399;
          }
        }
      }
    }
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
