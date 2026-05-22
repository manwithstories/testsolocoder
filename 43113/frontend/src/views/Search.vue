<template>
  <div class="search-page">
    <el-card>
      <template #header>
        <div class="search-header">
          <el-input
            v-model="keyword"
            placeholder="搜索问题..."
            size="large"
            @keyup.enter="search"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" size="large" @click="search">
            搜索
          </el-button>
        </div>
      </template>

      <div class="filters">
        <el-select v-model="categoryFilter" placeholder="选择分类" clearable @change="search">
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>

        <el-radio-group v-model="sortType" @change="search">
          <el-radio-button value="newest">最新</el-radio-button>
          <el-radio-button value="hot">热门</el-radio-button>
          <el-radio-button value="most_answers">回答最多</el-radio-button>
        </el-radio-group>
      </div>

      <el-table :data="questions" style="width: 100%">
        <el-table-column label="标题" prop="title">
          <template #default="{ row }">
            <span class="title-link" @click="goToDetail(row.id)">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            {{ row.category?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="浏览" prop="views" width="80" />
        <el-table-column label="回答" prop="answerCount" width="80" />
        <el-table-column label="提问者" width="150">
          <template #default="{ row }">
            {{ row.user?.nickname || row.user?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && questions.length === 0" description="暂无搜索结果" />

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="search"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { searchApi, questionApi } from '@/api'
import type { Question, Category } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const questions = ref<Question[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const keyword = ref(route.query.keyword as string || '')
const categoryFilter = ref<number | ''>('')
const sortType = ref('newest')

const fetchCategories = async () => {
  try {
    const res = await questionApi.getCategories()
    categories.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const search = async () => {
  loading.value = true
  try {
    const res = await searchApi.searchQuestions({
      keyword: keyword.value,
      categoryId: categoryFilter.value ? Number(categoryFilter.value) : undefined,
      page: page.value,
      pageSize: pageSize.value,
      sort: sortType.value
    })
    questions.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const goToDetail = (id: number) => {
  router.push(`/questions/${id}`)
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchCategories()
  if (keyword.value) {
    search()
  }
})
</script>

<style scoped lang="scss">
.search-page {
  .search-header {
    display: flex;
    gap: 12px;

    .el-input {
      flex: 1;
    }
  }

  .filters {
    display: flex;
    gap: 16px;
    margin: 20px 0;
  }

  .title-link {
    color: #409eff;
    cursor: pointer;

    &:hover {
      text-decoration: underline;
    }
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
