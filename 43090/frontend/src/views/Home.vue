<template>
  <div class="home">
    <section class="banner">
      <div class="banner-content">
        <h1>在线拍卖竞价管理系统</h1>
        <p>公平、透明、高效的在线拍卖平台</p>
        <el-button type="primary" size="large" @click="$router.push('/items')">立即参与</el-button>
      </div>
    </section>

    <section class="section">
      <div class="section-header">
        <h2>热门拍卖品</h2>
        <el-button link @click="$router.push('/items')">查看更多</el-button>
      </div>
      <el-row :gutter="20">
        <el-col v-for="item in hotItems" :key="item.id" :xs="24" :sm="12" :md="8" :lg="6">
          <ItemCard :item="item" />
        </el-col>
      </el-row>
    </section>

    <section class="section">
      <div class="section-header">
        <h2>正在进行的拍卖会</h2>
        <el-button link @click="$router.push('/sessions')">查看更多</el-button>
      </div>
      <el-row :gutter="20">
        <el-col v-for="session in activeSessions" :key="session.id" :xs="24" :sm="12" :md="8">
          <SessionCard :session="session" />
        </el-col>
      </el-row>
    </section>

    <section class="features">
      <el-row :gutter="40">
        <el-col :xs="24" :sm="8">
          <div class="feature-item">
            <el-icon size="48"><Lock /></el-icon>
            <h3>安全可靠</h3>
            <p>多重安全保障，确保交易安全</p>
          </div>
        </el-col>
        <el-col :xs="24" :sm="8">
          <div class="feature-item">
            <el-icon size="48"><Cpu /></el-icon>
            <h3>实时竞价</h3>
            <p>毫秒级实时响应，不错过任何出价</p>
          </div>
        </el-col>
        <el-col :xs="24" :sm="8">
          <div class="feature-item">
            <el-icon size="48"><Money /></el-icon>
            <h3>公平公正</h3>
            <p>透明的拍卖规则，公平竞争</p>
          </div>
        </el-col>
      </el-row>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { AuctionItem, AuctionSession } from '@/types'
import { itemApi, sessionApi } from '@/api'
import ItemCard from '@/components/ItemCard.vue'
import SessionCard from '@/components/SessionCard.vue'

const hotItems = ref<AuctionItem[]>([])
const activeSessions = ref<AuctionSession[]>([])

const fetchData = async () => {
  try {
    const [itemRes, sessionRes] = await Promise.all([
      itemApi.getList({ page: 1, page_size: 8, sort_by: 'bid_count', sort_order: 'desc' }),
      sessionApi.getActive(),
    ])
    hotItems.value = itemRes.list
    activeSessions.value = sessionRes
  } catch (e) {}
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 80px 20px;
  text-align: center;
  border-radius: 12px;
  margin-bottom: 40px;
}

.banner h1 {
  font-size: 48px;
  margin-bottom: 16px;
}

.banner p {
  font-size: 20px;
  margin-bottom: 30px;
  opacity: 0.9;
}

.section {
  margin-bottom: 50px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.features {
  background: white;
  padding: 60px 20px;
  border-radius: 12px;
  text-align: center;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.feature-item h3 {
  margin: 0;
  color: #303133;
}

.feature-item p {
  margin: 0;
  color: #909399;
}

.feature-item .el-icon {
  color: #409eff;
}
</style>
