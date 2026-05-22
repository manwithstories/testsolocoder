<template>
  <el-card class="item-card" shadow="hover" @click="$router.push(`/items/${item.id}`)">
    <div class="item-image">
      <img :src="getMainImage" :alt="item.title" />
      <div class="item-status" :class="statusClass">{{ statusText }}</div>
    </div>
    <div class="item-content">
      <h3 class="item-title">{{ item.title }}</h3>
      <div class="item-info">
        <div class="price">
          <span class="label">当前价</span>
          <span class="value">¥{{ item.current_price.toFixed(2) }}</span>
        </div>
        <div class="bids">
          <span>{{ item.bid_count }} 次出价</span>
        </div>
      </div>
      <div class="item-footer">
        <span class="views"><el-icon><View /></el-icon> {{ item.view_count }}</span>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AuctionItem } from '@/types'

const props = defineProps<{
  item: AuctionItem
}>()

const getMainImage = computed(() => {
  if (props.item.images && props.item.images.length > 0) {
    return props.item.images[0].url
  }
  return 'https://via.placeholder.com/300x200?text=No+Image'
})

const statusClass = computed(() => {
  const map: Record<number, string> = {
    0: 'status-draft',
    1: 'status-online',
    2: 'status-offline',
    3: 'status-sold',
  }
  return map[props.item.status] || ''
})

const statusText = computed(() => {
  const map: Record<number, string> = {
    0: '草稿',
    1: '拍卖中',
    2: '已下架',
    3: '已售出',
  }
  return map[props.item.status] || ''
})
</script>

<style scoped>
.item-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
  margin-bottom: 20px;
}

.item-card:hover {
  transform: translateY(-5px);
}

.item-image {
  position: relative;
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 8px 8px 0 0;
  margin: -20px -20px 15px;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-status {
  position: absolute;
  top: 10px;
  right: 10px;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  color: white;
}

.status-online {
  background: #67c23a;
}

.status-draft {
  background: #909399;
}

.status-offline {
  background: #f56c6c;
}

.status-sold {
  background: #409eff;
}

.item-title {
  font-size: 16px;
  margin: 0 0 12px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
  height: 44px;
}

.item-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.price .label {
  font-size: 12px;
  color: #909399;
  margin-right: 5px;
}

.price .value {
  font-size: 20px;
  font-weight: bold;
  color: #f56c6c;
}

.bids {
  font-size: 13px;
  color: #909399;
}

.item-footer {
  padding-top: 10px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  color: #909399;
  font-size: 13px;
}

.views {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
