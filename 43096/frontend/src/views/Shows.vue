<template>
  <Layout>
    <div class="shows-page">
      <div class="page-header">
        <h2>热门演出</h2>
        <div class="search-bar">
          <el-input
            v-model="keyword"
            placeholder="搜索演出名称或艺人"
            clearable
            @keyup.enter="fetchShows"
          >
            <template #append>
              <el-button @click="fetchShows">搜索</el-button>
            </template>
          </el-input>
        </div>
      </div>

      <div class="shows-grid">
        <el-card
          v-for="show in shows"
          :key="show.id"
          class="show-card"
          shadow="hover"
          @click="goToDetail(show.id)"
        >
          <img :src="show.poster || 'https://picsum.photos/400/250?random=' + show.id" class="show-poster" />
          <div class="show-info">
            <h3 class="show-name">{{ show.name }}</h3>
            <p class="show-artist">{{ show.artist }}</p>
            <p class="show-venue">{{ show.venue }}</p>
            <div class="show-sessions">
              <el-tag v-if="show.sessions && show.sessions.length > 0" type="info">
                {{ show.sessions.length }} 场次</el-tag>
              <el-tag type="success">可购票</el-tag>
            </div>
          </div>
        </el-card>
      </div>

      <el-pagination
        v-if="total > 0"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        class="pagination"
        @current-change="fetchShows"
      />

      <el-empty v-if="shows.length === 0" description="暂无演出数据" />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showApi } from '@/api'
import Layout from '@/components/Layout.vue'
import type { Show } from '@/types'

const router = useRouter()
const shows = ref<Show[]>([])
const keyword = ref('')
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)

async function fetchShows() {
  try {
    const res = await showApi.list({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value,
      status: 1
    })
    shows.value = res.list
    total.value = res.pagination.total
  } catch (err) {
    console.error(err)
  }
}

function goToDetail(id: number) {
  router.push(`/show/${id}`)
}

onMounted(() => {
  fetchShows()
})
</script>

<style lang="scss" scoped>
.shows-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;

  h2 {
    margin: 0 0 16px 0;
  }

  .search-bar {
    max-width: 400px;
  }
}

.shows-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.show-card {
  cursor: pointer;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-4px);
  }

  .show-poster {
    width: 100%;
    height: 180px;
    object-fit: cover;
    border-radius: 8px 8px 0 0;
  }

  .show-info {
    padding: 12px 0;

    .show-name {
      margin: 0 0 8px 0;
      font-size: 16px;
      font-weight: 600;
    }

    .show-artist {
      margin: 0 0 4px 0;
      color: #666;
      font-size: 14px;
    }

    .show-venue {
      margin: 0 0 8px 0;
      color: #999;
      font-size: 13px;
    }

    .show-sessions {
      display: flex;
      gap: 8px;
    }
  }
}

.pagination {
  margin-top: 30px;
  justify-content: center;
}
</style>
