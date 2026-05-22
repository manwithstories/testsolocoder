<template>
  <div class="favorites-page">
    <el-card>
      <template #header>
        <el-tabs v-model="activeTab">
          <el-tab-pane label="问题" name="question" />
          <el-tab-pane label="回答" name="answer" />
        </el-tabs>
      </template>

      <el-table :data="favorites" style="width: 100%">
        <el-table-column label="标题" prop="title">
          <template #default="{ row }">
            <span class="title-link" @click="goToDetail(row)">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.targetType === 'question' ? 'primary' : 'success'" size="small">
              {{ row.targetType === 'question' ? '问题' : '回答' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="收藏时间" prop="createdAt" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="danger" size="small" @click="removeFavorite(row)">
              取消收藏
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="favorites.length === 0" description="暂无收藏" />

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchFavorites"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { favoriteApi } from '@/api'
import type { Favorite } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()

const favorites = ref<(Favorite & { title?: string; targetType: string })[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeTab = ref('question')

const fetchFavorites = async () => {
  try {
    const res = await favoriteApi.getUserFavorites({
      page: page.value,
      pageSize: pageSize.value,
      targetType: activeTab.value
    })
    const list = res.data?.list || []
    favorites.value = list.map((f: Favorite) => ({
      ...f,
      title: f.targetType === 'question' ? '问题 #' + f.targetId : '回答 #' + f.targetId
    }))
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const removeFavorite = async (row: any) => {
  try {
    await favoriteApi.removeFavorite({
      targetType: row.targetType,
      targetId: row.targetId
    })
    fetchFavorites()
  } catch (e) {
    console.error(e)
  }
}

const goToDetail = (row: any) => {
  if (row.targetType === 'question') {
    router.push(`/questions/${row.targetId}`)
  }
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

watch(activeTab, () => {
  page.value = 1
  fetchFavorites()
})

onMounted(() => {
  fetchFavorites()
})
</script>

<style scoped lang="scss">
.favorites-page {
  .title-link {
    color: #409eff;
    cursor: pointer;

    &:hover {
      text-decoration: underline;
    }
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
