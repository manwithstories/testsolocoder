<template>
  <div class="reviews-page">
    <div class="page-header">
      <h2 class="page-title">评价管理</h2>
    </div>

    <div class="reviews-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="收到的评价" name="received">
          <div class="review-list" v-loading="loading">
            <div
              v-for="review in receivedReviews"
              :key="review.id"
              class="review-card card"
            >
              <div class="review-header">
                <el-avatar :size="40" :src="review.reviewer?.avatar">
                  {{ review.reviewer?.nickname?.charAt(0) || review.reviewer?.username?.charAt(0) || 'U' }}
                </el-avatar>
                <div class="reviewer-info">
                  <span class="reviewer-name">{{ review.reviewer?.nickname || review.reviewer?.username }}</span>
                  <span class="review-time">{{ formatTime(review.createdAt) }}</span>
                </div>
                <el-rate :model-value="review.rating" disabled />
              </div>
              <div class="review-content" v-if="review.content">
                {{ review.content }}
              </div>
              <div class="review-scores">
                <span v-if="review.qualityScore">商品质量：<el-rate :model-value="review.qualityScore" disabled size="small" /></span>
                <span v-if="review.serviceScore">服务态度：<el-rate :model-value="review.serviceScore" disabled size="small" /></span>
              </div>
            </div>
            <div class="empty-state" v-if="!loading && receivedReviews.length === 0">
              <el-empty description="暂无收到的评价" />
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="发出的评价" name="sent">
          <div class="review-list" v-loading="loading">
            <div
              v-for="review in sentReviews"
              :key="review.id"
              class="review-card card"
            >
              <div class="review-header">
                <el-avatar :size="40" :src="review.reviewee?.avatar">
                  {{ review.reviewee?.nickname?.charAt(0) || review.reviewee?.username?.charAt(0) || 'U' }}
                </el-avatar>
                <div class="reviewer-info">
                  <span class="reviewer-name">{{ review.reviewee?.nickname || review.reviewee?.username }}</span>
                  <span class="review-time">{{ formatTime(review.createdAt) }}</span>
                </div>
                <el-rate :model-value="review.rating" disabled />
              </div>
              <div class="review-content" v-if="review.content">
                {{ review.content }}
              </div>
            </div>
            <div class="empty-state" v-if="!loading && sentReviews.length === 0">
              <el-empty description="暂无发出的评价" />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { reviewApi } from '@/api'
import { useUserStore } from '@/store/user'
import type { Review } from '@/types'

const userStore = useUserStore()

const loading = ref(false)
const activeTab = ref('received')
const receivedReviews = ref<Review[]>([])
const sentReviews = ref<Review[]>([])

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

async function fetchReviews() {
  loading.value = true
  try {
    if (userStore.userInfo?.id) {
      const [receivedRes, sentRes] = await Promise.all([
        reviewApi.getList({
          revieweeId: userStore.userInfo.id,
          page: 1,
          pageSize: 20
        }),
        reviewApi.getList({
          page: 1,
          pageSize: 20
        })
      ])
      receivedReviews.value = receivedRes.data
      sentReviews.value = sentRes.data.filter(r => r.reviewerId === userStore.userInfo?.id)
    }
  } catch (error) {
    console.error('Failed to fetch reviews:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchReviews()
})
</script>

<style lang="scss" scoped>
.reviews-page {
  .review-list {
    .review-card {
      margin-bottom: 16px;

      .review-header {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 12px;

        .reviewer-info {
          flex: 1;

          .reviewer-name {
            display: block;
            font-weight: 500;
            margin-bottom: 4px;
          }

          .review-time {
            font-size: 12px;
            color: var(--text-lighter-color);
          }
        }
      }

      .review-content {
        padding: 12px;
        background: #f5f7fa;
        border-radius: 4px;
        margin-bottom: 12px;
        line-height: 1.6;
      }

      .review-scores {
        display: flex;
        gap: 20px;
        font-size: 13px;
        color: var(--text-lighter-color);
      }
    }
  }
}
</style>
