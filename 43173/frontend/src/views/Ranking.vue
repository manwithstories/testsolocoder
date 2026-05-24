<template>
  <div class="ranking">
    <div class="page-header">
      <h1 class="page-title">排行榜</h1>
    </div>
    
    <el-tabs v-model="activeTab" class="ranking-tabs">
      <el-tab-pane label="日榜" name="daily">
        <el-radio-group v-model="category" size="small" class="category-filter">
          <el-radio-button label="plays">播放量</el-radio-button>
          <el-radio-button label="hot">热度</el-radio-button>
          <el-radio-button label="likes">收藏</el-radio-button>
        </el-radio-group>
      </el-tab-pane>
      <el-tab-pane label="周榜" name="weekly">
        <el-radio-group v-model="category" size="small" class="category-filter">
          <el-radio-button label="plays">播放量</el-radio-button>
          <el-radio-button label="hot">热度</el-radio-button>
          <el-radio-button label="likes">收藏</el-radio-button>
        </el-radio-group>
      </el-tab-pane>
      <el-tab-pane label="月榜" name="monthly">
        <el-radio-group v-model="category" size="small" class="category-filter">
          <el-radio-button label="plays">播放量</el-radio-button>
          <el-radio-button label="hot">热度</el-radio-button>
          <el-radio-button label="likes">收藏</el-radio-button>
        </el-radio-group>
      </el-tab-pane>
    </el-tabs>
    
    <div class="ranking-list" v-loading="loading">
      <div 
        v-for="(item, index) in rankingList" 
        :key="item.work_id"
        class="ranking-item"
        @click="goToWork(item.work_id)"
      >
        <div class="rank-index" :class="{ top: index < 3 }">
          <span v-if="index < 3" class="medal">{{ ['🥇', '🥈', '🥉'][index] }}</span>
          <span v-else>{{ index + 1 }}</span>
        </div>
        <el-image :src="item.cover_url" class="cover" fit="cover">
          <template #error>
            <div class="image-placeholder">
              <el-icon :size="24"><Headset /></el-icon>
            </div>
          </template>
        </el-image>
        <div class="info">
          <div class="title">{{ item.title }}</div>
          <div class="artist">{{ item.artist_name }}</div>
        </div>
        <div class="score">
          <div class="value">{{ formatScore(item.score) }}</div>
          <div class="label">{{ categoryLabel }}</div>
        </div>
      </div>
      
      <el-empty v-if="rankingList.length === 0" description="暂无数据" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { rankingApi } from '@/api/ranking'
import type { RankingItem } from '@/api/ranking'

const router = useRouter()
const loading = ref(false)
const activeTab = ref('daily')
const category = ref('plays')
const rankingList = ref<RankingItem[]>([])

const categoryLabel = ref('播放量')

watch(category, (val) => {
  const labels: Record<string, string> = {
    plays: '播放量',
    hot: '热度',
    likes: '收藏'
  }
  categoryLabel.value = labels[val] || '播放量'
  loadRanking()
})

watch(activeTab, () => {
  loadRanking()
})

onMounted(() => {
  loadRanking()
})

async function loadRanking() {
  loading.value = true
  try {
    const res = await rankingApi.getRanking({
      type: activeTab.value,
      category: category.value,
      limit: 50
    })
    rankingList.value = res as unknown as RankingItem[]
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function formatScore(score: number) {
  if (score >= 10000) {
    return (score / 10000).toFixed(1) + 'w'
  }
  return Math.round(score).toLocaleString()
}

function goToWork(id: number) {
  router.push(`/work/${id}`)
}
</script>

<style scoped lang="scss">
.ranking {
  .category-filter {
    margin-bottom: 20px;
  }
  
  .ranking-list {
    .ranking-item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 16px;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.3s;
      
      &:hover {
        background: rgba(64, 158, 255, 0.05);
      }
      
      .rank-index {
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 18px;
        font-weight: 600;
        color: var(--text-light);
        
        &.top {
          .medal {
            font-size: 24px;
          }
        }
      }
      
      .cover {
        width: 60px;
        height: 60px;
        border-radius: 8px;
        overflow: hidden;
        flex-shrink: 0;
      }
      
      .info {
        flex: 1;
        min-width: 0;
        
        .title {
          font-weight: 500;
          margin-bottom: 4px;
        }
        
        .artist {
          font-size: 13px;
          color: var(--text-light);
        }
      }
      
      .score {
        text-align: right;
        
        .value {
          font-size: 18px;
          font-weight: 600;
          color: var(--primary-color);
        }
        
        .label {
          font-size: 12px;
          color: var(--text-light);
        }
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
