<template>
  <div class="home">
    <div class="banner">
      <div class="banner-content">
        <h1>发现独立音乐</h1>
        <p>探索独立音乐人的世界，发现不一样的声音</p>
        <el-button type="primary" size="large" @click="goToRanking">
          探索排行榜
        </el-button>
      </div>
    </div>
    
    <div class="section">
      <div class="section-header">
        <h2>热门排行</h2>
        <router-link to="/ranking" class="more">查看更多</router-link>
      </div>
      <div class="ranking-grid" v-loading="loading">
        <div 
          v-for="(item, index) in hotRanking.slice(0, 8)" 
          :key="item.work_id"
          class="ranking-item"
          @click="goToWork(item.work_id)"
        >
          <div class="rank" :class="{ top: index < 3 }">{{ index + 1 }}</div>
          <el-image :src="item.cover_url" class="cover" fit="cover">
            <template #error>
              <div class="image-placeholder">
                <el-icon :size="40"><Headset /></el-icon>
              </div>
            </template>
          </el-image>
          <div class="info">
            <div class="title text-ellipsis">{{ item.title }}</div>
            <div class="artist text-ellipsis">{{ item.artist_name }}</div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="section">
      <div class="section-header">
        <h2>新歌上架</h2>
        <router-link to="/works" class="more">查看更多</router-link>
      </div>
      <div class="works-grid" v-loading="loading">
        <div 
          v-for="work in newWorks" 
          :key="work.id"
          class="work-item"
          @click="goToWork(work.id)"
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
          </div>
        </div>
      </div>
    </div>
    
    <div class="section">
      <div class="section-header">
        <h2>推荐音乐人</h2>
        <router-link to="/artists" class="more">查看更多</router-link>
      </div>
      <div class="artists-grid" v-loading="loading">
        <div 
          v-for="artist in artists" 
          :key="artist.id"
          class="artist-item"
          @click="goToArtist(artist.id)"
        >
          <el-avatar :size="80" :src="artist.avatar">
            {{ artist.nickname?.charAt(0) || 'A' }}
          </el-avatar>
          <div class="name">{{ artist.nickname }}</div>
          <div class="desc text-ellipsis">{{ artist.artist_info?.artist_name || artist.bio || '独立音乐人' }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { rankingApi } from '@/api/ranking'
import { workApi } from '@/api/work'
import { userApi } from '@/api/auth'
import type { RankingItem } from '@/api/ranking'
import type { Work } from '@/api/work'
import type { UserInfo } from '@/api/auth'

const router = useRouter()
const loading = ref(false)
const hotRanking = ref<RankingItem[]>([])
const newWorks = ref<Work[]>([])
const artists = ref<UserInfo[]>([])

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    await Promise.all([
      loadHotRanking(),
      loadNewWorks(),
      loadArtists()
    ])
  } finally {
    loading.value = false
  }
}

async function loadHotRanking() {
  try {
    const res = await rankingApi.getHotRanking({ limit: 10 })
    hotRanking.value = res as unknown as RankingItem[]
  } catch (e) {
    console.error(e)
  }
}

async function loadNewWorks() {
  try {
    const res = await workApi.list({ page: 1, page_size: 8 })
    newWorks.value = res.list
  } catch (e) {
    console.error(e)
  }
}

async function loadArtists() {
  try {
    const res = await userApi.list({ page: 1, page_size: 8, role: 'artist' })
    artists.value = res.list
  } catch (e) {
    console.error(e)
  }
}

function goToRanking() {
  router.push('/ranking')
}

function goToWork(id: number) {
  router.push(`/work/${id}`)
}

function goToArtist(id: number) {
  router.push(`/artist/${id}`)
}
</script>

<style scoped lang="scss">
.home {
  .banner {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    padding: 60px 40px;
    color: #fff;
    margin-bottom: 30px;
    
    .banner-content {
      max-width: 600px;
      
      h1 {
        font-size: 36px;
        margin-bottom: 16px;
      }
      
      p {
        font-size: 18px;
        margin-bottom: 24px;
        opacity: 0.9;
      }
    }
  }
  
  .section {
    margin-bottom: 40px;
    
    .section-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 20px;
      
      h2 {
        font-size: 24px;
        font-weight: 600;
      }
      
      .more {
        color: var(--primary-color);
        font-size: 14px;
      }
    }
  }
  
  .ranking-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    
    .ranking-item {
      position: relative;
      cursor: pointer;
      transition: transform 0.3s;
      
      &:hover {
        transform: translateY(-5px);
        
        .cover {
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
        }
      }
      
      .rank {
        position: absolute;
        top: 8px;
        left: 8px;
        width: 28px;
        height: 28px;
        background: rgba(0, 0, 0, 0.6);
        color: #fff;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 600;
        z-index: 1;
        
        &.top {
          background: #f56c6c;
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
        }
      }
    }
  }
  
  .works-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
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
        }
      }
    }
  }
  
  .artists-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    
    .artist-item {
      text-align: center;
      cursor: pointer;
      padding: 20px;
      border-radius: 8px;
      transition: background 0.3s;
      
      &:hover {
        background: rgba(64, 158, 255, 0.05);
      }
      
      .name {
        margin-top: 12px;
        font-weight: 500;
      }
      
      .desc {
        margin-top: 4px;
        font-size: 13px;
        color: var(--text-light);
      }
    }
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
