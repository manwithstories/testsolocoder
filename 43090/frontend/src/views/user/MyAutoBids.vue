<template>
  <div class="my-auto-bids">
    <h3>自动出价管理</h3>
    <el-empty v-if="autoBids.length === 0" description="暂无自动出价设置" />
    <el-table :data="autoBids" v-else>
      <el-table-column label="拍卖品" min-width="200">
        <template #default="{ row }">
          <div class="item-title" @click="$router.push(`/items/${row.auction_item_id}`)">
            <img :src="getImage(row)" />
            <span>{{ row.auction_item?.title }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="current_bid" label="当前出价" width="120">
        <template #default="{ row }">¥{{ row.current_bid.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="max_price" label="最高上限" width="120">
        <template #default="{ row }">¥{{ row.max_price.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '运行中' : '已停止' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="设置时间" width="160">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button link type="danger" size="small" @click="cancelAutoBid(row)">取消</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import type { AutoBid } from '@/types'
import { bidApi } from '@/api'

const autoBids = ref<AutoBid[]>([])

const fetchAutoBids = async () => {
  try {
    const res = await bidApi.getMyAutoBids()
    autoBids.value = res
  } catch (e) {}
}

const getImage = (autoBid: AutoBid) => {
  if (autoBid.auction_item?.images?.[0]?.url) {
    return autoBid.auction_item.images[0].url
  }
  return 'https://via.placeholder.com/50x50?text=No'
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

const cancelAutoBid = async (autoBid: AutoBid) => {
  try {
    await ElMessageBox.confirm('确定要取消这个自动出价吗？', '提示', { type: 'warning' })
    await bidApi.cancelAutoBid(autoBid.id)
    ElMessage.success('已取消')
    fetchAutoBids()
  } catch (e) {}
}

onMounted(() => {
  fetchAutoBids()
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
</style>
