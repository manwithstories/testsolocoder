<template>
  <div class="session-list">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent>
        <el-form-item label="搜索">
          <el-input v-model="filters.keyword" placeholder="搜索拍卖会" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable>
            <el-option label="未开始" :value="0" />
            <el-option label="进行中" :value="1" />
            <el-option label="已结束" :value="2" />
            <el-option label="已取消" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchSessions">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-row :gutter="20">
      <el-col v-for="session in sessions" :key="session.id" :xs="24" :sm="12" :md="8">
        <SessionCard :session="session" />
      </el-col>
    </el-row>

    <el-empty v-if="sessions.length === 0" description="暂无拍卖会" />

    <el-pagination
      v-if="total > 0"
      class="pagination"
      v-model:current-page="filters.page"
      v-model:page-size="filters.page_size"
      :total="total"
      :page-sizes="[9, 18, 36]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="fetchSessions"
      @current-change="fetchSessions"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { AuctionSession } from '@/types'
import { sessionApi } from '@/api'
import SessionCard from '@/components/SessionCard.vue'

const sessions = ref<AuctionSession[]>([])
const total = ref(0)

const filters = reactive({
  page: 1,
  page_size: 9,
  keyword: '',
  status: undefined as number | undefined,
})

const fetchSessions = async () => {
  try {
    const res = await sessionApi.getList(filters)
    sessions.value = res.list
    total.value = res.total
  } catch (e) {}
}

const resetFilters = () => {
  filters.page = 1
  filters.keyword = ''
  filters.status = undefined
  fetchSessions()
}

onMounted(() => {
  fetchSessions()
})
</script>

<style scoped>
.filter-card {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
