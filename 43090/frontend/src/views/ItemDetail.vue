<template>
  <div v-if="item" class="item-detail">
    <el-row :gutter="30">
      <el-col :xs="24" :md="12">
        <el-carousel height="400px" class="image-carousel">
          <el-carousel-item v-for="(img, index) in item.images" :key="index">
            <img :src="img.url" :alt="item.title" />
          </el-carousel-item>
          <el-carousel-item v-if="!item.images || item.images.length === 0">
            <img src="https://via.placeholder.com/600x400?text=No+Image" />
          </el-carousel-item>
        </el-carousel>
      </el-col>

      <el-col :xs="24" :md="12">
        <div class="item-info">
          <h1>{{ item.title }}</h1>
          <div class="item-meta">
            <el-tag type="success" v-if="item.status === 1">拍卖中</el-tag>
            <el-tag type="warning" v-else-if="item.status === 0">草稿</el-tag>
            <el-tag type="danger" v-else-if="item.status === 2">已下架</el-tag>
            <el-tag type="info" v-else-if="item.status === 3">已售出</el-tag>
            <span class="meta-item"><el-icon><View /></el-icon> {{ item.view_count }} 浏览</span>
            <span class="meta-item"><el-icon><Cpu /></el-icon> {{ item.bid_count }} 次出价</span>
          </div>

          <div class="price-section">
            <div class="current-price">
              <span class="label">当前价格</span>
              <span class="price">¥{{ currentPrice.toFixed(2) }}</span>
            </div>
            <div class="price-info">
              <span>起拍价: ¥{{ item.start_price.toFixed(2) }}</span>
              <span>保留价: ¥{{ item.reserve_price.toFixed(2) }}</span>
            </div>
          </div>

          <div v-if="item.status === 1 && userStore.isLoggedIn" class="bid-section">
            <el-divider content-position="left">出价</el-divider>
            <div class="bid-form">
              <el-input-number
                v-model="bidAmount"
                :min="minBidAmount"
                :step="10"
                :precision="2"
                size="large"
                class="bid-input"
              />
              <el-button type="primary" size="large" @click="placeBid" :loading="bidding">
                出价
              </el-button>
            </div>
            <div class="auto-bid">
              <el-checkbox v-model="enableAutoBid">设置自动出价</el-checkbox>
              <el-input-number
                v-if="enableAutoBid"
                v-model="maxAutoBid"
                :min="bidAmount"
                :step="10"
                :precision="2"
                size="small"
                placeholder="最高出价上限"
              />
            </div>
          </div>

          <div v-if="item.status === 1 && !userStore.isLoggedIn" class="login-tip">
            <el-alert type="info" show-icon>
              <template #title>
                请 <el-link type="primary" @click="$router.push('/login')">登录</el-link> 后参与出价
              </template>
            </el-alert>
          </div>

          <el-divider />

          <div class="seller-info">
            <h4>卖家信息</h4>
            <el-avatar :size="48" :src="item.seller?.avatar">
              {{ item.seller?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <span class="seller-name">{{ item.seller?.username }}</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="30" class="detail-section">
      <el-col :xs="24" :md="16">
        <el-card>
          <template #header>
            <h3>商品详情</h3>
          </template>
          <div class="description" v-html="item.description"></div>
          <div class="item-details">
            <div class="detail-item">
              <span class="label">分类:</span>
              <span>{{ item.category?.name }}</span>
            </div>
            <div class="detail-item">
              <span class="label">所在地:</span>
              <span>{{ item.location }}</span>
            </div>
            <div class="detail-item">
              <span class="label">成色:</span>
              <span>{{ item.condition }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="8">
        <el-card>
          <template #header>
            <h3>出价记录 ({{ bidHistory.length }})</h3>
          </template>
          <div class="bid-history">
            <div
              v-for="bid in bidHistory"
              :key="bid.id"
              class="bid-record"
              :class="{ 'winning': bid.is_winning }"
            >
              <div class="bid-user">
                <el-avatar :size="24">{{ bid.user?.username?.charAt(0)?.toUpperCase() }}</el-avatar>
                <span>{{ bid.user?.username }}</span>
                <el-tag v-if="bid.is_auto_bid === 1" size="small" type="info">自动</el-tag>
              </div>
              <div class="bid-amount">¥{{ bid.amount.toFixed(2) }}</div>
              <div class="bid-time">{{ formatTime(bid.created_at) }}</div>
            </div>
            <el-empty v-if="bidHistory.length === 0" description="暂无出价记录" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { AuctionItem, Bid } from '@/types'
import { itemApi, bidApi } from '@/api'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()

const item = ref<AuctionItem | null>(null)
const bidHistory = ref<Bid[]>([])
const currentPrice = ref(0)
const bidding = ref(false)
const enableAutoBid = ref(false)

const bidForm = reactive({
  bidAmount: 0,
  maxAutoBid: 0,
})

const bidAmount = computed({
  get: () => bidForm.bidAmount,
  set: (val) => (bidForm.bidAmount = val),
})

const maxAutoBid = computed({
  get: () => bidForm.maxAutoBid,
  set: (val) => (bidForm.maxAutoBid = val),
})

const minBidAmount = computed(() => {
  return currentPrice.value + 10
})

const fetchItem = async () => {
  const id = Number(route.params.id)
  try {
    const res = await itemApi.getDetail(id)
    item.value = res
    currentPrice.value = res.current_price
    bidForm.bidAmount = currentPrice.value + 10
    fetchBidHistory()
  } catch (e) {}
}

const fetchBidHistory = async () => {
  const id = Number(route.params.id)
  try {
    const res = await itemApi.getBidHistory(id, { page: 1, page_size: 20 })
    bidHistory.value = res.list
  } catch (e) {}
}

const placeBid = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  bidding.value = true
  try {
    const id = Number(route.params.id)
    await bidApi.placeBid(id, {
      amount: bidForm.bidAmount,
      max_auto_bid: enableAutoBid.value ? bidForm.maxAutoBid : undefined,
    })
    ElMessage.success('出价成功')
    fetchItem()
  } finally {
    bidding.value = false
  }
}

const formatTime = (time: string) => {
  return dayjs(time).format('MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchItem()
})
</script>

<style scoped>
.item-detail {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.image-carousel img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.item-info h1 {
  margin: 0 0 15px;
  font-size: 24px;
  color: #303133;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #909399;
  font-size: 14px;
}

.price-section {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.current-price {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 10px;
}

.current-price .label {
  color: #606266;
  font-size: 14px;
}

.current-price .price {
  color: #f56c6c;
  font-size: 32px;
  font-weight: bold;
}

.price-info {
  display: flex;
  gap: 20px;
  color: #909399;
  font-size: 13px;
}

.bid-form {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.bid-input {
  flex: 1;
}

.auto-bid {
  display: flex;
  align-items: center;
  gap: 10px;
}

.login-tip {
  margin-bottom: 20px;
}

.seller-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.seller-info h4 {
  margin: 0 0 10px;
  width: 100%;
}

.seller-name {
  font-size: 16px;
  color: #303133;
}

.detail-section {
  margin-top: 30px;
}

.description {
  color: #606266;
  line-height: 1.8;
  margin-bottom: 20px;
}

.item-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.detail-item {
  display: flex;
  gap: 10px;
}

.detail-item .label {
  color: #909399;
  min-width: 60px;
}

.bid-history {
  max-height: 400px;
  overflow-y: auto;
}

.bid-record {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.bid-record.winning {
  background: #f0f9eb;
  border-radius: 4px;
  padding: 10px;
}

.bid-user {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.bid-amount {
  color: #f56c6c;
  font-weight: bold;
}

.bid-time {
  color: #909399;
  font-size: 12px;
}
</style>
