<template>
  <div class="review-container">
    <el-card v-if="showReviewForm" class="review-form-card">
      <template #header>
        <div class="card-header">
          <el-icon><Edit /></el-icon>
          <span>发表评价</span>
        </div>
      </template>
      <el-form :model="reviewForm" label-width="80px">
        <el-form-item label="评分">
          <el-rate v-model="reviewForm.rating" :max="5" />
        </el-form-item>
        <el-form-item label="评价内容">
          <el-input
            v-model="reviewForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入您的评价"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitReview" :loading="submitting">
            提交评价
          </el-button>
          <el-button @click="showReviewForm = false">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="review-list-card">
      <template #header>
        <div class="card-header">
          <el-icon><ChatDotRound /></el-icon>
          <span>我的评价</span>
        </div>
      </template>

      <div v-if="loading" class="loading">
        <el-skeleton :rows="3" animated />
      </div>

      <div v-else class="review-list">
        <div v-for="review in reviews" :key="review.id" class="review-item">
          <div class="review-header">
            <el-avatar :size="40" :src="review.reviewer?.avatar">
              {{ review.reviewer?.nickname?.charAt(0) }}
            </el-avatar>
            <div class="reviewer-info">
              <span class="reviewer-name">{{ review.reviewer?.nickname }}</span>
              <el-rate :model-value="review.rating" disabled size="small" />
            </div>
            <span class="review-time">{{ formatTime(review.created_at) }}</span>
          </div>
          <div v-if="review.content" class="review-content">
            {{ review.content }}
          </div>
          <div class="review-type">
            <el-tag size="small" :type="review.review_type === 'courier' ? 'primary' : 'success'">
              {{ review.review_type === 'courier' ? '评价跑腿员' : '评价发布者' }}
            </el-tag>
          </div>
        </div>
        <el-empty v-if="reviews.length === 0" description="暂无评价" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Edit, ChatDotRound } from '@element-plus/icons-vue'
import { reviewApi } from '@/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import type { Review } from '@/types'

const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const reviews = ref<Review[]>([])
const showReviewForm = ref(!!route.query.order_id)

const reviewForm = reactive({
  order_id: Number(route.query.order_id) || 0,
  review_type: 'courier',
  rating: 5,
  content: ''
})

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

const fetchReviews = async () => {
  loading.value = true
  try {
    const res = await reviewApi.list({
      user_id: userStore.userInfo?.id,
      page_size: 50
    })
    if (res.code === 200) {
      reviews.value = res.data.items || []
    }
  } catch (error) {
    console.error('Failed to fetch reviews:', error)
  } finally {
    loading.value = false
  }
}

const submitReview = async () => {
  if (!reviewForm.order_id) {
    ElMessage.warning('请选择要评价的订单')
    return
  }
  if (!reviewForm.rating) {
    ElMessage.warning('请选择评分')
    return
  }

  submitting.value = true
  try {
    const res = await reviewApi.create(reviewForm)
    if (res.code === 200) {
      ElMessage.success('评价提交成功')
      showReviewForm.value = false
      fetchReviews()
    }
  } catch (error) {
    console.error('Failed to submit review:', error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchReviews()
})
</script>

<style lang="scss" scoped>
.review-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.review-form-card,
.review-list-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.review-list {
  .review-item {
    padding: 16px 0;
    border-bottom: 1px solid #ebeef5;

    &:last-child {
      border-bottom: none;
    }

    .review-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 12px;
    }

    .reviewer-info {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 4px;

      .reviewer-name {
        font-weight: 500;
      }
    }

    .review-time {
      color: #909399;
      font-size: 12px;
    }

    .review-content {
      color: #606266;
      line-height: 1.6;
      margin-bottom: 12px;
      padding: 12px;
      background: #f5f7fa;
      border-radius: 6px;
    }
  }
}
</style>
