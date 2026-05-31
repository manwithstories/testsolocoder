<template>
  <div class="home-page">
    <section class="hero-section">
      <div class="hero-content">
        <h1>探索艺术与历史的数字世界</h1>
        <p>在线预约展览，沉浸式体验博物馆魅力，开启您的文化之旅</p>
        <div class="hero-buttons">
          <el-button type="primary" size="large" @click="$router.push('/exhibitions')">
            浏览展览
          </el-button>
          <el-button size="large" @click="$router.push('/collections')">
            探索藏品
          </el-button>
        </div>
      </div>
    </section>

    <section class="hot-exhibitions">
      <div class="section-header">
        <h2>热门展览</h2>
        <el-button text type="primary" @click="$router.push('/exhibitions')">
          查看全部 <el-icon><ArrowRight /></el-icon>
        </el-button>
      </div>
      <div v-loading="loading" class="exhibition-grid">
        <div
          v-for="exhibition in hotExhibitions"
          :key="exhibition.id"
          class="exhibition-card"
          @click="$router.push(`/exhibitions/${exhibition.id}`)"
        >
          <div class="card-image">
            <img :src="exhibition.image_url || '/placeholder.svg'" :alt="exhibition.title" />
            <div class="card-badge" v-if="exhibition.is_virtual">
              <el-icon><VideoCamera /></el-icon> 虚拟展厅
            </div>
          </div>
          <div class="card-content">
            <h3 class="card-title">{{ exhibition.title }}</h3>
            <p class="card-desc">{{ exhibition.description }}</p>
            <div class="card-meta">
              <span><el-icon><Calendar /></el-icon> {{ formatDate(exhibition.start_date) }} - {{ formatDate(exhibition.end_date) }}</span>
              <span><el-icon><Location /></el-icon> {{ exhibition.location }}</span>
            </div>
            <div class="card-footer">
              <span class="price" v-if="exhibition.ticket_price > 0">
                ¥{{ exhibition.ticket_price }}
              </span>
              <span class="price" v-else>免费</span>
              <span class="views">
                <el-icon><View /></el-icon> {{ exhibition.view_count }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="feature-section">
      <div class="features">
        <div class="feature-item">
          <div class="feature-icon icon-1">
            <el-icon size="32"><Tickets /></el-icon>
          </div>
          <h3>在线预约</h3>
          <p>快速预约展览，选择合适时段，支持多种导览服务</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon icon-2">
            <el-icon size="32"><Microphone /></el-icon>
          </div>
          <h3>智能导览</h3>
          <p>扫码获取藏品讲解，支持语音和文字多语言导览</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon icon-3">
            <el-icon size="32"><Picture /></el-icon>
          </div>
          <h3>高清藏品</h3>
          <p>浏览海量高清藏品图片，深度了解文物背后的故事</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon icon-4">
            <el-icon size="32"><DataAnalysis /></el-icon>
          </div>
          <h3>学术研究</h3>
          <p>学术研究申请，获取藏品高清资料用于研究用途</p>
        </div>
      </div>
    </section>

    <section class="stats-section">
      <div class="stats">
        <div class="stat-item">
          <div class="stat-number">{{ stats.exhibitions }}</div>
          <div class="stat-label">正在展览</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.collections }}</div>
          <div class="stat-label">藏品数量</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.visitors }}</div>
          <div class="stat-label">累计参观</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.researchers }}</div>
          <div class="stat-label">研究人员</div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as exhibitionApi from '@/api/exhibition'
import type { Exhibition } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const hotExhibitions = ref<Exhibition[]>([])
const stats = ref({
  exhibitions: 28,
  collections: 12580,
  visitors: '56.8万',
  researchers: 1200
})

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY.MM.DD')
}

const fetchHotExhibitions = async () => {
  try {
    loading.value = true
    const res = await exhibitionApi.getHotExhibitions()
    hotExhibitions.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchHotExhibitions()
})
</script>

<style scoped lang="scss">
.home-page {
  padding-bottom: 0;
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 80px 20px;
  text-align: center;
  color: #fff;

  .hero-content {
    max-width: 800px;
    margin: 0 auto;

    h1 {
      font-size: 48px;
      margin-bottom: 20px;
    }

    p {
      font-size: 18px;
      opacity: 0.9;
      margin-bottom: 32px;
    }

    .hero-buttons {
      display: flex;
      gap: 16px;
      justify-content: center;
    }
  }
}

.hot-exhibitions {
  max-width: 1400px;
  margin: 0 auto;
  padding: 60px 24px;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 32px;

    h2 {
      font-size: 28px;
    }
  }

  .exhibition-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
  }

  .exhibition-card {
    background: #fff;
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    }

    .card-image {
      position: relative;
      height: 200px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .card-badge {
        position: absolute;
        top: 12px;
        left: 12px;
        background: rgba(64, 158, 255, 0.9);
        color: #fff;
        padding: 4px 12px;
        border-radius: 20px;
        font-size: 12px;
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }

    .card-content {
      padding: 20px;

      .card-title {
        font-size: 18px;
        margin-bottom: 8px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .card-desc {
        color: #606266;
        font-size: 14px;
        line-height: 1.5;
        margin-bottom: 16px;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .card-meta {
        display: flex;
        flex-direction: column;
        gap: 8px;
        margin-bottom: 16px;
        font-size: 13px;
        color: #909399;

        span {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }

      .card-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .price {
          color: #f56c6c;
          font-weight: 600;
          font-size: 16px;
        }

        .views {
          display: flex;
          align-items: center;
          gap: 4px;
          color: #909399;
          font-size: 13px;
        }
      }
    }
  }
}

.feature-section {
  background: #fff;
  padding: 60px 24px;

  .features {
    max-width: 1200px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 40px;

    .feature-item {
      text-align: center;

      .feature-icon {
        width: 80px;
        height: 80px;
        border-radius: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 20px;
        color: #fff;

        &.icon-1 { background: linear-gradient(135deg, #667eea, #764ba2); }
        &.icon-2 { background: linear-gradient(135deg, #f093fb, #f5576c); }
        &.icon-3 { background: linear-gradient(135deg, #4facfe, #00f2fe); }
        &.icon-4 { background: linear-gradient(135deg, #43e97b, #38f9d7); }
      }

      h3 {
        font-size: 18px;
        margin-bottom: 12px;
      }

      p {
        color: #606266;
        font-size: 14px;
        line-height: 1.6;
      }
    }
  }
}

.stats-section {
  background: linear-gradient(135deg, #1f2937 0%, #111827 100%);
  padding: 60px 24px;

  .stats {
    max-width: 1200px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 40px;

    .stat-item {
      text-align: center;
      color: #fff;

      .stat-number {
        font-size: 48px;
        font-weight: 700;
        margin-bottom: 8px;
        background: linear-gradient(135deg, #667eea, #764ba2);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
      }

      .stat-label {
        font-size: 16px;
        color: #9ca3af;
      }
    }
  }
}

@media (max-width: 768px) {
  .features, .stats {
    grid-template-columns: repeat(2, 1fr) !important;
  }

  .hero-content h1 {
    font-size: 32px !important;
  }
}
</style>
