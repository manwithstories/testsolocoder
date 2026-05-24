<template>
  <div class="page-container">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h1 class="page-title">用户评价</h1>
    </div>

    <div v-loading="loading" class="card reviews-page">
      <div v-if="reviews.length > 0" class="reviews-list">
        <div v-for="review in reviews" :key="review.id" class="review-item">
          <div class="review-header">
            <el-avatar :size="40">
              {{ review.fromUserId }}
            </el-avatar>
            <div class="review-user-info">
              <el-rate :model-value="review.rating" disabled />
              <span class="review-date">{{ formatDate(review.createdAt) }}</span>
            </div>
          </div>
          <div v-if="review.content" class="review-content">
            {{ review.content }}
          </div>
        </div>
      </div>

      <div v-if="!loading && reviews.length === 0" class="empty-state">
        <el-icon><ChatDotRound /></el-icon>
        <p>暂无评价</p>
      </div>

      <div v-if="total > 0" class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[5, 10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadReviews"
          @current-change="loadReviews"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { reviewApi } from '@/api/order'
import type { Review } from '@/types'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const reviews = ref<Review[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const equipmentId = ref(0)

onMounted(() => {
  equipmentId.value = parseInt(route.params.equipmentId as string)
  loadReviews()
})

async function loadReviews() {
  loading.value = true
  try {
    const response = await reviewApi.getEquipmentReviews(equipmentId.value, page.value, pageSize.value)
    reviews.value = response.data.data || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load reviews:', error)
    ElMessage.error('加载评价列表失败')
  } finally {
    loading.value = false
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.review-item {
  padding: 20px 0;
  border-bottom: 1px solid #ebeef5;
}

.review-item:last-child {
  border-bottom: none;
}

.review-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.review-user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.review-date {
  font-size: 12px;
  color: #909399;
}

.review-content {
  color: #606266;
  line-height: 1.6;
  padding-left: 52px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
