<template>
  <div v-if="session" class="session-detail">
    <el-card class="session-header">
      <div class="header-info">
        <h1>{{ session.name }}</h1>
        <el-tag :type="statusType" size="large">{{ statusText }}</el-tag>
      </div>
      <p class="session-desc">{{ session.description }}</p>
      <div class="session-meta">
        <div class="meta-item">
          <el-icon><Clock /></el-icon>
          <span>开始时间: {{ formatTime(session.start_time) }}</span>
        </div>
        <div class="meta-item">
          <el-icon><Clock /></el-icon>
          <span>结束时间: {{ formatTime(session.end_time) }}</span>
        </div>
        <div class="meta-item">
          <el-icon><Money /></el-icon>
          <span>最小加价: ¥{{ session.min_increment.toFixed(2) }}</span>
        </div>
        <div class="meta-item">
          <el-icon><Clock /></el-icon>
          <span>延时规则: {{ session.extend_time }}秒</span>
        </div>
      </div>
    </el-card>

    <h3>拍卖品列表 ({{ session.auction_items?.length || 0 }})</h3>
    <el-row :gutter="20">
      <el-col
        v-for="itemSession in session.auction_items"
        :key="itemSession.id"
        :xs="24"
        :sm="12"
        :md="8"
        :lg="6"
      >
        <ItemCard :item="itemSession.auction_item!" />
      </el-col>
    </el-row>

    <el-empty v-if="!session.auction_items || session.auction_items.length === 0" description="暂无拍卖品" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import type { AuctionSession } from '@/types'
import { sessionApi } from '@/api'
import ItemCard from '@/components/ItemCard.vue'

const route = useRoute()
const session = ref<AuctionSession | null>(null)

const statusType = computed(() => {
  if (!session.value) return 'info'
  const map: Record<number, string> = { 0: 'info', 1: 'success', 2: 'info', 3: 'danger' }
  return map[session.value.status] || 'info'
})

const statusText = computed(() => {
  if (!session.value) return ''
  const map: Record<number, string> = { 0: '未开始', 1: '进行中', 2: '已结束', 3: '已取消' }
  return map[session.value.status] || ''
})

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

const fetchSession = async () => {
  const id = Number(route.params.id)
  try {
    const res = await sessionApi.getDetail(id)
    session.value = res
  } catch (e) {}
}

onMounted(() => {
  fetchSession()
})
</script>

<style scoped>
.session-header {
  margin-bottom: 30px;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.header-info h1 {
  margin: 0;
  font-size: 28px;
}

.session-desc {
  color: #606266;
  margin-bottom: 20px;
  font-size: 15px;
}

.session-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 30px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
}

h3 {
  margin: 0 0 20px;
  font-size: 20px;
}
</style>
