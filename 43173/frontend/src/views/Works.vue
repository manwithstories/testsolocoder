<template>
  <div class="works">
    <div class="page-header">
      <h1 class="page-title">作品</h1>
      <div class="search-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索作品"
          clearable
          :prefix-icon="Search"
          @keyup.enter="search"
        />
        <el-select v-model="genre" placeholder="风格" clearable class="genre-select">
          <el-option v-for="g in genres" :key="g" :label="g" :value="g" />
        </el-select>
        <el-button type="primary" @click="search">搜索</el-button>
      </div>
    </div>
    
    <div class="works-grid" v-loading="loading">
      <div 
        v-for="work in works" 
        :key="work.id"
        class="work-item"
        @click="goToDetail(work.id)"
      >
        <el-image :src="work.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="40"><Headset /></el-icon>
            </div>
          </template>
        </el-image>
        <div class="info">
          <div class="title text-ellipsis">{{ work.title }}</div>
          <div class="artist text-ellipsis">{{ work.artist_name }}</div>
          <div class="meta">
            <span><el-icon><Headset /></el-icon> {{ formatCount(work.play_count) }}</span>
          </div>
        </div>
      </div>
      
      <el-empty v-if="works.length === 0 && !loading" description="暂无作品" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 24, 48, 96]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadWorks"
        @size-change="loadWorks"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { workApi } from '@/api/work'
import type { Work } from '@/api/work'

const router = useRouter()
const loading = ref(false)
const keyword = ref('')
const genre = ref('')
const genres = ['流行', '摇滚', '电子', '民谣', '说唱', 'R&B', '古典', '爵士', '乡村', '金属', '朋克', '雷鬼']
const works = ref<Work[]>([])
const page = ref(1)
const pageSize = ref(24)
const total = ref(0)

onMounted(() => {
  loadWorks()
})

async function loadWorks() {
  loading.value = true
  try {
    const res = await workApi.list({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value,
      genre: genre.value,
      status: 2
    })
    works.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  loadWorks()
}

function formatCount(count: number) {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  return count
}

function goToDetail(id: number) {
  router.push(`/work/${id}`)
}
</script>

<style scoped lang="scss">
.works {
  .search-bar {
    display: flex;
    gap: 12px;
    
    .genre-select {
      width: 150px;
    }
  }
  
  .works-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 20px;
    
    .work-item {
      cursor: pointer;
      transition: transform 0.3s;
      
      &:hover {
        transform: translateY(-5px);
        
        .cover {
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
        }
      }
      
      .cover {
        width: 100%;
        aspect-ratio: 1;
        border-radius: 8px;
        overflow: hidden;
        transition: box-shadow 0.3s;
      }
      
      .info {
        padding: 12px 4px;
        
        .title {
          font-weight: 500;
          margin-bottom: 4px;
        }
        
        .artist {
          font-size: 13px;
          color: var(--text-light);
          margin-bottom: 4px;
        }
        
        .meta {
          font-size: 12px;
          color: var(--text-light);
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 30px;
  }
  
  .image-placeholder {
    width: 100%;
    height: 100%;
    background: #f5f7fa;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #c0c4cc;
  }
}
</style>
