<template>
  <div class="my-bids">
    <h3>我的出价</h3>
    <el-table :data="bids" v-loading="loading">
      <el-table-column label="拍卖品" min-width="200">
        <template #default="{ row }">
          <div class="item-title" @click="$router.push(`/items/${row.auction_item_id}`)">
            <img :src="getImage(row)" />
            <span>{{ row.auction_item?.title }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="出价金额" width="120">
        <template #default="{ row }">¥{{ row.amount.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="is_auto_bid" label="类型" width="80">
        <template #default="{ row }">
          <el-tag size="small" :type="row.is_auto_bid === 1 ? 'info' : 'primary'">
            {{ row.is_auto_bid === 1 ? '自动' : '手动' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_winning" label="状态" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.is_winning === 1" type="success" size="small">领先</el-tag>
          <el-tag v-else type="info" size="small">已被超越</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="出价时间" width="160">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchBids"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import type { Bid } from '@/types'
import { bidApi } from '@/api'

const bids = ref<Bid[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const fetchBids = async () => {
  loading.value = true
  try {
    const res = await bidApi.getMyBids({ page: page.value, page_size: pageSize.value })
    bids.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const getImage = (bid: Bid) => {
  if (bid.auction_item?.images?.[0]?.url) {
    return bid.auction_item.images[0].url
  }
  return 'https://via.placeholder.com/50x50?text=No'
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

onMounted(() => {
  fetchBids()
})
</script>

<style scoped>
h3 {
  margin: 0 0 20px;
}

.item-title {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.item-title img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 4px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
