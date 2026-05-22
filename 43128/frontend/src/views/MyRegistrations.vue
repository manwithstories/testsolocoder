<template>
  <div class="page">
    <div class="page-title">我的报名</div>
    <el-table :data="list" stripe>
      <el-table-column label="赛事">
        <template #default="{ row }">{{ row.event?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="项目">
        <template #default="{ row }">{{ row.event_item?.name || '-' }}</template>
      </el-table-column>
      <el-table-column prop="reg_type" label="类型" width="80">
        <template #default="{ row }">{{ row.reg_type === 'team' ? '团体' : '单人' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <span :class="['tag', statusClass(row.status)]">{{ statusText(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="候补位">
        <template #default="{ row }">
          <span v-if="row.status === 'waitlist'">第 {{ row.queue_position }} 位</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="费用" width="100">
        <template #default="{ row }">¥{{ row.amount }}</template>
      </el-table-column>
      <el-table-column prop="created_at" label="报名时间" />
    </el-table>
    <el-pagination
      style="margin-top:16px"
      layout="prev, pager, next"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="page"
      @current-change="fetch"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { regApi } from '@/api'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

function statusClass(s: string) {
  const map: Record<string, string> = {
    confirmed: 'tag-confirmed', waitlist: 'tag-waitlist', pending: 'tag-pending',
    rejected: 'tag-draft', cancelled: 'tag-draft',
  }
  return map[s] || 'tag-draft'
}
function statusText(s: string) {
  const map: Record<string, string> = {
    confirmed: '已确认', waitlist: '候补', pending: '待审核',
    rejected: '已拒绝', cancelled: '已取消',
  }
  return map[s] || s
}

async function fetch() {
  const res = await regApi.my({ page: page.value, page_size: pageSize.value })
  const data = res.data as any
  list.value = data.list || []
  total.value = data.total || 0
}

onMounted(fetch)
</script>
