<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">评价管理</h1>
    </div>

    <div class="search-bar">
      <el-input
        v-model="filters.keyword"
        placeholder="搜索评价内容"
        clearable
        style="width: 250px"
        @keyup.enter="loadReviews"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filters.minRating" placeholder="最低评分" clearable style="width: 120px">
        <el-option label="1星" :value="1" />
        <el-option label="2星" :value="2" />
        <el-option label="3星" :value="3" />
        <el-option label="4星" :value="4" />
        <el-option label="5星" :value="5" />
      </el-select>
      <el-button type="primary" @click="loadReviews">搜索</el-button>
    </div>

    <el-table :data="reviews" v-loading="loading" style="width: 100%">
      <el-table-column label="用户" width="120">
        <template #default="{ row }">
          {{ row.is_anonymous ? '匿名' : row.user?.username }}
        </template>
      </el-table-column>
      <el-table-column label="车辆" min-width="150">
        <template #default="{ row }">
          {{ row.car?.brand }} {{ row.car?.model }}
        </template>
      </el-table-column>
      <el-table-column label="评分" width="150">
        <template #default="{ row }">
          <el-rate :model-value="row.rating" disabled size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="content" label="内容" min-width="300" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_hidden ? 'info' : 'success'" size="small">
            {{ row.is_hidden ? '已隐藏' : '显示中' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="likes" label="点赞" width="80" />
      <el-table-column label="时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button
            type="warning"
            link
            size="small"
            @click="toggleHidden(row)"
          >
            {{ row.is_hidden ? '显示' : '隐藏' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadReviews"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { reviewApi } from '@/api'
import type { Review } from '@/types'

const reviews = ref<Review[]>([])
const loading = ref(false)

const filters = reactive({
  keyword: '',
  minRating: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadReviews()
})

const loadReviews = async () => {
  loading.value = true
  try {
    const res = await reviewApi.getAllReviews({
      page: pagination.page,
      page_size: pagination.pageSize,
      min_rating: filters.minRating
    })
    reviews.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const toggleHidden = async (row: Review) => {
  try {
    await reviewApi.toggleReviewHidden(row.id, !row.is_hidden)
    ElMessage.success('操作成功')
    loadReviews()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}
</script>
