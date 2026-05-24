<template>
  <div class="artists">
    <div class="page-header">
      <h1 class="page-title">音乐人</h1>
      <div class="search-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索音乐人"
          clearable
          :prefix-icon="Search"
          @keyup.enter="search"
        />
        <el-button type="primary" @click="search">搜索</el-button>
      </div>
    </div>
    
    <div class="artists-grid" v-loading="loading">
      <div 
        v-for="artist in artists" 
        :key="artist.id"
        class="artist-item"
        @click="goToDetail(artist.id)"
      >
        <el-avatar :size="100" :src="artist.avatar">
          {{ artist.nickname?.charAt(0) || 'A' }}
        </el-avatar>
        <div class="info">
          <div class="name">{{ artist.nickname }}</div>
          <div class="desc text-ellipsis">
            {{ artist.artist_info?.artist_name || '独立音乐人' }}
          </div>
          <div class="stats">
            <span>{{ formatCount(artist.artist_info?.total_followers || 0) }} 粉丝</span>
            <span>{{ formatCount(artist.artist_info?.total_works || 0) }} 作品</span>
          </div>
        </div>
      </div>
      
      <el-empty v-if="artists.length === 0 && !loading" description="暂无音乐人" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 24, 48]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadArtists"
        @size-change="loadArtists"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { userApi } from '@/api/auth'
import type { UserInfo } from '@/api/auth'

const router = useRouter()
const loading = ref(false)
const keyword = ref('')
const artists = ref<UserInfo[]>([])
const page = ref(1)
const pageSize = ref(24)
const total = ref(0)

onMounted(() => {
  loadArtists()
})

async function loadArtists() {
  loading.value = true
  try {
    const res = await userApi.list({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value,
      role: 'artist'
    })
    artists.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  loadArtists()
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}

function goToDetail(id: number) {
  router.push(`/artist/${id}`)
}
</script>

<style scoped lang="scss">
.artists {
  .search-bar {
    display: flex;
    gap: 12px;
  }
  
  .artists-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    
    .artist-item {
      text-align: center;
      cursor: pointer;
      padding: 24px;
      border-radius: 8px;
      transition: all 0.3s;
      
      &:hover {
        background: rgba(64, 158, 255, 0.05);
        transform: translateY(-5px);
        
        .el-avatar {
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
        }
      }
      
      .info {
        margin-top: 16px;
        
        .name {
          font-size: 16px;
          font-weight: 500;
          margin-bottom: 4px;
        }
        
        .desc {
          font-size: 13px;
          color: var(--text-light);
          margin-bottom: 8px;
        }
        
        .stats {
          font-size: 12px;
          color: var(--text-light);
          display: flex;
          justify-content: center;
          gap: 12px;
        }
      }
    }
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 30px;
  }
}
</style>
