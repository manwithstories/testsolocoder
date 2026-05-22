<template>
  <div>
    <div v-if="reviews.length > 0" class="review-list">
      <div v-for="review in reviews" :key="review.id" class="review-item">
        <div class="review-header">
          <el-avatar :size="40" :src="review.is_anonymous ? '' : review.user?.avatar">
            {{ review.is_anonymous ? '匿' : review.user?.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <div class="review-user">
            <span class="username">
              {{ review.is_anonymous ? '匿名用户' : review.user?.username }}
            </span>
            <el-rate :model-value="review.rating" disabled size="small" />
          </div>
          <span class="review-date">{{ formatDate(review.created_at) }}</span>
        </div>
        <div v-if="review.content" class="review-content">
          {{ review.content }}
        </div>
        <div class="review-actions">
          <el-button type="primary" link @click="handleLike(review.id)">
            <el-icon><Star /></el-icon>
            有用 ({{ review.likes }})
          </el-button>
        </div>
      </div>
    </div>

    <el-empty v-if="reviews.length === 0 && !loading" description="暂无评价" />

    <div v-if="reviews.length > 0" style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="prev, pager, next"
        @current-change="loadReviews"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Star } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { reviewApi } from '@/api'
import type { Review } from '@/types'

const props = defineProps<{
  carId: number
}>()

const reviews = ref<Review[]>([])
const loading = ref(false)

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
    const res = await reviewApi.getCarReviews(props.carId, {
      page: pagination.page,
      page_size: pagination.pageSize
    })
    reviews.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const handleLike = async (id: number) => {
  try {
    await reviewApi.likeReview(id)
    const review = reviews.value.find(r => r.id === id)
    if (review) review.likes++
  } catch {
    ElMessage.error('操作失败')
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}
</script>

<style scoped>
.review-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.review-item {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.review-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.review-user {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.username {
  font-weight: 500;
}

.review-date {
  margin-left: auto;
  color: #909399;
  font-size: 12px;
}

.review-content {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 12px;
}

.review-actions {
  display: flex;
  gap: 12px;
}
</style>
